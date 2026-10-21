# Using Azure Devops Projects as Workspaces

By default, Aikido treats your Azure Devops organization as a single environment. Some teams, however, use Azure Devops projects to organize their repositories and want each project to appear as its own workspace in Aikido. This guide explains how to set that up.

### Prerequisites

* You need to be an Azure Devops admin

{% hint style="warning" %}
The feature must be enabled by Aikido support. Contact us via chat to activate.
{% endhint %}

### Steps

1. **Create the global Azure Devops workspace in Aikido**
   * Go through the onboarding flow
   * No need to select any repositories here; this workspace will be used to split into other projects.
2. **Access the project-based onboarding link**
   * Go to <https://app.aikido.dev/onboarding/azure/connect-project-based-group>
3. **Enter your Azure Devops project name**
   * Provide the name of the project you want to create an Aikido workspace for.
4. **Select repositories**
   * Click **Next**, choose the repositories to include, and proceed.
5. **Finish setup**
   * Your project-based workspace is now created.
6. **Repeat for additional projects**
   * Switch back to your **global workspace** using the top-left dropdown.
   * Repeat from step 2 for each additional project.
