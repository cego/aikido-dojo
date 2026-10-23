# Pull request templates

AutoFix uses a pull request template from your repository to generate each pull request description. If no supported template is found, AutoFix uses its default template.

AutoFix checks these template locations:

* `pull_request_template.md`
* `PULL_REQUEST_TEMPLATE.md`
* `pull_request_template.txt`
* `docs/pull_request_template.md`
* `docs/PULL_REQUEST_TEMPLATE.md`
* `docs/pull_request_template.txt`
* `.github/PULL_REQUEST_TEMPLATE.md`
* `.github/pull_request_template.md`
* `.gitlab/merge_request_templates/Default.md`
* `.azuredevops/pull_request_template.md`
* `.azuredevops/pull_request_template.txt`
* `.vsts/pull_request_template.md`
* `.vsts/pull_request_template.txt`

Example pull request template:

{% code title="PULL\_REQUEST\_TEMPLATE.md" %}

```
# Summary

<!-- What does this PR do? -->

# Security Impact

<!-- Describe the security vulnerability being fixed and why it matters -->

# Fix Details

<!-- Technical details of the fix -->

# Warnings

<!-- Warnings about the remediation -->

# Related Tasks

<!-- Links to related tickets or issues -->

# Checklist

- [ ] Reviewed the changes
- [ ] Verified the fix does not introduce regressions
- [ ] Approved for merge
```

{% endcode %}
