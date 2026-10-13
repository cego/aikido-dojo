# Connect Azure DevOps Self-Managed Server

Aikido allows you to connect your Azure DevOps Self-Managed server instance to secure your code. To connect your Azure DevOps Self-Managed projects to Aikido, use the following procedure.

Note that each Azure Project with one or more repos maps to one Aikido workspace.

### Select Microsoft / Office 365 or Google to authenticate <a href="#select-microsoft--office-365-or-google-to-authenticate" id="select-microsoft--office-365-or-google-to-authenticate"></a>

**Step 1.** To connect your Azure DevOps project, you first need to authenticate via Microsoft / Office 365 or Google to create a user in Aikido. On the [signup screen](https://app.aikido.dev/login), click on 'Microsoft / Google' to continue.

> Depending on your organization's Microsoft Entra settings, you may need IT admin approval to authorize the Aikido application. [For more information about the approval process, check Azure docs.](https://learn.microsoft.com/en-us/entra/identity/enterprise-apps/grant-admin-consent?pivots=portal)

**Step 2.** Once you are authenticated via Google/Microsoft, select 'Connect Azure' on the page.

**Step 3.** Enter the details to connect the Azure Devops project of your choosing. We explain how to obtain the required information in the following section.

### Insert your collection name <a href="#insert-your-organizations-name" id="insert-your-organizations-name"></a>

Enter the name of the Azure DevOps collection/organization you want to connect. You can find this name by going to your Self-Managed server instance and copying it from the left-side navigation.

### Create a Personal Access Token <a href="#create-a-personal-access-token" id="create-a-personal-access-token"></a>

{% hint style="info" %}
The token holder must have at least Project Collection Administrator access.
{% endhint %}

Next up, you need to create a Personal Access Token to give us access to the resources you want.

1. Log in to your Azure DevOps Self-Managed instance
2. In the upper-right corner, click the user avatar. A dropdown will appear, like the following image:\
   \
   
3. Select **Security** in the dropdown.
4. Then select on the left **Personal Access Token**
5. Click on the **+New Token** button in the top right corner
6. Enter a name for the token, eg: 'Aikido Security Access Token'
7. Make sure to select the same organization/collection as the one you entered in the previous step
8. Make sure to select an expiration date in the future, the max should be 1 year.
9. Next, we need the following scopes to be selected (click 'show all scopes'):
   * **Code: Read**
   * **Member Entitlement Management: Read**
   * **Project and Team: Read**
   * **User Profile: Read**
10. Click the **Create** button at the bottom.

[image: Editing Azure DevOps personal access token with custom-defined code read permissions.]

11. Copy the token being shown on the following screen and enter it in the input field .

**Important**: You can't view the token value once you hit continue. Make sure you copy it first.

Aikido now checks the connection to your Azure DevOps instance. If the connection is not successful, make sure to double-check the organization name and personal access token you provided.

### Select the project and repos you'd like to secure <a href="#select-the-project-and-repos-youd-like-to-secure" id="select-the-project-and-repos-youd-like-to-secure"></a>

On the next screen, you can select which project you'd like to start with. You can always connect more of your projects to Aikido.

In the final step select all the repos you would like us to monitor.

#### Specific notes on TFVC Repos

Aikido supports the integration of both Git and Team Foundation Version Control (TFVC) repositories.

* For TFVC repositories, Aikido does not perform secret scanning.
* By default, Aikido tries to scan 'All Branches' for TFVC repos. This might result in slow scanning speed. **We recommend selecting the actual branch to scan.** This can be changed by clicking the branch label on the repo detail page.
