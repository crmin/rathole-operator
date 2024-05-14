package v1alpha1

import corev1 "k8s.io/api/core/v1"

type ResourceFrom struct {
	// +optional
	ConfigMapRef corev1.ConfigMapKeySelector `json:"configMapRef,omitempty"`
	// +optional
	SecretRef corev1.SecretKeySelector `json:"secretRef,omitempty"`
}

type RatholeConfigTarget struct {
	// +optional
	// +kubebuilder:validation:Enum=secret;Secret;configmap;Configmap;ConfigMap
	ResourceType string `json:"resourceType,omitempty" toml:"resource_type,omitempty"`
	// +optional
	Name string `json:"name,omitempty" toml:"name,omitempty"`
}
