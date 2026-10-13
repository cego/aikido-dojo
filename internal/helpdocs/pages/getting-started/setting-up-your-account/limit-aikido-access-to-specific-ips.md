# Limit Aikido Access to Specific IPs

IP restrictions let you limit access to your Aikido workspace so only users connecting from trusted IP addresses can use the dashboard. This adds a network-level control on top of your existing login method.

## Why use IP restrictions?

* ✨ Reduce risk from stolen credentials: leaked passwords or tokens are useless if the attacker can't connect from an approved network
* 🛡️ Enforce access from known locations such as your office network or VPN egress
* 🚦 Meet compliance requirements that call for IP-based access controls

## Who can configure IP restrictions?

Only **workspace admins** can view and manage IP restrictions. Non-admin users are blocked from the settings page and API.

## Set up IP restrictions

{% stepper %}
{% step %}

### Open IP restrictions

Go to [**IP Restrictions**](https://app.aikido.dev/settings/advanced/ip-locks) in your workspace settings. Use direct link to access.
{% endstep %}

{% step %}

### Add trusted IP addresses

Click **"Add IP restriction"** and enter one or more IPv4 addresses or CIDR ranges.

* Use a single IP (for example, `203.0.113.45`) for individual machines or VPN egress points
* Use a CIDR range (for example, `203.0.113.0/24`) to allow an entire subnet; netmasks must be between `/24` and `/32`

Click **"Save"** to apply the allowlist.
{% endstep %}

{% step %}

### Verify your current IP is included

Aikido automatically adds your current IP address when you save. This prevents you from locking yourself out of the workspace.

Your current IP is shown with a **"Current IP"** badge in the list. You can't remove your own IP while you're logged in.
{% endstep %}
{% endstepper %}

## How it works

When at least one IP address is on the allowlist, Aikido enables IP restrictions for your workspace. Any user connecting from an IP that isn't on the list sees a **403 Forbidden** error: *"Your IP is not allowed to access this workspace."*

**Disabling IP restrictions:** Remove all IP addresses from the allowlist. When you remove the last entry, IP restrictions are turned off automatically and all users can access the workspace again (subject to their normal login permissions).

**CIDR ranges:** When you add a CIDR range, Aikido expands it into individual IPv4 addresses behind the scenes. For example, `10.0.0.0/24` adds all 256 addresses in that subnet.

## What stays accessible

IP restrictions apply to workspace access through the Aikido dashboard and authenticated API calls. Aikido keeps the following integrations working so your scans, CI checks, and alerts aren't disrupted:

* CI/CD webhooks and check integrations (GitHub, GitLab, Bitbucket, Azure DevOps)
* Jira and Microsoft Teams integrations
* Aikido's background scanning and data-fetch services
* Zen Firewall runtime public API

Aikido Support can still access your workspace when providing support via impersonation.

## Supported formats

| Format       | Example         | Notes                               |
| ------------ | --------------- | ----------------------------------- |
| IPv4 address | `192.168.1.100` | Single trusted IP                   |
| CIDR range   | `10.0.0.0/24`   | Netmask must be `/24` through `/32` |
| IPv6         | Not applicable  | Not supported                       |

## Tips

{% hint style="warning" %}
Before enabling IP restrictions, make sure every team member can reach Aikido from an allowed network. Users on home networks or traveling without VPN access will be blocked.
{% endhint %}

{% hint style="info" %}
IP restrictions for Aikido are separate from IP allowlisting on your Git provider. If you also restrict access to GitHub, GitLab, or another SCM, see [IP Allowlisting for Your Git Provider](/miscellaneous-info/ip-allowlisting-for-your-git-provider.md).
{% endhint %}
