# Aikido Terraform Provider

Use the official Aikido Terraform provider to manage Aikido resources as code through the [management API](https://apidocs.aikido.dev/). This is useful if you already provision infrastructure with Terraform and want repository activation, AutoFix settings, and PR checks to follow the same workflow.

The provider is published on both the [Terraform Registry](https://registry.terraform.io/providers/AikidoTerraform/aikido/latest/docs) and the [OpenTofu Registry](https://search.opentofu.org/provider/aikidoterraform/aikido/latest), so it works with Terraform and OpenTofu. Full resource schemas, arguments, and import syntax live in the Terraform Registry documentation.

## What you can manage

The provider currently supports:

* **Repositories**: Activate or deactivate existing code repositories, and set sensitivity, connectivity, and labels
* **AutoFix settings**: Configure workspace AutoFix for dependencies, SAST, and Aikido Pentest findings
* **PR checks**: Configure pull request gating per repository (severity thresholds, scan types, inline comments, Deep Review)
* **Teams and Permissions**: setup teams and user permissions
* **Teams and resource mapping**: map repos and containers to teams
* **Repo/Team Mapping to Task Trackers**

{% hint style="info" %}
The provider manages configuration for repositories that already exist in Aikido (synced from your Git provider). It does not create Git repositories or connect new Git integrations.
{% endhint %}

## Requirements

* Terraform 1.0 or later, or OpenTofu
* An Aikido workspace with API access
* OAuth2 client credentials from the [Aikido REST API integration page](https://app.aikido.dev/settings/integrations/api/aikido/rest)

## Set up the provider

{% stepper %}
{% step %}

### Create API credentials

1. Go to the [Aikido REST API integration page](https://app.aikido.dev/settings/integrations/api/aikido/rest).
2. Create a client ID and client secret.
3. Store the secret securely. You will not be able to view it again after you leave the page.

See the [API authorization docs](https://apidocs.aikido.dev/reference/authorization) for more detail on how Aikido authenticates API clients.
{% endstep %}

{% step %}

### Declare the provider

Add the provider to your Terraform configuration, then run `terraform init`:

{% code title="main.tf" %}

```terraform
terraform {
  required_providers {
    aikido = {
      source  = "AikidoTerraform/aikido"
      version = "~> 1.1"
    }
  }
}

provider "aikido" {
  client_id     = var.aikido_client_id
  client_secret = var.aikido_client_secret
}
```

{% endcode %}

You can also pass credentials with environment variables instead of hardcoding them:

* `AIKIDO_CLIENT_ID`
* `AIKIDO_CLIENT_SECRET`

When those are set, you can use an empty provider block:

```terraform
provider "aikido" {}
```

{% endstep %}

{% step %}

### Target your Aikido region (optional)

By default, the provider uses the EU API (`https://app.aikido.dev/api`). If your workspace is hosted in another region, set `base_url` (or the `AIKIDO_BASE_URL` environment variable) to match the app URL you use to log in:

| Aikido login URL    | Example `base_url`                     |
| ------------------- | -------------------------------------- |
| `app.aikido.dev`    | `https://app.aikido.dev/api` (default) |
| `app.us.aikido.dev` | `https://app.us.aikido.dev/api`        |
| `app.au.aikido.dev` | `https://app.au.aikido.dev/api`        |
| `app.me.aikido.dev` | `https://app.me.aikido.dev/api`        |

```terraform
provider "aikido" {
  client_id     = var.aikido_client_id
  client_secret = var.aikido_client_secret
  base_url      = "https://app.us.aikido.dev/api"
}
```

{% endstep %}
{% endstepper %}

## Example configuration

```terraform
resource "aikido_repository" "app" {
  id     = "12345"
  active = true

  sensitivity  = "sensitive"
  connectivity = "connected"

  labels = ["payments", "production"]
}

resource "aikido_autofix_dependency_settings" "workspace" {
  enabled                      = true
  severity_filter              = "critical_and_high_only"
  repos_scope                  = "all"
  use_aikido_library_for_major = true
}

resource "aikido_repo_pr_checks_settings" "app" {
  code_repo_id = 12345

  minimum_severity                  = "high"
  fail_on_dependency_scan            = true
  fail_on_sast_scan                  = true
  fail_on_iac_scan                   = true
  fail_on_secrets_scan               = true
  fail_on_malware_scan               = true
  post_inline_comments_min_severity  = "critical"

  minimum_license_severity = "high"

  fail_on_code_quality_scan                      = true
  enable_code_quality_scan                       = true
  post_code_quality_inline_comments_min_severity = "medium"
}
```

For every argument, import ID, and resource-specific behavior, use the official docs:

<a href="https://registry.terraform.io/providers/AikidoTerraform/aikido/latest/docs" class="button primary" data-icon="book">Terraform Registry docs</a> <a href="https://search.opentofu.org/provider/aikidoterraform/aikido/latest" class="button secondary" data-icon="book">OpenTofu Registry</a>

## Related docs

* [Aikido management API](https://apidocs.aikido.dev/)
* [AutoFix](/autofix-and-remediation/overview-aikido-autofix.md)
* [PR gating](/pr-and-release-gating/aikido-ci-gating-functionality.md)
