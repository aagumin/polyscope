package v1alpha1

import v1 "k8s.io/api/rbac/v1"

type FromExist struct {
	RoleName string `json:"roleName"`
	// +optional
	RoleNamespace string `json:"roleNamespace,omitempty"`
}

type FromTemplate struct {
	RoleTemplate v1.Role `json:"roleTemplate"`
}
type RoleSpec struct {
	// +optional
	FromExist FromExist `json:"fromExist,omitempty"`
	// +optional
	FromTemplate FromTemplate `json:"fromTemplate,omitempty"`
}
