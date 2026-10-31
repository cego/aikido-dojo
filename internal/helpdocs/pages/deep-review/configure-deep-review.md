# Configure Deep Review

Enable Deep Review for one or more repositories from PR Checks. After enabling it, use each repository's settings to choose the strictness level and [link related repositories](/deep-review/link-repositories.md). The same card shows [Deep Review memory](/deep-review/memory.md).

Deep Review uses the same gating settings as the security review. It does not have a separate failure threshold.

{% hint style="info" %}
Deep Review costs **1 credit per pull request**. Additional reviews on the same PR don't cost extra.
{% endhint %}

## Enable Deep Review

{% stepper %}
{% step %}

### Open PR Checks

Go to **Repositories > Pull/Merge Requests > Manage PR/MR Checks**. You can also open [Manage PR/MR Checks](https://app.aikido.dev/repositories/prs) directly.
{% endstep %}

{% step %}

### Select repositories

Select one or more repositories. To enable them in bulk, select every repository you want to configure, then click **"Setup PR Scans"** in the bulk actions bar.
{% endstep %}

{% step %}

### Turn on Deep Review

Open the **Deep Review** section and turn on **Perform Deep Review of PRs**.
{% endstep %}

{% step %}

### Save your changes

Click **"Save"**. Aikido will review future pull requests for the selected repositories.
{% endstep %}
{% endstepper %}

<div data-with-frame="true">[image: Repository configuration with Deep Review enabled and the Perform Deep Review of PRs toggle turned on.]</div>

## PR check gating

Deep Review follows the gating settings of the security review for that repository. Change those settings in the **Vulnerability Scanning** section of the same PR Checks modal.

The **Scan failure for** threshold also applies to the **Aikido Security: Deep Review** check:

* Findings at or above the threshold fail the Deep Review check.
* Findings below the threshold are still reported and commented on, but they do not fail the check.
* **Always Pass** keeps the Deep Review check green. Aikido still reviews the PR and posts comments.

If **Always make the PR check green** is enabled for the workspace, the Deep Review check also stays green.

See [PR Gating Overview](/pr-and-release-gating/aikido-ci-gating-functionality.md) for how to set the failure threshold.

<div data-with-frame="true"></div>

## Open the Deep Review settings for a repository

Strictness, linked repositories, and memory are shown for each repository:

1. Open [**Repositories**](https://app.aikido.dev/repositories) and select the repository you want to configure.
2. Open its **Repository Settings**.
3. Find the **Deep PR Review** card.
4. Click **Edit** next to **Strictness** or **Linked Repositories**. To see what Aikido has learned from replies, click **View Memory**. See [Deep Review Memory](/deep-review/memory.md).

<div data-with-frame="true">[image: Deep PR Review card in repository settings, with Edit buttons for strictness and linked repositories, and a View Memory button.]</div>

## Choose a strictness level

Strictness controls how much feedback Deep Review posts:

* **Chill**: Keeps reviews focused with fewer comments.
* **Normal**: Provides balanced coverage without overwhelming reviewers.
* **Assertive**: Flags more potential issues and improvements, resulting in more comments.

## Exclude paths and files

Aikido already ignores a large set of paths out of the box, including dependency folders, build output, generated files, lock files, binaries, and media. You can exclude more paths with a [`.aikido` file](/code-scanning/scanning-practices/ignore-via-code-with-aikido-files.md). If a pull request only changes excluded files, Aikido skips the review and you are not charged.

See [How Deep Review Works](/deep-review/how-deep-review-works.md#exclude-paths-and-files) for more detail.

## Credits

Deep Review costs **1 credit per pull request**.

Further reviews on the same PR (for example after new commits) don't consume extra credits. Credits are charged per PR, not per commit or per individual review run. Learn more about [Aikido Credits](/miscellaneous-info/wallet-and-credits.md).

{% hint style="info" %}
Curious what Deep Review would cost for your team?

Our [estimator](https://app.aikido.dev/repositories/deep-review) gives you an idea based on your activity of the last 30 days.
{% endhint %}

## Set a monthly spending limit

Workspace admins can cap Deep Review spend from [Wallet settings](https://app.aikido.dev/settings/wallet/settings).

{% stepper %}
{% step %}

### Open Wallet settings

Go to [Wallet settings](https://app.aikido.dev/settings/wallet/settings). You can also open **Organization Settings** and click **Manage Wallet**.
{% endstep %}

{% step %}

### Configure the Deep Review limit

Click **Configure Limits**. Turn on the **Deep Review** toggle, enter a **credits/month** cap, and click **Save**.
{% endstep %}
{% endstepper %}

<div data-with-frame="true">[image: Spending Limits dialog with Deep Review enabled and a 100 credits per month cap.]</div>

When usage reaches the cap, Aikido skips Deep Reviews that would charge a credit until the next calendar month. Further reviews on the same PR still run.

Turn the toggle off if you don't want a monthly cap.
