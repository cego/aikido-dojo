# Required permissions for Aikido AutoFix Github App

On Monday 24th February 2025 we added additional permissions to the Aikido Autofix GitHub App. The following permissions were added:

* Code: read & write
* PRs: read & write
* (new) Workflows: read & write
* (new) Checks: read
* (new) Actions: read

All new installations of the Autofix GitHub app include these permissions. Users that installed the app before February 18th, received an email asking to confirm these new permissions. Until the user confirms the new permissions can be granted, the app will remain on the initial permissions.

The new permissions will allow Aikido to create more and better autofixes. With the added workflow permissions, Aikido is able to create PRs that fix Github actions, for example to pin 3rd party actions you are using, mitigating potential supply chain attacks.

The additional read permissions for Checks and Actions enable Aikido to retrieve the status and logs for GitHub Actions. If an Autofix is created that breaks the build, the output from the build will be used to fix the PR and create a working build.
