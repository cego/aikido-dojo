# Add Extra Code Context

Code Context helps Aikido understand your team’s unique coding practices and exceptions. By providing context about your codebase, you can fine-tune how Code Quality analyzes your pull requests, reducing noise and making reviews more relevant to your team’s standards.

The same guidance can live in the repository. Commit `AGENTS.md` at the root when it should apply to every review. See [Guide Agents with Repository Files](/ai-and-dev-tools/guide-an-audit-with-repository-files.md).

### What is Code Context?

Code Context allows you to teach Aikido about your team’s specific coding conventions, architectural decisions, and acceptable exceptions to general best practices. This ensures Code Quality comments are relevant and actionable rather than generic.

Think of Code Context as a way to document your team’s “tribal knowledge” - those unwritten rules and exceptions that every team member knows but aren’t captured in traditional linting rules.

### Common Code Context examples

#### Performance and architecture decisions

{% code overflow="wrap" %}

```
Our data-export-service repository is our nightly batch processor for archived data. We intentionally avoid database indexes here because the tables are append-only and indexes would slow down our bulk insert operations.
```

{% endcode %}

#### Tool-specific exceptions

{% code overflow="wrap" %}

```
Debug logging and console.log statements are acceptable in our troubleshooting-utils repository as these are CLI tools designed for debugging production issues.
```

{% endcode %}

#### Relaxed standards for internal tools

{% code overflow="wrap" %}

```
For data migration scripts in the migrations/ folder, we exceed line count limits multiple times. This is acceptable as these are one-time scripts. We also don't enforce strict error handling in internal-facing scripts because our team knows how to use them.
```

{% endcode %}

#### Import conventions

{% code overflow="wrap" %}

```
Wildcard imports (import * as) are acceptable when importing from our @company/ui-components library as we intentionally export all components as a single namespace.
```

{% endcode %}

#### Documentation requirements

{% code overflow="wrap" %}

```
Relaxed JSDoc requirements apply to internal helper functions in the utils/ folder. Full documentation is only required for exported public APIs.
```

{% endcode %}

#### Performance optimizations

{% code overflow="wrap" %}

```
Performance-critical functions in our real-time-processing module can skip certain TypeScript strict checks when properly commented with @performance-critical.
```

{% endcode %}

### How to Add Extra Code Context

**Step 1.** Go to [Code Quality Checks](https://app.aikido.dev/code-quality/checks) page and click Code Context

**Step 2.** Click **Add Code Context** to add new code context.

**Step 3.** Define the Code Context

* Write a clear description explaining your team’s practice or exception
* Choose the scope:
  * **All repositories** - Applies to your entire organization
  * **Selected repositories** - Only applies to specific repos

{% hint style="info" %}
**Tips on Writing Effective Context**

**Do:**

* Be specific about when and why exceptions apply
* Reference specific folders, file patterns, or repositories
* Explain the reasoning behind architectural decisions
* Use clear, concise language

**Don’t:**

* Create overly broad exceptions that defeat the purpose of quality checks
* Duplicate what’s already handled by toggling checks on/off
* Write vague statements without clear scope
  {% endhint %}

### Code Context vs Custom Code Rules

Here some small guidelines to understand when to use which functionality

#### Use Code Context when:

* You have **acceptable exceptions** to best practices in certain areas
* You need to **adjust existing rules** for specific situations
* You want to **provide explanations** for architectural decisions

#### Use Custom Code Rules when:

* You need to **enforce new standards** not covered by default checks
* You want to **detect specific patterns** unique to your codebase
* You need **hard enforcement** with consistent detection
* You want to **block PRs** for specific violations

## Use feedback to generate Code Context

You can add extra code context straight from your pull request. Simply post a comment to @AikidoSec with your feedback or explanation, and Aikido will automatically process it to improve future review insights.
