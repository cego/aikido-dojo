# How Deep Review Works

Deep Review is Aikido's agentic pull request reviewer. Instead of treating the diff as the full story, Aikido explores relevant code, reasons about the wider impact of a change, and posts feedback directly on your pull request.

It covers **both security and quality**. It can flag a missing access check, and a refactor that breaks pagination on large result sets, as long as the issue shows up in the impact of this pull request.

## Deep Review vs other PR checks

Pull requests can already receive comments from more than one Aikido check. They answer different questions:

| Check                                                                                                                                                                       | What it looks at                                                        | Cadence                   |
| --------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------- | ------------------------- |
| [SAST](/code-scanning/scanning-practices/sast-by-aikido-supported-languages-and-security-focus.md) in [PR gating](/pr-and-release-gating/aikido-ci-gating-functionality.md) | Known security patterns on the diff                                     | On new code, once enabled |
| [Code Quality](/code-quality/code-quality-overview.md)                                                                                                                      | Quality rules on the changed lines                                      | On new code, once enabled |
| Deep Review                                                                                                                                                                 | Security and quality impact of the change, with context beyond the diff | Credits per PR            |

Deep Review does not replace SAST or Code Quality. Enable those in PR Checks and keep them on for every commit. Use Deep Review when you want to find complex issues that need more context than the diff.

For the full map, see [Which Aikido product should I use?](/getting-started/which-aikido-product-should-i-use.md).

## Context beyond the diff

A diff shows what changed, but not always how that change interacts with the rest of your application. Deep Review pulls in relevant context from outside the changed lines before deciding whether to comment, so it can review changes that depend on existing behavior elsewhere in the codebase.

## Repository instructions

Deep Review checks `AGENTS.md`, `CLAUDE.md`, `.claude/`, and `.cursor/` for accepted risk. A risk those files already accept stays unreported, unless the pull request widens that risk.

Put standing guidance in `AGENTS.md` at the repository root: trust boundaries, intentional behavior, and accepted risks. See [Guide Agents with Repository Files](/ai-and-dev-tools/guide-an-audit-with-repository-files.md).

Aikido also builds [memory](/deep-review/memory.md) from replies on pull requests. Use `AGENTS.md` when you want to write the guidance yourself.

## Cross-repository context

If your application spans multiple repositories, you can [link related repositories](/deep-review/link-repositories.md) to the repository under review. Deep Review can then use those repositories as additional context when reasoning about shared libraries, APIs, and services.

The number of repositories you can link depends on your Aikido plan.

## Feedback in your pull request

Aikido posts comments directly on your pull request when Deep Review has feedback to share. You don't need to leave your existing review workflow to see or discuss the results.

Deep Review also posts its own status check (**Aikido Security: Deep Review**). Whether that check fails follows the same [gating settings](/deep-review/configure-deep-review.md#pr-check-gating) as the security review: the failure threshold, **Always Pass**, and **Always make the PR check green**.

When you fix an issue that a comment identifies, Aikido automatically resolves the comment. [Learn how Aikido manages comments after an issue is fixed](/pr-and-release-gating/aikido-ci-gating-functionality/autofix-suggestions-and-inline-commenting-for-pr-checks.md#what-happens-to-comments-after-you-fix-an-issue).

### Reply in natural language

Reply to any Deep Review comment in plain language. You don't need to tag Aikido or use a special command format, and you can stay in the pull request without switching to the Aikido UI.

| What you say                   | Example                                                  | Outcome                                                                                                                  |
| ------------------------------ | -------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------ |
| **Ignore a finding**           | "This is a false positive, please ignore."               | Aikido ignores the finding.                                                                                              |
| **Share context or push back** | "This is intentional because we validate on the server." | Aikido stores this in [repository memory](/deep-review/memory.md) for later reviews.                                     |
| **Ask for a fix**              | "Can you fix this?"                                      | Aikido opens a pull request with the fix when [AutoFix](/autofix-and-remediation/overview-aikido-autofix.md) is enabled. |
| **Ask for clarification**      | "Can you explain why this is a problem?"                 | Aikido replies to your question.                                                                                         |

## Exclude paths and files

Aikido already ignores a large set of paths out of the box. Deep Review skips common noise such as dependency folders, build output, generated files, lock files, binaries, and media, so you usually do not need extra configuration.

You can exclude more paths with a [`.aikido` file](/code-scanning/scanning-practices/ignore-via-code-with-aikido-files.md). Deep Review respects those excluded paths as well.

If a pull request only changes excluded files, Aikido skips Deep Review and you are not charged.

## Credits

Deep Review costs **1 credit per pull request**. Additional reviews on the same PR (for example after new commits) don't cost extra. When you're out of credits, Deep Reviews are skipped.

You can also set a [monthly spending limit](/deep-review/configure-deep-review.md#set-a-monthly-spending-limit) for Deep Review from [Wallet settings](https://app.aikido.dev/settings/wallet/settings). When the limit is reached, Deep Reviews that would charge a credit are skipped until the next calendar month.

{% hint style="info" %}
Curious what Deep Review would cost for your team? Our [estimator](https://app.aikido.dev/repositories/deep-review) gives you an idea based on your activity of the last 30 days.
{% endhint %}

Learn how to [configure Deep Review](/deep-review/configure-deep-review.md).
