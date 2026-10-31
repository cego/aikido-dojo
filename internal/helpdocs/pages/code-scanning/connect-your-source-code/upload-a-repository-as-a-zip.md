# Upload a Repository as a Zip

When you cannot connect a git provider, upload a static `.zip` as a snapshot of the code. Aikido adds that snapshot as a repository so you can run AI Code Analysis, whitebox pentests, and other features that read the code. Actions that need a live git connection, such as an AutoFix pull request, are not available yet.

{% hint style="info" %}
Zip upload is the exception to Aikido's usual rule of not keeping your code. Aikido stores this archive with the repository so the snapshot stays available for AI Code Analysis, whitebox pentests, and other features that read it. You can delete the repository at any time.
{% endhint %}

## Prerequisites

* You are a workspace admin.
* The archive is a `.zip` of one repository.
* Include the code you want assessed: source, lockfiles, and any dependency folders or build output that should be part of the snapshot.

## How it works

Aikido stores the zip and uses it as the repository contents. Features that read the snapshot can use it. Features that need a live git connection cannot, and Aikido is adding support for more of them.

Uploading a new version of an existing repository is not supported yet. Support for that is coming soon.

## Upload the repository

{% stepper %}
{% step %}

### Create the zip

From the repository root, run:

```shell
git archive --format=zip --output=my-app.zip HEAD
```

`git archive` packs tracked files at `HEAD`. Zip the folder from your file manager when the snapshot should also include untracked files, such as installed dependencies or build output. Files can sit at the root of the archive or inside one top-level folder.
{% endstep %}

{% step %}

### Open repository settings

Go to [Repository settings](https://app.aikido.dev/settings/integrations/repositories). From **Repositories**, click **Manage Repos**.
{% endstep %}

{% step %}

### Choose upload source code

Click **Add repo**. Select **Upload source code**, then click **Continue**.

<div data-with-frame="true">[image: Add a repository dialog with Upload source code selected]</div>
{% endstep %}

{% step %}

### Name the repository and add the zip

Enter a **Repository name**. Click the upload area, or drag the `.zip` onto it.
{% endstep %}

{% step %}

### Upload

Click **Upload**. Aikido adds the repository, confirms with **Repository uploaded**, and starts the scan automatically.
{% endstep %}

{% step %}

### Select the repository in a setup flow

Open [AI Code Analysis](/ai-code-analysis/ai-code-analysis.md) or [set up a pentest](/pentests/configure-a-pentest/starting-an-assessment.md). When the setup flow asks you to choose repositories, select this one. Other features that read repository code work the same way.
{% endstep %}
{% endstepper %}

## Troubleshooting

| Problem                     | What to do                                                                                                          |
| --------------------------- | ------------------------------------------------------------------------------------------------------------------- |
| The archive is too large    | Remove files you can spare, then zip the repository again.                                                          |
| The file is not a valid zip | Create the archive with `git archive` or your file manager. Renaming another file to `.zip` does not make it a zip. |

## FAQ

### How do I upload a newer copy of the code?

Uploading a new version of an existing repository is not supported yet. Support for that is coming soon.

### Which features are unavailable?

The zip is a static snapshot, so anything that needs a live git connection is limited. Opening an AutoFix pull request is one example. AI Code Analysis, whitebox pentests, and other features that read the code can use the snapshot. Connect a git provider when you need the full set.

## Related

{% content-ref url="/pages/5VOmohZt3vDXzoKPHmDk" %}
[Connect Your Repositories](/getting-started/setting-up-your-account/create-account-and-connect-your-repositories.md)
{% endcontent-ref %}

{% content-ref url="/pages/KVMq059PBtgOs1XY2NCO" %}
[Local Code Scanning](/code-scanning/local-code-scanning.md)
{% endcontent-ref %}
