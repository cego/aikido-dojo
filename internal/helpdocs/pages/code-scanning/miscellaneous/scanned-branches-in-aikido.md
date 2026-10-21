# Scanned Branches in Aikido

### Single-Branch Scanning <a href="#single-branch-scanning" id="single-branch-scanning"></a>

By default, Aikido will scan the default branch in your repository for dependency and code issues. The default branch is usually main/master. **Aikido scans your default branch on every push via SCM webhooks, plus a full re-scan of every repo, cloud, container and domain runs on a daily schedule.**

If you update the default branch, Aikido will pick up this change during its nightly sync, so no changes are needed on Aikido's end.

In some cases, it might be interesting to scan another branch instead of the default branch. This could be the case if you are developing on another branch for an extended period of time. You can overwrite the branch being scanned on the repository detail page in Aikido by clicking the branch button at the top.

[image: GitHub admin-panel repository showing master branch and issue status.]

A modal will appear that will allow you to change the branch that is scanned on the scheduled sync.

[image: Dialog box to update and save the branch name for scanning.]

## Multi-Branch Scanning <a href="#multi-branch-scanning" id="multi-branch-scanning"></a>

If you are looking to scan multiple branches of the same repo on a regular basis, you can contact us to enable Multi-Branch Scanning for you. This is only available for Paid Plans. Please check whether multi-branch scanning is right for you [by checking the use cases in this article](/code-scanning/miscellaneous/support-for-multi-branch-scanning.md).
