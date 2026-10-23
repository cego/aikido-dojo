# Using Aikido MCP

The tools the Aikido MCP server exposes, and example prompts for asking your agent to use them.

Once the MCP server is installed, you don't call these tools yourself. You ask in plain language and your agent picks the right one.

## Available tools

### aikido\_scan\_paths

Scans files that are saved on disk for vulnerabilities (SAST), hardcoded secrets and IaC misconfigurations. This is the preferred scan tool.

* Provide the files to scan in `paths`, absolute or relative to `root`. Directories are rejected, so list each file
* Set `root` to the workspace or repository root so findings are reported relative to it. Paths that are missing or outside `root` are skipped and listed in `skipped_paths`
* Optionally pass `repository_name` when the files come from a repository, so they are linked to that repository in Aikido

### aikido\_full\_scan

Scans code that is passed in as file content, for the cases where the code is not on disk yet: generated but unsaved code, a pasted snippet, or a proposed change. Use `aikido_scan_paths` whenever the files do exist on disk.

* Provide each file in `files_to_scan` as a `relativeFilePath` and its `content`

### aikido\_oss\_dependency\_scan

Scans dependency manifests and lockfiles that are saved on disk for vulnerable, outdated or malicious packages.

* Provide the manifests and lockfiles to scan in `paths`, for example `package.json` and `package-lock.json`, `requirements.txt`, `go.mod` and `go.sum`, or `pom.xml`. Directories are rejected
* Set `root` to the workspace or repository root so findings are reported relative to it
* Optionally pass `repository_name` when the files come from a repository

### aikido\_issues\_list

Fetches security issues from your Aikido feed.

* Filter by scope: `repo_name`, `cloud_name`, `vm_name`, `domain_name`, `container_name`, `workspace_name` or `team_name`
* Optionally narrow a repo scope to a single branch with `repo_branch_name`
* Pick one or more issue types: `sast`, `leaked_secret`, `iac`, `open_source`, `cloud`, `cloud_instance`, `docker_container`, `malware`, `eol`, `mobile`, `surface_monitoring`, `scm_security`, `license`, `ai_pentest`, `ai_code_analysis`
* Filter by `labels` to only return issues on repositories or container images carrying any of the given labels
* Filter by `severity` to only return issues of a certain severity or above a given threshold
* Page through results with `page` (zero-based)
* Filter on SLA by choosing the following filters: `out_of_sla` or `sla_due_soon`
* Returns each issue with title, type, severity, and remediation steps

#### Looking up a single issue

Pass `issue_id` to fetch one issue in detail. It comes back with more than the list view shows:

* `score_reasons`: every reason the issue got its severity, each with a readable explanation and how much it moved the score, for example a known exploit in the wild
* `pentest`: for AI pentest findings, the attack analysis with its summary, risk, reproduction steps, root cause and remediation
* `licenses`: for license issues, the detected license with its legal risk and why

### aikido\_workspaces\_list

Lists the Aikido workspaces you have access to. Your agent uses it to find the exact workspace name before filtering issues by `workspace_name`.

### aikido\_ignore\_issue

Ignores a security issue in the feed. Requires `issue_id` and a `reason`.

### aikido\_unignore\_issue

Unignores a previously ignored issue. Requires `issue_id` and a `reason`.

### aikido\_snooze\_issue

Temporarily hides a security issue in the feed. Requires `issue_id`, a `reason` and a `duration`: `1_day`, `1_week`, `1_month`, `3_months`, `6_months` or `1_year`.

### aikido\_unsnooze\_issue

Unsnoozes a previously snoozed issue. Requires `issue_id`.

### aikido\_adjust\_issue\_severity

Raises or lowers the severity of a security issue. Requires `issue_id`, the new `adjusted_severity` (`critical`, `high`, `medium` or `low`) and a `reason`.

### aikido\_login

Starts the Aikido sign-in flow; returns sign-in URL or confirms you're already signed in.

* Sign-in redirects you to your account's region automatically. Only set `region` (`eu`, `us`, `me`, `au`) if you land in the wrong one
* Use `force_reauth` to switch accounts or sign in again

{% hint style="info" %}
Not all MCP tools are enabled by default. Admins can enable them on the [permissions page](https://app.aikido.dev/settings/integrations/ide/mcp/permissions) for everyone
{% endhint %}

## Example prompts

### Scanning code

* "Use Aikido to scan this file for security issues"
* "Run an Aikido scan on my staged changes to check for secrets before I commit"
* "Scan the files I just edited with Aikido and link them to the `payments-api` repo"

### Reviewing issues by repo

* "Show me all critical Aikido issues in `payments-api`"
* "List any leaked secrets in `frontend-web` from Aikido"
* "What open source vulnerabilities does Aikido see in `api-gateway`?"
* "Show SAST and IaC issues in `infra-core` from Aikido"

### Reviewing issues by cloud, VM, or container

* "List all Aikido cloud issues in `prod-aws`"
* "Show malware findings on `web-server-01` from Aikido"
* "What end-of-life software is running in the `nginx-proxy` container per Aikido?"
* "Show me surface monitoring issues for `example.com` in Aikido"

### Digging into one issue

* "Give me the full details of Aikido issue `12345`"
* "Why is Aikido issue `12345` rated critical? Is there a known exploit?"
* "Walk me through the reproduction steps for the Aikido pentest finding `12345`"

### Managing issues

* "Snooze Aikido issue `12345` for a month, we're fixing it in the next release"
* "Lower the severity of Aikido issue `12345` to low, the endpoint is internal only"
* "Unignore Aikido issue `12345`, the mitigation was removed"
* "Which Aikido workspaces do I have access to?"

### Combined workflows

* "Use Aikido to scan my current changes, then show existing critical issues in the same repo"
* "Check this PR with Aikido and compare against open SAST issues in the repo"

## Controlling when scans run

Aikido does not decide when to scan. Your agent does. A scan happens when the agent calls `aikido_scan_paths`.

So anything you tell the agent about scanning takes effect straight away, for the rest of the session.

### Pause scanning

* "Don't run any Aikido scans until I ask for one."
* "Pause Aikido scanning while I iterate on this."
* "Skip the Aikido scan for now, I'll ask for one before we commit."

### Scan once at the end instead of after every edit

* "Don't scan after each edit. Run one Aikido scan when the whole task is done."
* "Run an Aikido scan on my staged changes right before I commit, not before that."

### Keep scans small

Scan time scales with how much code goes in, so narrowing the scope is usually faster than switching scanning off.

* "Only scan the files you changed, not the whole directory."
* "Scan `src/payments/` only."
* Prefer `aikido_scan_paths` over `aikido_full_scan`. It reads the files from disk instead of passing their content through the conversation.

### Make it your default

A prompt instruction lasts for one session. To change the default for every session, put the same instruction in your agent's rules file (`CLAUDE.md`, `.cursor/rules`, `AGENTS.md`, and so on):

```
Run an Aikido scan only when I explicitly ask for one, or before a commit.
```
