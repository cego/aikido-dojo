# Main Feed

Aikido's Feed is meticulously designed to streamline security management: it groups similar issues for clarity, allows for customizable views to highlight what's currently important. The main feed focusses on open issues only. Snoozed and ignored issues are accessible on a separate page for a clean giving a focused and clean overview.

### Aikido's Main Feed <a href="#aikidos-main-feed" id="aikidos-main-feed"></a>

The **Feed** is the central hub for monitoring and managing security issues. A similar view can be found all over Aikido when going into details of repositories, clouds and domains.

By default, Aikido enables the **Focus** view, containing all issues that are important to follow-up. By toggling to **All** we will also show issues that have been auto-triaged and ignored.

*Status Column*

Shows the status of the particular issue.

* *New*: issues that have been around for less than 7 days
* *To Do:* all other issues that are still unresolved
* *Task Open:* only when a task manager is linked
* *PR Open:* clicking will take you to PR / branch scan overview.<br>

  

### Sidebar <a href="#sidebar" id="sidebar"></a>

Aikido's Feed is equipped with a detailed sidebar, designed to provide users with comprehensive information and actionable options for each security issue. The image below shows the main important actions one can take.

**#1 Group Actions:** these are actions that can be taken on group level of an issue.

**#2 Subissue Actions:** you can also take actions on a subissue separately. This can be important when certain subissues want to be snoozed or ignored, but the overall issue group should remain in the main feed.

**#3 Reachability Analysis:** visualizes accessible code paths for that specific issue.

**#4 Detail Screen**: Separate detail screen with even more information (e.g., CVE information).

[image: Security dashboard showing a critical vulnerability in the 'bson' package with actionable steps.]

**Reachability Analysis**

### Filter Issues <a href="#filter-issues" id="filter-issues"></a>

You can easily adjust the view in order to filter on those issues that are most important to you in the moment. These filters can be found above the issue table in every feed view. You can filter on

* **Issue Type**
* **Predefined filtered views**
  * If you are using the [**SLA functionality**](https://help.aikido.dev/en/articles/8926339-enabling-slas-in-aikido), you will be able to see all upcoming and out of SLA issues

### Export Issues to CSV or PDF <a href="#export-issues-to-csv-or-pdf" id="export-issues-to-csv-or-pdf"></a>

You can export a CSV or PDF of issues. You can configure which issues to export exactly. This can be triggered from the Actions menu in the top right.

[image: Export issues: Filter by status, type, team, and format for CSV or PDF download.]

### Show issues per team or project <a href="#show-issues-per-team-or-project" id="show-issues-per-team-or-project"></a>

**Aikido's Feed** features a **team filter** at the top of the page. This filter allows users to tailor the feed to display only the issues relevant to selected teams ([📖 Set Up Teams](https://help.aikido.dev/en/articles/9005606-using-teams-for-repository-and-user-management)).

[image: Team selector dropdown with search and statistics on solved and new issues.]
