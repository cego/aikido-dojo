# JumpCloud: Login with SAML

{% hint style="info" %}
This feature is only available on a **paying** plan and is not enabled by default. **Pro plans and up** can enable SAML [via this link](https://app.aikido.dev/settings/saml/enable).
{% endhint %}

> If you switch to SAML Login instead of auto-onboarding via your Git provider, team import from GitHub, Bitbucket, or Azure DevOps will no longer work. You will need to manage your teams manually moving forward, either through the Aikido UI or [Access Profiles.](/getting-started/automated-user-management/saml-login/saml-user-rights-access-profiles-recommended.md)

## Setting up SAML in your account <a href="#setting-up-saml-in-your-account" id="setting-up-saml-in-your-account"></a>

**Step 1.** Go to [**General Settings**](https://app.aikido.dev/settings/account) and click '**Enable SAML Authentication'**

[image: Workspace info page with option to enable SAML authentication.]

**Step 2.** Copy **all details** to your identity provider. See steps below.

[image: SAML authentication setup screen with required URLs and Name ID format.]

### Continue in JumpCloud <a href="#continue-in-jumpcloud" id="continue-in-jumpcloud"></a>

**Step 1.** Go to **User Authentication** > **SSO Applications** in the JumpCloud Admin Portal navigation.

**Step 2.** Click the **Add New Application** and search for **SAML 2.0**. Click **Next**.

**Step 3.** Choose a **Display Label** and click **Save Application**.

**Step 4.** Click **Configure Application**.

**Step 5.** Click on the **SSO** tab and fill following fields:

* **Idp Entity ID:** `https://console.jumpcloud.com/<appname>`
* **SP Entity ID:** `https://app.aikido.dev/saml`
* **ACS URLs** - **Default URL:** `https://app.aikido.dev/api/saml/saml_auth?samlClientId=...` (As shown in Aikido)
* **SAMLSubject NameID:** `email`
* **SAMLSubject NameID Format:** `urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress`
* **Signature Algorithm:** `RSA-SHA256`
* **Default RelayState:** You can leave this empty

[image: SAML 2.0 single sign-on configuration settings for secure identity management integration.]

**Step 6.** We'll continue in Aikido, but you might as well click **Save** and come back to this screen.

### Go back to Aikido <a href="#go-back-to-aikido" id="go-back-to-aikido"></a>

* Fill in the:
  * **Entity ID / Issuer:** `https://console.jumpcloud.com/<appname>` (Make sure this matches what you've entered as **Idp Entity ID** in JumpCloud. If you're having issues with this, see the Troubleshooting section at the bottom)
  * **Single Sign-On URL:** as shown in JumpCloud under **IDP URL**. (looks like `https://sso.jumpcloud.com/saml2/<appname>`)
  * **X.509 Certificate**: This can be fetched in different ways. One way is to click **Export Metadata** in JumpCloud the config and open the downloaded xml. You'll find your certificate between the `ds:X509Certificate` tags.
* Also fill out the **Company Domain** to make sure people can log in without the need of a Single Sign-On URL.

[image: SAML Authentication setup form for configuring identity provider details.]

> Success! People having access to your JumpCloud SAML app will now be able to auto-onboard to your Aikido workspace.

### 2 options for users to login using your SAML client <a href="#id-2-options-for-users-to-login-using-your-saml-client" id="id-2-options-for-users-to-login-using-your-saml-client"></a>

**Option 1. Using SSO Link Directly**

Copy the Login Link and share this internally with other users.

[image: SAML Authentication settings with options to manage or copy the login link.]

**Option 2.** Going to the Aikido login screen, selecting **Login Via SSO** and filling in the email address **Important**: the email needs to contain the company domain that has been set up.

[image: One-click login and sign-up with Google, Microsoft, or SSO, with no credit card needed.]

[image: Login screen with Google, Microsoft, and email login options.]

### Troubleshooting <a href="#troubleshooting" id="troubleshooting"></a>

**Error**

[image: SAML client already exists error; prompts user to contact support for linking organization.]

**Solution**

Make sure the **Idp Entity ID** is unique. Perhaps you could change it to `https://console.jumpcloud.com/<samlClientId>`. Note that you'll also need to change it in Aikido in **Entity ID / Issuer** as these should match.
