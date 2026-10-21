# Cloud Discovery: Search Asset Inventory

### How to use <a href="#how-to-use" id="how-to-use"></a>

[Search your cloud inventory](https://app.aikido.dev/clouds/assets) by describing what you are looking for **in natural language** and let the system figure out how to find the relevant asset. Use the examples below to get started or to get inspiration.

#### Simple prompts <a href="#simple-prompts" id="simple-prompts"></a>

| Prompt                                | Why It Matters                                                                                                                               |
| ------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------- |
| public s3 buckets                     | Public S3 buckets are frequently misconfigured and can lead to data exposure or leaks. There are also multiple ways to make a bucket public. |
| buckets outside eu                    | Helps enforce data residency compliance (e.g., GDPR), ensuring sensitive data doesn't leave allowed regions.                                 |
| users without mfa                     | Accounts without Multi-Factor Authentication are vulnerable to account takeover via credential theft.                                        |
| users with programmatic access        | Identifies users who can interact with the cloud via API keys. These credentials are a common target for attackers.                          |
| databases without deletion protection | Prevents accidental or malicious deletion of critical databases.                                                                             |

#### Networking prompts <a href="#networking-prompts" id="networking-prompts"></a>

| Prompt                                            | Why It Matters                                                                                                |
| ------------------------------------------------- | ------------------------------------------------------------------------------------------------------------- |
| EC2 instances with open management ports          | Ports like SSH (22) and RDP (3389) open to the internet are major attack vectors for unauthorized access.     |
| RDS databases allowing traffic from ec2 instances | Helps identify trust relationships and lateral movement paths that attackers could exploit.                   |
| Lambda functions not running in VPCs              | Functions outside VPCs may lack network controls and expose sensitive traffic to the public internet.         |
| ec2 instances that might host databases           | Helps identify data stores that may need additional protection or monitoring, even if not explicitly labeled. |
| lambdas with access to VPC endpoints              | Misused Lambda functions with VPC access can interact with sensitive internal services or databases.          |

#### IAM prompts <a href="#iam-prompts" id="iam-prompts"></a>

| Prompt                                   | Why It Matters                                                                                  |
| ---------------------------------------- | ----------------------------------------------------------------------------------------------- |
| ec2 instances with access to s3 buckets  | Detects possible data exfiltration paths via overly-permissive IAM roles.                       |
| lambdas that can create users            | Functions with privilege to create users can be abused for persistence or privilege escalation. |
| iam roles accessible from other accounts | Cross-account access increases your attack surface and may be unmonitored.                      |
| users with admin privileges              | Overprivileged users are a primary cause of security misconfigurations and insider threats.     |
| overprivileged IAM roles                 | Detects roles with excessive permissions that exceed least-privilege best practices.            |

#### CVEs/EOL issues <a href="#cveseol-issues" id="cveseol-issues"></a>

| Prompt                                           | Why It Matters                                                                                                                                                                                                        |
| ------------------------------------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| ec2 instances vulnerable to CVE-2025-21613       | Allows targeted remediation of known, high-risk vulnerabilities in your infrastructure.                                                                                                                               |
| ec2 instances running outdated OS                | Legacy systems often lack critical security patches and support, increasing risk.                                                                                                                                     |
| VMs with outdated python                         | Outdated runtimes can be vulnerable and incompatible with modern security libraries.                                                                                                                                  |
| VM with critical vulnerabilities                 | Prioritizes the remediation of VMs with the highest likelihood of being exploited, based on vulnerability severity.                                                                                                   |
| ec2 instances vulnerable to log4shell            | Specific vulnerability targeting ensures you can patch critical zero-days quickly and thoroughly.                                                                                                                     |
| containers with high or critical vulnerabilities | There are multiple ways to run containers in the cloud, including ECS, Lambda, App Runner, and others. With this, you can identify and prioritize containers impacted by significant vulnerabilities.                 |
| containers with outdated python                  | Quickly identify containers across all services that use outdated packages (python in this case).                                                                                                                     |
| containers with eol distros                      | Containers can also run unmaintained operating systems, typically part of the base image. Updating these will significantly increase the security of your containers, likely addressing many vulnerabilities as well. |

#### Advanced: combine them <a href="#advanced-combine-them" id="advanced-combine-them"></a>

| Prompt                                                                              | Why It Matters                                                                                                                  |
| ----------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------- |
| show me public ec2 instances vulnerable to CVE-2025-21613 with access to s3 buckets | Models an end-to-end attack chain: public exposure + vulnerability + access to sensitive data.                                  |
| lambda functions created manually                                                   | Manual provisioning can bypass IaC guardrails or compliance checks.                                                             |
| functions exposed to the internet with admin permissions                            | Serverless resources with admin rights and internet exposure can be abused for privilege escalation and data exfiltration.      |
| my riskiest datastores                                                              | Lets Aikido surface the highest-risk data assets based on exposure, vulnerability, and privilege to help prioritize protection. |

### Prompting Best Practices <a href="#prompting-best-practices" id="prompting-best-practices"></a>

* **Don't struggle to find the perfect terms:** There are no predefined prompts, terms, or rules to follow. You can describe anything you want to see from your cloud environment, and let Aikido figure out what it needs to search and generate the queries.
* **Prompt engineering:** If you are unsatisfied with the result, try to be more specific about what you want. Example

  ❌functions with access to rds

  ✅lambda functions with network access to rds instances
* **Use single keywords for broad discovery:** Enter a single keyword, and Aikido will run a broad text search across your assets. For example, searching for a user’s name will return:
  * user with the same or a similar name
  * any groups they belong to,\
    assets where their name appears in tags,\
    and policies that mention them.

**Tip: Use keywords such as usernames, instance names, or tag values to quickly find relevant assets and permissions.** For example, typing `alice` might show you her user profile, groups she's in, EC2 instances tagged with her name, and IAM policies that mention her, all in one search.

### How It Works (Extended) <a href="#how-it-works-extended" id="how-it-works-extended"></a>

After you [connect your cloud environments](https://help.aikido.dev/cloud-scanning/connect-your-cloud), Aikido automatically builds an asset inventory, giving you visibility across your accounts and regions. Aikido synchronizes the inventory as part of the cloud scans (scheduled and manually triggered) by calling the cloud provider's APIs. Additionally, it enriches the graph with data from the Aikido platform, such as vulnerabilities and end-of-life issues found on [virtual machines](https://help.aikido.dev/virtual-machine-scanning/virtual-machines-overview) and [containers](https://help.aikido.dev/container-image-scanning/container-image-scanning-overview).

You can access the cloud asset inventory by going to **Clouds -> Discovery** tab. This shows you the assets from all your connected clouds. Alternatively, you can navigate to a specific cloud and, on the Discovery tab, you will see only the assets from that cloud.

To search your cloud inventory, you describe what you want in **natural language** and let the system figure out how to find the relevant assets. Aikido translates the prompt into one or more steps, depending on the complexity of the prompt. It then shows you the intermediate results as it implements the steps to achieve the final result.

In the example from above, for the prompt "show me EC2 instances with access to S3 buckets", Aikido looked for EC2 instances with IAM roles (attached through instance profiles), found the IAM roles with access to S3 buckets (whether granted through inline or attached policies), found the bucket policies granting access to IAM roles, and combined these in the final result.

You can see all steps, their summary and a list of assets in the first tab. If you have multiple cloud providers connected (AWS, Azure, and Google Cloud), multiple tabs related to the specific provider will be shown.

The final results will be shown in the tab on the right and will include the full list of assets that apply to your search. This is also the case for multi-cloud search - all AWS, Azure, and Google Cloud assets together.

#### Extra notes <a href="#extra-notes" id="extra-notes"></a>

* If you search for one word, Aikido performs a text search, allowing you to find assets by name or other fields. For example, if you search for the name of a user, Aikido will return that user and any group the user is a member of, assets for which the user is mentioned in the tags, as well as policies referencing the user by ARN.
* Aikido caches the prompts, serving subsequent searches much quicker. It also displays your previous searches (only for your user).
