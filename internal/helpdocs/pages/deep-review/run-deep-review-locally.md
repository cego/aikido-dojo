# Run Deep Review Locally

You can run Deep Review on the changes in your local git repository, before you open a pull request (PR). Your AI agent (in your IDE or a CLI agent such as Claude Code, Codex CLI or Gemini CLI) starts the review through the [Aikido MCP server](https://help.aikido.dev/ai-and-dev-tools/aikido-mcp#installation) and shows you the findings in the conversation.

Use it when you want feedback on a branch while you're still working on it. Once you open a PR, Deep Review on the PR works as described in [How Deep Review Works](https://help.aikido.dev/deep-review/how-deep-review-works).

## How to Install / Prerequisites

* The Aikido MCP server is installed in your agent. See [AI Coding Assistants (MCP)](https://help.aikido.dev/ai-and-dev-tools/aikido-mcp#installation)
* Credits in your Aikido Wallet. See [Credits](#credits)

## Run a review

Ask your agent in plain language:

* "Run an Aikido Deep Review on the changes in my current branch"
* "Deep Review only my uncommitted changes with Aikido"
* "Run an Aikido Deep Review on my branch and include the new files I haven't added to git yet"

A review usually takes 2 to 4 minutes. The agent starts it with `aikido_deep_review`, then keeps checking with `aikido_deep_review_status` until the result is ready. You don't need to do anything while it waits.

## What gets reviewed

By default, Deep Review looks at all changes on your branch: commits that aren't on the default branch yet, plus your uncommitted edits. Ask your agent to narrow it down:

| Ask for                    | What gets reviewed                                                                     |
| -------------------------- | -------------------------------------------------------------------------------------- |
| All changes (default)      | Commits that aren't on the default branch yet, plus uncommitted edits                  |
| Only committed changes     | Commits that aren't on the default branch yet                                          |
| Only uncommitted changes   | Uncommitted edits on top of your latest commit. That commit must be pushed             |

To compare against a branch other than the default one, name it in your prompt, for example "Deep Review my changes against `origin/release`".

New files that you haven't added to git are skipped, and the agent lists which ones. Ask for them to be included when they're part of your change. Check first that they don't contain scratch files or secrets.

## Credits

Each local review costs **1 credit**. Unlike Deep Review on a PR, every run is a new review, even for a branch you reviewed before. Ask for a review when your change is ready, not after every edit.

You're not charged when:

* The review fails
* The changes only touch files Deep Review ignores, such as lock files, generated code, or vendored code

{% hint style="info" %}
Tell your agent when a change is worth a Deep Review, for example when it touches authentication, access control, or core business logic. Put that in your agent's rules file (`AGENTS.md`, `CLAUDE.md`, `.cursor/rules`, and so on), for example: "Run an Aikido Deep Review before I commit changes to auth or payment code."
{% endhint %}

## Troubleshooting

### The review was skipped

Your changes only touch files Deep Review ignores, so there was nothing to review. No credit was charged.

### The review failed

Deep Review could not run. No credit was charged. Ask your agent to start a new review.

### A review of uncommitted changes doesn't start

A review of only uncommitted changes needs your latest commit on the remote. Push your branch, then ask again.
