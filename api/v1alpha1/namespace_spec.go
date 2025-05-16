package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type NamespaceSelectionSpec struct {
	// +optional
	Name []string `json:"forName,omitempty"`
	// +optional
	Labels metav1.LabelSelector `json:"forLabels,omitempty"`
}
