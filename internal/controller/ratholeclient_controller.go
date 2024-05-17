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
	"k8s.io/client-go/tools/record"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	ratholev1alpha1 "github.com/crmin/rathole-operator/api/v1alpha1"
)

// RatholeClientReconciler reconciles a RatholeClient object
type RatholeClientReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

const clientFinalizerName = "rathole.superclass.io/client"

//+kubebuilder:rbac:groups=rathole.superclass.io,resources=ratholeclients,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rathole.superclass.io,resources=ratholeclients/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=rathole.superclass.io,resources=ratholeclients/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the RatholeClient object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.0/pkg/reconcile
func (r *RatholeClientReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var client_ ratholev1alpha1.RatholeClient
	if err := r.Get(ctx, req.NamespacedName, &client_); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if client_.ObjectMeta.DeletionTimestamp.IsZero() {
		// Add finalizer
		if !containsString(client_.ObjectMeta.Finalizers, clientFinalizerName) {
			client_.ObjectMeta.Finalizers = append(client_.ObjectMeta.Finalizers, clientFinalizerName)
			if err := r.Update(ctx, &client_); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// Remove finalizer for deletion
		if containsString(client_.ObjectMeta.Finalizers, clientFinalizerName) {
			client_.ObjectMeta.Finalizers = removeString(client_.ObjectMeta.Finalizers, clientFinalizerName)
			if err := r.Update(ctx, &client_); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	// Skip if already reconciled. Spec hasn't changed
	if client_.Status.Condition.ObservedGeneration == client_.Generation {
		return ctrl.Result{}, nil
	}

	if err := ReconcileClient(r, ctx, &client_); err != nil {
		client_.Status.Condition.Status = "Error"
		client_.Status.Condition.Reason = err.Error()
		if err := r.Status().Update(ctx, &client_); err != nil {
			return ctrl.Result{}, err
		}
		// Retry after 10 seconds
		return ctrl.Result{RequeueAfter: 10}, nil
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RatholeClientReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ratholev1alpha1.RatholeClient{}).
		Complete(r)
}
