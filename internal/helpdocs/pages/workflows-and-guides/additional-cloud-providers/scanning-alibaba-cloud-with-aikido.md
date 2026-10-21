# Scanning Alibaba Cloud with Aikido

Aikido fully supports protecting workloads on Alibaba Cloud through specific integrations. A native integration may be added in the future, but you can already achieve full coverage by combining:<br>

* [**Container image scanning**](#container-image-scanning) for Alibaba Cloud registry or any other OCI-compatible registry
* [**Kubernetes cluster scanning**](#kubernetes-cluster-scanning) for Alibaba Cloud Managed Kubernetes or self-managed clusters
* [**Kubernetes cluster image scanning**](#kubernetes-cluster-image-scanning) for Alibaba Cloud Managed Kubernetes or self-managed clusters
* [**Virtual machine scanning**](#virtual-machine-scanning) via the Local VM Scanner on Alibaba Cloud instances

## Features

### Container image scanning

Alibaba Cloud’s container registry and most third-party registries you use from Alibaba Cloud are OCI-compatible, so they can be scanned using Aikido.

Create a read-only or pull-only user in Alibaba Cloud registry: <https://www.alibabacloud.com/help/en/acr/user-guide/configure-access-credentials?spm=a2c63.p38356.0.i1#task617>

Follow the OCI guide below to configure container image scanning

{% content-ref url="/pages/WbT1GJFncHssApSoIXH7" %}
[Generic OCI-Compatible Registry](/container-image-scanning/standalone-registries/generic-oci-compatible-registry.md)
{% endcontent-ref %}

### Kubernetes cluster scanning

If you use Alibaba Cloud Managed Kubernetes or run your own Kubernetes clusters on Alibaba Cloud VMs, you can connect them as generic Kubernetes clusters.

{% content-ref url="/pages/HJqrXQMWoUZAKUrWnpkb" %}
[Kubernetes Cluster Scanning](/cloud-scanning/kubernetes-cluster-scanning.md)
{% endcontent-ref %}

### Kubernetes cluster image scanning

If you use Alibaba Cloud Managed Kubernetes or run your own Kubernetes clusters on Alibaba Cloud VMs, you can scan the images of running containers with Kubernetes image scanning.

{% content-ref url="/pages/rSbyejklgqsaULug65TZ" %}
[Kubernetes In-Cluster Image Scanning](/cloud-scanning/kubernetes-cluster-scanning/kubernetes-in-cluster-image-scanning.md)
{% endcontent-ref %}

### Virtual Machine scanning

To scan Virtual Machines on Alibaba Cloud, use the Local VM Scanner. It inspects packages, system dependencies and configuration directly on the instance.

{% content-ref url="/pages/07U0LWsYijoM6UzZxn2i" %}
[Local VM Scanning](/virtual-machine-scanning/local-vm-scanning.md)
{% endcontent-ref %}

You can roll this out centrally using your usual automation tooling (e.g. Ansible, Terraform-provisioned scripts, or cloud-init) so that new Alibaba Cloud instances are automatically enrolled.
