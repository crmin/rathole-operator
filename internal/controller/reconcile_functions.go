package controller

import (
	"context"
	"fmt"
	ratholev1alpha1 "github.com/crmin/rathole-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

type Reconciler interface {
	Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error)
	Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error
	List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error
	Status() client.StatusWriter
	Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error
	Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error
}

func ReconcileServer(r Reconciler, ctx context.Context, server *ratholev1alpha1.RatholeServer) error {
	if server.Status.Condition.ObservedGeneration == server.Generation { // Skip if already reconciled. Spec hasn't changed
		return nil
	}

	// Get services linked to this server
	var services ratholev1alpha1.RatholeServiceList
	if err := r.List(ctx, &services, client.InNamespace(server.Namespace)); err != nil {
		return err
	}

	server.Spec.Services = make(map[string]*ratholev1alpha1.RatholeServiceSpec)

	for _, service := range services.Items {
		fmt.Printf("############ service=%s\n", service.ObjectMeta.Name)
		if service.Spec.ServerRef.Name != server.ObjectMeta.Name {
			fmt.Printf("service ref=%s, server=%s\n", service.Spec.ServerRef.Name, server.ObjectMeta.Name)
			continue
		}
		// remove client options
		service.Spec.LocalAddr = ""
		service.Spec.RetryInterval = 0

		server.Spec.Services[service.ObjectMeta.Name] = &service.Spec
	}

	if len(server.Spec.Services) == 0 {
		// Create dummy service
		serviceSpec := &ratholev1alpha1.RatholeServiceSpec{
			ServerRef: ratholev1alpha1.RatholeServiceResourceRef{
				Name: server.ObjectMeta.Name,
			},
			Type:     "tcp",
			Token:    "default",
			BindAddr: "0.0.0.0:3030",
		}
		server.Spec.Services["dummy"] = serviceSpec
	}

	// Set default token if set .Spec.DefaultTokenFrom and .Spec.DefaultToken is empty
	// TODO: If both .Spec.DefaultTokenFrom and .Spec.DefaultToken are set, an error should occur through the webhook validate
	if server.Spec.DefaultToken == "" {
		if server.Spec.DefaultTokenFrom.ConfigMapRef.Name != "" { // Use ConfigMap
			var (
				configMap corev1.ConfigMap
				ok        bool
			)
			if err := r.Get(ctx, client.ObjectKey{Namespace: server.Namespace, Name: server.Spec.DefaultTokenFrom.ConfigMapRef.Name}, &configMap); err != nil {
				return err
			}
			server.Spec.DefaultToken, ok = configMap.Data[server.Spec.DefaultTokenFrom.ConfigMapRef.Key]
			if !ok {
				return fmt.Errorf("key %s not found in configmap %s", server.Spec.DefaultTokenFrom.ConfigMapRef.Key, server.Spec.DefaultTokenFrom.ConfigMapRef.Name)
			}
		} else if server.Spec.DefaultTokenFrom.SecretRef.Name != "" { // Use Secret
			var (
				secret        corev1.Secret
				secretContent []byte
				ok            bool
			)
			if err := r.Get(ctx, client.ObjectKey{Namespace: server.Namespace, Name: server.Spec.DefaultTokenFrom.SecretRef.Name}, &secret); err != nil {
				return err
			}
			secretContent, ok = secret.Data[server.Spec.DefaultTokenFrom.SecretRef.Key]
			if !ok {
				return fmt.Errorf("key %s not found in secret %s", server.Spec.DefaultTokenFrom.SecretRef.Key, server.Spec.DefaultTokenFrom.SecretRef.Name)
			}
			server.Spec.DefaultToken = string(secretContent)
		}
	}

	server.Status.Condition.ObservedGeneration = server.Generation
	server.Status.Condition.Status = "Synced"
	server.Status.Condition.Reason = "Reconciled"
	// TODO: update last synced time
	if err := r.Status().Update(ctx, server); err != nil {
		return err
	}

	tomlParent := "server"
	config, err := ConvertSpecToToml(&tomlParent, server.Spec)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Config: %s\n", config)

	// upsert config to configTarget
	configResourceType := strings.ToLower(server.Spec.ConfigTarget.ResourceType)

	ownerRefs := []metav1.OwnerReference{
		{
			APIVersion: server.APIVersion,
			Kind:       server.Kind,
			Name:       server.Name,
			UID:        server.UID,
		},
	}

	switch configResourceType {
	case "secret":
		var (
			secret corev1.Secret
			ok     = true
		)

		// search secret
		if err := r.Get(ctx, client.ObjectKey{Namespace: server.Namespace, Name: server.Spec.ConfigTarget.Name}, &secret); err != nil {
			if client.IgnoreNotFound(err) != nil {
				return err
			}
			ok = false
		}

		// if not exist, create secret
		if !ok {
			secret = corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:            server.Spec.ConfigTarget.Name,
					Namespace:       server.Namespace,
					OwnerReferences: ownerRefs,
				},
				StringData: map[string]string{
					"config.toml": config,
				},
			}
			if err := r.Create(ctx, &secret); err != nil {
				return err
			}
		} else {
			secret.StringData["config.toml"] = config
			if err := r.Update(ctx, &secret); err != nil {
				return err
			}
		}

	case "configmap":
		var (
			configMap corev1.ConfigMap

			ok = true
		)

		// search configmap
		if err := r.Get(ctx, client.ObjectKey{Namespace: server.Namespace, Name: server.Spec.ConfigTarget.Name}, &configMap); err != nil {
			if client.IgnoreNotFound(err) != nil {
				return err
			}
			ok = false
		}

		// if not exist, create configmap
		if !ok {
			configMap = corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:            server.Spec.ConfigTarget.Name,
					Namespace:       server.Namespace,
					OwnerReferences: ownerRefs,
				},
				Data: map[string]string{
					"config.toml": config,
				},
			}
			if err := r.Create(ctx, &configMap); err != nil {
				return err
			}
		} else {
			configMap.Data["config.toml"] = config
			if err := r.Update(ctx, &configMap); err != nil {
				return err
			}
		}
	}
	return nil
}

