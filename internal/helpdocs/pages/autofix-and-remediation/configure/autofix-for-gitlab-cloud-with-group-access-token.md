# AutoFix for Gitlab Cloud with Group Access Token

> It is recommended to use a Group Access Token to setup the Gitlab AutoFix integration. If you don't have GitLab Premium or require [signed commits](/autofix-and-remediation/configure/signed-commits.md), then a [Personal Access Token](https://gitlab.com/-/user_settings/personal_access_tokens) (PAT) should be used.

### Choose Authorize via Access Token in Aikido <a href="#go-to-pat-management-page-in-aikido" id="go-to-pat-management-page-in-aikido"></a>

In Aikido, go to the [**AutoFix settings** page](https://app.aikido.dev/issues/fix/settings) and click **Authorize** (on initial setup) or **Manage Personal Access Token** (when a token is already set). Choose **Authorize Via Access Token**.

**Next**, follow the instructions below to generate a Group Access Token in Gitlab.

## Generate a "Group Access Token" in Gitlab <a href="#generate-a-group-access-token-in-gitlab" id="generate-a-group-access-token-in-gitlab"></a>

1. Navigate to '**Settings**' > '**Access tokens**', via the navigation of your Group's page
2. You should end up on a page similar to this. You can navigate to this page yourself by opening your group's settings menu and selecting "**Access tokens**"

   
3. Click on "**Add** **new token**"
4. Enter a name for token, we suggest something like: "**Aikido Autofix**"
5. For the expiration, use custom defined and set it to 1 year from now, or remove it.
6. The token should have at least the role of "**Developer**"
7. Next you need to select the following permissions for the token: **api** and **write\_repository**

   
8. Click on "**Create group access token**"
9. Copy the token on the next screen and insert it on the access token management page in Aikido
