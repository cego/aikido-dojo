# Microsoft Azure: Custom Attributes with SAML /Entra ID

First, make sure you have SAML login working using following guide:

[https://help.aikido.dev/doc/microsoft-azure-login-with-saml--entra-id/doc74BfKR60Z](/getting-started/automated-user-management/saml-login/microsoft-azure-login-with-saml-entra-id.md)

### Setting up Azure Group based SAML custom attributes <a href="#setting-up-azure-group-based-saml-custom-attributes" id="setting-up-azure-group-based-saml-custom-attributes"></a>

1. Go to the application registration

   [image: Azure portal view for managing Aikido-SSO enterprise application properties and settings.]
2. Create an app role.**value** here should be the value of the claim. In this example, we're setting up for `aikido_role`, so valid values for this are `admin`, `default`, `team_only`.

   [image: Creating a new app role "Aikido Admin" in Microsoft Azure for Aikido-SSO.]
3. After saving, go back to the app settings, and add a group to 'Users and Groups'

   [image: Azure portal: Assign users or groups to the Aikido-SSO enterprise application.]
4. Add the Entra group you'd like to give admin access (in this case) and add the role we created in step 2.

   [image: Azure Add Assignment: Select users, groups, and roles for directory permissions.]
5. Back in the Single Sign-on settings of the app, go to the Attributes & Claims -> Edit

   [image: Azure portal SAML-based single sign-on configuration for Aikido-SSO application.]
6. Click 'Add new claim'

   [image: Azure portal Attributes & Claims page for adding and managing SAML claims.]
7. Fill in the attribute name & user.assignedroles as source attribute. (this is the `admin` value we set up in step 2)

   [image: Azure claim setup: mapping "aikido\_role" to "user.assignedroles" attribute.]
8. All done. On SAML login, these changes will take effect.
