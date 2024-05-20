package controller

import (
	"context"
	"encoding/base64"
	"fmt"
	ratholev1alpha1 "github.com/crmin/rathole-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

const ratholeSecretRoot = "/var/run/secrets/rathole"

type Reconciler interface {
	Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error)
	Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error
	List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error
	Status() client.StatusWriter
	Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error
	Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error
}

func ReconcileServer(r Reconciler, ctx context.Context, server *ratholev1alpha1.RatholeServer) error {
	// Get services linked to this server
	var services ratholev1alpha1.RatholeServiceList
	if err := r.List(ctx, &services, client.InNamespace(server.Namespace)); err != nil {
		return err
	}

	server.Spec.Services = make(map[string]*ratholev1alpha1.RatholeServiceSpec)

	for _, service := range services.Items {
		if service.Spec.ServerRef.Name != server.ObjectMeta.Name {
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
	if server.Spec.DefaultToken == "" {
		var err error
		if server.Spec.DefaultToken, err = ReadConfig(r, ctx, server.Namespace, server.Spec.DefaultTokenFrom); err != nil {
			return err
		}
	}

	// Read PKCS12 content from PKCS12From field and set to PKCS12 field
	if server.Spec.Transport.Type == "tls" {
		var (
			pkcsContent string
			err         error
		)
		//pkcs12 -> filename
		fileDir := fmt.Sprintf("%s/%s/%s", ratholeSecretRoot, server.Namespace, server.Name)
		filePath := fmt.Sprintf("%s/%s", fileDir, server.Spec.Transport.TLS.PKCS12From.SecretRef.Name)
		if pkcsContent, err = ReadConfig(r, ctx, server.Namespace, server.Spec.Transport.TLS.PKCS12From); err != nil {
			return err
		}

		// Create directory if not exist
		if _, err := os.Stat(fileDir); os.IsNotExist(err) {
			if err := os.MkdirAll(fileDir, 0755); err != nil {
				return err
			}
		}

		// Write pkcs12 content to file
		if err := os.WriteFile(filePath, []byte(pkcsContent), 0644); err != nil {
			return err
		}

		server.Spec.Transport.TLS.PKCS12 = filePath

		if server.Spec.Transport.TLS.PKCS12Password == "" {
			if server.Spec.Transport.TLS.PKCS12Password, err = ReadConfig(r, ctx, server.Namespace, server.Spec.Transport.TLS.PKCS12PasswordFrom); err != nil {
				return err
			}
		}
	} else if server.Spec.Transport.Type == "noise" {
		var err error
		if server.Spec.Transport.Noise.LocalPrivateKey == "" && (server.Spec.Transport.Noise.LocalPrivateKeyFrom.ConfigMapRef.Name != "" || server.Spec.Transport.Noise.LocalPrivateKeyFrom.SecretRef.Name != "") {
			if server.Spec.Transport.Noise.LocalPrivateKey, err = ReadConfig(r, ctx, server.Namespace, server.Spec.Transport.Noise.LocalPrivateKeyFrom); err != nil {
				return err
			}
		}
		if server.Spec.Transport.Noise.RemotePublicKey == "" && (server.Spec.Transport.Noise.RemotePublicKeyFrom.ConfigMapRef.Name != "" || server.Spec.Transport.Noise.RemotePublicKeyFrom.SecretRef.Name != "") {
			if server.Spec.Transport.Noise.RemotePublicKey, err = ReadConfig(r, ctx, server.Namespace, server.Spec.Transport.Noise.RemotePublicKeyFrom); err != nil {
				return err
			}
		}

		// If exist localPrivateKey and remotePublicKey, encode to base64
		if server.Spec.Transport.Noise.LocalPrivateKey != "" {
			server.Spec.Transport.Noise.EncodedLocalPrivateKey = base64.StdEncoding.EncodeToString([]byte(server.Spec.Transport.Noise.LocalPrivateKey))
		}
		if server.Spec.Transport.Noise.RemotePublicKey != "" {
			server.Spec.Transport.Noise.EncodedRemotePublicKey = base64.StdEncoding.EncodeToString([]byte(server.Spec.Transport.Noise.RemotePublicKey))
		}
	}

	server.Status.Condition.ObservedGeneration = server.Generation
	server.Status.Condition.Status = "Synced"
	server.Status.Condition.Reason = "Reconciled"
	// TODO: update last synced time
	if err := r.Status().Update(ctx, server.DeepCopy()); err != nil {
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
			secret.StringData = make(map[string]string)
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
	// Get services linked to this server
	var services ratholev1alpha1.RatholeServiceList
	if err := r.List(ctx, &services, client.InNamespace(client_.Namespace)); err != nil {
		return err
	}

	client_.Spec.Services = make(map[string]*ratholev1alpha1.RatholeServiceSpec)

	for _, service := range services.Items {
		if service.Spec.ClientRef.Name != client_.ObjectMeta.Name {
			continue
		}
		// remove server options
		service.Spec.BindAddr = ""

		client_.Spec.Services[service.ObjectMeta.Name] = &service.Spec
		// POINT-1
	}

	if len(client_.Spec.Services) == 0 {
		// Create dummy service
		serviceSpec := &ratholev1alpha1.RatholeServiceSpec{
			ServerRef: ratholev1alpha1.RatholeServiceResourceRef{
				Name: client_.ObjectMeta.Name,
			},
			Type:      "tcp",
			Token:     "default",
			LocalAddr: "127.0.0.1:3030",
		}
		client_.Spec.Services["dummy"] = serviceSpec
	}

	// Set default token if set .Spec.DefaultTokenFrom and .Spec.DefaultToken is empty
	// TODO: If both .Spec.DefaultTokenFrom and .Spec.DefaultToken are set, an error should occur through the webhook validate
	if client_.Spec.DefaultToken == "" {
		var err error
		if client_.Spec.DefaultToken, err = ReadConfig(r, ctx, client_.Namespace, client_.Spec.DefaultTokenFrom); err != nil {
			return err
		}
	}

	client_.Status.Condition.ObservedGeneration = client_.Generation
	client_.Status.Condition.Status = "Synced"
	client_.Status.Condition.Reason = "Reconciled"
	// TODO: update last synced time
	if err := r.Status().Update(ctx, client_.DeepCopy()); err != nil {
		return err
	}
	// POINT-2

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
			secret.StringData = make(map[string]string)
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

func ReconcileService(r Reconciler, ctx context.Context, service *ratholev1alpha1.RatholeService) error {
	service.Status.Condition.ObservedGeneration = service.Generation
	service.Status.Condition.Status = "Synced"
	service.Status.Condition.Reason = "Reconciled"

	if err := r.Status().Update(ctx, service.DeepCopy()); err != nil {
		return err
	}

	if service.Spec.ServerRef.Name != "" {
		var server ratholev1alpha1.RatholeServer
		if err := r.Get(ctx, client.ObjectKey{Namespace: service.Namespace, Name: service.Spec.ServerRef.Name}, &server); err != nil {
			return err
		}
		if err := ReconcileServer(r, ctx, &server); err != nil {
			return err
		}
	}

	if service.Spec.ClientRef.Name != "" {
		var client_ ratholev1alpha1.RatholeClient
		if err := r.Get(ctx, client.ObjectKey{Namespace: service.Namespace, Name: service.Spec.ClientRef.Name}, &client_); err != nil {
			return err
		}
		if err := ReconcileClient(r, ctx, &client_); err != nil {
			return err
		}
	}

	return nil
}

func ReadConfig(r Reconciler, ctx context.Context, namespace string, resourceFrom ratholev1alpha1.ResourceFrom) (string, error) {
	if resourceFrom.SecretRef.Name != "" {
		var (
			secret        corev1.Secret
			secretContent []byte
			ok            bool
		)
		if err := r.Get(ctx, client.ObjectKey{Namespace: namespace, Name: resourceFrom.SecretRef.Name}, &secret); err != nil {
			return "", err
		}
		secretContent, ok = secret.Data[resourceFrom.SecretRef.Key]
		if !ok {
			return "", fmt.Errorf("key %s not found in secret %s", resourceFrom.SecretRef.Key, resourceFrom.SecretRef.Name)
		}
		return string(secretContent), nil
	} else if resourceFrom.ConfigMapRef.Name != "" {
		var (
			configMap     corev1.ConfigMap
			configContent string
			ok            bool
		)
		if err := r.Get(ctx, client.ObjectKey{Namespace: namespace, Name: resourceFrom.ConfigMapRef.Name}, &configMap); err != nil {
			return "", err
		}
		configContent, ok = configMap.Data[resourceFrom.ConfigMapRef.Key]
		if !ok {
			return "", fmt.Errorf("key %s not found in configmap %s", resourceFrom.ConfigMapRef.Key, resourceFrom.ConfigMapRef.Name)
		}
		return configContent, nil
	}
	return "", fmt.Errorf("resourceFrom is not set")
}
