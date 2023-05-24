/*
Copyright 2023 The Nephio Authors.

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
	"reflect"

	allocv1alpha1 "github.com/nokia/k8s-ipam/apis/alloc/common/v1alpha1"
	"github.com/nokia/k8s-ipam/pkg/meta"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// NodeSpec defines the desired state of Node
type NodeSpec struct {
	// UserDefinedLabels define metadata  associated to the resource.
	// defined in the spec to distingiush metadata labels from user defined labels
	allocv1alpha1.UserDefinedLabels `json:",inline" yaml:",inline"`

	// Location provider the location information where this resource is located
	Location *Location `json:"location,omitempty" yaml:"location,omitempty"`
	
	// ParametersRef points to the vendor or implementation specific params for the
	// network.
	// +optional
	ParametersRef *corev1.ObjectReference `json:"parametersRef,omitempty" yaml:"parametersRef,omitempty"`

	// Provider specifies the provider implementing this network.
	Provider string `json:"provider" yaml:"provider"`
}

type Location struct {
	Latitude  *string `json:"latitude,omitempty" yaml:"latitude,omitempty"`
	Longitude *string `json:"longitude,omitempty" yaml:"longitude,omitempty"`
}

// NodeStatus defines the observed state of Node
type NodeStatus struct {
	// ConditionedStatus provides the status of the Node allocation using conditions
	// 2 conditions are used:
	// - a condition for the reconcilation status
	// - a condition for the ready status
	// if both are true the other attributes in the status are meaningful
	allocv1alpha1.ConditionedStatus `json:",inline" yaml:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:resource:categories={nephio,inv}
// Node is the Schema for the vlan API
type Node struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	Spec   NodeSpec   `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status NodeStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

//+kubebuilder:object:root=true

// NodeList contains a list of Nodes
type NodeList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Items           []Node `json:"items" yaml:"items"`
}

func init() {
	SchemeBuilder.Register(&Node{}, &NodeList{})
}

var (
	NodeKind             = reflect.TypeOf(Node{}).Name()
	NodeGroupKind        = schema.GroupKind{Group: GroupVersion.Group, Kind: NodeKind}.String()
	NodeKindAPIVersion   = NodeKind + "." + GroupVersion.String()
	NodeGroupVersionKind = GroupVersion.WithKind(NodeKind)
	NodeKindGVKString    = meta.GVKToString(schema.GroupVersionKind{
		Group:   GroupVersion.Group,
		Version: GroupVersion.Version,
		Kind:    NodeKind,
	})
)