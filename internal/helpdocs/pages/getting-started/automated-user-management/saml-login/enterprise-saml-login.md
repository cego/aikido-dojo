# Enterprise SAML Login

Create a company SAML client from the enterprise login page, then link it to the new workspace before anyone else signs in.

Use the enterprise SAML login page when the workspace does not exist yet. The page creates a SAML client for your company domain and signs you in.

{% hint style="warning" %}
Link that client to the workspace before anyone else signs in. Until it is linked, SSO opens the create-workspace page, including for people you invited.
{% endhint %}

Open the page on the region where the workspace should live. The values you copy into your identity provider use that region's address.

| Region      | Enterprise SAML login page                        |
| ----------- | ------------------------------------------------- |
| EU          | <https://app.aikido.dev/login/enterprise/saml>    |
| US          | <https://app.us.aikido.dev/login/enterprise/saml> |
| Australia   | <https://app.au.aikido.dev/login/enterprise/saml> |
| Middle East | <https://app.me.aikido.dev/login/enterprise/saml> |
| US Gov      | <https://app.aikidogov.us/login/enterprise/saml>  |

On an on-prem installation, open `/login/enterprise/saml` on that installation. The page uses the installation's own address.

If the workspace already exists, turn on SAML from [General](https://app.aikido.dev/settings/account) under **Settings**. That path links the client while you create it.

## Create the SAML client

{% stepper %}
{% step %}

### Open the enterprise page

Open the enterprise SAML login page for your region, enter your work email, and click **Continue**.

Aikido takes the company domain from that email. `name@company.com` gives the domain `company.com`. Teammates later sign in with **Continue with SSO** using an email on that domain.

Use a company address such as `company.com` or `company.co.uk`. Public inboxes such as Gmail, Outlook, Yahoo, and iCloud are rejected. Click **Change email** to start over.
{% endstep %}

{% step %}

### Copy the values into your identity provider

Copy these values into a new SAML 2.0 application:

| Field on the enterprise page | Paste it into your identity provider as   |
| ---------------------------- | ----------------------------------------- |
| Single Sign-On URL / ACS URL | ACS URL, Reply URL, or Single sign-on URL |
| Audience URI / Entity ID     | Audience URI, Entity ID, or SP Entity ID  |
| Name ID Format               | `EmailAddress`                            |

The ACS URL includes your company domain. On the EU app it looks like this:

`https://app.aikido.dev/api/saml/saml_auth?samlClientDomain=company.com`

The identity provider clicks start at these sections. Stop when the application exists there, then return to this page.

* [Continue in Okta](/getting-started/automated-user-management/saml-login/okta-login-with-saml.md#continue-in-okta)
* [Continue in JumpCloud](/getting-started/automated-user-management/saml-login/jumpcloud-login-with-saml.md#continue-in-jumpcloud)
* [Continue in Google](/getting-started/automated-user-management/saml-login/google-workspaces-login-with-saml.md#continue-in-google)
* [Continue in Azure](/getting-started/automated-user-management/saml-login/microsoft-azure-login-with-saml-entra-id.md#continue-in-azure)
  {% endstep %}

{% step %}

### Paste the identity provider values

On the enterprise page, fill in the values your identity provider shows:

* **Entity ID / Issuer**
* **Single Sign-On URL.** Use your identity provider's own domain, for example `your-org.okta.com` or `login.microsoftonline.com`. A custom domain can be rejected. If Aikido returns "IdP domain not allowed. Please contact support.", [contact support](/miscellaneous-info/contact-support.md).
* **X.509 Certificate**

**Company Email Domain** is filled from your work email and stays read-only.

Keep the certificate and the company domain. Linking asks for both, and Aikido does not show the certificate again. If another person creates the workspace, send them those two values.

A few providers expect an exact format:

* **Microsoft Entra ID.** Entity ID is `https://sts.windows.net/<tenant-id>/`. Single Sign-On URL is `https://login.microsoftonline.com/<tenant-id>/saml2`. The tenant id in both values is the same. In Entra ID these are the **Microsoft Entra Identifier** and the **Login URL**.
* **Okta.** The Entity ID starts with `http://www.okta.com/`, with `http`, as Okta shows it.
  {% endstep %}

{% step %}

### Create the workspace

Click **Continue**. Aikido creates the SAML client, signs you in through your identity provider, and opens workspace setup.

Connect your source code manager, then go to [Link the client to the workspace](#link-the-client-to-the-workspace) before you invite anyone.

If Aikido says this SAML client already exists, that domain or issuer is already registered. [Contact support](/miscellaneous-info/contact-support.md) to connect your organization to it. Creating another workspace leaves the client unlinked.
{% endstep %}
{% endstepper %}

## Link the client to the workspace

A workspace admin does this, signed in with the SSO login from the enterprise page. If someone else created the identity provider application, use the company domain and certificate they saved.

{% stepper %}
{% step %}

### Open SAML Setup

In the Aikido app you just signed in to, go to **Settings** → **General**. On the EU app that page is [General](https://app.aikido.dev/settings/account). On another region, open `/settings/account` on that region's host.

Find the **SAML Setup** card.
{% endstep %}

{% step %}

### Attach the SAML client

Match the buttons on the **SAML Setup** card:

* **Link existing SAML app** is shown. The line under it says "This option copies SAML settings from another Aikido workspace." Click the button. Enter **Company Domain** and the same **X.509 Certificate**, then click **Link SAML app**.
* The card only shows **Authentication**, the status "SAML authentication is disabled", and **Enable SAML Authentication**. Leave that button. [Contact support](/miscellaneous-info/contact-support.md), send the company domain, and ask them to link this SAML client to the workspace.

**Enable SAML Authentication** starts a second client. Aikido then stops with "This SAML client already exists. Please contact support if you want to connect the current organization to an existing SAML client."
{% endstep %}

{% step %}

### Turn on auto-onboarding

After the client is linked, the card shows that auto-onboarding is disabled. Click **Manage SAML Authentication**, set **Auto-Onboarding** to **Enabled**, and click **Save Settings**.

People assigned to the SAML app can then open this workspace without an invite. Leave auto-onboarding disabled when you want to invite each person yourself.
{% endstep %}
{% endstepper %}

**SAML Setup** then shows **Copy Login link**. Teammates can also click **Continue with SSO** on the Aikido login screen and enter an email on the company domain.

## Troubleshooting

### Colleagues land on the create-workspace page

The SAML client is not linked to this workspace yet. A workspace admin who is signed in with that SSO login completes [Link the client to the workspace](#link-the-client-to-the-workspace).

### Aikido says this SAML client already exists

On the enterprise page, the domain or issuer is already registered. [Contact support](/miscellaneous-info/contact-support.md). Creating another workspace leaves the client unlinked.

The same message appears after **Enable SAML Authentication**. That button starts a new client. Use **Link existing SAML app**, or contact support when that button is not on the card.

### The link is rejected

The workspace admin has to be signed in with the SAML app being linked. Sign out, click **Continue with SSO**, enter an email on the company domain, open the workspace, and try **Link existing SAML app** again.
