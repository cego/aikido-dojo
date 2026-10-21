# GitLab Integration: Authentication and User Management

### Introduction <a href="#introduction" id="introduction"></a>

When you create an Aikido workspace based on GitLab Cloud version, you select a GitLab Group. The Aikido workspace will consist of the main group and all its subgroups and underlying repositories.

### Benefits <a href="#benefits" id="benefits"></a>

* **Automatic Onboarding:** From that point on, other GitLab members of this group will be able to onboard Aikido automatically. Note: this only on the condition that they have at the least 'reporter', 'developer', 'maintainer' or 'owner' level of access to this group. 'Guests' will not be able to get in.
* **Effortless Offboarding:** When a member is removed from your GitLab Instance, they automatically lose access to the Aikido workspace. Aikido will sync access every night and auto-offboard members that are no longer active in your GitLab instance.

### Adding multiple GitLab Groups in a single Aikido Workspace <a href="#adding-multiple-gitlab-groups-in-a-single-aikido-workspace" id="adding-multiple-gitlab-groups-in-a-single-aikido-workspace"></a>

By default, a new workspace will be created for every GitLab group. If you want to have multiple GitLab groups inside 1 Aikido workspace, you need to create a root-level group in GitLab that contains all the subgroups.

### Inviting users to Aikido without an organisational GitLab account <a href="#inviting-users-to-aikido-without-an-organisational-gitlab-account" id="inviting-users-to-aikido-without-an-organisational-gitlab-account"></a>

You can invite users via Google, Microsoft or using a personal GitLab Account. More information can be [found here.](/getting-started/automated-user-management/invite-users-to-aikido-without-a-git-account.md)
