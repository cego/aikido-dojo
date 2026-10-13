# Assign Team Responsibilities with Gitlab Topics

Aikido allows you to use **Gitlab Project Topics** to automatically manage team assignments in Aikido. This streamlines issue management in large organizations and enables more granular reporting without manual setup.

### Use Cases <a href="#use-cases" id="use-cases"></a>

* **Clear ownership across many repos**: In large orgs, each repo has a defined owning team, making it easier to track who is responsible.
* **Automatic onboarding**: When new repos are created with the right topic, they’re instantly assigned to the correct team in Aikido.

### Configuration <a href="#assigning-specific-paths" id="assigning-specific-paths"></a>

{% hint style="warning" %}
This feature is **not enabled by default** and needs to be activated for your account. Please contact **Aikido support** for assistance.
{% endhint %}

For each repository you want managed by Aikido:

1. In GitLab, [go to your repository’s Settings → General → Topics](https://docs.gitlab.com/user/project/project_topics/).
2. Add a topic like: `owner:Internal Tooling Team`.
3. Aikido will automatically create (or update) the team and assign the repo.
4. If the topic is removed or changed later, Aikido will unassign or reassign the repo automatically.
