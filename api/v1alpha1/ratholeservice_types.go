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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RatholeServiceSpec defines the desired state of RatholeService
type RatholeServiceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	ServerRef RatholeServiceResourceRef `json:"serverRef" toml:"-"`
	// +optional
	ClientRef RatholeServiceResourceRef `json:"clientRef,omitempty" toml:"-"`

	// +optional
	// +kubebuilder:validation:Enum=tcp;udp
	Type string `json:"type,omitempty" toml:"type,omitempty"` // default=tcp
	// +optional
	Token     string `json:"token,omitempty" toml:"token,omitempty"` // default=client.default_token TODO: Write hook for validate; Token must set if client.default_token was not set
	LocalAddr string `json:"localAddr" toml:"local_addr,omitempty"`  // necessary, for client
	BindAddr  string `json:"bindAddr" toml:"bind_addr,omitempty"`    // necessary, for server
	// +optional
	Nodelay bool `json:"nodelay,omitempty" toml:"nodelay,omitempty"` // Override client.transport.nodelay
	// +optional
	RetryInterval uint `json:"retryInterval,omitempty" toml:"retry_interval,omitempty"` // Override client.retry_interval, for client
}

// RatholeServiceStatus defines the observed state of RatholeService
type RatholeServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// +optional
	Condition RatholeServiceStatusCondition `json:"condition,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RatholeService is the Schema for the ratholeservices API
type RatholeService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RatholeServiceSpec   `json:"spec,omitempty"`
	Status RatholeServiceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RatholeServiceList contains a list of RatholeService
type RatholeServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RatholeService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RatholeService{}, &RatholeServiceList{})
}

type RatholeServiceResourceRef struct {
	Name string `json:"name" toml:"-"`
}

type RatholeServiceStatusCondition struct {
	Status string `json:"status,omitempty"`
	Reason string `json:"reason,omitempty"`
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}
