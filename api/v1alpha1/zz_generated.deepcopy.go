//go:build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RatholeClient) DeepCopyInto(out *RatholeClient) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RatholeClient.
func (in *RatholeClient) DeepCopy() *RatholeClient {
	if in == nil {
		return nil
	}
	out := new(RatholeClient)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RatholeClient) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RatholeClientList) DeepCopyInto(out *RatholeClientList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]RatholeClient, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RatholeClientList.
func (in *RatholeClientList) DeepCopy() *RatholeClientList {
	if in == nil {
		return nil
	}
	out := new(RatholeClientList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RatholeClientList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RatholeClientSpec) DeepCopyInto(out *RatholeClientSpec) {
	*out = *in
	out.ConfigTarget = in.ConfigTarget
	in.DefaultTokenFrom.DeepCopyInto(&out.DefaultTokenFrom)
	in.Transport.DeepCopyInto(&out.Transport)
	if in.Services != nil {
		in, out := &in.Services, &out.Services
		*out = make(map[string]*RatholeServiceSpec, len(*in))
		for key, val := range *in {
			var outVal *RatholeServiceSpec
			if val == nil {
				(*out)[key] = nil
			} else {
				inVal := (*in)[key]
				in, out := &inVal, &outVal
				*out = new(RatholeServiceSpec)
				**out = **in
			}
			(*out)[key] = outVal
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RatholeClientSpec.
func (in *RatholeClientSpec) DeepCopy() *RatholeClientSpec {
	if in == nil {
		return nil
	}
	out := new(RatholeClientSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RatholeClientSpecTransport) DeepCopyInto(out *RatholeClientSpecTransport) {
	*out = *in
	if in.TCP != nil {
		in, out := &in.TCP, &out.TCP
		*out = new(RatholeClientSpecTransportTCP)
		**out = **in
	}
	if in.TLS != nil {
		in, out := &in.TLS, &out.TLS
		*out = new(RatholeClientSpecTransportTLS)
		(*in).DeepCopyInto(*out)
	}
	if in.Noise != nil {
		in, out := &in.Noise, &out.Noise
		*out = new(RatholeClientSpecTransportNoise)
		(*in).DeepCopyInto(*out)
	}
	if in.Websocket != nil {
		in, out := &in.Websocket, &out.Websocket
		*out = new(RatholeClientSpecTransportWebsocket)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RatholeClientSpecTransport.
func (in *RatholeClientSpecTransport) DeepCopy() *RatholeClientSpecTransport {
	if in == nil {
		return nil
	}
	out := new(RatholeClientSpecTransport)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RatholeClientSpecTransportNoise) DeepCopyInto(out *RatholeClientSpecTransportNoise) {
	*out = *in
	in.LocalPrivateKeyFrom.DeepCopyInto(&out.LocalPrivateKeyFrom)
	in.RemotePublicKeyFrom.DeepCopyInto(&out.RemotePublicKeyFrom)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RatholeClientSpecTransportNoise.
func (in *RatholeClientSpecTransportNoise) DeepCopy() *RatholeClientSpecTransportNoise {
	if in == nil {
		return nil
	}
	out := new(RatholeClientSpecTransportNoise)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RatholeClientSpecTransportTCP) DeepCopyInto(out *RatholeClientSpecTransportTCP) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RatholeClientSpecTransportTCP.
func (in *RatholeClientSpecTransportTCP) DeepCopy() *RatholeClientSpecTransportTCP {
	if in == nil {
		return nil
	}
	out := new(RatholeClientSpecTransportTCP)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RatholeClientSpecTransportTLS) DeepCopyInto(out *RatholeClientSpecTransportTLS) {
	*out = *in
	in.TrustedRootFrom.DeepCopyInto(&out.TrustedRootFrom)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RatholeClientSpecTransportTLS.
func (in *RatholeClientSpecTransportTLS) DeepCopy() *RatholeClientSpecTransportTLS {
	if in == nil {
		return nil
	}
	out := new(RatholeClientSpecTransportTLS)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RatholeClientSpecTransportWebsocket) DeepCopyInto(out *RatholeClientSpecTransportWebsocket) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RatholeClientSpecTransportWebsocket.
func (in *RatholeClientSpecTransportWebsocket) DeepCopy() *RatholeClientSpecTransportWebsocket {
	if in == nil {
		return nil
	}
	out := new(RatholeClientSpecTransportWebsocket)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RatholeClientStatus) DeepCopyInto(out *RatholeClientStatus) {
	*out = *in
	in.Condition.DeepCopyInto(&out.Condition)
	out.ConfigTarget = in.ConfigTarget
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RatholeClientStatus.
func (in *RatholeClientStatus) DeepCopy() *RatholeClientStatus {
	if in == nil {
		return nil
	}
	out := new(RatholeClientStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RatholeClientStatusCondition) DeepCopyInto(out *RatholeClientStatusCondition) {
	*out = *in
	if in.LastSyncedTime != nil {
		in, out := &in.LastSyncedTime, &out.LastSyncedTime
		*out = (*in).DeepCopy()
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RatholeClientStatusCondition.
func (in *RatholeClientStatusCondition) DeepCopy() *RatholeClientStatusCondition {
	if in == nil {
		return nil
	}
	out := new(RatholeClientStatusCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RatholeConfigTarget) DeepCopyInto(out *RatholeConfigTarget) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RatholeConfigTarget.
func (in *RatholeConfigTarget) DeepCopy() *RatholeConfigTarget {
	if in == nil {
		return nil
	}
	out := new(RatholeConfigTarget)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RatholeServer) DeepCopyInto(out *RatholeServer) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RatholeServer.
func (in *RatholeServer) DeepCopy() *RatholeServer {
	if in == nil {
		return nil
	}
	out := new(RatholeServer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RatholeServer) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RatholeServerList) DeepCopyInto(out *RatholeServerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]RatholeServer, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RatholeServerList.
func (in *RatholeServerList) DeepCopy() *RatholeServerList {
	if in == nil {
		return nil
	}
	out := new(RatholeServerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RatholeServerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RatholeServerSpec) DeepCopyInto(out *RatholeServerSpec) {
	*out = *in
	out.ConfigTarget = in.ConfigTarget
	in.DefaultTokenFrom.DeepCopyInto(&out.DefaultTokenFrom)
	in.Transport.DeepCopyInto(&out.Transport)
	in.Deployment.DeepCopyInto(&out.Deployment)
	if in.Services != nil {
		in, out := &in.Services, &out.Services
		*out = make(map[string]*RatholeServiceSpec, len(*in))
		for key, val := range *in {
			var outVal *RatholeServiceSpec
			if val == nil {
				(*out)[key] = nil
			} else {
				inVal := (*in)[key]
				in, out := &inVal, &outVal
				*out = new(RatholeServiceSpec)
				**out = **in
			}
			(*out)[key] = outVal
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RatholeServerSpec.
func (in *RatholeServerSpec) DeepCopy() *RatholeServerSpec {
	if in == nil {
		return nil
	}
	out := new(RatholeServerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RatholeServerSpecDeployment) DeepCopyInto(out *RatholeServerSpecDeployment) {
	*out = *in
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.NodeAffinity.DeepCopyInto(&out.NodeAffinity)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RatholeServerSpecDeployment.
func (in *RatholeServerSpecDeployment) DeepCopy() *RatholeServerSpecDeployment {
	if in == nil {
		return nil
	}
	out := new(RatholeServerSpecDeployment)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RatholeServerSpecTransport) DeepCopyInto(out *RatholeServerSpecTransport) {
	*out = *in
	if in.TCP != nil {
		in, out := &in.TCP, &out.TCP
		*out = new(RatholeServerSpecTransportTCP)
		**out = **in
	}
	if in.TLS != nil {
		in, out := &in.TLS, &out.TLS
		*out = new(RatholeServerSpecTransportTLS)
		(*in).DeepCopyInto(*out)
	}
	if in.Noise != nil {
		in, out := &in.Noise, &out.Noise
		*out = new(RatholeServerSpecTransportNoise)
		(*in).DeepCopyInto(*out)
	}
	if in.Websocket != nil {
		in, out := &in.Websocket, &out.Websocket
		*out = new(RatholeServerSpecTransportWebsocket)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RatholeServerSpecTransport.
func (in *RatholeServerSpecTransport) DeepCopy() *RatholeServerSpecTransport {
	if in == nil {
		return nil
	}
	out := new(RatholeServerSpecTransport)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RatholeServerSpecTransportNoise) DeepCopyInto(out *RatholeServerSpecTransportNoise) {
	*out = *in
	in.LocalPrivateKeyFrom.DeepCopyInto(&out.LocalPrivateKeyFrom)
	in.RemotePublicKeyFrom.DeepCopyInto(&out.RemotePublicKeyFrom)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RatholeServerSpecTransportNoise.
func (in *RatholeServerSpecTransportNoise) DeepCopy() *RatholeServerSpecTransportNoise {
	if in == nil {
		return nil
	}
	out := new(RatholeServerSpecTransportNoise)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RatholeServerSpecTransportTCP) DeepCopyInto(out *RatholeServerSpecTransportTCP) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RatholeServerSpecTransportTCP.
func (in *RatholeServerSpecTransportTCP) DeepCopy() *RatholeServerSpecTransportTCP {
	if in == nil {
		return nil
	}
	out := new(RatholeServerSpecTransportTCP)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RatholeServerSpecTransportTLS) DeepCopyInto(out *RatholeServerSpecTransportTLS) {
	*out = *in
	in.PKCS12From.DeepCopyInto(&out.PKCS12From)
	in.PKCS12PasswordFrom.DeepCopyInto(&out.PKCS12PasswordFrom)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RatholeServerSpecTransportTLS.
func (in *RatholeServerSpecTransportTLS) DeepCopy() *RatholeServerSpecTransportTLS {
	if in == nil {
		return nil
	}
	out := new(RatholeServerSpecTransportTLS)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RatholeServerSpecTransportWebsocket) DeepCopyInto(out *RatholeServerSpecTransportWebsocket) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RatholeServerSpecTransportWebsocket.
func (in *RatholeServerSpecTransportWebsocket) DeepCopy() *RatholeServerSpecTransportWebsocket {
	if in == nil {
		return nil
	}
	out := new(RatholeServerSpecTransportWebsocket)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RatholeServerStatus) DeepCopyInto(out *RatholeServerStatus) {
	*out = *in
	in.Condition.DeepCopyInto(&out.Condition)
	out.ConfigTarget = in.ConfigTarget
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RatholeServerStatus.
func (in *RatholeServerStatus) DeepCopy() *RatholeServerStatus {
	if in == nil {
		return nil
	}
	out := new(RatholeServerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RatholeServerStatusCondition) DeepCopyInto(out *RatholeServerStatusCondition) {
	*out = *in
	if in.LastSyncedTime != nil {
		in, out := &in.LastSyncedTime, &out.LastSyncedTime
		*out = (*in).DeepCopy()
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RatholeServerStatusCondition.
func (in *RatholeServerStatusCondition) DeepCopy() *RatholeServerStatusCondition {
	if in == nil {
		return nil
	}
	out := new(RatholeServerStatusCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RatholeService) DeepCopyInto(out *RatholeService) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RatholeService.
func (in *RatholeService) DeepCopy() *RatholeService {
	if in == nil {
		return nil
	}
	out := new(RatholeService)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RatholeService) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RatholeServiceList) DeepCopyInto(out *RatholeServiceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]RatholeService, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RatholeServiceList.
func (in *RatholeServiceList) DeepCopy() *RatholeServiceList {
	if in == nil {
		return nil
	}
	out := new(RatholeServiceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RatholeServiceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RatholeServiceResourceRef) DeepCopyInto(out *RatholeServiceResourceRef) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RatholeServiceResourceRef.
func (in *RatholeServiceResourceRef) DeepCopy() *RatholeServiceResourceRef {
	if in == nil {
		return nil
	}
	out := new(RatholeServiceResourceRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RatholeServiceSpec) DeepCopyInto(out *RatholeServiceSpec) {
	*out = *in
	out.ServerRef = in.ServerRef
	out.ClientRef = in.ClientRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RatholeServiceSpec.
func (in *RatholeServiceSpec) DeepCopy() *RatholeServiceSpec {
	if in == nil {
		return nil
	}
	out := new(RatholeServiceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RatholeServiceStatus) DeepCopyInto(out *RatholeServiceStatus) {
	*out = *in
	out.Condition = in.Condition
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RatholeServiceStatus.
func (in *RatholeServiceStatus) DeepCopy() *RatholeServiceStatus {
	if in == nil {
		return nil
	}
	out := new(RatholeServiceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RatholeServiceStatusCondition) DeepCopyInto(out *RatholeServiceStatusCondition) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RatholeServiceStatusCondition.
func (in *RatholeServiceStatusCondition) DeepCopy() *RatholeServiceStatusCondition {
	if in == nil {
		return nil
	}
	out := new(RatholeServiceStatusCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceFrom) DeepCopyInto(out *ResourceFrom) {
	*out = *in
	in.ConfigMapRef.DeepCopyInto(&out.ConfigMapRef)
	in.SecretRef.DeepCopyInto(&out.SecretRef)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceFrom.
func (in *ResourceFrom) DeepCopy() *ResourceFrom {
	if in == nil {
		return nil
	}
	out := new(ResourceFrom)
	in.DeepCopyInto(out)
	return out
}
