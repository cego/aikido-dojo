# Kubernetes Cluster Scanning

This page describes the scanning capabilities for your Kubernetes clusters. It applies to managed Kubernetes environments (e.g., EKS, AKS, and GKE) as well as self-managed/on-prem/OpenShift clusters.

## Why connect your Kubernetes clusters?

While Aikido [scans your infrastructure as code files](https://help.aikido.dev/general-information/sast-by-aikido-supported-languages-and-security-focus), including Kubernetes manifests and Helm charts, connecting your cluster enables Aikido to perform more (and more powerful) checks, as well as validating your actual environments.

Aikido assesses your Kubernetes resources against [a set of rules](https://app.aikido.dev/clouds/checks?cloudCheckType=kubernetes), generating issues for any deviation. You can see the generated issues [in your feed](https://app.aikido.dev/queue?issue_type_filter=kubernetes), by sorting for "Kubernetes Configurations".

Moreover, Aikido can scan the images deployed in your cluster. More on this [below](#image-scanning).

## Getting Started

Go to [the Clouds page](https://app.aikido.dev/clouds), click "Connect Cloud", and choose "Kubernetes". This applies to all Kubernetes environments, including AWS EKS, Azure Kubernetes Service, Google Kubernetes Engine, and other managed Kubernetes services, as well as self-managed on-premises or cloud-based deployments.

In the first step, you will provide the following:

* **Cluster name**: This is used only in Aikido. You can provide any name you see fit, but it must be unique within the Aikido workspace.
* **Excluded namespaces**: Optionally, you can exclude the collection of resources from specific Kubernetes namespaces. An option to quickly add commonly excluded namespaces is available. (e.g., kube-system, kube-public, cloud provider-specific namespaces, etc.).
* **Image scanning**: An additional component will be installed in your cluster responsible for generating SBOMs from your container images. These SBOM are then scanned by Aikido, and you will see any issues impacting these images. See [Kubernetes In-Cluster Image Scanning](https://help.aikido.dev/~/revisions/zVer63pN7HO83qF6ZCQx/cloud-scanning/kubernetes-cluster-scanning/kubernetes-in-cluster-image-scanning).
* **Environment**: This is identical to the purpose of cloud, allowing Aikido to adjust the severity of the findings.

Next, you will install the Aikido agent in your cluster. You will need:

* [The Helm CLI](https://helm.sh/docs/intro/install/)
* Access to the cluster and permissions to install a Helm chart

Follow the in-app steps to install the chart that will deploy the Aikido agent as a deployment (currently running as a single pod). The agent should start in a few seconds and you should be able to finalize the onboarding.

### Updating An Existing Cluster

You can update the excluded namespaces and enable/disable the in-cluster image scanning on a cluster that's already connected.

{% hint style="info" %}
When changing the in-cluster image scanning setting, you must ensure the SBOM collector is enabled or disabled accordingly in the [Kubernetes agent Helm release](https://github.com/AikidoSec/helm-charts/blob/main/kubernetes-agent/values.yaml#L27). This setting is not updated automatically.

In most cases, you can apply the change with:

```shellscript
helm upgrade aikido -n aikido \
  aikido/kubernetes-agent \
  --reuse-values \
  --set sbomCollector.enabled=true # or false
```

{% endhint %}

## OpenShift Support

Aikido fully supports **OpenShift** clusters, including OpenShift Container Platform running on-prem or in managed environments. The agent and SBOM collector work the same way as on other Kubernetes distributions, with a few OpenShift-specific integrations and considerations outlined below:

### Image mirror resolution

OpenShift commonly uses image mirrors (for example via Quay or internal registries). Aikido automatically resolves these mappings so that images are reported and correlated using their original source image, not the mirrored reference.

To achieve this, the agent reads OpenShift image mirror configuration from the cluster, including:

* `ImageContentSourcePolicies`
* `ImageDigestMirrorSets`
* `ImageTagMirrorSets`

This ensures accurate image identification, vulnerability correlation, and avoids duplicated or misleading image entries in the UI.

### Custom CA certificates

Some OpenShift clusters rely on custom or internal certificate authorities. Aikido supports this by allowing you to mount additional CA certificates.

The Helm chart [allows you to](https://github.com/AikidoSec/helm-charts?tab=readme-ov-file#custom-ca-certificates):

* Add extra volumes and volume mounts
* Point the agent to the custom trust bundle using standard environment variables such as `SSL_CERT_FILE` or `SSL_CERT_DIR`

### Image pull secrets

By default, Aikido can access container images using RBAC permissions, which is the simplest setup and does not require copying registry credentials into the Aikido namespace.

For environments where RBAC-based access is not desired or not possible, the agent also supports [mounting image pull secrets directly](https://github.com/AikidoSec/helm-charts/blob/main/kubernetes-agent/values.yaml#L106). This allows explicit control over registry credentials and can be useful in more restricted OpenShift setups.

Both approaches are supported and functionally equivalent. Choose the one that best fits your cluster’s security and operational model.

## FIPS-compliant images

Aikido publishes FIPS-compliant variants of the Kubernetes agent and SBOM collector for FedRAMP environments:

**Kubernetes agent (FIPS)** `public.ecr.aws/aikido-cloud/kubernetes-agent-fips:v1.19.0`

**SBOM collector (FIPS)** `public.ecr.aws/aikido-cloud/kubernetes-sbom-collector-fips:v1.9.0`

To use them, append `-fips` to the repo values in your Helm chart.

## FAQs

#### 1. What can the Aikido agent do in my cluster?

The Aikido agent can only read your Kubernetes resources. Additionally, it has minimal permissions, limited to its own namespace, for self-management. You can inspect the provided permissions [in the Helm chart](https://github.com/AikidoSec/helm-charts/blob/main/kubernetes-agent/templates/agent/cluster_rbac.yaml).

#### 2. What performance impact will the Aikido agent have in my cluster?

The agent's sole purpose is to collect the resources from your Kubernetes cluster. As such, the agent's resource utilization will be negligible. You can customize the [resource requests and limits](https://github.com/AikidoSec/helm-charts/blob/main/kubernetes-agent/values.yaml#L17) when installing the Helm chart.

#### 3. Does my cluster need to be public?

No. The traffic only flows from the agent to the Aikido platform.

#### 4. Is there any maintenance for the Aikido agent?

No. When we release new versions of the agent, it will update itself automatically. Note that the agent can only [update/patch its own deployment and secret](https://github.com/AikidoSec/helm-charts/blob/main/kubernetes-agent/templates/agent/rbac.yaml#L14).

#### 5. Does the Aikido agent read my Kubernetes secrets?

No. It does not have access to Kubernetes secrets, except for image pull secrets (for image scanning), which you can [limit only to specific secrets](https://github.com/AikidoSec/helm-charts/blob/main/kubernetes-agent/values.yaml#L99).

#### 6. What if the Aikido agent cannot communicate with the platform? Will it try indefinitely?

The agent uses an exponential backoff mechanism for resource collection. If it cannot establish the connection with the Aikido platform, it will sit idle, without consuming any additional resources.
