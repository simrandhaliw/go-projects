/*

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
// +kubebuilder:docs-gen:collapse=Apache License

/*
Since we're in a v1alpha2 package, controller-gen will assume this is for the v1alpha2
version automatically.  We could override that with the [`+versionName`
marker](/reference/markers/crd.md).
*/
package v1alpha2

/*
 */
import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// +kubebuilder:docs-gen:collapse=Imports

// MemcachedSpec defines the desired state of Memcached
// +k8s:openapi-gen=true
type MemcachedSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Price is a field representing price per GB for a disk.
	Price Price `json:"price"`

	// +kubebuilder:validation:Minimum=0
	// Size is the size of the memcached deployment
	Size int32 `json:"size"`

	// This flag tells the controller to suspend subsequent executions, it does
	// not apply to already started executions.  Defaults to false.
	// +optional
	Suspend *bool `json:"suspend,omitempty"`
}

// Price represents a generic price value that has amount and currency.
type Price struct {
	// specifies the amount value.
	// +optional
	Amount int64 `json:"amount"`
	// specifies the curreny type.
	// +optional
	Currency string `json:"currency"`
}

// MemcachedStatus defines the observed state of Memcached
// +k8s:openapi-gen=true
type MemcachedStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// Nodes are the names of the memcached pods
	Nodes []string `json:"nodes"`
}

// +kubebuilder:object:root=true

// Memcached is the Schema for the memcacheds API
// +kubebuilder:subresource:status
// k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// Memcached is the Schema for the memcached API
type Memcached struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MemcachedSpec   `json:"spec,omitempty"`
	Status MemcachedStatus `json:"status,omitempty"`
}

/*
 */
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.
// +kubebuilder:object:root=true

// MemcachedList contains a list of Memcached
type MemcachedList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Memcached `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Memcached{}, &MemcachedList{})
}
