# Connect AWS Organization

Connect your entire AWS Organization to Aikido by onboarding just the management account. Automatically discover and connect all member accounts, simplify setup, and ensure full coverage.

{% hint style="info" %}
This functionality is available only for **Pro** and **Advanced** plans. **Contact us** via chat for more information.
{% endhint %}

If you have an AWS organization with many member accounts, you can connect just the [management account](https://docs.aws.amazon.com/organizations/latest/userguide/orgs_getting-started_concepts.html#management-account) and let Aikido automatically discover and connect the rest of the AWS accounts.

## Why Connect AWS Organization?

By onboarding at the organization level, you benefit from:

* **Faster setup**: You only need to connect the management account.
* **Automatic account discovery**: New member accounts are automatically added to Aikido, including accounts you will create in the future.

## Prerequisites

* You are on the Pro, Advanced, or Enterprise plan.
* You have access to the AWS management account of your AWS organization.
* You have [enabled Trusted Access for AWS CloudFormation StackSets](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/stacksets-orgs-activate-trusted-access.html).

## Getting Started

To connect your AWS organization, select the 'Full AWS Organization' option in the [AWS connection wizard](https://app.aikido.dev/clouds/add/aws). You will need to provide the following:

* **Organization ID**: It should look like 'o-wma21z4agr'.
* **Root ID** or one or more **Organization Unit (OU) IDs**, separated by commas: This option allows you to connect the entire organization (by providing the root ID) or specific parts of it (for example, you may want to connect only the production and staging OUs).
* **Excluded Account IDs**: Optionally, you can exclude specific AWS accounts from being added to Aikido.

You can obtain this information from your [AWS Organization page](https://us-east-1.console.aws.amazon.com/organizations/v2/home/accounts).

Once the setup is complete, you will see all your AWS accounts connected to Aikido within a few minutes.

### Cloud Purpose Determination

The purpose/environment of your AWS accounts is automatically determined based on the name of the parent OU. Aikido looks for terms such as "production", "staging", "uat", etc., and sets the cloud purpose accordingly. If it doesn't find any match, the purpose will be "mixed". You can manually update the purpose of each cloud connection using the "Configure" button.

### ECR/EBS Scanning for Member Accounts

Aikido can automatically configure ECR scanning and/or EBS (EC2) scanning for your AWS member accounts. Enabling these will lead to additional IAM resources being deployed with the CloudFormation stack set (an IAM role and custom policy for each feature) - ensure you set these options according to your needs before opening the CloudFormation page in your AWS account.

If you haven't connected your AWS organization yet, you have the option to enable these during the initial setup. If you don't see the "Enable EBS Scanning" option, please reach out to us.

#### Editing Organization Details Post Onboarding

To update the organization-specific values (e.g., excluded AWS accounts) **after** you have onboarded your AWS organization, you need to:

1. Update the parameters in your CloudFormation stack.
   1. Open the `aikido-security-readonly-org` CloudFormation stack.
   2. "Update stack" -> "Make a direct update".
   3. Go with the default "Use existing template" option.
   4. Update the parameters according to your needs.
   5. On the next page, acknowledge the creation of IAM resources.
   6. Press "Submit". The change should take a few minutes to deploy.
2. In Aikido, go to "Clouds", select "Configure" on the cloud corresponding to your AWS management account, and update the settings.

3. Trigger a scan of the management account in Aikido. This will discover any new AWS accounts. If you've excluded AWS accounts (or removed OUs), please delete the corresponding cloud(s) from Aikido.

## FAQs

1. **Is it secure?**

Yes. Connecting your AWS Organization relies on the same setup we use for connecting individual AWS accounts, with a least-privilege IAM role that requires an external ID. In fact, it is the same template, just deployed using CloudFormation StackSets in each one of your AWS accounts.

2. **If I add an AWS account to my organization, will it appear in Aikido?**

Yes. Aikido scans your AWS organization every time it scans your management account, and automatically connects new AWS accounts. This process is facilitated by AWS CloudFormation StackSets that automatically creates the required IAM role and policy in your AWS accounts.

3. **I just added a new AWS account to my org, and it did not show up in Aikido.**

If Aikido scanned your AWS management account (you can manually scan it) and the new account still does not appear, you may have reached the limit of cloud accounts for your plan. Contact us to increase your limit.

4. **What happens if I suspend an AWS account or remove it from the AWS org?**

Aikido will detect that the account is no longer active or part of the organization and will mark the corresponding connection as "not reachable". You will see this on the [clouds page](https://app.aikido.dev/clouds).

5. **Can I use Terraform to set up the required AWS resources?**

Yes. You can use [our Terraform module](https://github.com/AikidoSec/aws-native-terraform-module). The primary variant creates an AWS CloudFormation StackSet. If StackSets is not a option for you, consider using [the IAM submodule](https://github.com/AikidoSec/aws-native-terraform-module/tree/main/modules/iam-roles) to create the role in each AWS account. Note that you still need to follow the in-app steps (except the CloudFormation part).
