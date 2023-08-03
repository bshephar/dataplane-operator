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

package deployment

import (
	"context"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v3"

	dataplanev1 "github.com/openstack-k8s-operators/dataplane-operator/api/v1beta1"
	infranetworkv1 "github.com/openstack-k8s-operators/infra-operator/apis/network/v1beta1"
	"github.com/openstack-k8s-operators/lib-common/modules/ansible"
	"github.com/openstack-k8s-operators/lib-common/modules/common/configmap"
	"github.com/openstack-k8s-operators/lib-common/modules/common/helper"
	utils "github.com/openstack-k8s-operators/lib-common/modules/common/util"
)

// GenerateRoleInventory yields a parsed Inventory for role
func GenerateRoleInventory(ctx context.Context, helper *helper.Helper,
	instance *dataplanev1.OpenStackDataPlaneNodeSet,
	allIPSets map[string]infranetworkv1.IPSet, dnsAddresses []string) (string, error) {
	var err error
	var configMaps []string

	inventory := ansible.MakeInventory()
	roleNameGroup := inventory.AddGroup(instance.Name)
	err = resolveAnsibleVars(&instance.Spec.NodeTemplate, &ansible.Host{}, &roleNameGroup)
	if err != nil {
		return "", err
	}

	configMaps = append(configMaps, fmt.Sprintf("dataplanerole-%s", instance.Name))
	for nodeName, node := range instance.Spec.NodeTemplate.Nodes {
		// Add the node name to the list of configMaps required.
		// This is required for compatibility with NovaExternalCompute.
		configMaps = append(configMaps, fmt.Sprintf("dataplanenode-%s", nodeName))
		host := roleNameGroup.AddHost(nodeName)
		var dnsSearchDomains []string

		// Use ansible_host if provided else use hostname. Fall back to
		// nodeName if all else fails.
		if node.AnsibleHost != "" {
			host.Vars["ansible_host"] = node.AnsibleHost
		} else if node.HostName != "" {
			host.Vars["ansible_host"] = node.HostName
		} else {
			host.Vars["ansible_host"] = nodeName
		}

		ipSet, ok := allIPSets[nodeName]
		if ok {
			pupulateInventoryFromIPAM(&ipSet, host, dnsAddresses)
			for _, res := range ipSet.Status.Reservation {
				// Build the vars for ips/routes etc
				switch n := res.Network; n {
				case CtlPlaneNetwork:
					host.Vars["ctlplane_ip"] = res.Address
					_, ipnet, err := net.ParseCIDR(res.Cidr)
					if err == nil {
						netCidr, _ := ipnet.Mask.Size()
						host.Vars["ctlplane_subnet_cidr"] = netCidr
					}
					host.Vars["ctlplane_mtu"] = res.MTU
					host.Vars["gateway_ip"] = res.Gateway
					host.Vars["ctlplane_dns_nameservers"] = dnsAddresses
					host.Vars["ctlplane_host_routes"] = res.Routes
					dnsSearchDomains = append(dnsSearchDomains, res.DNSDomain)
				default:
					entry := toSnakeCase(string(n))
					host.Vars[entry+"_ip"] = res.Address
					_, ipnet, err := net.ParseCIDR(res.Cidr)
					if err == nil {
						netCidr, _ := ipnet.Mask.Size()
						host.Vars[entry+"_cidr"] = netCidr
					}
					host.Vars[entry+"_vlan_id"] = res.Vlan
					host.Vars[entry+"_mtu"] = res.MTU
					//host.Vars[string.Join(entry, "_gateway_ip")] = res.Gateway
					host.Vars[entry+"_host_routes"] = res.Routes
					dnsSearchDomains = append(dnsSearchDomains, res.DNSDomain)
				}
				networkConfig := getAnsibleNetworkConfig(instance, nodeName)

				if networkConfig.Template != "" {
					host.Vars["edpm_network_config_template"] = NicConfigTemplateFile
				}

				host.Vars["ansible_user"] = getAnsibleUser(instance, nodeName)
				host.Vars["ansible_port"] = getAnsiblePort(instance, nodeName)
				host.Vars["management_network"] = getAnsibleManagementNetwork(instance, nodeName)
				host.Vars["networks"] = getAnsibleNetworks(instance, nodeName)

				ansibleVarsData, err := getAnsibleVars(helper, instance, nodeName)
				if err != nil {
					return "", err
				}
				for key, value := range ansibleVarsData {
					host.Vars[key] = value
				}
				host.Vars["dns_search_domains"] = dnsSearchDomains
			}
		}

		err = resolveNodeAnsibleVars(&node, &host, &ansible.Group{})
		if err != nil {
			return "", err
		}
	}

	invData, err := inventory.MarshalYAML()
	if err != nil {
		utils.LogErrorForObject(helper, err, "Could not parse Role inventory", instance)
		return "", err
	}
	cmData := map[string]string{
		"inventory": string(invData),
		"network":   string(instance.Spec.NodeTemplate.NetworkConfig.Template),
	}

	var configMapName string
	for _, configMapName = range configMaps {
		cms := []utils.Template{
			// ConfigMap
			{
				Name:         configMapName,
				Namespace:    instance.Namespace,
				Type:         utils.TemplateTypeNone,
				InstanceType: instance.Kind,
				CustomData:   cmData,
				Labels:       instance.ObjectMeta.Labels,
			},
		}
		err = configmap.EnsureConfigMaps(ctx, helper, instance, cms, nil)
	}

	return configMapName, err
}

