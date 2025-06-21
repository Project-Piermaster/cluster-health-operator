/*
Copyright 2025.

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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterIssueSpec defines the desired state of ClusterIssue
type ClusterIssueSpec struct {
	PodRef          corev1.ObjectReference `json:"podRef"`
	Reason          string                 `json:"reason"`
	Message         string                 `json:"message"`
	DiagnosedAt     metav1.Time            `json:"diagnosedAt"`
	NodeSuggestions []string               `json:"nodeSuggestions,omitempty"`
}

// ClusterIssueStatus defines the observed state of ClusterIssue
type ClusterIssueStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Resolved        bool   `json:"resolved,omitempty"`
	ResolutionNotes string `json:"resolutionNotes,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,shortName=ci

// ClusterIssue is the Schema for the clusterissues API
type ClusterIssue struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterIssueSpec   `json:"spec,omitempty"`
	Status ClusterIssueStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterIssueList contains a list of ClusterIssue
type ClusterIssueList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterIssue `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterIssue{}, &ClusterIssueList{})
}
