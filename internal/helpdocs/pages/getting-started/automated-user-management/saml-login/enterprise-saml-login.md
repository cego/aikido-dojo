# Enterprise SAML Login

Use the enterprise SAML login page to create an Aikido workspace and sign in with SAML. You do not need an existing Aikido account.

Open the page on the region where the workspace should live. The values you copy into your identity provider use that region's address.

| Region      | Enterprise SAML login page                        |
| ----------- | ------------------------------------------------- |
| EU          | <https://app.aikido.dev/login/enterprise/saml>    |
| US          | <https://app.us.aikido.dev/login/enterprise/saml> |
| Australia   | <https://app.au.aikido.dev/login/enterprise/saml> |
| Middle East | <https://app.me.aikido.dev/login/enterprise/saml> |
| US Gov      | <https://app.aikidogov.us/login/enterprise/saml>  |

On an on-prem installation, open `/login/enterprise/saml` on that installation. The page uses the installation's own address.

If the workspace already exists, turn on SAML from [General Settings](https://app.aikido.dev/settings/account) and follow the identity provider guide below.

## Set up your identity provider

Create a SAML 2.0 application in your identity provider and paste in the three values from the enterprise page. The provider guides below walk through that setup. On this page, start at the identity provider steps in the guide. The **Single Sign-On URL / ACS URL**, **Audience URI / Entity ID**, and **Name ID Format** come from the enterprise page.

* [Okta: Login with SAML](/getting-started/automated-user-management/saml-login/okta-login-with-saml.md)
* [JumpCloud: Login with SAML](/getting-started/automated-user-management/saml-login/jumpcloud-login-with-saml.md)
* [Google Workspaces: Login with SAML](/getting-started/automated-user-management/saml-login/google-workspaces-login-with-saml.md)
* [Microsoft Azure: Login with SAML / Entra ID](/getting-started/automated-user-management/saml-login/microsoft-azure-login-with-saml-entra-id.md)

After the workspace exists, set roles with [SAML Access Profiles](/getting-started/automated-user-management/saml-login/saml-user-rights-access-profiles-recommended.md) or [custom attributes](/getting-started/automated-user-management/saml-login/saml-user-rights-using-custom-attributes-advanced.md). For Entra ID claims, see [Microsoft Azure: Custom Attributes with SAML / Entra ID](/getting-started/automated-user-management/saml-login/microsoft-azure-custom-attributes-with-saml-entra-id.md).

## Create the workspace

**Step 1.** Open the enterprise SAML login page for your region and enter your work email, then click **Continue**.

Aikido takes the company domain from that email. Someone at `name@company.com` gets the domain `company.com`. That domain is how teammates later sign in with **Login via SSO**.

Use a company address such as `company.com` or `company.co.uk`. Public inboxes such as Gmail, Outlook, Yahoo, and iCloud are rejected. Click **Change email** if you need to start over.

**Step 2.** Copy these values into the SAML application in your identity provider:

| Field on the enterprise page | Paste it into your identity provider as   |
| ---------------------------- | ----------------------------------------- |
| Single Sign-On URL / ACS URL | ACS URL, Reply URL, or Single sign-on URL |
| Audience URI / Entity ID     | Audience URI, Entity ID, or SP Entity ID  |
| Name ID Format               | `EmailAddress`                            |

The ACS URL includes your company domain. On the EU app it looks like this:

`https://app.aikido.dev/api/saml/saml_auth?samlClientDomain=company.com`

**Step 3.** Back on the enterprise page, fill in the values your identity provider shows:

* **Entity ID / Issuer**
* **Single Sign-On URL.** Use your identity provider's own domain, for example `your-org.okta.com` or `login.microsoftonline.com`. A custom domain can be rejected. If Aikido returns "IdP domain not allowed", contact support.
* **X.509 Certificate**

**Company Email Domain** is filled from your work email and stays read-only.

A few providers expect an exact format:

* **Microsoft Entra ID.** Entity ID is `https://sts.windows.net/<tenant-id>/`. Single Sign-On URL is `https://login.microsoftonline.com/<tenant-id>/saml2`. The tenant id in both values is the same. In Entra ID these are the **Microsoft Entra Identifier** and the **Login URL**.
* **Okta.** The Entity ID starts with `http://www.okta.com/`, with `http`, as Okta shows it.

**Step 4.** Click **Continue**. Aikido creates the SAML client, signs you in through your identity provider, and opens onboarding.

If Aikido says this SAML client already exists, that domain or issuer is already in use. Contact support to connect your organization to it.

## Link the SAML client to the workspace

The enterprise page creates the SAML client for your company. It does not attach that client to a workspace yet. Other people who sign in with the same identity provider can reach this workspace only after you link it.

**Step 5.** Finish onboarding and connect your source code manager.

**Step 6.** In the new workspace, open [General Settings](https://app.aikido.dev/settings/account) and find **SAML Setup**. Click **Link existing SAML app**.

Enter the company domain from the enterprise page and the same **X.509 Certificate**, then click **Link SAML app**. Stay signed in with this SAML login, and be an admin of the workspace.

**Step 7.** Click **Manage SAML Authentication** and set **Auto-Onboarding** to **Enabled**. People assigned to the SAML app can then enter the workspace without an invite. With auto-onboarding left disabled, invite each person yourself.

**General Settings** then shows a login link you can share. Teammates can also use **Login via SSO** on the Aikido login screen and enter an email on the company domain. Both options are described at the end of each identity provider guide above.
