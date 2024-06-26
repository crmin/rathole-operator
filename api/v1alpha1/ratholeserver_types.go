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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RatholeServerSpec defines the desired state of RatholeServer
type RatholeServerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// +optional
	ConfigTarget RatholeConfigTarget `json:"configTarget,omitempty" toml:"-"` // TODO: If not set, create random name secret

	BindAddr string `json:"bindAddr" toml:"bind_addr"`
	// +optional
	DefaultToken string `json:"defaultToken,omitempty" toml:"default_token,omitempty"`
	// +optional
	DefaultTokenFrom ResourceFrom `json:"defaultTokenFrom,omitempty" toml:"-"`
	// +optional
	HeartbeatInterval uint `json:"heartbeatInterval,omitempty" toml:"heartbeat_interval,omitempty"`
	// +optional
	Transport RatholeServerSpecTransport `json:"transport,omitempty" toml:"transport,omitempty"`

	// For server deployment
	// +optional
	Deployment RatholeServerSpecDeployment `json:"deployment,omitempty" toml:"-"`

	// +optional
	Services map[string]*RatholeServiceSpec `json:"-" toml:"services,omitempty"`

	// For deployment and service
	// +optional
	// +kubebuilder:validation:Enum=ClusterIP;NodePort;LoadBalancer
	ServiceType           v1.ServiceType    `json:"serviceType,omitempty" toml:"-"`           // default=ClusterIP
	ServiceAnnotations    map[string]string `json:"serviceAnnotations,omitempty" toml:"-"`    // default={}
	DeploymentAnnotations map[string]string `json:"deploymentAnnotations,omitempty" toml:"-"` // default={}
}

// RatholeServerStatus defines the observed state of RatholeServer
type RatholeServerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// +optional
	Condition RatholeServerStatusCondition `json:"condition,omitempty"`
	// +optional
	ConfigTarget RatholeConfigTarget `json:"configTarget,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RatholeServer is the Schema for the ratholeservers API
type RatholeServer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RatholeServerSpec   `json:"spec,omitempty"`
	Status RatholeServerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RatholeServerList contains a list of RatholeServer
type RatholeServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RatholeServer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RatholeServer{}, &RatholeServerList{})
}

type RatholeServerSpecTransport struct {
	// +optional
	// +kubebuilder:validation:Enum=tcp;tls;noise;websocket
	Type string `json:"type,omitempty" toml:"type,omitempty"`
	// +optional
	TCP *RatholeServerSpecTransportTCP `json:"tcp,omitempty" toml:"tcp,omitempty"`
	// +optional
	TLS *RatholeServerSpecTransportTLS `json:"tls,omitempty" toml:"tls,omitempty"`
	// +optional
	Noise *RatholeServerSpecTransportNoise `json:"noise,omitempty" toml:"noise,omitempty"`
	// +optional
	Websocket *RatholeServerSpecTransportWebsocket `json:"websocket,omitempty" toml:"websocket,omitempty"`
}

type RatholeServerSpecTransportTCP struct {
	// Optional, also affects `noise` and `tls`
	// +optional
	Nodelay bool `json:"nodelay,omitempty" toml:"nodelay,omitempty"`
	// +optional
	KeepaliveSecs uint `json:"keepaliveSecs,omitempty" toml:"keepalive_secs,omitempty"`
	// +optional
	KeepaliveInterval uint `json:"keepaliveInterval,omitempty" toml:"keepalive_interval,omitempty"`
}

type RatholeServerSpecTransportTLS struct {
	// If .Spec.Transport.Type is "tls", this field must be set.
	PKCS12From ResourceFrom `json:"pkcs12From" toml:"-"` // necessary, PKCS12 is file path
	// +optional
	PKCS12Password     string       `json:"pkcs12Password,omitempty" toml:"pkcs12_password,omitempty"`
	PKCS12PasswordFrom ResourceFrom `json:"pkcs12PasswordFrom,omitempty" toml:"-"`
	// TODO: Write hook for Validate; One of PKCS12 or PKCS12From must be set.

	// Field ignored in CRD generation. Used for internal logic.
	// +kubebuilder:skipversion
	PKCS12 string `json:"-" toml:"pkcs12"` // Make temp file using PKCS12From and set temp file path
}

type RatholeServerSpecTransportNoise struct {
	// If .Spec.Transport.Type is "noise", this field must be set.
	// +optional
	Pattern string `json:"pattern,omitempty" toml:"pattern,omitempty"`
	// If not set localPrivateKey and remotePublicKey, do not use auth
	// +optional
	LocalPrivateKey string `json:"localPrivateKey,omitempty" toml:"-"` // Set plain text, not base64 encoded
	// +optional
	LocalPrivateKeyFrom ResourceFrom `json:"localPrivateKeyFrom,omitempty" toml:"-"`
	// +optional
	RemotePublicKey string `json:"remotePublicKey,omitempty" toml:"-"` // Set plain text, not base64 encoded
	// +optional
	RemotePublicKeyFrom ResourceFrom `json:"remotePublicKeyFrom,omitempty" toml:"-"`

	// Field ignored in CRD generation. Used for internal logic.
	// +kubebuilder:skipversion
	EncodedLocalPrivateKey string `json:"-" toml:"local_private_key,omitempty"` // Base64 encoded value. Make temp file using LocalPrivateKeyFrom and set temp file path
	// +kubebuilder:skipversion
	EncodedRemotePublicKey string `json:"-" toml:"remote_public_key,omitempty"` // Base64 encoded value. Make temp file using RemotePublicKeyFrom and set temp file path
}

type RatholeServerSpecTransportWebsocket struct {
	// If .Spec.Transport.Type is "websocket", this field must be set.
	TLS bool `json:"tls,omitempty" toml:"tls,omitempty"` // necessary
}

type RatholeServerSpecDeployment struct {
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty" toml:"-"`
	// +optional
	NodeAffinity v1.NodeAffinity `json:"nodeAffinity,omitempty" toml:"-"`
}

type RatholeServerStatusCondition struct {
	// +optional
	Status string `json:"status,omitempty"`
	Reason string `json:"reason,omitempty"`
	// +optional
	LastSyncedTime *metav1.Time `json:"lastSyncedTime,omitempty"`
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}
