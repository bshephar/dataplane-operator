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
	"github.com/openstack-k8s-operators/lib-common/modules/common/condition"
)

const (
	// DataPlaneNodeReadyMessage ready
	DataPlaneNodeReadyMessage = "DataPlaneNode ready"

	// DataPlaneNodeReadyWaitingMessage ready
	DataPlaneNodeReadyWaitingMessage = "DataPlaneNode not yet ready"

	// DataPlaneNodeErrorMessage error
	DataPlaneNodeErrorMessage = "DataPlaneNode error occurred %s"

	// DataPlaneNodeSetReadyMessage ready
	DataPlaneNodeSetReadyMessage = "DataPlaneNodeSet ready"

	// DataPlaneNodeSetReadyWaitingMessage ready
	DataPlaneNodeSetReadyWaitingMessage = "DataPlaneNodeSet not yet ready"

	// DataPlaneNodeSetErrorMessage error
	DataPlaneNodeSetErrorMessage = "DataPlaneNodeSet error occurred %s"

	// DataPlaneReadyMessage ready
	DataPlaneReadyMessage = "DataPlane ready"

	// DataPlaneReadyWaitingMessage ready
	DataPlaneReadyWaitingMessage = "DataPlane not yet ready"

	// DataPlaneErrorMessage error
	DataPlaneErrorMessage = "DataPlane error occurred %s"

	// ServiceReadyCondition Status=True condition indicates if the
	// service is finished and successful.
	ServiceReadyCondition string = "%s service ready"

	// ServiceReadyMessage ready
	ServiceReadyMessage = "%s service ready"

	// ServiceReadyWaitingMessage not yet ready
	ServiceReadyWaitingMessage = "%s service not yet ready"

	// ServiceErrorMessage error
	ServiceErrorMessage = "Service error occurred %s"

	// SetupReadyCondition - Overall setup condition
	SetupReadyCondition condition.Type = "SetupReady"

	// ConfigureCephClientReadyCondition Status=True condition indicates if the
	// Ceph client configuration is finished and successful.
	ConfigureCephClientReadyCondition condition.Type = "ConfigureCephClientReady"

	// ConfigureCephClientReadyMessage ready
	ConfigureCephClientReadyMessage = "ConfigureCephClient ready"

	// ConfigureCephClientReadyWaitingMessage not yet ready
	ConfigureCephClientReadyWaitingMessage = "ConfigureCephClient not yet ready"

	// ConfigureCephClientErrorMessage error
	ConfigureCephClientErrorMessage = "ConfigureCephClient error occurred %s"

	// NovaComputeReadyCondition Status=True condition indicates nova-compute
	// has been deployed and is ready
	NovaComputeReadyCondition condition.Type = "NovaComputeReady"

	// NovaComputeReadyMessage ready
	NovaComputeReadyMessage = "NovaComputeReady ready"

	// NovaComputeReadyWaitingMessage not yet ready
	NovaComputeReadyWaitingMessage = "NovaComputeReady not yet ready"

	// NovaComputeErrorMessage error
	NovaComputeErrorMessage = "NovaCompute error occurred"

	// NodeSetBareMetalProvisionReadyCondition Status=True condition indicates
	// all baremetal nodes provisioned for the NodeSet.
	NodeSetBareMetalProvisionReadyCondition condition.Type = "NodeSetBaremetalProvisionReady"

	// NodeSetBaremetalProvisionReadyMessage ready
	NodeSetBaremetalProvisionReadyMessage = "NodeSetBaremetalProvisionReady ready"

	// NodeSetBaremetalProvisionReadyWaitingMessage not yet ready
	NodeSetBaremetalProvisionReadyWaitingMessage = "NodeSetBaremetalProvisionReady not yet ready"

	// NodeSetBaremetalProvisionErrorMessage error
	NodeSetBaremetalProvisionErrorMessage = "NodeSetBaremetalProvisionReady error occurred"

	// NodeSetIPReservationReadyCondition Status=True condition indicates
	// IPSets reserved for all nodes in a NodeSet.
	NodeSetIPReservationReadyCondition condition.Type = "NodeSetIPReservationReady"

	// NodeSetIPReservationReadyMessage ready
	NodeSetIPReservationReadyMessage = "NodeSetIPReservationReady ready"

	// NodeSetIPReservationReadyWaitingMessage not yet ready
	NodeSetIPReservationReadyWaitingMessage = "NodeSetIPReservationReady not yet ready"

	// NodeSetIPReservationReadyErrorMessage error
	NodeSetIPReservationReadyErrorMessage = "NodeSetIPReservationReady error occurred"

	// NodeSetDNSDataReadyCondition Status=True condition indicates
	// DNSData created for the NodeSet.
	NodeSetDNSDataReadyCondition condition.Type = "NodeSetDNSDataReady"

	// NodeSetDNSDataReadyMessage ready
	NodeSetDNSDataReadyMessage = "NodeSetDNSDataReady ready"

	// NodeSetDNSDataReadyWaitingMessage not yet ready
	NodeSetDNSDataReadyWaitingMessage = "NodeSetDNSDataReady not yet ready"

	// NodeSetDNSDataReadyErrorMessage error
	NodeSetDNSDataReadyErrorMessage = "NodeSetDNSDataReady error occurred"

	// InputReadyWaitingMessage not yet ready
	InputReadyWaitingMessage = "Waiting for input %s, not yet ready"
)
