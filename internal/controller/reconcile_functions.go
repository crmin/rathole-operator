package controller

import (
	"context"
	"encoding/base64"
	"fmt"
	ratholev1alpha1 "github.com/crmin/rathole-operator/api/v1alpha1"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
	"strings"
	"time"
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
	if err := ValidateServer(server); err != nil {
		return err
	}

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
			err      error
			filePath string
		)
		//pkcs12 -> filename
		if server.Spec.Transport.TLS.PKCS12From.SecretRef.Name != "" {
			filePath = fmt.Sprintf("%s/%s", ratholeSecretRoot, server.Spec.Transport.TLS.PKCS12From.SecretRef.Name)
		} else if server.Spec.Transport.TLS.PKCS12From.ConfigMapRef.Name != "" {
			filePath = fmt.Sprintf("%s/%s", ratholeSecretRoot, server.Spec.Transport.TLS.PKCS12From.ConfigMapRef.Name)
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

	// If not set .Spec.ConfigTarget, set random name of secret
	if server.Spec.ConfigTarget.ResourceType == "" {
		server.Spec.ConfigTarget.ResourceType = "secret"
	}
	if server.Spec.ConfigTarget.Name == "" {
		server.Spec.ConfigTarget.Name = fmt.Sprintf("rathole-server-config-%s", GetSuffix(5))
	}

	server.Status.Condition.ObservedGeneration = server.Generation
	server.Status.Condition.Status = "Synced"
	server.Status.Condition.Reason = "Reconciled"
	server.Status.Condition.LastSyncedTime = &metav1.Time{Time: time.Now()}
	server.Status.ConfigTarget = *server.Spec.ConfigTarget.DeepCopy()
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

	if err := CreateServerDeployment(r, ctx, server); err != nil {
		return err
	}

	return nil
}

