# Assign Team Responsibilities with Code Owners

Aikido allows you to use **Code Owners** to automatically manage team assignments in Aikido. This streamlines issue management in large organizations and enables more granular reporting without manual setup.

{% hint style="info" %}
Code Owners is independent from source code manager. It can easily be used with GitHub, Gitlab, Bitbucket and even our Local Scanning workspaces.
{% endhint %}

### Use Cases <a href="#use-cases" id="use-cases"></a>

* **Clear ownership across many repos**: In large orgs, each repo has a defined owning team, making it easier to track who is responsible.
* **Path-level ownership:** For codebases where multiple teams own different directories or paths (e.g. monorepos), Aikido maps CODEOWNERS paths directly to teams, no manual assignment needed.
* **Automatic onboarding**: When new repos are created with the right topic, they’re instantly assigned to the correct team in Aikido.

### Configuration <a href="#assigning-specific-paths" id="assigning-specific-paths"></a>

{% hint style="warning" %}
This feature is **not enabled by default** and needs to be activated for your account. Please contact **Aikido support** for assistance.
{% endhint %}

1. In Github/Gitlab/Bitbucket, create a new file called `CODEOWNERS.` See [GitHub docs](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-code-owners), [Gitlab docs](https://docs.gitlab.com/user/project/codeowners/) and [Bitbucket docs](https://support.atlassian.com/bitbucket-cloud/docs/set-up-and-use-code-owners/). Aikido assumes there is only 1 codeowner file and will use the first one it finds.
2. Aikido will automatically create (or update) the team and assign the repo.
3. If the Codeowners file is changed later, Aikido will unassign or reassign the repo automatically.
