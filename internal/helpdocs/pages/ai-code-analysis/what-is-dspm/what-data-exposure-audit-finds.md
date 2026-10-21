# What Data Exposure Audit Finds

What Data Exposure Audit checks for in your source code.

A **Data Exposure Audit** is Aikido's code-based DSPM. It reasons about how sensitive data enters, moves, is stored, and leaves your application: from APIs and ORM models to logging, queues, third-party integrations, and AI workflows. You don't need to connect databases or cloud data stores. Each finding ships with code-based **Evidence**. See [Track Progress and Findings](/ai-code-analysis/ai-code-audit-overview/track-progress-and-findings.md) for how to read it.

For how this differs from traditional DSPM, see [Data Exposure Audit (DSPM)](/ai-code-analysis/what-is-dspm.md). To run an audit, see [Run a Data Exposure Audit](/ai-code-analysis/what-is-dspm/data-exposure-audit.md).

## Coverage at a glance

Data Exposure Audit focuses on sensitive data handling in code, including:

* **Plaintext and weak storage**: secrets, tokens, or regulated personal data (PII, PHI, payment card data) stored without encryption, or with weak substitutes like MD5, SHA-1, or base64 treated as encryption.
* **Credentials in URLs**: tokens, passwords, or API keys embedded in clone URLs, query strings, redirects, or other request paths where they can leak through logs, proxies, or browser history.
* **API over-exposure**: sensitive fields returned more broadly than needed, and missing field-level or tenant ownership checks on data access.
* **Unsafe logging, caching, and exports**: unredacted PII or secrets in application logs, traces, caches, dumps, or data export paths.
* **Over-broad third-party and AI sharing**: customer data or secrets sent to vendors, analytics, support tools, LLM prompts, or vector stores beyond what's required.
* **Incomplete data deletion**: delete flows that leave downstream records, caches, or related stores intact when a user or retention policy expects removal.
* **Shadow data and weak backups**: unintended copies of sensitive data (staging restores, snapshots, side stores) and backup or snapshot paths without adequate protection or masking.
* **Over-permissioned access**: cloud IAM roles or third-party OAuth scopes that grant more data access than the code needs.
* **Unauthenticated or unaudited data stores**: data store access without authentication, or audit logging disabled on stores that hold sensitive data.
* **Identity and source exposure**: account or identity enumeration risks, and unnecessary exposure of source IPs or similar identifiers tied to sensitive flows.

## Data types in scope

Findings can involve any sensitive data the code handles, including:

* Personally identifiable information (PII)
* Protected health information (PHI)
* Payment card and financial data (PCI)
* Secrets, tokens, and credentials
* Customer content and business data
* Data passed into AI prompts, embeddings, or related workflows

## Beyond what traditional DSPM sees

Traditional DSPM often tells you sensitive data exists in a store. Data Exposure Audit connects the exposure back to the code path that created it, so you get the destination, the risk, and a remediation you can ship:

* **Logging that leaks PII**: for example, a checkout path that logs a full billing object (email, phone, address, card last four) into your observability stack. The fix is to redact those fields and keep only IDs and error codes.
* **Token mishandling end to end**: where a token enters (OAuth, GitHub App, PAT), what scopes it has, where it's stored or cached, whether it's embedded in URLs, and whether it can leak through logs, traces, errors, or process arguments.
* **API responses that over-share**: serializers or handlers that return more sensitive fields than the caller should see, without a matching authorization check.
* **AI and vendor pipelines**: prompts or integrations that pull customer data into third-party or LLM workflows without minimization.

## What it doesn't replace

Data Exposure Audit reasons from **source**. It doesn't connect to your live datastores, so it won't inventory what's already sitting in buckets, warehouses, or SaaS apps the way traditional DSPM does. Use it to find and fix how your application exposes data; pair it with cloud and storage reviews when you need estate-wide inventory of data at rest.

It also isn't a substitute for a [Code Security Audit](/ai-code-analysis/ai-code-audit-overview/what-ai-code-audit-finds.md) or a live [Aikido Pentest](/pentests/aikido-pentest.md). Those cover vulnerability classes and runtime exploitability; Data Exposure Audit focuses on sensitive data handling and exposure paths in code.

## FAQ

### Can you confirm coverage before I start?

Yes. [Contact support](/miscellaneous-info/contact-support.md) to check a specific data type, framework, or finding class before you spend credits.

## Related

{% content-ref url="/pages/Io2yJZC8MsZHFIvz5M38" %}
[Data Exposure Audit (DSPM)](/ai-code-analysis/what-is-dspm.md)
{% endcontent-ref %}
