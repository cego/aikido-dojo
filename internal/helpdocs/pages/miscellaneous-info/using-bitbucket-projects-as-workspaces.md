# Using Bitbucket Projects as Workspaces

By default, Aikido treats your Bitbucket workspace as a single environment. Some teams, however, use Bitbucket projects to organize their repositories and want each project to appear as its own workspace in Aikido. This guide explains how to set that up.

### Prerequisites

* Currently supported only for Bitbucket.
* You need to be a Bitbucket admin

{% hint style="warning" %}
The feature must be enabled by Aikido support. Contact us via chat to activate.
{% endhint %}

### Steps

1. **Create the global Bitbucket workspace in Aikido**
   * Go through the onboarding flow
   * No need to select any repositories here; this workspace will be used to split into other projects.
2. **Access the project-based onboarding link**
   * Go to [https://app.aikido.dev/onboarding/bitbucket/connect-project-based-group](https://app.aikido.dev/onboarding/bitbucket/connect-project-based-group?utm_source=chatgpt.com)
3. **Enter your Bitbucket project key**
   * Provide the key of the project you want to create an Aikido workspace for.
4. **Select repositories**
   * Click **Next**, choose the repositories to include, and proceed.
5. **Finish setup**
   * Your project-based workspace is now created.
6. **Repeat for additional projects**
   * Switch back to your **global workspace** using the top-left dropdown.
   * Repeat from step 2 for each additional project.
