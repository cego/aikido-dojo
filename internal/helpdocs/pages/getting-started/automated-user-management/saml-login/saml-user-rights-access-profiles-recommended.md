# SAML User Rights: Access Profiles (Recommended)

SAML Access Profiles allow you to define user access rights based on SAML attributes. You can manage these profiles in the app under:

**Settings > General > SAML Setup > Add SAML Profile**<https://app.aikido.dev/settings/account>

## Configuring SAML Access Profiles <a href="#configuring-saml-access-profiles" id="configuring-saml-access-profiles"></a>

When adding a new SAML Profile, you can define the following settings:

### 1. Profile Name <a href="#id-1-profile-name" id="id-1-profile-name"></a>

* The name that should be passed as the `aikido_access_profile` SAML claim.

### 2. Role <a href="#id-2-role" id="id-2-role"></a>

* Defines the user's role:
  * **Admin**
  * **Default**
  * **Team Only**

### 3. Edit Rights <a href="#id-3-edit-rights" id="id-3-edit-rights"></a>

* Determines the user's edit capabilities:
  * **Standard**
  * **Read Only**

### 4. Can Ignore <a href="#id-4-can-ignore" id="id-4-can-ignore"></a>

* Specifies whether the user can ignore issues:
  * **Yes**
  * **No**

### 5. Can Snooze <a href="#id-5-can-snooze" id="id-5-can-snooze"></a>

* Specifies whether the user can snooze issues:
  * **Yes**
  * **No**

### 6. Can Change Severity <a href="#id-6-can-change-severity" id="id-6-can-change-severity"></a>

* Defines if the user can change the severity of issues:
  * **Yes**
  * **No**

### 7. Can Export Data <a href="#id-7-can-manage-teams" id="id-7-can-manage-teams"></a>

* Defines if the user can export data:
  * **Yes**
  * **No**

### 8. Can Manage Teams <a href="#id-7-can-manage-teams" id="id-7-can-manage-teams"></a>

* Defines if the user can manage teams:
  * **Yes**
  * **No**

### 9. Can Manage Clouds <a href="#id-7-can-manage-teams" id="id-7-can-manage-teams"></a>

* Defines if the user can manage clouds:
  * **Yes**
  * **No**

### 9. Can Manage Containers <a href="#id-7-can-manage-teams" id="id-7-can-manage-teams"></a>

* Defines if the user can manage containers:
  * **Yes**
  * **No**

### 10. Can Manage Domains <a href="#id-7-can-manage-teams" id="id-7-can-manage-teams"></a>

* Defines if the user can manage domains:
  * **Yes**
  * **No**

### 11. Can Manage Pentests <a href="#id-7-can-manage-teams" id="id-7-can-manage-teams"></a>

* Defines if the user can manage pentests:
  * **Yes**
  * **No**

### 12. Can Manage Code Quality Rules <a href="#id-7-can-manage-teams" id="id-7-can-manage-teams"></a>

* Defines if the user can manage code quality rules:
  * **Yes**
  * **No**

### 13. Member of Teams <a href="#id-8-member-of-teams" id="id-8-member-of-teams"></a>

* A comma-separated list of team names the user belongs to.
* Matches the existing `aikido_teams` [SAML claim](/getting-started/automated-user-management/saml-login/saml-user-rights-using-custom-attributes-advanced.md).

### 14. Workspace IDs <a href="#id-9-workspace-ids" id="id-9-workspace-ids"></a>

* A comma-separated list of workspace IDs where the user has access.
* Matches the existing `aikido_workspace_ids` [SAML claim](/getting-started/automated-user-management/saml-login/saml-user-rights-using-custom-attributes-advanced.md).
* If left empty, the profile grants access to all workspaces linked to the SAML client.

## Using SAML Access Profiles <a href="#using-saml-access-profiles" id="using-saml-access-profiles"></a>

Once a profile is created, you can set up a custom SAML claim `aikido_access_profile` with the profile name as value. **If set**, users who authenticate via SAML will receive access based on the profile associated with this claim. Ensure that the correct claims are configured in your Identity Provider (IdP) to match the assigned profiles.

> **Note**
>
> When using the `aikido_access_profile` in combination with other [custom SAML claims](/getting-started/automated-user-management/saml-login/saml-user-rights-using-custom-attributes-advanced.md), those other claims will take precedence.
