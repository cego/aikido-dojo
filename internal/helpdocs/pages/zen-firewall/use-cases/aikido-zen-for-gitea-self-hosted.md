# Aikido Zen for Gitea Self-Hosted

[Gitea](https://about.gitea.com/) is a self-hosted Git service for hosting source code. It typically includes pull requests, issues, releases, and webhooks.

That often means it holds very sensitive data: private repositories, credentials committed by mistake, deploy keys, access tokens, and internal engineering discussions.

If attackers abuse a vulnerability in your Gitea stack, they can exfiltrate source code, steal secrets, or gain a foothold inside your network. Aikido Zen helps reduce this risk by blocking common web attacks at runtime inside the application code.

## What is Aikido Zen and why it matters

Aikido Zen is an in app firewall. It embeds into your runtime and instruments key functions like file access, outbound HTTP calls and database queries. When user input reaches these sensitive sinks in a dangerous way, Zen throws an exception and stops the operation before it executes.

For a Gitea admin this matters because:

* Gitea relies on a database (SQLite/MySQL/MariaDB/PostgreSQL) plus a filesystem repo store for Git data and attachments.
* A single injection or path traversal bug can expose private repos or leak secrets at scale.
* Gitea is often internet facing, which makes it a frequent target for bots and known threat actors.

Zen adds an extra safety net on the server where Gitea runs, without changing how users access the platform.

## How Zen protects Gitea

**Mitigating abusive traffic**

Zen includes traffic control features that help slow down or block brute force style attempts and cut noise from scanning bots before they reach application logic.

**Blocking unsafe database operations**

Zen inspects database queries at runtime and can stop patterns that look like SQL injection or harmful query manipulation. This helps contain issues coming from custom code or third-party modules that interact with your data models.

**Preventing unauthorized file system access**

Many applications store configuration, attachments, and temporary files on disk. Zen monitors sensitive file operations and blocks path traversal or unexpected file reads and writes outside allowed directories.

**Protecting against unsafe outbound requests and SSRF**

Many applications integrate with external services via webhooks and APIs. Zen checks outbound HTTP calls and can block unsafe requests created from untrusted input. This prevents SSRF attempts and limits access to internal network targets.

**Monitoring and controlling outbound domains**

Zen shows every domain your server connects to. You can block or allow domains to control where data is sent. This is useful for detecting unexpected calls from custom modules, preventing data exfiltration and keeping a tight boundary around integrations your team approves.

**Local only inspection**

Zen makes security decisions locally on the server without sending request data to Aikido. This fits well with deployments where critical business data and privacy is a priority.

## Next step

Beta available soon, contact Aikido support for more information
