# Scanning OVH cloud with Aikido

Aikido fully supports protecting workloads on OVHcloud through specific integrations. A native integration may be added in the future, but you can already achieve full coverage by combining:<br>

* [**Container image scanning**](https://help.aikido.dev/workflows-and-guides/additional-cloud-providers/scanning-ovh-cloud-with-aikido#container-image-scanning) for OVHcloud’s registry or any other OCI-compatible registry
* [**Kubernetes cluster scanning**](https://help.aikido.dev/workflows-and-guides/additional-cloud-providers/scanning-ovh-cloud-with-aikido#kubernetes-cluster-scanning) for OVHcloud Managed Kubernetes or self-managed clusters
* [**Kubernetes cluster image scanning**](https://help.aikido.dev/workflows-and-guides/additional-cloud-providers/scanning-ovh-cloud-with-aikido#kubernetes-cluster-image-scanning) for OVHcloud Managed Kubernetes or self-managed clusters
* [**Virtual machine scanning**](https://help.aikido.dev/workflows-and-guides/additional-cloud-providers/scanning-ovh-cloud-with-aikido#virtual-machine-scanning) via the Local VM Scanner on OVHcloud instances

## Features

### Container image scanning

OVHcloud’s container registry and most third-party registries you use from OVHcloud are OCI-compatible, so they can be scanned using Aikido.

Create a read-only or pull-only user in OVH registry: [https://help.ovhcloud.com/csm/en-public-cloud-private-registry-creation?id=kb\_article\_view\&sysparm\_article=KB0050328](https://help.ovhcloud.com/csm/en-public-cloud-private-registry-creation?id=kb_article_view\&sysparm_article=KB0050328)

Follow the OCI guide below to configure container image scanning

{% content-ref url="../../container-image-scanning/standalone-registries/generic-oci-compatible-registry.md" %}
[generic-oci-compatible-registry.md](https://help.aikido.dev/container-image-scanning/standalone-registries/generic-oci-compatible-registry)
{% endcontent-ref %}

### Kubernetes cluster scanning

If you use OVHcloud Managed Kubernetes or run your own Kubernetes clusters on OVHcloud VMs, you can connect them as generic Kubernetes clusters.

{% content-ref url="../../cloud-scanning/kubernetes-cluster-scanning/" %}
[kubernetes-cluster-scanning](https://help.aikido.dev/cloud-scanning/kubernetes-cluster-scanning)
{% endcontent-ref %}

### Kubernetes cluster image scanning&#x20;

If you use OVHcloud Managed Kubernetes or run your own Kubernetes clusters on OVHcloud VMs, you can scan the images of running containers with Kubernetes image scanning.

{% content-ref url="../../cloud-scanning/kubernetes-cluster-scanning/kubernetes-in-cluster-image-scanning.md" %}
[kubernetes-in-cluster-image-scanning.md](https://help.aikido.dev/cloud-scanning/kubernetes-cluster-scanning/kubernetes-in-cluster-image-scanning)
{% endcontent-ref %}

### Virtual Machine scanning

To scan Virtual Machines on OVHcloud, use the Local VM Scanner. It inspects packages, system dependencies and configuration directly on the instance.

{% content-ref url="../../virtual-machine-scanning/local-vm-scanning/" %}
[local-vm-scanning](https://help.aikido.dev/virtual-machine-scanning/local-vm-scanning)
{% endcontent-ref %}

You can roll this out centrally using your usual automation tooling (e.g. Ansible, Terraform-provisioned scripts, or cloud-init) so that new OVHcloud instances are automatically enrolled.
