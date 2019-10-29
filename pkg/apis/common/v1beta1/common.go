// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package v1beta1

import (
	"reflect"

	v1 "k8s.io/api/core/v1"
	"k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// ReconcilerStatus represents status information about desired/available nodes.
type ReconcilerStatus struct {
	AvailableNodes int `json:"availableNodes,omitempty"`
}

// SecretRef is a reference to a secret that exists in the same namespace.
type SecretRef struct {
	// SecretName is the name of the secret.
	SecretName string `json:"secretName,omitempty"`
}

// ObjectSelector defines a reference to a Kubernetes object.
type ObjectSelector struct {
	// Name of the Kubernetes object.
	Name string `json:"name"`
	// Namespace of the Kubernetes object. If empty, defaults to the current namespace.
	Namespace string `json:"namespace,omitempty"`
}

// NamespacedName is a convenience method to turn an ObjectSelector into a NamespacedName.
func (s ObjectSelector) NamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Name:      s.Name,
		Namespace: s.Namespace,
	}
}

// IsDefined checks if the object selector is not nil and has a name.
// Namespace is not mandatory as it may be inherited by the parent object.
func (s *ObjectSelector) IsDefined() bool {
	return s != nil && s.Name != ""
}

// HTTPConfig holds the HTTP layer configuration for resources.
type HTTPConfig struct {
	// Service defines the template for the associated Kubernetes Service object.
	Service ServiceTemplate `json:"service,omitempty"`
	// TLS defines options for configuring TLS for HTTP.
	TLS TLSOptions `json:"tls,omitempty"`
}

// Scheme returns the scheme for this HTTP config
func (http HTTPConfig) Scheme() string {
	if http.TLS.Enabled() {
		return "https"
	}
	return "http"
}

// TLSOptions holds TLS configuration options.
type TLSOptions struct {
	// SelfSignedCertificate allows configuring the self-signed certificate generated by the operator.
	SelfSignedCertificate *SelfSignedCertificate `json:"selfSignedCertificate,omitempty"`

	// Certificate is a reference to a Kubernetes secret that contains the certificate and private key for enabling TLS.
	// The referenced secret should contain the following:
	//
	// - `ca.crt`: The certificate authority (optional).
	// - `tls.crt`: The certificate (or a chain).
	// - `tls.key`: The private key to the first certificate in the certificate chain.
	Certificate SecretRef `json:"certificate,omitempty"`
}

// Enabled returns true when TLS is enabled based on this option struct.
func (tls TLSOptions) Enabled() bool {
	selfSigned := tls.SelfSignedCertificate
	return selfSigned == nil || !selfSigned.Disabled || tls.Certificate.SecretName != ""
}

// SelfSignedCertificate holds configuration for the self-signed certificate generated by the operator.
type SelfSignedCertificate struct {
	// SubjectAlternativeNames is a list of SANs to include in the generated HTTP TLS certificate.
	SubjectAlternativeNames []SubjectAlternativeName `json:"subjectAltNames,omitempty"`
	// Disabled indicates that the provisioning of the self-signed certifcate should be disabled.
	Disabled bool `json:"disabled,omitempty"`
}

// SubjectAlternativeName represents a SAN entry in a x509 certificate.
type SubjectAlternativeName struct {
	// DNS is the DNS name of the subject.
	DNS string `json:"dns,omitempty"`
	// IP is the IP address of the subject.
	IP string `json:"ip,omitempty"`
}

// ServiceTemplate defines the template for a Kubernetes Service.
type ServiceTemplate struct {
	// ObjectMeta is the metadata of the service.
	// The name and namespace provided here are managed by ECK and will be ignored.
	// +kubebuilder:validation:Optional
	ObjectMeta metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the specification of the service.
	// +kubebuilder:validation:Optional
	Spec v1.ServiceSpec `json:"spec,omitempty"`
}

// DefaultPodDisruptionBudgetMaxUnavailable is the default max unavailable pods in a PDB.
var DefaultPodDisruptionBudgetMaxUnavailable = intstr.FromInt(1)

// PodDisruptionBudgetTemplate defines the template for creating a PodDisruptionBudget.
type PodDisruptionBudgetTemplate struct {
	// ObjectMeta is the metadata of the PDB.
	// The name and namespace provided here are managed by ECK and will be ignored.
	// +kubebuilder:validation:Optional
	ObjectMeta metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the specification of the PDB.
	// +kubebuilder:validation:Optional
	Spec v1beta1.PodDisruptionBudgetSpec `json:"spec,omitempty"`
}

// IsDisabled returns true if the PodDisruptionBudget is explicitly disabled (not nil, but empty).
func (p *PodDisruptionBudgetTemplate) IsDisabled() bool {
	return reflect.DeepEqual(p, &PodDisruptionBudgetTemplate{})
}

// SecretSource defines a data source based on a Kubernetes Secret.
type SecretSource struct {
	// SecretName is the name of the secret.
	SecretName string `json:"secretName"`
	// Entries define how to project each key-value pair in the secret to filesystem paths.
	// If not defined, all keys will be projected to similarly named paths in the filesystem.
	// If defined, only the specified keys will be projected to the corresponding paths.
	// +kubebuilder:validation:Optional
	Entries []KeyToPath `json:"entries,omitempty"`
}

// KeyToPath defines how to map a key in a Secret object to a filesystem path.
type KeyToPath struct {
	// Key is the key contained in the secret.
	Key string `json:"key"`

	// Path is the relative file path to map the key to.
	// Path must not be an absolute file path and must not contain any ".." components.
	// +kubebuilder:validation:Optional
	Path string `json:"path,omitempty"`
}