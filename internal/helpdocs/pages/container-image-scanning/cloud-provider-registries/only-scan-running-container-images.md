# Only Scan Running Container Images

Aikido cuts noise by automatically deactivating container images it can't find running in your connected cloud or Kubernetes environments, so you're not scanning and getting findings for old, unused images.

#### **Use Cases**

* Registries full of stale builds: Old feature-branch images, superseded releases, and abandoned CI artifacts pile up in registries over time. Each one keeps generating findings even though nothing depends on it anymore.
* Alert fatigue from irrelevant risk: When most of your open findings come from images nobody runs, it's harder to spot and prioritize the ones that actually matter.
* Migrations and registry cleanup: Teams moving between registries or consolidating infrastructure often leave old copies behind - this keeps them from cluttering your active issue count.

#### How It Works

* Aikido matches each image against your connected runtime and cloud environments.
* Images it can't match to a running deployment are deactivated: they stop generating new issues, but stay visible in Aikido with scan history.
* If Aikido later detects an image running again, it's reactivated automatically.
* Manually activated or deactivated images are unaffected - this setting won't override a state you've set yourself.

#### Prerequisites

* Connect your cloud and Kubernetes environments before enabling this setting - without them, Aikido has no runtime data to match images against, and can't tell which images are safe to deactivate.

{% hint style="info" %}
If you're on AWS and connected your account before July 13th, 2026, also [update your CloudFormation stack](https://help.aikido.dev/cloud-scanning/connect-your-cloud/aws/connect-aws-account-to-aikido/update-aws-cloudformation-stack)
{% endhint %}

#### Enabling **Only Scan Running Images**

* Navigate to [**Settings** > **Containers**](https://app.aikido.dev/settings/container-image-registry)
* Click the triple dots on the registry you want to configure, and click '**Edit Registry Configuration**'
* In the **Registry Settings** modal, toggle on **Only scan running images**
* Click **Save**

#### **Notes**

* If the same image is discovered through both Kubernetes scanning and registry discovery, the Kubernetes-reported version takes precedence. The registry-discovered version still exists in Aikido, but gets automatically deactivated to avoid duplicate scanning - you may see it listed as a separate, deactivated container.
* If enabling this would deactivate the majority of the registry's running images (Aikido uses a 75% threshold), you'll see a warning first - this usually means your cloud or Kubernetes environments aren't fully connected yet. Connect them for accurate results, or enable it anyway.
