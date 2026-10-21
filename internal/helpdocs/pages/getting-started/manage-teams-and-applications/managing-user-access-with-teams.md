# Managing User Access with Teams

## Introduction <a href="#introduction" id="introduction"></a>

Aikido lets you create teams, connect multiple repositories and clouds, and manage access using RBAC (Role-Based Access Control) for better security. This article focuses on Managing of User Access. If you are looking to use Teams to group your resources into Projects or Apps, please [**click here**](/getting-started/manage-teams-and-applications/manage-and-view-your-apps-and-projects-via-our-teams-feature.md).

## Use Cases <a href="#use-cases" id="use-cases"></a>

* **Selective Access Control:** Assign repository access exclusively to designated team members.
* **Filter Repositories Quickly:** In companies with multiple teams, each team may have access to the entire codebase, but their primary responsibilities are often limited to specific repositories or projects. By creating teams in Aikido, team members **can have a focused overview** of only the repositories relevant to their tasks.\
  ​

  
* **Better overview when using a monorepo split.** You can assign team members to specific directories of your monorepo, improving their overview. More information on splitting monorepositories can be found [here](https://help.aikido.dev/en/articles/9026666-splitting-up-your-monorepo-per-directory).

## How To Create Teams <a href="#how-to-create-teams" id="how-to-create-teams"></a>

**Step 1:** Navigate to **Settings -> Teams**\
​

[image: Team dashboard showing user "mdschuym" with 4 active repositories and 1 member.]

**Step 2:** Click **Create Team** and give your team a name

**Step 3:** **Add team members** to the newly created team.

[image: User management interface showing available users and team assignment options.]

**Step 4:** Define the **team's responsibilities** by adding **resources via the responsibility tab.** You can add different resources such as clouds, repositories, containers, domains and zen apps.

> If you want to link specific domains to a team, you can set this up by linking your domains to a repo or container. It will automatically inherit access permissions.

**Step 5.** Go back to your feed and filter on a specific team. You should only see issues that are related to those repositories and clouds that were attached to the team.

## Syncing with GitHub, Bitbucket or Azure DevOps <a href="#syncing-with-github-bitbucket-or-azure-devops" id="syncing-with-github-bitbucket-or-azure-devops"></a>

If you have **existing teams** set up in GitHub, Bitbucket or Azure DevOps, Aikido will import them and maintain synchronization on a nightly basis. This ensures that any changes in team structures or access rights managed in GitHub/Bitbucket are accurately reflected in Aikido. Any new teams that are created in GitHub will appear in Aikido. The same applies to when you remove a team in GitHub: Aikido will pick this up and remove the team too. Any repos that are part of the team, will be synced too.

{% hint style="info" %}
It's important to note that in this scenario, GitHub/Bitbucket/Azure DevOps acts as the source of truth for access rights, and all management should be conducted within those platforms. This also means that **no extra users can be added to scm-linked teams** inside Aikid&#x6F;**.**
{% endhint %}

Aikido makes it clear which teams have been imported from your SCM.

[image: Team roles and users overview with DevOps group imported from GitHub.]

## Syncing with Backstage.io <a href="#syncing-with-backstageio" id="syncing-with-backstageio"></a>

Aikido integrates seamlessly with repositories containing `catalog-info.yaml` files for [Backstage.io](https://backstage.io). This allows for the automatic importing of teams, taking into account the path of where the file is located.

{% hint style="info" %}
Please note: This feature needs to be enabled by Aikido manually. Please reach out to support to enable this.
{% endhint %}

#### How It Works <a href="#how-it-works" id="how-it-works"></a>

1. Aikido scans repositories for `catalog-info.yaml` files.
2. Aikido looks for the `spec->owner` field in the file and imports this as team.
3. Aikido records the exact path of each `catalog-info.yaml` file, ensuring the team is responsible for those specific paths (and repositories).

## Syncing with Port.io <a href="#syncing-with-portio" id="syncing-with-portio"></a>

Aikido integrates seamlessly with repositories containing `port.yaml` files for [Port.io](https://port.io). This allows for the automatic importing of teams, taking into account the path of where the file is located.

#### How It Works <a href="#how-it-works" id="how-it-works"></a>

1. Aikido scans repositories for `port.yaml` files.
2. Aikido looks for the `team` field in the file and imports this as team.
3. Aikido records the exact path of each `port.yaml` file, ensuring the team is responsible for those specific paths (and repositories).

## How to select your team in UI <a href="#how-to-select-your-team-in-ui" id="how-to-select-your-team-in-ui"></a>

**Aikido's Feed** features a **team filter** at the top of the page. This filter allows users to tailor the feed to display only the issues relevant to selected teams. This filter can be used on basically every page in Aikido (feed, reports, settings etc).
