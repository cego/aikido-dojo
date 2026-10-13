# GitLab MR Gating

Aikido GitLab MR gating runs security checks on merge requests (MRs). It focuses on **newly introduced** issues in the MR and shows which existing issues the MR solves.

Use it to catch risky changes before they land on your main branch. You decide what should fail the gate.

### How it works

* You authorize Aikido to access your GitLab group.
* You pick repositories and scan types in Aikido.
* On each MR update, Aikido scans the MR diff.
* Aikido publishes results back to GitLab as pipeline checks and comments.
* You can make failed checks block merging in GitLab.

{% hint style="info" %}
If you want a broader view of the MR gating funtionality, see [PR Gating Overview](/pr-and-release-gating/aikido-ci-gating-functionality.md).
{% endhint %}

### Pick the right guide

Start with the page that matches your GitLab plan. Then configure repositories in the dashboard.

* **GitLab Premium, Ultimate, or GitLab Server (self-managed)**
  * Use a **Service Account** plus an Access Token.
  * Follow: [GitLab Premium, Ultimate & Server: MR Scans Setup](/pr-and-release-gating/gitlab-mr-gating/gitlab-server-ci-mr-gating-via-aikido-dashboard-with-a-service-account-token.md).
* **GitLab Free**
  * Use a dedicated GitLab user plus a Personal Access Token (PAT).
  * Follow: [GitLab Free: MR Scans Setup](/pr-and-release-gating/gitlab-mr-gating/gitlab-free-mr-scans-setup.md).

### What to expect after setup

* Aikido runs MR scans when commits are pushed to the MR branch.
* Results show up in the MR under **Pipelines**.
* If you enable **Pipelines must succeed**, a failing scan can block the merge.

{% hint style="warning" %}
Start with one repository first. Confirm checks run and comments appear as desired. Then roll out to more repos.
{% endhint %}
