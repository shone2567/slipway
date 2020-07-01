/*
Copyright (c) 2020 Facebook, Inc. and its affiliates.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Important: Run "make" to regenerate code after modifying this file

// ImageMirrorSpec defines the desired state of ImageMirror
type ImageMirrorSpec struct {
	// SourceRepository is a URL resource, including scheme (optional),
	// registry host, and registry organization (e.g. docker.io/dwat/) which
	// will be used to pull images to mirror. NOTE: This must not include
	// the container image name or any tags.
	SourceRepository string `json:"source_repository,required"`

	// DestRepository is a URL resource as above, which is used to
	// push mirrored container images.
	DestRepository string `json:"dest_repository,required"`

	// ImageName is the name of the image without tag (e.g. cuda).
	ImageName string `json:"image_name,required"`

	// Pattern is a regex matching the tags which should be mirrored. We
	// follow Flux (who of course copied as well) because compatibility
	// is nice for our users. Cf.
	// https://github.com/fluxcd/flux/blob/v1.19.0/pkg/policy/pattern.go
	Pattern string `json:"tag_regex,omitempty"`
}

// ImageMirrorStatus defines the observed state of ImageMirror
type ImageMirrorStatus struct {
	// MirroredTags is a slice of tags which have already been mirrored.
	MirroredTags []string `json:"mirrored_tags"`
}

// +kubebuilder:object:root=true

// ImageMirror is the Schema for the imagemirrors API
type ImageMirror struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ImageMirrorSpec   `json:"spec,omitempty"`
	Status ImageMirrorStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ImageMirrorList contains a list of ImageMirror
type ImageMirrorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ImageMirror `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ImageMirror{}, &ImageMirrorList{})
}
