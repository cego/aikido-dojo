# Configure Team Sync Settings

Team sync controls how Aikido imports teams and members from your source control manager (GitHub, Bitbucket, or Azure DevOps), and whether repositories are automatically linked to those teams.

{% hint style="info" %}
This is an admin-only setting. Only workspace **admins** can change the team sync mode.
{% endhint %}

### **Different Sync Settings**

#### **Full team sync (default)**

Teams and members are imported from your SCM and kept up to date. Repositories are automatically linked to teams based on your SCM's team-to-repo mapping, so ownership and coverage mirror your SCM structure exactly.

**When to use:** Use this when you want your SCM to drive both group membership and repository ownership. This is the lowest-maintenance option.

#### **Team sync without automatic repo linking**

Teams and members are still imported and kept in sync from your SCM. Repositories are not linked automatically. You manage repo-team assignments directly in Aikido.

**When to use:** Use this when you want accurate team membership in Aikido but need to control repo-team assignments manually, for example when your SCM team structure does not reflect how you assign code ownership.

#### **Team sync off**

No teams are imported from your SCM going forward. Existing SCM-imported teams will be removed when you turn this off.

**When to use**: use this when you are not using SCM teams in Aikido at all, or when you are migrating to a different team management approach.

{% hint style="warning" %}
Turning off team sync removes teams that were created and maintained by the import process. Teams created manually in Aikido are not affected.
{% endhint %}

#### How to configure sync settings

**Step 1.** Go to [**settings > Teams**](https://app.test.aikido.dev/settings/teams) and click **Settings** in the Action menu

**Step 2.** Choose the preferred sync configuration and press save.