func ReconcileClient(r Reconciler, ctx context.Context, client_ *ratholev1alpha1.RatholeClient) error {
	if client_.Status.Condition.ObservedGeneration == client_.Generation { // Skip if already reconciled. Spec hasn't changed
		return nil
	}

	// Get services linked to this server
	var services ratholev1alpha1.RatholeServiceList
	if err := r.List(ctx, &services, client.InNamespace(client_.Namespace)); err != nil {
		return err
	}

	client_.Spec.Services = make(map[string]*ratholev1alpha1.RatholeServiceSpec)

	for _, service := range services.Items {
		for _, clientRef := range service.Spec.ClientRefs {
			if clientRef.Name != client_.ObjectMeta.Name {
				continue
			}
			// remove server options
			service.Spec.BindAddr = ""

			client_.Spec.Services[service.ObjectMeta.Name] = &service.Spec
		}
	}

	// NOTE: Client does not need default service
	//if len(client_.Spec.Services) == 0 {
	//	// Create dummy service
	//	serviceSpec := &ratholev1alpha1.RatholeServiceSpec{
	//		ServerRef: ratholev1alpha1.RatholeServiceResourceRef{
	//			Name: client_.ObjectMeta.Name,
	//		},
	//		Type:     "tcp",
	//		Token:    "default",
	//		BindAddr: "0.0.0.0:3030",
	//	}
	//	client_.Spec.Services["dummy"] = serviceSpec
	//}

	// Set default token if set .Spec.DefaultTokenFrom and .Spec.DefaultToken is empty
	// TODO: If both .Spec.DefaultTokenFrom and .Spec.DefaultToken are set, an error should occur through the webhook validate
	if client_.Spec.DefaultToken == "" {
		if client_.Spec.DefaultTokenFrom.ConfigMapRef.Name != "" { // Use ConfigMap
			var (
				configMap corev1.ConfigMap
				ok        bool
			)
			if err := r.Get(ctx, client.ObjectKey{Namespace: client_.Namespace, Name: client_.Spec.DefaultTokenFrom.ConfigMapRef.Name}, &configMap); err != nil {
				return err
			}
			client_.Spec.DefaultToken, ok = configMap.Data[client_.Spec.DefaultTokenFrom.ConfigMapRef.Key]
			if !ok {
				return fmt.Errorf("key %s not found in configmap %s", client_.Spec.DefaultTokenFrom.ConfigMapRef.Key, client_.Spec.DefaultTokenFrom.ConfigMapRef.Name)
			}
		} else if client_.Spec.DefaultTokenFrom.SecretRef.Name != "" { // Use Secret
			var (
				secret        corev1.Secret
				secretContent []byte
				ok            bool
			)
			if err := r.Get(ctx, client.ObjectKey{Namespace: client_.Namespace, Name: client_.Spec.DefaultTokenFrom.SecretRef.Name}, &secret); err != nil {
				return err
			}
			secretContent, ok = secret.Data[client_.Spec.DefaultTokenFrom.SecretRef.Key]
			if !ok {
				return fmt.Errorf("key %s not found in secret %s", client_.Spec.DefaultTokenFrom.SecretRef.Key, client_.Spec.DefaultTokenFrom.SecretRef.Name)
			}
			client_.Spec.DefaultToken = string(secretContent)
		}
	}

	client_.Status.Condition.ObservedGeneration = client_.Generation
	client_.Status.Condition.Status = "Synced"
	client_.Status.Condition.Reason = "Reconciled"
	// TODO: update last synced time
	if err := r.Status().Update(ctx, client_); err != nil {
		return err
	}

	tomlParent := "client"
	config, err := ConvertSpecToToml(&tomlParent, client_.Spec)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Config: %s\n", config)

	// upsert config to configTarget
	configResourceType := strings.ToLower(client_.Spec.ConfigTarget.ResourceType)

	ownerRefs := []metav1.OwnerReference{
		{
			APIVersion: client_.APIVersion,
			Kind:       client_.Kind,
			Name:       client_.Name,
			UID:        client_.UID,
		},
	}

	switch configResourceType {
	case "secret":
		var (
			secret corev1.Secret
			ok     = true
		)

		// search secret
		if err := r.Get(ctx, client.ObjectKey{Namespace: client_.Namespace, Name: client_.Spec.ConfigTarget.Name}, &secret); err != nil {
			if client.IgnoreNotFound(err) != nil {
				return err
			}
			ok = false
		}

		// if not exist, create secret
		if !ok {
			secret = corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:            client_.Spec.ConfigTarget.Name,
					Namespace:       client_.Namespace,
					OwnerReferences: ownerRefs,
				},
				StringData: map[string]string{
					"config.toml": config,
				},
			}
			if err := r.Create(ctx, &secret); err != nil {
				return err
			}
		} else {
			secret.StringData["config.toml"] = config
			if err := r.Update(ctx, &secret); err != nil {
				return err
			}
		}

	case "configmap":
		var (
			configMap corev1.ConfigMap

			ok = true
		)

		// search configmap
		if err := r.Get(ctx, client.ObjectKey{Namespace: client_.Namespace, Name: client_.Spec.ConfigTarget.Name}, &configMap); err != nil {
			if client.IgnoreNotFound(err) != nil {
				return err
			}
			ok = false
		}

		// if not exist, create configmap
		if !ok {
			configMap = corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:            client_.Spec.ConfigTarget.Name,
					Namespace:       client_.Namespace,
					OwnerReferences: ownerRefs,
				},
				Data: map[string]string{
					"config.toml": config,
				},
			}
			if err := r.Create(ctx, &configMap); err != nil {
				return err
			}
		} else {
			configMap.Data["config.toml"] = config
			if err := r.Update(ctx, &configMap); err != nil {
				return err
			}
		}
	}
	return nil
}
