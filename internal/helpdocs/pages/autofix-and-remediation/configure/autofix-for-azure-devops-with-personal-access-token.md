# AutoFix for Azure DevOps with Personal Access Token

> When the installation of the Azure DevOps Autofix integration does not work via the OAuth flow, you can also configure it to work with a Personal Access Token (PAT).

> Please note that for the integration to work properly, you need to be able to have write access in all of the projects in your organization. This is usually a reserved permission for workspace admins.

**First**, go to the PAT management page in Aikido, which can be [found here](https://app.aikido.dev/settings/integrations/azure-devops/autofix/personal-access-token).

**Next**, follow the instructions below to generate a PAT in Azure DevOps.

## Generating a PAT <a href="#generating-a-pat" id="generating-a-pat"></a>

1. Navigate to '**User Settings**' > '**Personal Access Tokens**', via the header icon shown here

   

   Or use the shortcut link '**Generate Access Token here**' on the personal access token management page in Aikido.

   
2. You should end up on a page similar to this

   
3. Make sure you have selected the same organization as the one in your Aikido workspace via the '**Access scope**' button on the top right
4. Click on "**New token**"
5. Next you need to select the following permissions for the token: **Code (Read & Write)** and **Graph (Read).**\
   If you don't see the **Graph** scope appear, you will have to click the blue link at the bottom "Show more scopes"

   
6. For the expiration, use custom defined and set it to 1 year from now, this is the maximum duration that Azure DevOps allows the token to be valid for.
7. Click on '**Create**' to create the access token
8. Copy the token on the next screen and insert it on the personal access token management page in Aikido.
