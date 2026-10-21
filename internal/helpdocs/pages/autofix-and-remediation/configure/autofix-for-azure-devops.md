# AutoFix for Azure DevOps using OAuth

Aikdo Autofix is a tool you can use to have Aikido fix vulnerabilities in 3rd party dependencies in your projects. It will do this by creating pull requests that remove the vulnerability via package updates or by other means. In some cases an Aikido Autofix can remove a whole class of vulnerabilities instead of just 1 issue.

**Setup Autofix for Azure DevOps**

> All users within your workspace will need to setup Autofix individually.

By default, Aikido only has read access on your Azure DevOps instance. To use Aikido Autofix, a separate access token with write access is required. Please make sure that "Third-party application access via Oauth" is enabled for your organization, by going to "Organization settings" and then clicking "Policies".

[image: Organization security policies for application access, public projects, and user invitations.]

**Step 1. Enable Autofix on the** [**Autofix Settings**](https://app.aikido.dev/settings/integrations/autofix) **page or go to** [**Autofix Page**](https://app.aikido.dev/issues/fix) **and click on Enable Autofix.**

**Step 2.** Click **Authorize**

**Step 3.** Grant Aikido permissions to **Write**

**Step 4.** Click save and you are all set. You will now be able to execute autofix PRs from the [Autofix page](https://app.aikido.dev/issues/fix) or from the action menu for subissues in the sidebar ([read more here](/autofix-and-remediation/scope/autofix-for-open-source-dependencies.md)).
