# Microsoft Azure: Login with SAML/ Entra ID

{% hint style="info" %}
This feature is only available on a **paying** plan and is not enabled by default. **Pro plans and up** can enable SAML [via this link](https://app.aikido.dev/settings/saml/enable).
{% endhint %}

> If you switch to SAML Login instead of auto-onboarding via your Git provider, team import from GitHub, Bitbucket, or Azure DevOps will no longer work. You will need to manage your teams manually moving forward, either through the Aikido UI or [Access Profiles.](/getting-started/automated-user-management/saml-login/saml-user-rights-access-profiles-recommended.md)

## Setting up SAML in your account <a href="#setting-up-saml-in-your-account" id="setting-up-saml-in-your-account"></a>

**Step 1.** Go to [**General Settings**](https://app.aikido.dev/settings/account) and click '**Enable SAML Authentication'**

[image: Workspace info screen with option to enable SAML authentication for GitHub account.]

**Step 2.** Copy **all details** to your identity provider. See steps below.

[image: SAML Authentication setup screen showing required URLs and Name ID format for configuration.]

### Continue in Azure <a href="#continue-in-azure" id="continue-in-azure"></a>

**Step 1.** Go to **Microsoft Entra ID**.

**Step 2.** Click the **Add** dropdown and select **Enterprise application**.

[image: Adding a new enterprise application in Microsoft Azure Active Directory.]

**Step 3.** Click **Create your own application**, choose a name for your app and select 'Non-gallery'.

[image: Creating a custom non-gallery application named "Aikido-SSO" for integration purposes.]

**Step 4.** Select **Set up single sign on**.

[image: Aikido-SSO application setup: Assign users and configure single sign-on in Microsoft Entra.]

**Step 5.** Click the **SAML** option.

[image: Enable SAML single sign-on for secure application authentication in Aikido-SSO.]

**Step 6.** On **step 1**, click **Edit.**

[image: Basic SAML Configuration: Edit required Identifier and Reply URL fields.]

**Step 7.** Fill in the **Entity ID** and **ACS URL** as shown in Aikido.

[image: Configuration screen for SAML SSO with Entity ID and Reply URL fields specified.]

**Step 8.** At **step 2**, click **Edit.**

[image: User attributes and claims mapping with editable options highlighted.]

**Step 9.** Click the **Unique User Identifier (Name ID)**.\
Optional: clicking 'Add new claim' at the top of this page allows you to add [custom attributes](/getting-started/automated-user-management/saml-login/saml-user-rights-using-custom-attributes-advanced.md) to SAML. More info [here](https://help.aikido.dev/doc/microsoft-azure-custom-attributes-with-saml--entra-id/docFaysVwVZy).

[image: Highlighted SAML claim: Unique User Identifier (Name ID) for user identification.]

**Step 10.** Make sure to set **Source attribute** to `user.mail` here.

[image: Configuring a SAML claim for user email as the name identifier in Azure AD.]

**Step 11.** At step 3 you can download the **Certificate (Base64)** & at step 4 you'll see the **Login URL** and **Mircosoft Entra Identifier**. These should be copy and pasted to Aikido.

### Go back to Aikido <a href="#go-back-to-aikido" id="go-back-to-aikido"></a>

* Fill in the **Entity ID / Issuer** using Azure's **Microsoft Entra Identifier**, and the **Single Sign-On URL** using Azure's **Login URL** (both from Step 11), plus the **X.509 Certificate** as shown in Azure.
* Also fill out the **Company Domain** to make sure people can log in without the need of a Single Sign-On URL.

[image: SAML Authentication setup form for configuring Single Sign-On (SSO) credentials.]

> Success! People having access to your Azure SAML app will now be able to auto-onboard to your Aikido workspace.

### 2 options for users to login using your SAML client <a href="#id-2-options-for-users-to-login-using-your-saml-client" id="id-2-options-for-users-to-login-using-your-saml-client"></a>

**Option 1. Using SSO Link Directly**

Copy the Login Link and share this internally with other users.

[image: SAML Authentication settings with options to manage or copy the login link.]

**Option 2.** Going to the Aikido login screen, selecting **Login Via SSO** and filling in the email address **Important**: the email needs to contain the company domain that has been set up.

[image: One-click login and sign-up with Google, Microsoft, or SSO; no credit card needed.]

[image: Login screen offering Google, Microsoft, or email sign-in options.]
