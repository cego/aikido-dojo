# Code Quality Overview

<table data-view="cards"><thead><tr><th>Section</th><th>What you will find</th><th data-hidden data-card-target data-type="content-ref">Target</th></tr></thead><tbody><tr><td><strong>Set Up PR Checks</strong></td><td>Quality checks on <strong>new code</strong>: scan changed lines on every pull request.</td><td><a href="https://help.aikido.dev/code-quality/code-quality-setup">code-quality-setup.md</a></td></tr><tr><td><strong>Repository Scans</strong></td><td>Quality checks on <strong>existing code</strong>: full-repo scores, findings, and hotspots.</td><td><a href="https://help.aikido.dev/code-quality/repository-scans">repository-scans</a></td></tr><tr><td><strong>Code Coverage</strong></td><td>Upload coverage reports from CI, track coverage per repository and file, and fail pull requests below a threshold.</td><td><a href="https://help.aikido.dev/code-quality/code-coverage">code-coverage</a></td></tr><tr><td><strong>Add Custom Code Rules</strong></td><td>Define team-specific checks beyond the defaults.</td><td><a href="https://help.aikido.dev/code-quality/add-custom-code-rules">add-custom-code-rules.md</a></td></tr><tr><td><strong>Add Extra Code Context</strong></td><td>Give Aikido background on your codebase and conventions.</td><td><a href="https://help.aikido.dev/code-quality/add-extra-code-context">add-extra-code-context.md</a></td></tr><tr><td><strong>Ignore Specific Lines</strong></td><td>Suppress findings on lines you've reviewed and accepted.</td><td><a href="https://help.aikido.dev/code-quality/ignoring-specific-lines">ignoring-specific-lines.md</a></td></tr></tbody></table>

## What is Code Quality?

Code Quality helps you maintain high standards across the **new code** you ship and the **existing code** already in your repos.

It is the quality counterpart to [SAST](https://help.aikido.dev/code-scanning/scanning-practices/sast-by-aikido-supported-languages-and-security-focus): checks on new code as you add it, once you enable them. It does not replace [Deep Review](https://help.aikido.dev/deep-review/how-deep-review-works). Deep Review covers **both security and quality**, with context beyond the changed lines.

Aikido applies AI-powered checks against your team's enabled rules. Findings include plain-language explanations and fix guidance, so developers can act on them without being security experts.

## How it works

Code Quality splits across these modes:

<table><thead><tr><th width="188.30859375">Mode</th><th>What it covers</th></tr></thead><tbody><tr><td><a href="https://help.aikido.dev/code-quality/code-quality-setup"><strong>PR checks</strong></a></td><td><strong>New code:</strong> lines added or changed in a pull request</td></tr><tr><td><a href="https://help.aikido.dev/code-quality/repository-scans"><strong>Repository scans</strong></a></td><td><strong>Existing code:</strong> the full codebase as it stands today</td></tr><tr><td><a href="https://help.aikido.dev/code-quality/code-coverage"><strong>Code Coverage</strong></a></td><td><strong>Test coverage:</strong> reports uploaded from CI, with optional <a href="https://help.aikido.dev/code-quality/code-coverage/pr-gating-on-code-coverage">PR gating</a> on a minimum percentage</td></tr></tbody></table>

Use either quality mode on its own, or both together. PR checks stop quality from getting worse; repository scans show you what's already there and where to focus cleanup. Add [Code Coverage](https://help.aikido.dev/code-quality/code-coverage) to see how much of that code your tests actually exercise, and [fail pull requests](https://help.aikido.dev/code-quality/code-coverage/pr-gating-on-code-coverage) that fall short of a coverage threshold.

### PR checks for new code

On every pull request, Aikido scans only the lines you changed and leaves inline comments (or fails the gate, if you configure it that way). This catches regressions before they merge and keeps new contributions aligned with your standards.

See [Set Up PR Checks](https://help.aikido.dev/code-quality/code-quality-setup) for the full setup flow.

### Repository scans for existing code

A full-codebase scan scores each repository across four quality categories, lists every open finding in code that's already merged, and highlights hotspots where high churn meets low quality.

See [Repository Scans](https://help.aikido.dev/code-quality/repository-scans) for setup, pricing, and how to read results.

## Use cases

* **Guard new code at review time:** PR checks on changed lines, with optional gate failures
* **Understand existing code health:** Category scores across the full codebase per repository
* **Prioritize legacy cleanup:** Scatter plot hotspots where frequently changed files carry the most debt
* **Enforce team standards:** Custom rules and code context tailored to your conventions
* **Cover both layers:** PR checks on every merge; repository scans with incremental rescans during your yearly subscription
* **Track test coverage:** Upload coverage from CI, see it per repository, file, and line, and fail pull requests below a threshold

## Supported languages

Code Quality supports a wide and growing set of languages and frameworks. Click below for the full list.

{% hint style="info" %}
Don't see your language listed? Reach out in intercom. We can look at adding it.
{% endhint %}

<details>

<summary>Supported languages</summary>

* ABAP (SAP)
* angular
* Ansible
* Apex (Salesforce)
* C/C++
* C#
* Clojure
* COBOL
* ColdFusion
* Dart
* Delphi
* Docker
* Elixir
* Erlang
* go
* Groovy
* Haskell
* Java
* JS
* JSON
* Julia
* Kotlin
* Kubernetes / Helm charts
* Lua
* MySQL
* OCaml
* OpenEdge
* Pascal
* perl
* PHP
* PL/SQL
* PowerShell
* PY
* R
* react
* Ruby
* Rust
* Scala
* shell
* Solidity
* SQL
* svelte
* Swift/Objective-C
* Terraform
* TS
* Visual Basic
* vue
* YAML
* Zig

</details>

For ignoring findings via `.aikido` configuration files in your repository, see [Ignore via Code with Aikido Files](https://help.aikido.dev/code-scanning/scanning-practices/ignore-via-code-with-aikido-files).
