/*
Copyright 2023.

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

package v1beta1

import (
	"fmt"

	condition "github.com/openstack-k8s-operators/lib-common/modules/common/condition"
	"github.com/openstack-k8s-operators/lib-common/modules/storage"
	baremetalv1 "github.com/openstack-k8s-operators/openstack-baremetal-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OpenStackDataPlaneRoleSpec defines the desired state of OpenStackDataPlaneRole
type OpenStackDataPlaneRoleSpec struct {
	// +kubebuilder:validation:Optional
	// DataPlane name of OpenStackDataPlane for this role
	DataPlane string `json:"dataPlane,omitempty"`
	
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:io.kubernetes:Secret"}
	// AnsibleSSHPrivateKeySecret Private SSH Key secret containing private SSH
	// key for connecting to node. Must be of the form:
	// Secret.data.ssh-privatekey: <base64 encoded private key contents>
	// <https://kubernetes.io/docs/concepts/configuration/secret/#ssh-authentication-secrets>
	AnsibleSSHPrivateKeySecret string `json:"ansibleSSHPrivateKeySecret,omitempty"`

	// +kubebuilder:validation:Optional
	// BaremetalSetTemplate Template for BaremetalSet for the Role
	BaremetalSetTemplate baremetalv1.OpenStackBaremetalSetSpec `json:"baremetalSetTemplate,omitempty"`

	// +kubebuilder:validation:Required
	// NodeTemplate - node attributes specific to nodes defined by this resource. These
	// attributes can be overriden at the individual node level, else take their defaults
	// from valus in this section.
	NodeTemplate NodeTemplate `json:"nodes"`

	// +kubebuilder:validation:Optional
	//
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	// PreProvisioned - Whether the nodes are actually pre-provisioned (True) or should be
	// preprovisioned (False)
	PreProvisioned bool `json:"preProvisioned,omitempty"`

	// Env is a list containing the environment variables to pass to the pod
	Env []corev1.EnvVar `json:"env,omitempty"`

	// +kubebuilder:validation:Optional
	// DeployStrategy section to control how the node is deployed
	DeployStrategy DeployStrategySection `json:"deployStrategy,omitempty"`

	// +kubebuilder:validation:Optional
	// NetworkConfig - Network configuration details. Contains os-net-config
	// related properties.
	NetworkConfig NetworkConfigSection `json:"networkConfig"`
	
	// +kubebuilder:validation:Optional
	// NetworkAttachments is a list of NetworkAttachment resource names to pass to the ansibleee resource
	// which allows to connect the ansibleee runner to the given network
	NetworkAttachments []string `json:"networkAttachments,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default={configure-network,validate-network,install-os,configure-os,run-os}
	// Services list
	Services []string `json:"services"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+operator-sdk:csv:customresourcedefinitions:displayName="OpenStack Data Plane Role"
// +kubebuilder:resource:shortName=osdprole;osdproles
//+kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.conditions[0].status",description="Status"
//+kubebuilder:printcolumn:name="Message",type="string",JSONPath=".status.conditions[0].message",description="Message"

// OpenStackDataPlaneRole is the Schema for the openstackdataplaneroles API
type OpenStackDataPlaneRole struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenStackDataPlaneRoleSpec `json:"spec,omitempty"`
	Status OpenStackDataPlaneStatus   `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// OpenStackDataPlaneRoleList contains a list of OpenStackDataPlaneRole
type OpenStackDataPlaneRoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenStackDataPlaneRole `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OpenStackDataPlaneRole{}, &OpenStackDataPlaneRoleList{})
}

// IsReady - returns true if the DataPlane is ready
func (instance OpenStackDataPlaneRole) IsReady() bool {
	return instance.Status.Conditions.IsTrue(condition.DeploymentReadyCondition)
}

// InitConditions - Initializes Status Conditons
func (instance *OpenStackDataPlaneRole) InitConditions() {
	instance.Status.Conditions = condition.Conditions{}

	cl := condition.CreateList(
		condition.UnknownCondition(condition.DeploymentReadyCondition, condition.InitReason, condition.InitReason),
		condition.UnknownCondition(condition.InputReadyCondition, condition.InitReason, condition.InitReason),
		condition.UnknownCondition(SetupReadyCondition, condition.InitReason, condition.InitReason),
		condition.UnknownCondition(RoleBareMetalProvisionReadyCondition, condition.InitReason, condition.InitReason),
		condition.UnknownCondition(RoleIPReservationReadyCondition, condition.InitReason, condition.InitReason),
		condition.UnknownCondition(RoleDNSDataReadyCondition, condition.InitReason, condition.InitReason),
	)

	if instance.Spec.Services != nil {
		for _, service := range instance.Spec.Services {
			readyCondition := condition.Type(fmt.Sprintf(ServiceReadyCondition, service))
			cl = append(cl, *condition.UnknownCondition(readyCondition, condition.InitReason, condition.InitReason))
		}
	}

	haveCephSecret := false
	for _, node := range instance.Spec.NodeTemplate.Nodes {
		for _, extraMount := range node.ExtraMounts {
			if extraMount.ExtraVolType == "Ceph" {
				haveCephSecret = true
				break
			}
		}
	}

	if haveCephSecret {
		cl.Set(condition.UnknownCondition(ConfigureCephClientReadyCondition, condition.InitReason, condition.InitReason))

	}
	instance.Status.Conditions.Init(&cl)
	instance.Status.Deployed = false
}

// GetAnsibleEESpec - get the fields that will be passed to AEE
func (instance OpenStackDataPlaneRole) GetAnsibleEESpec() AnsibleEESpec {
	var extraMounts []storage.VolMounts 
	for _, node := range instance.Spec.NodeTemplate.Nodes {
		for _, extraMount := range node.ExtraMounts {
			extraMounts = append(extraMounts, extraMount)
			}
		}
	
	return AnsibleEESpec{
		NetworkAttachments: instance.Spec.NetworkAttachments,
		AnsibleTags:        instance.Spec.DeployStrategy.AnsibleTags,
		AnsibleLimit:       instance.Spec.DeployStrategy.AnsibleLimit,
		AnsibleSkipTags:    instance.Spec.DeployStrategy.AnsibleSkipTags,
		ExtraMounts:        extraMounts,
		Env:                instance.Spec.Env,
	}
}
