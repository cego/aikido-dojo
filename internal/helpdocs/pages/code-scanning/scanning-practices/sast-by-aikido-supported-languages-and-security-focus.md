# SAST/IaC: Supported Languages and Security Focus

SAST (static analysis) checks your source code for known dangerous patterns, such as untrusted input reaching a database query. It runs continuously on connected repositories and on pull requests through [PR gating](https://help.aikido.dev/pr-and-release-gating/aikido-ci-gating-functionality). It does not need a running app.

SAST will not catch a logic or authorization bug that does not look like a known pattern. For the impact of a pull request, use [Deep Review](https://help.aikido.dev/deep-review/how-deep-review-works). For existing source, use [AI Code Analysis](https://help.aikido.dev/ai-code-analysis/ai-code-analysis). The quality counterpart to SAST is [Code Quality](https://help.aikido.dev/code-quality/code-quality-overview).

See [Which Aikido product should I use?](https://help.aikido.dev/getting-started/which-aikido-product-should-i-use).

### How Aikido SAST currently works <a href="#how-aikido-sast-currently-works" id="how-aikido-sast-currently-works"></a>

Aikido’s SAST engine is built to find and prioritize security issues in your code. No noise, just the vulnerabilities you need to fix.

Aikido SAST engine is based on our **custom risk categorisation model**. Some of these categorisation: -

* Aikido removes findings that are not related to security (eg opinionated code styling rules).
* Findings that reside in repositories that a user categorized as sensitive will get upgraded.
* Findings inside of files that are not intended for production (eg unit tests or functions that aren't used in production) might get downgraded and so on.

Our SAST engine also leverage some of the best open-source engines out there, which we have significantly customized and fine-tuned to provide you sharper, relevant results over the years.

To view all individual rules that are active per language, check out our [SAST Checks](https://app.aikido.dev/repositories/sast?regionRedirect=1) or [Infrastructure as Code](https://app.aikido.dev/repositories/iac?regionRedirect=1) checks to view the rules per language.

### Language support <a href="#language-support" id="language-support"></a>

Aikido is not sensitive to the versions of languages. By default, we support all versions.\
Aikido supports tracking tainted user input from top-level controllers to other files where dangerous functions are used for a growing set of languages.

| **Language**                                                                 | **Base engine**                                    | **Taint analysis**    |
| ---------------------------------------------------------------------------- | -------------------------------------------------- | --------------------- |
| JavaScript                                                                   | Aikido Engine + Opengrep                           | Across multiple files |
| Typescript                                                                   | Aikido Engine + Opengrep                           | Across multiple files |
| PHP                                                                          | Aikido Engine + Opengrep                           | Across multiple files |
| .NET/C#                                                                      | Aikido Engine + Opengrep                           | Across multiple files |
| Java                                                                         | Aikido Engine + Opengrep                           | Across multiple files |
| Rust                                                                         | Aikido Engine + Opengrep                           | Across multiple files |
| Go                                                                           | Aikido Engine + Opengrep                           | Across multiple files |
| Ruby                                                                         | Aikido Engine + Opengrep                           | Across multiple files |
| Python                                                                       | Aikido Engine + Opengrep                           | Across multiple files |
| Scala                                                                        | Aikido Engine + Opengrep                           | Across multiple files |
| C/C++                                                                        | Aikido Engine + Opengrep                           | Across multiple files |
| Visual Basic                                                                 | Aikido Engine + Opengrep                           | Across multiple files |
| Swift                                                                        | Aikido Engine + Opengrep                           | Within files          |
| Android                                                                      | Aikido Engine + Opengrep                           | Within files          |
| Kotlin                                                                       | Aikido Engine + Opengrep                           | Within files          |
| Dart                                                                         | Aikido Engine + Opengrep                           | Within files          |
| Elixir                                                                       | Aikido Engine + Opengrep                           | Within files          |
| Apex                                                                         | Aikido Engine + Opengrep                           | Within files          |
| Clojure                                                                      | Aikido Engine + Opengrep                           | Within files          |
| Infrastructure-as-code files (Terraform, Cloudformation, Docker, Pulumi,...) | Aikido Engine + Checkov                            | Not applicable        |
| Exposed secret discovery in all files inside of Git history                  | Aikido Base Engine with Liveness Checks + Gitleaks | Not applicable        |
