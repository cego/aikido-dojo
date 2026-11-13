# Assigning Resources to Teams

Assigning resources to teams defines ownership inside Aikido. It ensures that findings are routed to the right people, dashboards stay relevant, and teams only see what they are responsible for.

This page explains all available ways to assign resources, including individual and bulk options.

### Assign from the Team Page

1. Go to [Settings → Teams](https://app.aikido.dev/settings/teams?regionRedirect=1)
2. Select a team
3.  Open the Responsibility tab

    
4. Click "Link Resource"
5.  Select resource type, and search and add resources<br>

    

### Assign from the Resource Page

You can also assign a team while viewing the resource itself.

#### Repositories

Open a Repository → Configure → Teams responsible

#### Cloud Accounts

Open a Cloud Connection → Configure → Teams responsible

### Domains & API's

Open the Front-end, REST, or GraphQL scan configuration from the Settings page (or the action menu in the list view) and link the domain to an asset.

You can link it to either a repository or a container. Once linked, the scan inherits the team permissions from that asset.

### Bulk Assignment

Bulk actions are available on resource list pages in [settings](https://app.aikido.dev/settings/account?regionRedirect=1).

{% hint style="info" %}
If bulk actions are not visible, assignment must be done individually.
{% endhint %}

**Repositories**

[Settings → Repositories](https://app.aikido.dev/settings/integrations/repositories?regionRedirect=1) → Select multiple → Bulk Actions → Assign to team

**Containers**

[Settings → Containers](https://app.aikido.dev/settings/container-image-registry?regionRedirect=1) → Select multiple → Bulk Actions → Assign to team

**Cloud**

[Settings -> Clouds](https://app.aikido.dev/settings/integrations/clouds?regionRedirect=1) -> Select multiple -> Bulk Actions -> Assign to team
