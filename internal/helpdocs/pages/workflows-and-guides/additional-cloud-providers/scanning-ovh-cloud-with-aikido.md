# Scanning OVH cloud with Aikido

Aikido fully supports protecting workloads on OVHcloud through specific integrations. A native integration may be added in the future, but you can already achieve full coverage by combining:<br>

* [**Container image scanning**](#container-image-scanning) for OVHcloud’s registry or any other OCI-compatible registry
* [**Kubernetes cluster scanning**](#kubernetes-cluster-scanning) for OVHcloud Managed Kubernetes or self-managed clusters
* [**Kubernetes cluster image scanning**](#kubernetes-cluster-image-scanning) for OVHcloud Managed Kubernetes or self-managed clusters
* [**Virtual machine scanning**](#virtual-machine-scanning) via the Local VM Scanner on OVHcloud instances

## Features

### Container image scanning

OVHcloud’s container registry and most third-party registries you use from OVHcloud are OCI-compatible, so they can be scanned using Aikido.

Create a read-only or pull-only user in OVH registry: <https://help.ovhcloud.com/csm/en-public-cloud-private-registry-creation?id=kb_article_view&sysparm_article=KB0050328>

Follow the OCI guide below to configure container image scanning

{% content-ref url="/pages/WbT1GJFncHssApSoIXH7" %}
[Generic OCI-Compatible Registry](/container-image-scanning/standalone-registries/generic-oci-compatible-registry.md)
{% endcontent-ref %}

### Kubernetes cluster scanning

If you use OVHcloud Managed Kubernetes or run your own Kubernetes clusters on OVHcloud VMs, you can connect them as generic Kubernetes clusters.

{% content-ref url="/pages/HJqrXQMWoUZAKUrWnpkb" %}
[Kubernetes Cluster Scanning](/cloud-scanning/kubernetes-cluster-scanning.md)
{% endcontent-ref %}

### Kubernetes cluster image scanning

If you use OVHcloud Managed Kubernetes or run your own Kubernetes clusters on OVHcloud VMs, you can scan the images of running containers with Kubernetes image scanning.

{% content-ref url="/pages/rSbyejklgqsaULug65TZ" %}
[Kubernetes In-Cluster Image Scanning](/cloud-scanning/kubernetes-cluster-scanning/kubernetes-in-cluster-image-scanning.md)
{% endcontent-ref %}

### Virtual Machine scanning

To scan Virtual Machines on OVHcloud, use the Local VM Scanner. It inspects packages, system dependencies and configuration directly on the instance.

{% content-ref url="/pages/07U0LWsYijoM6UzZxn2i" %}
[Local VM Scanning](/virtual-machine-scanning/local-vm-scanning.md)
{% endcontent-ref %}

You can roll this out centrally using your usual automation tooling (e.g. Ansible, Terraform-provisioned scripts, or cloud-init) so that new OVHcloud instances are automatically enrolled.
