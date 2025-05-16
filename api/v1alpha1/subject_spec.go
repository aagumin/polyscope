package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
)

type SubjectSpec struct {
	// +optional
	SaName string `json:"saName"`
	// +optional
	SaTemplate v1.ServiceAccount `json:"saTemplate"`
}
