# Forge Integration: Troubleshooting

**Table of contents:**

* [1. Missing permissions](#1-missing-permissions)
* [2. Repository not active in Aikido](#2-repository-not-active-in-aikido)
* [3. Setup for multiple Gits](#3-setup-for-multiple-gits)

## Aikido - Forge Integration: Troubleshooting

#### 1. Missing permissions <a href="#id-1-missing-permissions" id="id-1-missing-permissions"></a>

> A person has a user on aikido, but have not given permission to the current repository in the workspace

All permissions for Aikido can be managed inside your source code manager. This often can be found at App Permissions. A quick way to **grant permissions and activating repos** is through the Aikido UI

**Steps to give permissions and activate repos**

* Go to the [Repository Settings page ](https://app.aikido.dev/settings/integrations/repositories)and click **Add Repo**. This will open up the respective source code manager
* Select repos you want to Aikido to have access to

  [image: Repository access settings: choose all or specific repositories for permissions.]
* Select which repos Aikido should **activate** and scan

  > This is important so Aikido can scan for vulnerabilities and send information to Forge.

  [image: Aikido interface for selecting repositories to scan for vulnerabilities.]

#### 2. Repository not active in Aikido <a href="#id-2-repository-not-active-in-aikido" id="id-2-repository-not-active-in-aikido"></a>

> The repository has permission, but was not selected as one of the repos to scan.

It is possible that you have granted permission to Aikido to all repositories, but that only a few were selected for being nightly scanned. By default, Aikido disables these repositories for scanning.

**Steps to activate a single repository in Aikido**

* Go to the [Repository Settings page](https://app.aikido.dev/settings/integrations/repositories)
* Go to the repository that you want to activate and click the triple dots on the right. In the action menu, select 'Activate'.

  [image: Dropdown menu showing "Configure" and "Activate" actions for a Javaspringvulny entry.]

If you need to enable a big number of repos at once, you can click 'add repo' which will take you back through the [initial setup](https://app.aikido.dev/onboarding/github/connect-repositories) page. There you can multiselect which repos you would like to be scanned / activated.

[image: Select repositories to scan for vulnerabilities using Aikido's platform.]

#### 3. Setup for multiple Gits <a href="#id-3-setup-for-multiple-gits" id="id-3-setup-for-multiple-gits"></a>

> A user has a workspace, but the repository is not in that workspace. This can be the case when a repository is another GitHub organisation.

Aikido supports [connecting multiple Gits](/getting-started/setting-up-your-account/account-setup-with-multiple-gits.md) (e.g., GitHub combined with Bitbucket) or different organisations within the same Git. If you want to connect a new GitHub organisation to an existing account, **you will need to create a new workspace.**

**Steps to setup a new workspace**

**Step 1:** Add a new workspace by clicking your profile icon in the top right corner.

[image: User profile dropdown with option to add another workspace highlighted.]

**Step 2:** Continue by selecting "Add a New Workspace."

[image: Aikido dashboard: create new workspace or logout option available.]

**Step 3:** Switch between different Git accounts using the workspace switcher located in the top left corner.

[image: Dashboard showing GitHub organization selection and security vulnerability overview by severity.]
