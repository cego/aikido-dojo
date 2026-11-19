# Scanning Hetzner with Aikido

Aikido fully supports protecting workloads on Hetzner through specific integrations. A native integration may be added in the future, but you can already achieve full coverage by combining:<br>

* [**Container image scanning**](https://help.aikido.dev/workflows-and-guides/additional-cloud-providers/scanning-hetzner-with-aikido#container-image-scanning) for OCI-compatible registries hosted on Hetzner
* [**Kubernetes cluster scanning**](https://help.aikido.dev/workflows-and-guides/additional-cloud-providers/scanning-hetzner-with-aikido#kubernetes-cluster-scanning) for self-managed clusters
* [**Kubernetes cluster image scanning**](https://help.aikido.dev/workflows-and-guides/additional-cloud-providers/scanning-hetzner-with-aikido#kubernetes-cluster-image-scanning) for self-managed clusters
* [**Virtual machine scanning**](https://help.aikido.dev/workflows-and-guides/additional-cloud-providers/scanning-hetzner-with-aikido#virtual-machine-scanning) via the Local VM Scanner on Hetzner instances

## Features

### Container image scanning

Scan third-party registries that are OCI-compatible and hosted on Hetzner. Follow the OCI guide below to configure container image scanning

{% content-ref url="../../container-image-scanning/standalone-registries/generic-oci-compatible-registry.md" %}
[generic-oci-compatible-registry.md](https://help.aikido.dev/container-image-scanning/standalone-registries/generic-oci-compatible-registry)
{% endcontent-ref %}

### Kubernetes cluster scanning

If you run your own Kubernetes clusters on Hetzner VMs, you can connect them as generic Kubernetes clusters.

{% content-ref url="../../cloud-scanning/kubernetes-cluster-scanning/" %}
[kubernetes-cluster-scanning](https://help.aikido.dev/cloud-scanning/kubernetes-cluster-scanning)
{% endcontent-ref %}

### Kubernetes cluster image scanning&#x20;

If you run your own Kubernetes clusters on Hetzner VMs, you can scan the images of running containers with Kubernetes image scanning.

{% content-ref url="../../cloud-scanning/kubernetes-cluster-scanning/kubernetes-in-cluster-image-scanning.md" %}
[kubernetes-in-cluster-image-scanning.md](https://help.aikido.dev/cloud-scanning/kubernetes-cluster-scanning/kubernetes-in-cluster-image-scanning)
{% endcontent-ref %}

### Virtual Machine scanning

To scan Virtual Machines on Hetzner, use the Local VM Scanner. It inspects packages, system dependencies and configuration directly on the instance.

{% content-ref url="../../virtual-machine-scanning/local-vm-scanning/" %}
[local-vm-scanning](https://help.aikido.dev/virtual-machine-scanning/local-vm-scanning)
{% endcontent-ref %}

You can roll this out centrally using your usual automation tooling (e.g. Ansible, Terraform-provisioned scripts, or cloud-init) so that new Hetzner instances are automatically enrolled.
