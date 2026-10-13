# Why does AutoFix close PRs?

We believe that having a pile of open Pull Request with overlapping dependency fixes is not helping you fix your dependency issues. Therefor AutoFix closes Pull Requests to keep the amount of open AutoFix Pull Requests manageable.

Pull requests are closed for several reasons:

1. The issue has been **fixed** in a different way
2. A **new AutoFix Pull Request** with a superset of fixes has been created
3. The pull request was **conflicted** and AutoFix recreated the Pull Request

### Issues already fixed

When an issue is closed in Aikido, we will check for open Pull Requests for the issue. When an open Pull Request for the issue is found, AutoFix will close the open Pull Request.

This can happen when a manual fix is performed while an Aikido Pull Request is open.

### New AutoFix Pull Requests

When AutoFix creates a Pull Request (through the UI or automated job), it checks for open pull requests that contain a subset of the changes in the newly created pull request. When AutoFix finds a pull request that contains a subset of fixes, it closed that pull request in favour of the new Pull Request.

As an example, if PR 1 fixes issue A en PR 2 fixes issue A and B, PR1 will be closed.

### Conflicted Pull Request

When the automated AutoFix job detects a pull request is open but it contains merge conflicts, it closes the Pull Request to create a new pull request to mitigated the conflict.
