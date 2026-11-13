# Guide Agents with Repository Files

Aikido agents read markdown in your repositories. This applies to all AI features, such as pentests, Code Quality, Deep Review, and AI Code Analysis. Use the files to explain how the system is built, which behavior is intentional, and which risks you already accept.

Commit that guidance in the repository when it should apply every time Aikido reviews the code. Use notes on a single run when the context is missing from the source and only applies to that run.

## Add a markdown file

Put standing instructions in `AGENTS.md` at the repository root. Aikido also checks these files for accepted risk:

* `CLAUDE.md`
* `.claude/`
* `.cursor/`

A risk these files already accept stays unreported, unless a change widens that risk.

Keep `AGENTS.md` at the repository root. `README.md`, `ARCHITECTURE.md`, and `SECURITY.md` still help for architecture and threat-model detail. Link them from `AGENTS.md` when that detail lives in another file.

On **Exclude Paths**, exclude unrelated applications or services. Leave documentation included so Aikido can read it.

{% code title="AGENTS.md" overflow="wrap" %}
```markdown
# Security context

## Trust boundaries

The public API is `api/`. Handlers under `internal/` are only reachable from the service mesh. Missing authentication on `internal/` is accepted. Report the same pattern under `api/`.

## Authentication

End-user requests are authenticated at the gateway. This service trusts the `X-User-Id` header set by that gateway.

## Accepted risks

`scripts/migrations/` is one-off operational code. Do not report missing error handling there.
```
{% endcode %}

## Write instructions agents can follow

Name the path, role, or flow, and say what the agent should treat as intended. Also say what should still be reported.

* **Be specific.** "Admins in `billing/refunds` may read any tenant. Report the same check missing on any other route."
* **Separate intended behavior from bugs.** "Webhook signatures are verified in the gateway, before traffic reaches this service."
* **Use headings.** One topic per section, so an agent can find the relevant part.
* **Keep secrets out.** Do not commit credentials, tokens, customer data, or production connection strings.
* **Keep exceptions narrow.** "Ignore security issues" hides real findings. Name the folder or the behavior.

## Choose where the context goes

| Where                      | When to use it                                                                                                                                        |
| -------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- |
| Markdown in the repository | Context that should apply every time Aikido reviews this repository: architecture, trust boundaries, accepted risks, and how a flow is meant to work. |
| Notes on a single run      | Context that is missing from the source and applies to that run only. Up to 2000 characters. Runtime configuration, a proxy, or a one-off exception.  |
| **Exclude Paths**          | Drop an unrelated application in a monorepo. Leave documentation included.                                                                            |

Notes apply to that run. A markdown file stays in the repository, so the next review can use it too. Put the durable guidance in the file. Use notes for what the file does not cover.

## What to include

* How requests enter, and which folders are production code.
* Trust boundaries, authentication, authorization, and multi-tenant rules.
* Accepted risks, and what should still be reported.
* Data that is intentionally public or internal, and where a gateway redacts fields.
* Trusted admin, keeper, and oracle roles when the repository contains contracts. Put target chains and deployed addresses in the run notes when they are not in the repository.
* How a cloud account is meant to be segmented, and which exceptions are intentional. Markdown in the infrastructure as code (IaC) repositories explains how a live setting is defined.
* Quality conventions the team already accepts, such as where strict error handling does not apply.

## FAQ

### Does Aikido load every markdown file?

Aikido checks `AGENTS.md`, `CLAUDE.md`, `.claude/`, and `.cursor/` for accepted risk. Keep `AGENTS.md` short and at the repository root, and link longer write-ups from there.

### Do these files replace rules and ignores?

These files guide the agents. [Static application security testing (SAST)](https://help.aikido.dev/code-scanning/scanning-practices/sast-by-aikido-supported-languages-and-security-focus) rules, [custom code rules](https://help.aikido.dev/code-quality/add-custom-code-rules), and [`.aikido` ignores](https://help.aikido.dev/code-scanning/scanning-practices/ignore-via-code-with-aikido-files) still apply on their own.

## Related

{% content-ref url="../pentests/configure-a-pentest/leveraging-code-and-documentation.md" %}
[leveraging-code-and-documentation.md](https://help.aikido.dev/pentests/configure-a-pentest/leveraging-code-and-documentation)
{% endcontent-ref %}

{% content-ref url="../code-quality/add-extra-code-context.md" %}
[add-extra-code-context.md](https://help.aikido.dev/code-quality/add-extra-code-context)
{% endcontent-ref %}

{% content-ref url="../deep-review/how-deep-review-works.md" %}
[how-deep-review-works.md](https://help.aikido.dev/deep-review/how-deep-review-works)
{% endcontent-ref %}

{% content-ref url="../ai-code-analysis/ai-code-analysis.md" %}
[ai-code-analysis.md](https://help.aikido.dev/ai-code-analysis/ai-code-analysis)
{% endcontent-ref %}

{% content-ref url="../ai-code-analysis/ai-code-audit-overview/create-an-ai-code-audit.md" %}
[create-an-ai-code-audit.md](https://help.aikido.dev/ai-code-analysis/ai-code-audit-overview/create-an-ai-code-audit)
{% endcontent-ref %}

{% content-ref url="../ai-code-analysis/what-is-dspm/data-exposure-audit.md" %}
[data-exposure-audit.md](https://help.aikido.dev/ai-code-analysis/what-is-dspm/data-exposure-audit)
{% endcontent-ref %}

{% content-ref url="../ai-code-analysis/cloud-security-audit/run-a-cloud-security-audit.md" %}
[run-a-cloud-security-audit.md](https://help.aikido.dev/ai-code-analysis/cloud-security-audit/run-a-cloud-security-audit)
{% endcontent-ref %}

{% content-ref url="../ai-code-analysis/smart-contract-audit/run-a-smart-contract-audit.md" %}
[run-a-smart-contract-audit.md](https://help.aikido.dev/ai-code-analysis/smart-contract-audit/run-a-smart-contract-audit)
{% endcontent-ref %}
