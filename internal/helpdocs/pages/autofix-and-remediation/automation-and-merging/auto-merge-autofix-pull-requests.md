# Auto-Merge Autofix Pull Requests

## When to use auto-merge?

The auto-merge feature streamlines the development workflow by automatically merging Pull Requests once all defined conditions are met.

This can significantly improve workflow speed, by eliminating the need to manually track and merge pull requests. Developers no longer have to constantly check the conditions of each PR or periodically initiate merges.

Aikido relies on the built-in Github, Gitlab and Azure DevOps auto-merge features. In the Autofix configuration you can enable auto-merge on the PRs created Aikido.

## GitHub

Enable auto-merge for pull request in your repository: [Managing auto-merge for pull requests in your repository - Github Docs](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/configuring-pull-request-merges/managing-auto-merge-for-pull-requests-in-your-repository?utm_source=chatgpt.com).[\ <br>](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/configuring-pull-request-merges/managing-auto-merge-for-pull-requests-in-your-repository?utm_source=chatgpt.com)We strongly recommend you set up branch protection rules for the target branch: [Managing a branch protection rule - GitHub Docs](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches/managing-a-branch-protection-rule). This allows you to define the conditions that need to be met before the auto-merge is triggered.

## GitLab (Cloud & Server)

Instruction on how to enable auto-merge on your GitLab repo can be found here: [Auto-merge | GitLab Docs](https://docs.gitlab.com/user/project/merge_requests/auto_merge/)

Gitlab requires a branch policy to be setup in order for Aikido to create PRs with auto-merge enabled: [Protected branches | GitLab Docs](https://docs.gitlab.com/user/project/repository/branches/protected/)

Lastly, a pipeline has to be linked to the PR in order for Aikido to create PRs with auto-merge enabled: [CI/CD pipelines | GitLab Docs](https://docs.gitlab.com/ci/pipelines/)

Once the pipeline succeeds and the merge checks have been completed the PR will auto merge.

## Azure DevOps

We strongly recommend enabling branch policies for target branch to define when the PR should be auto-completed: [Git branch policies and settings - Azure Repos | Microsoft Learn](https://learn.microsoft.com/en-us/azure/devops/repos/git/branch-policies?view=azure-devops\&tabs=browser)

## Bitbucket

Bitbucket does not have an auto-merge functionality similar to the ones above. When using the auto-merge functionality provided by Bitbucket, the user must still manually click the merge button. [Merge a pull request | Bitbucket Data Center 10.0 | Atlassian Documentation](https://confluence.atlassian.com/bitbucketserver/merge-a-pull-request-808488562.html)

<br>
