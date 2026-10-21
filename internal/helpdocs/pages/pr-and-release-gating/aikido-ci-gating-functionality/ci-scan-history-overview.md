# CI Scan History Overview

The CI Scan History feature in Aikido provides a detailed overview of all past CI scans, offering insights into their status, newly introduced issues, and resolved issues. This is important to keep an overview of all executed scans and allow you to quickly check the diff page from previous executed scans.

### How to check your CI Scan History <a href="#how-to-check-your-ci-scan-history" id="how-to-check-your-ci-scan-history"></a>

**Step 1.** Go to the CI section on the [Integrations page](https://app.aikido.dev/settings/integrations) and enter the [GitHub](https://app.aikido.dev/settings/integrations/github/checks)/[Azure](https://app.aikido.dev/settings/integrations/azure-devops/checks) CI dashboard **OR** - if you are running your CI via code - go to the [CI Settings Page](https://app.aikido.dev/settings/integrations/continuous-integration).

**Step 2.** Click on "CI Scan History" located at the top right of the page.

[image: GitHub CI dashboard showing repository count and a button for CI scan history.]

**Step 3.** Review the overview page to see all previous CI scans, including their status, newly introduced issues, and resolved issues.

[image: CI scan history table showing pull requests, issue counts, scan dates, and gate statuses.]

**Step 4.** Click one of the PRs to get a full overview of the executed scan.

[image: CI check failed: 2 secrets exposed, 1 vulnerability solved in code review.]

### Filter by scan date

You can filter the CI scan history to a specific date by adding a `scan_date` query parameter to the URL:

```
https://app.aikido.dev/settings/integrations/ci/checks/scans?scan_date=2026-05-10
```

Use the `YYYY-MM-DD` format. This is handy when you need to pull up scans from a particular day during an audit or compliance review, such as when you need to show what ran on a specific release date or during a past incident.