func ReconcileClient(r Reconciler, ctx context.Context, client_ *ratholev1alpha1.RatholeClient) error {
	if err := ValidateClient(client_); err != nil {
		return err
	}

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
	if client_.Spec.DefaultToken == "" {
		var err error
		if client_.Spec.DefaultToken, err = ReadConfig(r, ctx, client_.Namespace, client_.Spec.DefaultTokenFrom); err != nil {
			return err
		}
	}

	// Read PKCS12 content from PKCS12From field and set to PKCS12 field
	if client_.Spec.Transport.Type == "tls" {
		var confName string

		if client_.Spec.Transport.TLS.TrustedRoot == "" {
			if client_.Spec.Transport.TLS.TrustedRootFrom.SecretRef.Name != "" {
				confName = client_.Spec.Transport.TLS.TrustedRootFrom.SecretRef.Name
			} else if client_.Spec.Transport.TLS.TrustedRootFrom.ConfigMapRef.Name != "" {
				confName = client_.Spec.Transport.TLS.TrustedRootFrom.ConfigMapRef.Name
			}

			// trustedRoot -> filename
			filePath := fmt.Sprintf("%s/%s", ratholeSecretRoot, confName)
			client_.Spec.Transport.TLS.TrustedRoot = filePath
		}
	} else if client_.Spec.Transport.Type == "noise" {
		var err error
		if client_.Spec.Transport.Noise.LocalPrivateKey == "" && (client_.Spec.Transport.Noise.LocalPrivateKeyFrom.ConfigMapRef.Name != "" || client_.Spec.Transport.Noise.LocalPrivateKeyFrom.SecretRef.Name != "") {
			if client_.Spec.Transport.Noise.LocalPrivateKey, err = ReadConfig(r, ctx, client_.Namespace, client_.Spec.Transport.Noise.LocalPrivateKeyFrom); err != nil {
				return err
			}
		}
		if client_.Spec.Transport.Noise.RemotePublicKey == "" && (client_.Spec.Transport.Noise.RemotePublicKeyFrom.ConfigMapRef.Name != "" || client_.Spec.Transport.Noise.RemotePublicKeyFrom.SecretRef.Name != "") {
			if client_.Spec.Transport.Noise.RemotePublicKey, err = ReadConfig(r, ctx, client_.Namespace, client_.Spec.Transport.Noise.RemotePublicKeyFrom); err != nil {
				return err
			}
		}

		// If exist localPrivateKey and remotePublicKey, encode to base64
		if client_.Spec.Transport.Noise.LocalPrivateKey != "" {
			client_.Spec.Transport.Noise.EncodedLocalPrivateKey = base64.StdEncoding.EncodeToString([]byte(client_.Spec.Transport.Noise.LocalPrivateKey))
		}
		if client_.Spec.Transport.Noise.RemotePublicKey != "" {
			client_.Spec.Transport.Noise.EncodedRemotePublicKey = base64.StdEncoding.EncodeToString([]byte(client_.Spec.Transport.Noise.RemotePublicKey))
		}
	}

	if client_.Spec.ConfigTarget.ResourceType == "" {
		client_.Spec.ConfigTarget.ResourceType = "secret"
	}
	if client_.Spec.ConfigTarget.Name == "" {
		client_.Spec.ConfigTarget.Name = fmt.Sprintf("rathole-client-config-%s", GetSuffix(5))
	}

	client_.Status.Condition.ObservedGeneration = client_.Generation
	client_.Status.Condition.Status = "Synced"
	client_.Status.Condition.Reason = "Reconciled"
	client_.Status.Condition.LastSyncedTime = &metav1.Time{Time: time.Now()}
	client_.Status.ConfigTarget = *client_.Spec.ConfigTarget.DeepCopy()
	if err := r.Status().Update(ctx, client_.DeepCopy()); err != nil {
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
	var (
		server  ratholev1alpha1.RatholeServer
		client_ ratholev1alpha1.RatholeClient
	)

	if service.Spec.ServerRef.Name != "" {
		if err := r.Get(ctx, client.ObjectKey{Namespace: service.Namespace, Name: service.Spec.ServerRef.Name}, &server); err != nil {
			return err
		}
	}

	if service.Spec.ClientRef.Name != "" {
		if err := r.Get(ctx, client.ObjectKey{Namespace: service.Namespace, Name: service.Spec.ClientRef.Name}, &client_); err != nil {
			return err
		}
	}

	if err := ValidateService(service, &server, &client_); err != nil {
		return err
	}

	service.Status.Condition.ObservedGeneration = service.Generation
	service.Status.Condition.Status = "Synced"
	service.Status.Condition.Reason = "Reconciled"
	if err := r.Status().Update(ctx, service.DeepCopy()); err != nil {
		return err
	}

	if err := ReconcileServer(r, ctx, &server); err != nil {
		return err
	}
	if err := ReconcileClient(r, ctx, &client_); err != nil {
		return err
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

func CreateServerDeployment(r Reconciler, ctx context.Context, server *ratholev1alpha1.RatholeServer) error {
	var replicas int32 = 1

	var (
		serverConfVolumeSrc corev1.VolumeSource
	)
	if server.Spec.ConfigTarget.ResourceType == "secret" {
		serverConfVolumeSrc = corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: server.Spec.ConfigTarget.Name,
			},
		}
	} else if server.Spec.ConfigTarget.ResourceType == "configmap" {
		serverConfVolumeSrc = corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: server.Spec.ConfigTarget.Name,
				},
			},
		}
	}
	serverConfVolume := corev1.Volume{
		Name:         "config",
		VolumeSource: serverConfVolumeSrc,
	}

	ownerRefs := []metav1.OwnerReference{
		{
			APIVersion: server.APIVersion,
			Kind:       server.Kind,
			Name:       server.Name,
			UID:        server.UID,
		},
	}

	deploy := v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:            server.ObjectMeta.Name,
			Namespace:       server.Namespace,
			OwnerReferences: ownerRefs,
		},
		Spec: v1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":     server.ObjectMeta.Name,
					"service": "rathole-server",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":     server.ObjectMeta.Name,
						"service": "rathole-server",
					},
				},
				Spec: corev1.PodSpec{
					HostNetwork: true,
					Containers: []corev1.Container{
						{
							Name:  "rathole-server",
							Image: "crmin/rathole:v0.5.0-debug",
							Args: []string{
								"--server",
								"/var/run/config.toml",
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config",
									MountPath: "/var/run/config.toml",
									SubPath:   "config.toml",
								},
							},
						},
					},
					NodeSelector: server.Spec.Deployment.NodeSelector,
					Affinity: &corev1.Affinity{
						NodeAffinity: &server.Spec.Deployment.NodeAffinity,
					},
					Volumes: []corev1.Volume{
						serverConfVolume,
					},
				},
			},
		},
	}

	if err := r.Create(ctx, deploy.DeepCopy()); err != nil {
		return err
	}

	serverPort, err := strconv.Atoi(strings.Split(server.Spec.BindAddr, ":")[1])
	if err != nil {
		return err
	}

	service := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            server.ObjectMeta.Name,
			Namespace:       server.Namespace,
			OwnerReferences: ownerRefs,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app":     server.ObjectMeta.Name,
				"service": "rathole-server",
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "rathole",
					Port:       int32(serverPort),
					TargetPort: intstr.IntOrString{IntVal: int32(serverPort)},
				},
			},
		},
	}

	if err := r.Create(ctx, service.DeepCopy()); err != nil {
		return err
	}

	return nil
}
