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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	ratholev1alpha1 "github.com/crmin/rathole-operator/api/v1alpha1"
)

// RatholeServerReconciler reconciles a RatholeServer object
type RatholeServerReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

const serverFinalizerName = "rathole.superclass.io/server"

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
		if !containsString(server.ObjectMeta.Finalizers, serverFinalizerName) {
			log.Log.Info("Adding finalizer for server")
			server.ObjectMeta.Finalizers = append(server.ObjectMeta.Finalizers, serverFinalizerName)
			if err := r.Update(ctx, &server); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// Remove finalizer for deletion
		if containsString(server.ObjectMeta.Finalizers, serverFinalizerName) {
			log.Log.Info("Removing finalizer for server")
			server.ObjectMeta.Finalizers = removeString(server.ObjectMeta.Finalizers, serverFinalizerName)
			if err := r.Update(ctx, &server); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	if server.Status.Condition.ObservedGeneration == server.Generation {
		log.Log.Info("Skipping reconciliation: server spec hasn't changed")
		return ctrl.Result{}, nil
	}

	if err := ReconcileServer(r, ctx, &server); err != nil {
		log.Log.Error(err, "Failed to reconcile server")
		server.Status.Condition.Status = "Error"
		server.Status.Condition.Reason = err.Error()
		if err := r.Status().Update(ctx, &server); err != nil {
			log.Log.Error(err, "Failed to update server status")
			return ctrl.Result{}, err
		}
		// Retry after 10 seconds
		return ctrl.Result{RequeueAfter: 10}, nil
	}

	log.Log.Info("Rathole server reconciled successfully")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RatholeServerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ratholev1alpha1.RatholeServer{}).
		Complete(r)
}
