# Manage and View Your Apps and Projects via Our Teams Feature

## Introduction <a href="#introduction" id="introduction"></a>

Aikido lets you create teams, connect multiple repositories and clouds to have a clear overview of your apps and projects. This project/app view is available throughout the entire Aikido app, going from feed, alerting and reporting. This article focuses on using Teams in order to group resources into Project and Apps. Managing of User Access. If you are looking to manage User Access with teams, please [**click here**](/getting-started/manage-teams-and-applications/managing-user-access-with-teams.md).

## Use Cases <a href="#use-cases" id="use-cases"></a>

* **Agency Project Management:** Agencies frequently handle multiple clients, necessitating a structured approach to manage each client's repositories separately. Aikido facilitates this by allowing the creation of teams for each client, making it straightforward to organize and access client-specific repositories. Additionally, this allows for easy generation of client-specific reports.

  [image: Dropdown menu for selecting "All teams" or individual client teams.]
* **Specify resources connect to an app:** you can use the team functionality to combine resources that are part of the same app (combination of repos, containers and clouds). This also allows you to view all reports
* **Better overview when using a monorepo split.** You can assign team members to specific directories of your monorepo, improving their overview. More information on splitting monorepositories can be found [here](https://help.aikido.dev/en/articles/9026666-splitting-up-your-monorepo-per-directory).

{% hint style="info" %}
Repo linking is only available for manually created teams (not imported from SCM). Other resources like containers, clouds & Zen apps can be linked to both manually created and SCM-imported teams.
{% endhint %}

## How To Create Teams <a href="#how-to-create-teams" id="how-to-create-teams"></a>

**Step 1:** Navigate to **Settings -> Teams**\
​

[image: Team dashboard for “mdschuym” showing 4 active repositories and management options.]

**Step 2:** Click **Create Team** and give your team a name\
​

[image: Form to create a new team by entering a team name for vulnerability tracking.]

**Step 3 (optional):** **Add team members** to the newly created team.

[image: User management interface for assigning people to the "Architects" team.]

**Step 4:** **Link different resources** via the **responsibility tab.** You can add different resources such as repositories, clouds, containers, domains and zen apps. You can add repositories and containers **in bulk** via the [repository](https://app.aikido.dev/settings/integrations/repositories) and [container settings](https://app.aikido.dev/settings/container-image-registry) screen.

> If you want to link specific domains to a team, you can set this up by linking your domains to a repo or container. It will automatically inherit access permissions.

[image: Interface for linking repositories, cloud, or container resources to the Backend team.]

**Step 5.** Go back to your feed and filter on a specific team. You should only see issues that are related to those repositories and clouds that were attached to the team.

## How to select your team in UI <a href="#how-to-select-your-team-in-ui" id="how-to-select-your-team-in-ui"></a>

**Aikido's Feed** features a **team filter** at the top of the page. This filter allows users to tailor the feed to display only the issues relevant to selected teams. This filter can be used on basically every page in Aikido (feed, reports, settings etc).

[image: Team filter dropdown with stats for solved and newly detected issues this week.]

## Parent teams

A parent team groups existing teams and takes on their resources. Select the parent team in the team filter to see issues for every child team.

{% hint style="info" %}
Parent teams are available on request. [Contact Support](/miscellaneous-info/contact-support.md) to enable them for your workspace.
{% endhint %}

Create a parent team with the [Create parent team](https://apidocs.aikido.dev/reference/createparentteam) endpoint. Pass the parent name and the child team IDs. Replace the child teams later with [Update child teams for parent team](https://apidocs.aikido.dev/reference/updatechildteamsforparentteam). Both calls need the `teams:write` scope.

Aikido copies repositories, clouds, containers, domains, and Zen apps from the child teams onto the parent. Link and unlink those resources on the child teams, not on the parent. Members of the parent team stay separate from the child teams.

The team filter nests child teams under each parent team. Aikido refreshes the parent's resources when you open it under **Settings** → **Teams**, and when you select it in the team filter. Selecting it refreshes the copy at most once an hour.
