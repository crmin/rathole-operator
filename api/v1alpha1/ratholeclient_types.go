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

// RatholeClientSpec defines the desired state of RatholeClient
type RatholeClientSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// +optional
	ConfigTarget RatholeConfigTarget `json:"configTarget,omitempty" toml:"config_target,omitempty"` // If not set, create random name secret

	RemoteAddr string `json:"remoteAddr" toml:"remote_addr"`
	// +optional
	DefaultToken string `json:"defaultToken,omitempty" toml:"default_token,omitempty"`
	// +optional
	HeartbeatTimeout uint `json:"heartbeatTimeout,omitempty" toml:"heartbeat_timeout,omitempty"`
	// +optional
	RetryInterval uint `json:"retryInterval,omitempty" toml:"retry_interval,omitempty"`
	// +optional
	Transport RatholeClientSpecTransport `json:"transport,omitempty" toml:"transport,omitempty"`

	// TODO: DefaultTokenFrom 필드 생성 필요
}

// RatholeClientStatus defines the observed state of RatholeClient
type RatholeClientStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// +optional
	Condition RatholeClientStatusCondition `json:"condition,omitempty"`
	// +optional
	ConfigTarget RatholeConfigTarget `json:"configTarget,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RatholeClient is the Schema for the ratholeclients API
type RatholeClient struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RatholeClientSpec   `json:"spec,omitempty"`
	Status RatholeClientStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RatholeClientList contains a list of RatholeClient
type RatholeClientList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RatholeClient `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RatholeClient{}, &RatholeClientList{})
}

type RatholeClientSpecTransport struct {
	// +optional
	// +kubebuilder:validation:Enum=tcp;tls;noise;websocket
	Type string `json:"type,omitempty" toml:"type,omitempty"`
	// +optional
	TCP *RatholeClientSpecTransportTCP `json:"tcp,omitempty" toml:"tcp,omitempty"`
	// +optional
	TLS *RatholeClientSpecTransportTLS `json:"tls,omitempty" toml:"tls,omitempty"`
	// +optional
	Noise *RatholeClientSpecTransportNoise `json:"noise,omitempty" toml:"noise,omitempty"`
	// +optional
	Websocket *RatholeClientSpecTransportWebsocket `json:"websocket,omitempty" toml:"websocket,omitempty"`
}

type RatholeClientSpecTransportTCP struct {
	// Optional, also affects `noise` and `tls`
	// +optional
	Proxy string `json:"proxy,omitempty" toml:"proxy,omitempty"`
	// +optional
	Nodelay bool `json:"nodelay,omitempty" toml:"nodelay,omitempty"`
	// +optional
	KeepaliveSecs uint `json:"keepaliveSecs,omitempty" toml:"keepalive_secs,omitempty"`
	// +optional
	KeepaliveInterval uint `json:"keepaliveInterval,omitempty" toml:"keepalive_interval,omitempty"`
}

type RatholeClientSpecTransportTLS struct {
	// If .Spec.Transport.Type is "tls", this field must be set.
	TrustedRootFrom ResourceFrom `json:"trustedRootFrom,omitempty" toml:"-"`
	// +optional
	Hostname string `json:"hostname,omitempty" toml:"hostname,omitempty"`

	// Field ignored in CRD generation. Used for internal logic.
	// +kubebuilder:skipversion
	TrustedRoot string `json:"-" toml:"trusted_root"` // Make temp file using TrustedRootFrom and set temp file path
}

type RatholeClientSpecTransportNoise struct {
	// If .Spec.Transport.Type is "noise", this field must be set.
	// +optional
	Pattern string `json:"pattern,omitempty" toml:"pattern,omitempty"`
	// +optional
	LocalPrivateKey ResourceFrom `json:"localPrivateKey,omitempty" toml:"-"` // Set plain text, not base64 encoded
	// +optional
	LocalPrivateKeyFrom ResourceFrom `json:"localPrivateKeyFrom,omitempty" toml:"-"`
	// +optional
	RemotePublicKey ResourceFrom `json:"remotePublicKey,omitempty" toml:"-"` // Set plain text, not base64 encoded
	// +optional
	RemotePublicKeyFrom ResourceFrom `json:"remotePublicKeyFrom,omitempty" toml:"-"`

	// TODO: Write hook for Validate; One of LocalPrivateKey or LocalPrivateKeyFrom must be set.
	// TODO: If EncodedLocalPrivateKey was set, set LocalPrivateKey after base64 encoding.
	// TODO: If LocalPrivateKeyFrom was set, read value LocalPrivateKeyFrom and encode to base64, set EncodedLocalPrivateKey.

	// TODO: Write hook for Validate; One of RemotePublicKey or RemotePublicKeyFrom must be set.
	// TODO: If EncodedRemotePublicKey was set, set RemotePublicKey after base64 encoding.
	// TODO: If RemotePublicKeyFrom was set, read value RemotePublicKeyFrom and encode to base64, set EncodedRemotePublicKey.

	// Field ignored in CRD generation. Used for internal logic.
	// +kubebuilder:skipversion
	EncodedLocalPrivateKey string `json:"-" toml:"local_private_key,omitempty"` // Make temp file using LocalPrivateKeyFrom and set temp file path
	// +kubebuilder:skipversion
	EncodedRemotePublicKey string `json:"-" toml:"remote_public_key,omitempty"` // Make temp file using RemotePublicKeyFrom and set temp file path
}

type RatholeClientSpecTransportWebsocket struct {
	// If .Spec.Transport.Type is "websocket", this field must be set.
	TLS bool `json:"tls,omitempty" toml:"tls,omitempty"` // necessary
}

type RatholeClientStatusCondition struct {
	// +optional
	Status string `json:"status,omitempty"`
	Reason string `json:"reason,omitempty"`
	// +optional
	LastSyncedTime *metav1.Time `json:"lastSyncedTime,omitempty"`
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}
