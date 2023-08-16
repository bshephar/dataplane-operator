
### Custom Resources

* [OpenStackDataPlaneRole](#openstackdataplanerole)

### Sub Resources

* [AnsibleEESpec](#ansibleeespec)
* [DeployStrategySection](#deploystrategysection)
* [NetworkConfigSection](#networkconfigsection)
* [NodeTemplate](#nodetemplate)
* [OpenStackDataPlaneRoleList](#openstackdataplanerolelist)
* [OpenStackDataPlaneRoleSpec](#openstackdataplanerolespec)

#### AnsibleEESpec

AnsibleEESpec is a specification of the ansible EE attributes

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| networkAttachments | NetworkAttachments is a list of NetworkAttachment resource names to pass to the ansibleee resource which allows to connect the ansibleee runner to the given network | []string | false |
| openStackAnsibleEERunnerImage | OpenStackAnsibleEERunnerImage image to use as the ansibleEE runner image | string | false |
| ansibleTags | AnsibleTags for ansible execution | string | false |
| ansibleLimit | AnsibleLimit for ansible execution | string | false |
| ansibleSkipTags | AnsibleSkipTags for ansible execution | string | false |
| extraMounts | ExtraMounts containing files which can be mounted into an Ansible Execution Pod | []storage.VolMounts | false |
| env | Env is a list containing the environment variables to pass to the pod | []corev1.EnvVar | false |
| dnsConfig | DNSConfig for setting dnsservers | *corev1.PodDNSConfig | false |

[Back to Custom Resources](#custom-resources)

#### DeployStrategySection

DeployStrategySection for fields controlling the deployment

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| deploy | Deploy boolean to trigger ansible execution | bool | true |
| ansibleTags | AnsibleTags for ansible execution | string | false |
| ansibleLimit | AnsibleLimit for ansible execution | string | false |
| ansibleSkipTags | AnsibleSkipTags for ansible execution | string | false |

[Back to Custom Resources](#custom-resources)

#### NetworkConfigSection

NetworkConfigSection is a specification of the Network configuration details

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| template | Template - Contains a Ansible j2 nic config template to use when applying node network configuration | string | false |

[Back to Custom Resources](#custom-resources)

#### NodeTemplate

NodeTemplate is a specification of the node attributes

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| hostName | HostName - node name | string | false |
| networkAttachments | NetworkAttachments is a list of NetworkAttachment resource names to pass to the ansibleee resource which allows to connect the ansibleee runner to the given network | []string | false |
| networkConfig | NetworkConfig - Network configuration details. Contains os-net-config related properties. | [NetworkConfigSection](#networkconfigsection) | false |
| networks | Networks - Instance networks | []infranetworkv1.IPSetNetwork | false |
| managementNetwork | ManagementNetwork - Name of network to use for management (SSH/Ansible) | string | false |
| ansibleUser | AnsibleUser SSH user for Ansible connection | string | false |
| ansiblePort | AnsiblePort SSH port for Ansible connection | int | false |
| ansibleVars | AnsibleVars for configuring ansible | map[string]json.RawMessage | false |
| ansibleHost | AnsibleHost host used for Ansible connections | string | false |
| ansibleSSHPrivateKeySecret | AnsibleSSHPrivateKeySecret Private SSH Key secret containing private SSH key for connecting to node. Must be of the form: Secret.data.ssh-privatekey: <base64 encoded private key contents> <https://kubernetes.io/docs/concepts/configuration/secret/#ssh-authentication-secrets> | string | false |
| extraMounts | ExtraMounts containing files which can be mounted into an Ansible Execution Pod | []storage.VolMounts | false |
| userData | UserData  node specific user-data | *corev1.SecretReference | false |
| networkData | NetworkData  node specific network-data | *corev1.SecretReference | false |

[Back to Custom Resources](#custom-resources)

#### OpenStackDataPlaneRole

OpenStackDataPlaneRole is the Schema for the openstackdataplaneroles API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | metav1.ObjectMeta | false |
| spec |  | [OpenStackDataPlaneRoleSpec](#openstackdataplanerolespec) | false |
| status |  | OpenStackDataPlaneStatus | false |

[Back to Custom Resources](#custom-resources)

#### OpenStackDataPlaneRoleList

OpenStackDataPlaneRoleList contains a list of OpenStackDataPlaneRole

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | metav1.ListMeta | false |
| items |  | [][OpenStackDataPlaneRole](#openstackdataplanerole) | true |

[Back to Custom Resources](#custom-resources)

#### OpenStackDataPlaneRoleSpec

OpenStackDataPlaneRoleSpec defines the desired state of OpenStackDataPlaneRole

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| dataPlane | DataPlane name of OpenStackDataPlane for this role | string | false |
| baremetalSetTemplate | BaremetalSetTemplate Template for BaremetalSet for the Role | baremetalv1.OpenStackBaremetalSetSpec | false |
| nodeTemplate | NodeTemplate - node attributes specific to this roles | [NodeTemplate](#nodetemplate) | false |
| nodes | Nodes - Map of Node Names and node specific data. Values here override defaults in the NodeTemplate section. | map[string][NodeTemplate](#nodetemplate) | true |
| preProvisioned | \n\nPreProvisioned - Whether the nodes are actually pre-provisioned (True) or should be preprovisioned (False) | bool | false |
| env | Env is a list containing the environment variables to pass to the pod | []corev1.EnvVar | false |
| deployStrategy | DeployStrategy section to control how the node is deployed | [DeployStrategySection](#deploystrategysection) | false |
| networkAttachments | NetworkAttachments is a list of NetworkAttachment resource names to pass to the ansibleee resource which allows to connect the ansibleee runner to the given network | []string | false |
| services | Services list | []string | true |

[Back to Custom Resources](#custom-resources)