// pupulateInventoryFromIPAM populates inventory from IPAM
func pupulateInventoryFromIPAM(
	ipSet *infranetworkv1.IPSet, host ansible.Host,
	dnsAddresses []string) {
	var dnsSearchDomains []string
	for _, res := range ipSet.Status.Reservation {
		// Build the vars for ips/routes etc
		switch n := res.Network; n {
		case CtlPlaneNetwork:
			host.Vars["ctlplane_ip"] = res.Address
			_, ipnet, err := net.ParseCIDR(res.Cidr)
			if err == nil {
				netCidr, _ := ipnet.Mask.Size()
				host.Vars["ctlplane_subnet_cidr"] = netCidr
			}
			host.Vars["ctlplane_mtu"] = res.MTU
			host.Vars["gateway_ip"] = res.Gateway
			host.Vars["ctlplane_dns_nameservers"] = dnsAddresses
			host.Vars["ctlplane_host_routes"] = res.Routes
			dnsSearchDomains = append(dnsSearchDomains, res.DNSDomain)
		default:
			entry := toSnakeCase(string(n))
			host.Vars[entry+"_ip"] = res.Address
			_, ipnet, err := net.ParseCIDR(res.Cidr)
			if err == nil {
				netCidr, _ := ipnet.Mask.Size()
				host.Vars[entry+"_cidr"] = netCidr
			}
			host.Vars[entry+"_vlan_id"] = res.Vlan
			host.Vars[entry+"_mtu"] = res.MTU
			host.Vars[entry+"_gateway_ip"] = res.Gateway
			host.Vars[entry+"_host_routes"] = res.Routes
			dnsSearchDomains = append(dnsSearchDomains, res.DNSDomain)
		}
		host.Vars["dns_search_domains"] = dnsSearchDomains
	}
}

// getAnsibleUser returns the string value from the template unless it is set in the node
func getAnsibleUser(instance *dataplanev1.OpenStackDataPlaneNodeSet, nodeName string) string {
	if instance.Spec.NodeTemplate.Nodes[nodeName].AnsibleUser != "" {
		return instance.Spec.NodeTemplate.Nodes[nodeName].AnsibleUser
	}
	return instance.Spec.NodeTemplate.AnsibleUser
}

// getAnsiblePort returns the string value from the template unless it is set in the node
func getAnsiblePort(instance *dataplanev1.OpenStackDataPlaneNodeSet, nodeName string) string {
	if instance.Spec.NodeTemplate.Nodes[nodeName].AnsiblePort > 0 {
		return strconv.Itoa(instance.Spec.NodeTemplate.Nodes[nodeName].AnsiblePort)
	}
	return strconv.Itoa(instance.Spec.NodeTemplate.AnsiblePort)
}

// getAnsibleManagementNetwork returns the string value from the template unless it is set in the node
func getAnsibleManagementNetwork(
	instance *dataplanev1.OpenStackDataPlaneNodeSet,
	nodeName string) string {
	if instance.Spec.NodeTemplate.Nodes[nodeName].ManagementNetwork != "" {
		return instance.Spec.NodeTemplate.Nodes[nodeName].ManagementNetwork
	}
	return instance.Spec.NodeTemplate.ManagementNetwork
}

// getAnsibleNetworkConfig returns a JSON string value from the template unless it is set in the node
func getAnsibleNetworkConfig(instance *dataplanev1.OpenStackDataPlaneNodeSet, nodeName string) dataplanev1.NetworkConfigSection {
	if instance.Spec.NodeTemplate.Nodes[nodeName].NetworkConfig.Template != "" {
		return instance.Spec.NodeTemplate.Nodes[nodeName].NetworkConfig
	}
	return instance.Spec.NodeTemplate.NetworkConfig
}

