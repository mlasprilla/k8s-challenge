/*
Copyright 2023.

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

package controllers

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	examplecomv1alpha1 "example.com/api/v1alpha1"
)

// PythonHTTPServerConfigReconciler reconciles a PythonHTTPServerConfig object
type PythonHTTPServerConfigReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=example.com.example.com,resources=pythonhttpserverconfigs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com.example.com,resources=pythonhttpserverconfigs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com.example.com,resources=pythonhttpserverconfigs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PythonHTTPServerConfig object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
// func (r *PythonHTTPServerConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

// _ = log.FromContext(ctx)

// TODO(user): your logic here
func (r *PythonHTTPServerConfigReconciler) Reconcile(ctx context.Context, req types.Request) (ctrlResult sdk.Result, err error) {
	log, _ := logr.FromContext(ctx)

	// Fetch the PythonHTTPServerConfig instance
	var instance examplecomv1alpha1.PythonHTTPServerConfig
	if err := r.Get(ctx, req.NamespacedName, &instance); err != nil {
		// Error reading the object - requeue the request.
		return sdk.Result{}, sdk.IgnoreNotFound(err)
	}

	// Check if the object is being deleted
	if !instance.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.Delete(ctx, &instance)
	}

	// Add finalizer if it doesn't exist
	if !containsString(instance.ObjectMeta.Finalizers, finalizerName) {
		instance.ObjectMeta.Finalizers = append(instance.ObjectMeta.Finalizers, finalizerName)
		if err := r.Update(ctx, &instance); err != nil {
			return sdk.Result{}, err
		}
	}

	// Check if the instance was created or updated
	if reflect.DeepEqual(instance.Status, examplecomv1alpha1.PythonHTTPServerConfigStatus{}) {
		log.Info(fmt.Sprintf("PythonHTTPServerConfig %s/%s created", instance.Namespace, instance.Name))
	} else {
		if instance.Spec.Replicas != instance.Status.Replicas {
			log.Info(fmt.Sprintf("PythonHTTPServerConfig %s/%s updated - replicas: %d -> %d", instance.Namespace, instance.Name, instance.Status.Replicas, instance.Spec.Replicas))
		}
	}

	// Create or update the service for this instance
	if err := r.reconcileService(ctx, &instance); err != nil {
		return sdk.Result{}, err
	}

	// Create or update the deployment for this instance
	if err := r.reconcileDeployment(ctx, &instance); err != nil {
		return sdk.Result{}, err
	}

	return sdk.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PythonHTTPServerConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplecomv1alpha1.PythonHTTPServerConfig{}).
		Complete(r)
}
