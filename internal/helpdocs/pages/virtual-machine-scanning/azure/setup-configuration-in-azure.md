# Azure VM Access Configuration

This page describes how to grant access to Aikido to scan your Azure Virtual Machines.

## Entra ID App Registration

Log into your [**Azure Portal**](https://portal.azure.com/) and navigate to the **Microsoft Entra ID service**.

Click on **Add** and select **App registration**

[image: Azure portal: Add new user, group, enterprise app, or app registration in Default Directory.]

Give the application a meaningful name, we need this name later.

Leave the **Supported account types** default: **Accounts in this organizational directory only**.

Click on **Register**.

[image: Azure portal screen for registering a new application named "AikidoSecurity" with account type selection.]

You get redirected to the detail page of the newly created application. Here you can find and copy the **Application (client) ID** and the **Directory (tenant) ID**

[image: Azure portal showing AikidoSecurity application's client ID and essential overview details.]

At the client credentials field, click "Add a certificate or secret"

[image: Azure portal displaying AikidoSecurity app overview and client credential configuration options.]

Click the "New client secret"-button, give a description for the secret and set the expiration date to 2 years (730 days / 24 months)

[image: Creating a new client secret for application authentication in Azure AD.]

Copy the **Secret's Value**

[image: Azure portal client secret management screen showing secret description, expiry date, and copy option.]

You now have all the required values to add the Azure Cloud via the Public API once the application setup is complete in Azure Portal.

## Azure RBAC Role Assignment

Go to the subscription detail page. Now we need to make sure we grant access to the roles we need.

Navigate to **Subscriptions**, find the relevant Subscription for your Virtual Machines

Click on **"Access Control (IAM)"**.

[image: Azure IAM access control panel for managing roles and permissions in a subscription.]

Go to the Role assignments tab & Click on **"Add"**, then **"Add role assignment"**.

[image: Azure portal interface for adding role assignments and managing classic administrators.]

In the **"Role"** tab, search and select **"VM Scanner Operator"** & Click **"Next"**.

[image: Assigning "VM Scanner Operator" role for disk snapshot security analysis in Azure.]

Leave the **"Assign access to**" default value.

Click on **"Select Members"**, search for the name of the app registration (e.g. "AikidoSecurity") you created and select it.

Click **"Select"**

Click **"Review + assign"** twice

[image: Assigning a role to a member in Azure subscription using Access Control (IAM) settings.]

Repeat the role assignment process for the role **"Disk Snapshot Contributor"**.

The App Registration now has the required roles to scan your virtual machines.
