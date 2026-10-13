# Bitbucket Cloud: Assigning Repositories to Teams via Projects

### Introduction

In Bitbucket, repositories can be linked to a project, and that project can be linked to a team. Aikido now supports automatically linking those repositories to the team their Bitbucket project belongs to.

{% hint style="info" %}
This feature is currently behind a feature flag. Please contact [Aikido Support](mailto:support@aikido.dev) to have it enabled for your workspace.
{% endhint %}

### Setup

Once support has enabled the feature for your workspace:

1. Go to [**Settings → Bitbucket → Update Workspace Access Token**](https://app.aikido.dev/onboarding/bitbucket/update-workspace-access-token)
2. Scroll to the **Show Extra Tokens** section and add a new workspace token.
3. Make sure the token includes the following scopes:
   * `repository:admin`
   * `project:admin`
4. Save the token.

Aikido will now use your Bitbucket project structure to assign repositories to the correct teams during the next team sync.

{% hint style="info" %}
These additional scopes are required by the Bitbucket v2 API to read project-level access.
{% endhint %}
