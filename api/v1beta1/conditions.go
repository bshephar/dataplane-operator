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
	// DataPlaneNodeSetErrorMessage error
	DataPlaneNodeSetErrorMessage = "DataPlaneNodeSet error occurred %s"

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

	// RoleBareMetalProvisionReadyCondition Status=True condition indicates
	// all baremetal nodes provisioned for the Role.
	RoleBareMetalProvisionReadyCondition condition.Type = "RoleBaremetalProvisionReady"

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

	// NodeSetDeploymentReadyCondition Status=True condition indicates if the
	// NodeSet Deployment is finished and successful.
	NodeSetDeploymentReadyCondition condition.Type = "NodeSetDeploymentReady"

	// NodeSetDeploymentReadyMessage ready
	NodeSetDeploymentReadyMessage = "%s Deployment ready"

	// NodeSetDeploymentReadyWaitingMessage not yet ready
	NodeSetDeploymentReadyWaitingMessage = "%s Deployment not yet ready"

	// NodeSetDeploymentErrorMessage error
	NodeSetDeploymentErrorMessage = "%s Deployment error occurred %s"

	// NodeSetServiceDeploymentReadyCondition Status=True condition indicates if the
	// NodeSet Deployment is finished and successful.
	NodeSetServiceDeploymentReadyCondition string = "%s %s Deployment ready"

	// NodeSetServiceDeploymentReadyMessage ready
	NodeSetServiceDeploymentReadyMessage = "%s %s Deployment ready"

	// NodeSetServiceDeploymentReadyWaitingMessage not yet ready
	NodeSetServiceDeploymentReadyWaitingMessage = "%s %s Deployment not yet ready"

	// NodeSetServiceDeploymentErrorMessage error
	NodeSetServiceDeploymentErrorMessage = "%s %s Deployment error occurred"
)
