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

// RatholeServiceReconciler reconciles a RatholeService object
type RatholeServiceReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

const serviceFinalizerName = "rathole.superclass.io/service"

//+kubebuilder:rbac:groups=rathole.superclass.io,resources=ratholeservices,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rathole.superclass.io,resources=ratholeservices/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=rathole.superclass.io,resources=ratholeservices/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the RatholeService object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.0/pkg/reconcile
func (r *RatholeServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var service ratholev1alpha1.RatholeService
	if err := r.Get(ctx, req.NamespacedName, &service); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if service.ObjectMeta.DeletionTimestamp.IsZero() {
		// Add finalizer
		if !containsString(service.ObjectMeta.Finalizers, clientFinalizerName) {
			service.ObjectMeta.Finalizers = append(service.ObjectMeta.Finalizers, clientFinalizerName)
			if err := r.Update(ctx, &service); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// Remove finalizer for deletion
		if containsString(service.ObjectMeta.Finalizers, clientFinalizerName) {
			service.ObjectMeta.Finalizers = removeString(service.ObjectMeta.Finalizers, clientFinalizerName)
			if err := r.Update(ctx, &service); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	// Skip if already reconciled. Spec hasn't changed
	if service.Status.Condition.ObservedGeneration == service.Generation {
		return ctrl.Result{}, nil
	}

	if err := ReconcileService(r, ctx, &service); err != nil {
		service.Status.Condition.Status = "Error"
		service.Status.Condition.Reason = err.Error()
		if err := r.Status().Update(ctx, &service); err != nil {
			return ctrl.Result{}, err
		}
		// Retry after 10 seconds
		return ctrl.Result{RequeueAfter: 10}, nil
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RatholeServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ratholev1alpha1.RatholeService{}).
		Complete(r)
}
