# Setup Custom Role in Azure

This document guides you through creating a new App Registration with client credentials, and a custom RBAC role for Aikido VM scanning.

The credentials will be used by Aikido to make the necessary API requests to scan your Virtual Machines. Access to scan your Virtual Machines will be granted to the App Registration using a custom role.

These permissions are limited to the minimum required for Virtual Machine scanning:

Use the following steps to create

Log into your [**Azure Portal**](https://portal.azure.com/) and navigate to the **Microsoft Entra ID service**.

Click on **Add** and select **App registration**

[image: Azure Portal: Adding a new app registration under Default Directory Overview.]

Give the application a meaningful name, we need this name later.

Leave the **Supported account types** default: **Accounts in this organizational directory only**.

Click on **Register**.

[image: Registering a new application named "AikidoSecurity" in Microsoft Azure Active Directory.]

You get redirected to the detail page of the newly created application. Here you can find and copy the **Application (client) ID** and the **Directory (tenant) ID**

[image: Azure portal displaying AikidoSecurity application overview and client ID information.]

At the client credentials field, click "Add a certificate or secret"

[image: Azure portal app registration overview, highlighting "Add a certificate or secret" for client credentials.]

Click the "New client secret"-button, give a description for the secret and set the expiration date to 2 years (730 days / 24 months)

[image: Creating a new client secret for application authentication in Azure Portal.]

Copy the **Secret's Value**

[image: Azure portal showing an active client secret for application authentication with expiration date.]

You now have all the required values to configure VM scanning in Aikido, once the application setup is complete in Azure Portal. Next, we need to make sure we grant the application access for VM scanning.

Navigate to **Subscriptions**, find the relevant Subscription for your Virtual Machines

Click on **"Access Control (IAM)"**.

[image: Azure IAM portal for managing role assignments, permissions, and access levels on subscriptions.]

Click on the **"Add"** button.

Select **"Add custom role"**

[image: Azure IAM: Add role assignment, co-administrator, or custom role for subscription access control.]

Go to the **"JSON"** tab and open the editor by clicking on **"Edit"**

[image: Azure portal interface for creating and downloading a custom role JSON definition.]

Copy generated JSON config from the Aikido setup screen, paste it into the editor

Click **"Save"**

[image: Creating a custom Azure IAM role using JSON permissions template.]

At the bottom, click **"Review + assign"**, then **"Create"**

[image: Azure custom role creation: review of permissions and assignable scopes for AikidoVMScanner.]

Now that the custom role is created, we can assign it to the App Registration we created at the start.

Navigate to **Subscriptions**, find the relevant Subscription for your Virtual Machines

Click on **"Access Control (IAM)"**.

[image: Azure portal IAM: Manage access roles and permissions for subscription resources.]

Go to the Role assignments tab & Click on **"Add"**, then **"Add role assignment"**.

[image: Role assignment menu in Azure portal with options to add or download roles.]

In the **"Role"** tab, search and select the custom role you created ("Aikido VM Scanner") & Click **"Next"**.

[image: Assigning the "Aikido VM Scanner" role in Azure for virtual machine access.]

Leave the **"Assign access to"** default value.

Click on **"Select Members"**, search for the name of the app registration (e.g. "AikidoSecurity") you created and select it.

Click **"Select"**

Click **"Review + assign"** twice

[image: Assigning user, group, or application roles in Azure subscription via role assignment interface.]

The App Registration now has the required permissions to scan your Azure Virtual Machines.
