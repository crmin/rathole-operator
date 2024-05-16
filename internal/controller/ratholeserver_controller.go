/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	ratholev1alpha1 "github.com/crmin/rathole-operator/api/v1alpha1"
)

// RatholeServerReconciler reconciles a RatholeServer object
type RatholeServerReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	Recorder  record.EventRecorder
	Clientset *kubernetes.Clientset
}

const finalizerName = "rathole.superclass.io/server"

//+kubebuilder:rbac:groups=rathole.superclass.io,resources=ratholeservers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rathole.superclass.io,resources=ratholeservers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=rathole.superclass.io,resources=ratholeservers/finalizers,verbs=update
//+kubebuilder:rbac:groups=rathole.superclass.io,resources=ratholeservices,verbs=get;list;watch
//+kubebuilder:rbac:groups=,resources=configmap,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=,resources=secret,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the RatholeServer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.0/pkg/reconcile
func (r *RatholeServerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var server ratholev1alpha1.RatholeServer
	if err := r.Get(ctx, req.NamespacedName, &server); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if server.ObjectMeta.DeletionTimestamp.IsZero() {
		// Add finalizer
		if !containsString(server.ObjectMeta.Finalizers, finalizerName) {
			server.ObjectMeta.Finalizers = append(server.ObjectMeta.Finalizers, finalizerName)
			if err := r.Update(ctx, &server); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// Remove finalizer for deletion
		if containsString(server.ObjectMeta.Finalizers, finalizerName) {
			server.ObjectMeta.Finalizers = removeString(server.ObjectMeta.Finalizers, finalizerName)
			if err := r.Update(ctx, &server); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	if err := r.ReconcileServer(ctx, &server); err != nil {
		server.Status.Condition.Status = "Error"
		server.Status.Condition.Reason = err.Error()
		if err := r.Status().Update(ctx, &server); err != nil {
			return ctrl.Result{}, err
		}
		// Retry after 10 seconds
		return ctrl.Result{RequeueAfter: 10}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RatholeServerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ratholev1alpha1.RatholeServer{}).
		Complete(r)
}

func (r *RatholeServerReconciler) ReconcileServer(ctx context.Context, server *ratholev1alpha1.RatholeServer) error {
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
	// If both .Spec.DefaultTokenFrom and .Spec.DefaultToken are set, an error should occur through the webhook validate
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
			ok     bool
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
			if _, err := r.Clientset.CoreV1().Secrets(server.Namespace).Create(ctx, &secret, metav1.CreateOptions{}); err != nil {
				return err
			}
		} else {
			secret.StringData["config.toml"] = config
			if _, err := r.Clientset.CoreV1().Secrets(server.Namespace).Update(ctx, &secret, metav1.UpdateOptions{}); err != nil {
				return err
			}
		}

	case "configmap":
		var (
			configMap corev1.ConfigMap
			ok        bool
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
			if _, err := r.Clientset.CoreV1().ConfigMaps(server.Namespace).Create(ctx, &configMap, metav1.CreateOptions{}); err != nil {
				return err
			}
		} else {
			configMap.Data["config.toml"] = config

			if _, err := r.Clientset.CoreV1().ConfigMaps(server.Namespace).Update(ctx, &configMap, metav1.UpdateOptions{}); err != nil {
				return err
			}
		}
	}
	return nil
}
