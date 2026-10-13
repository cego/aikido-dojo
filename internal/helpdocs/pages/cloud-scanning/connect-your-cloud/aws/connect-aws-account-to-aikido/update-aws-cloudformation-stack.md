# Update AWS CloudFormation Stack

{% hint style="info" %}
This page is only relevant if you connected your AWS account **before** July 13th, 2026. Accounts connected after this date already have the latest permissions.
{% endhint %}

Aikido has updated the AWS CloudFormation template with expanded permissions. If your account was connected before this change, you'll need to update your existing stack to stay current.

The updated template enables:

* Broader AWS service coverage – visibility across more resource types
* Deeper inspection of managed services – including EKS, Lambda, and Batch
* [Container and workload network reachability analysis](https://help.aikido.dev/container-image-scanning/container-reachability-analysis) – tracing network paths from the internet to your containers

All permissions are **read-only**. `rds:DownloadDBLogFilePortion` is explicitly denied to prevent access to database logs.

Aikido automatically masks sensitive values in your AWS resources before indexing them. This includes secrets in Lambda environment variables, EC2 user data, and similar configurations.

### How to update your stack

{% stepper %}
{% step %}
**Get the updated template from Aikido**

Go to [Clouds → Add AWS](https://app.aikido.dev/clouds/add/aws). Select **"Single"** or **"Full Organization"** and copy the CloudFormation template link.
{% endstep %}

{% step %}
**Open CloudFormation in AWS Console**

Navigate to **CloudFormation** and locate your existing Aikido stack.

If you're updating permissions for an organization connection, make sure you're in your **AWS organization management account**.
{% endstep %}

{% step %}
**Update the stack**

Click **"Update stack"** and choose one of:

* **Create a change set** (recommended) – lets you preview all changes before applying
* **Make a direct update** – faster, but applies changes immediately
  {% endstep %}

{% step %}
**Replace the template**

Select **"Replace existing template"**, then either:

* Paste the template URL from Aikido
* or Upload the file if you downloaded it
  {% endstep %}

{% step %}
**Complete the deployment**

Proceed with the update and finish the stack deployment. If you chose to create a change set, remember to **execute it** after reviewing.
{% endstep %}
{% endstepper %}