// getAnsibleNetworks returns a JSON string mapping fixedIP and/or network name to their valules
func getAnsibleNetworks(instance *dataplanev1.OpenStackDataPlaneNodeSet, nodeName string) []infranetworkv1.IPSetNetwork {
	if len(instance.Spec.NodeTemplate.Nodes[nodeName].Networks) > 0 {
		return instance.Spec.NodeTemplate.Nodes[nodeName].Networks
	}
	return instance.Spec.NodeTemplate.Networks
}

// getAnsibleVars returns ansible vars for a node
func getAnsibleVars(
	helper *helper.Helper, instance *dataplanev1.OpenStackDataPlaneNodeSet, nodeName string) (map[string]interface{}, error) {
	// Merge the ansibleVars from the role into the value set on the node.
	// Top level keys set on the node ansibleVars should override top level keys from the role AnsibleVars.
	// However, there is no "deep" merge of values. Only top level keys are comvar matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")

	// Unmarshal the YAML strings into two maps
	var role, node map[string]interface{}
	roleYamlError := yaml.Unmarshal([]byte(instance.Spec.NodeTemplate.AnsibleVars), &role)
	if roleYamlError != nil {
		utils.LogErrorForObject(
			helper,
			roleYamlError,
			fmt.Sprintf("Failed to unmarshal YAML data from role AnsibleVars '%s'",
				instance.Spec.NodeTemplate.AnsibleVars), instance)
		return nil, roleYamlError
	}
	nodeYamlError := yaml.Unmarshal([]byte(instance.Spec.NodeTemplate.Nodes[nodeName].AnsibleVars), &node)
	if nodeYamlError != nil {
		utils.LogErrorForObject(
			helper,
			nodeYamlError,
			fmt.Sprintf("Failed to unmarshal YAML data from node AnsibleVars '%s'",
				instance.Spec.NodeTemplate.Nodes[nodeName].AnsibleVars), instance)
		return nil, nodeYamlError
	}

	if role == nil && node != nil {
		return node, nil
	}
	if role != nil && node == nil {
		return role, nil
	}

	// Merge the two maps
	for k, v := range node {
		role[k] = v
	}
	return role, nil
}

func resolveAnsibleVars(node *dataplanev1.NodeTemplate, host *ansible.Host, group *ansible.Group) error {
	ansibleVarsData := make(map[string]interface{})

	if node.AnsibleUser != "" {
		ansibleVarsData["ansible_user"] = node.AnsibleUser
	}
	if node.AnsiblePort > 0 {
		ansibleVarsData["ansible_port"] = node.AnsiblePort
	}
	if node.ManagementNetwork != "" {
		ansibleVarsData["management_network"] = node.ManagementNetwork
	}
	if node.NetworkConfig.Template != "" {
		ansibleVarsData["edpm_network_config_template"] = NicConfigTemplateFile
	}
	if len(node.Networks) > 0 {
		ansibleVarsData["networks"] = node.Networks
	}

	err := yaml.Unmarshal([]byte(node.AnsibleVars), ansibleVarsData)
	if err != nil {
		return err
	}

	if host.Vars != nil {
		for key, value := range ansibleVarsData {
			host.Vars[key] = value
		}
	}

	if group.Vars != nil {
		for key, value := range ansibleVarsData {
			group.Vars[key] = value
		}
	}

	return nil
}

func resolveNodeAnsibleVars(node *dataplanev1.NodeSection, host *ansible.Host, group *ansible.Group) error {
	ansibleVarsData := make(map[string]interface{})

	if node.AnsibleUser != "" {
		ansibleVarsData["ansible_user"] = node.AnsibleUser
	}
	if node.AnsiblePort > 0 {
		ansibleVarsData["ansible_port"] = node.AnsiblePort
	}
	if node.ManagementNetwork != "" {
		ansibleVarsData["management_network"] = node.ManagementNetwork
	}
	if node.NetworkConfig.Template != "" {
		ansibleVarsData["edpm_network_config_template"] = NicConfigTemplateFile
	}
	if len(node.Networks) > 0 {
		ansibleVarsData["networks"] = node.Networks
	}
	var err error
	for key, val := range node.AnsibleVars {
		var v interface{}
		err = yaml.Unmarshal(val, &v)
		if err != nil {
			return err
		}
		ansibleVarsData[key] = v
	}

	if host.Vars != nil {
		for key, value := range ansibleVarsData {
			host.Vars[key] = value
		}
	}

	if group.Vars != nil {
		for key, value := range ansibleVarsData {
			group.Vars[key] = value
		}
	}

	return nil
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
