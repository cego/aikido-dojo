# IP Allowlisting for Your Git Provider

Your source code repositories contain some of your most sensitive assets: code, secrets history, CI/CD configurations, and deployment workflows. By default, anyone with valid credentials can access them from anywhere in the world.

IP allowlisting adds a network-level control on top of credential-based authentication. You define a list of trusted IP addresses (your office network, VPN egress, CI/CD servers), and access from any other location is blocked -- even with a valid token or SSH key.

**Why this matters:**

* A stolen token or leaked credential becomes useless if the attacker can't connect from a trusted network
* It limits your attack surface to known, controlled locations
* It satisfies access control requirements for SOC 2, ISO 27001, and similar compliance frameworks

**What to add to your allowlist:**

* Your office network IP range
* Your VPN egress IP(s)
* Your CI/CD runner IPs (required -- pipelines will break otherwise)
* Any IP used by third-party integrations like Aikido

**Note:** This feature requires a paid plan on most platforms and works best when your team connects through a consistent network or VPN.

***

**How to set it up on your platform:**

* **GitHub** (Enterprise Cloud): [Managing allowed IP addresses for your organization](https://docs.github.com/en/enterprise-cloud@latest/organizations/keeping-your-organization-secure/managing-security-settings-for-your-organization/managing-allowed-ip-addresses-for-your-organization)
* **GitLab** (Premium/Ultimate): [Restrict group access by IP address](https://docs.gitlab.com/user/group/access_and_permissions/)
* **Bitbucket Cloud** (Premium): [Control access to private content in a workspace](https://support.atlassian.com/bitbucket-cloud/docs/control-access-to-private-content-in-a-workspace/)
* **Azure DevOps**: [Conditional Access policies](https://learn.microsoft.com/en-us/azure/devops/organizations/accounts/conditional-access-policies?view=azure-devops) (via Microsoft Entra ID)
