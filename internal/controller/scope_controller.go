/*
Copyright 2025 Arsen Gumin.

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
	v1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	polyscopev1alpha1 "github.com/aagumin/polyscope/api/v1alpha1"
)

var log = logf.Log.WithName("controller")

// ScopeReconciler reconciles a Scope object
type ScopeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=polyscope.my.domain,resources=scopes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=polyscope.my.domain,resources=scopes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=polyscope.my.domain,resources=scopes/finalizers,verbs=update

// +kubebuilder:rbac:groups=rbac.authorization.k8s.io/v1,resources=Role,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io/v1,resources=RoleBinding,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Scope object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.4/pkg/reconcile
func (r *ScopeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// update -> create
	log.V(1).Info("reconciling scope object")

	var scope polyscopev1alpha1.Scope

	if err := r.Get(ctx, req.NamespacedName, &scope); err != nil {
		return ctrl.Result{}, err
	}

	roles, err := r.listAllRoles(ctx, &scope, req)

	scope

	if err != nil {
		return ctrl.Result{}, err
	}
	log.V(1).Info("found roles", "roles", roles)
	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

func (r *ScopeReconciler) listAllRoles(ctx context.Context, scope *polyscopev1alpha1.Scope, req ctrl.Request) (map[string][]v1.Role, error) {
	roles := make(map[string][]v1.Role)

	// Get roles from specific namespaces listed in the scope
	if len(scope.Spec.NamespaceSelectionSpec.Name) > 0 {
		for _, namespace := range scope.Spec.NamespaceSelectionSpec.Name {
			var namespaceRoles v1.RoleList
			if err := r.List(ctx, &namespaceRoles, client.InNamespace(namespace)); err != nil {
				log.Error(err, "unable to list roles", "namespace", namespace)
				return nil, err
			}
			roles[namespace] = namespaceRoles.Items
		}
	}

	// Get roles with matching labels across namespaces
	if scope.Spec.NamespaceSelectionSpec.Labels.MatchLabels != nil {
		// If specific namespaces are also provided, limit the search to those namespaces
		if len(scope.Spec.NamespaceSelectionSpec.Name) > 0 {
			for _, namespace := range scope.Spec.NamespaceSelectionSpec.Name {
				var labeledRoles v1.RoleList
				if err := r.List(ctx, &labeledRoles,
					client.InNamespace(namespace),
					client.MatchingLabels(scope.Spec.NamespaceSelectionSpec.Labels.MatchLabels)); err != nil {
					log.Error(err, "unable to list roles with labels", "namespace", namespace)
					return nil, err
				}
				// Append to existing roles if the namespace already has entries
				if existingRoles, ok := roles[namespace]; ok {
					roles[namespace] = append(existingRoles, labeledRoles.Items...)
				} else {
					roles[namespace] = labeledRoles.Items
				}
			}
		} else {
			// If no specific namespaces are provided, search across all namespaces
			// Note: This will only work if your controller has cluster-wide permissions
			var labeledRoles v1.RoleList
			if err := r.List(ctx, &labeledRoles,
				client.MatchingLabels(scope.Spec.NamespaceSelectionSpec.Labels.MatchLabels)); err != nil {
				log.Error(err, "unable to list roles with labels across all namespaces")
				return nil, err
			}

			// Group the roles by namespace
			for _, role := range labeledRoles.Items {
				namespace := role.Namespace
				if existingRoles, ok := roles[namespace]; ok {
					roles[namespace] = append(existingRoles, role)
				} else {
					roles[namespace] = []v1.Role{role}
				}
			}
		}
	}

	return roles, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ScopeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&polyscopev1alpha1.Scope{}).
		Named("scope").
		Complete(r)
}
