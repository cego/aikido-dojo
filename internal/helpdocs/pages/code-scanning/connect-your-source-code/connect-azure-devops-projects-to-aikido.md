# Connect Azure DevOps Projects

Aikido allows you to connect your Azure DevOps projects to secure your code. To connect your Azure DevOps projects to Aikido, you will need to follow the steps below.

Note that each Azure Project with one or more repos will map to one Aikido workspace.

### Select Microsoft / Office 365 or Google to authenticate <a href="#select-microsoft--office-365-or-google-to-authenticate" id="select-microsoft--office-365-or-google-to-authenticate"></a>

**Step 1.** To connect your Azure DevOps project, you first need to authenticate via Microsoft / Office 365 or Google to create a user in Aikido. On the [signup screen](https://app.aikido.dev/login), click on 'Microsoft / Google' to continue.

> Depending on your organization's Microsoft Entra settings, you may need IT admin approval to authorize the Aikido application. [For more information about the approval process, check Azure docs.](https://learn.microsoft.com/en-us/entra/identity/enterprise-apps/grant-admin-consent?pivots=portal)

**Step 2.** Once you are authenticated via Google/Microsoft, you can go ahead and select 'Azure DevOps' on the page below.

**Step 3.** Enter the details to connect the Azure Devops project of your choosing. We'll explain how to obtain the required information right below.<br>

### Insert your organization's name <a href="#insert-your-organizations-name" id="insert-your-organizations-name"></a>

Enter the name of the Azure DevOps organization you'd like to connect. You can find this name by going to [https://dev.azure.com](https://dev.azure.com/) and copying it from the left-side navigation.

### Create a Personal Access Token <a href="#create-a-personal-access-token" id="create-a-personal-access-token"></a>

{% hint style="info" %}
If you're using EntraID groups as members of Devops teams, make sure the user who generates this token has Entra Read permissions.
{% endhint %}

{% hint style="success" %}
Once your workspace is connected, you can replace the personal access token with a Microsoft Entra service principal that never expires and stores no secret. See [Code Scanning for Azure DevOps with an Entra Service Principal](/code-scanning/connect-your-source-code/code-scanning-for-azure-devops-with-entra-service-principal.md).
{% endhint %}

Next up, you need to create a Personal Access Token to give us access to the resources you want.

* Log in to your Azure DevOps account
* In the upper-right corner, click the user settings icon next to your avatar. It looks like this:

  [image: User profile icon with a settings gear for account management options.]
* Select **Personal Access Token** in the dropdown.
* Click on the **+New Token** button in the top left corner
* Enter a name for the token, eg: 'Aikido Security Access Token'
* Make sure to select the same organization as the one you entered in the previous step
* Make sure to select an expiration date in the future, the max should be 1 year.
* Next, we need the following scopes to be selected (click 'show all scopes'):
  * **Code: Read**
  * **Member Entitlement Management: Read**
  * **Project and Team: Read**
  * **User Profile: Read**
* Click the **Create** button at the bottom.

[image: Editing Azure DevOps personal access token with custom-defined code read permissions.]

* Copy the token being shown on the following screen and enter it in the input field .

**Important**: You will no longer be able to view the value of the token once you hit continue. Make sure you copied it first.

Aikido will now check the connection to your Azure DevOps organization. If the connection was not successful, make sure to double-check the organization name and personal access token you provided.

### Select the project and repos you'd like to secure <a href="#select-the-project-and-repos-youd-like-to-secure" id="select-the-project-and-repos-youd-like-to-secure"></a>

On the next screen, you can select which project you'd like to start with. You'll always be able to connect more of your projects to Aikido.

In the final step you can select all the repos you would like us to monitor.

#### Specific notes on TFVC Repos

Aikido supports the integration of both Git and Team Foundation Version Control (TFVC) repositories.

* For TFVC repositories, Aikido **does not** perform secret scanning.
* By default, Aikido tries to scan 'All Branches' for TFVC repos. This might result in slow scanning speed. **It is recommended to select the actual branch to scan.** This can be changed by clicking the branch label on the repo detail page.
