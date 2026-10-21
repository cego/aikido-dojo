# Pull request templates

AutoFix uses a pull request template from your repository to generate each pull request description. If no supported template is found, AutoFix uses its default template.

AutoFix checks these template locations:

* `.github/PULL_REQUEST_TEMPLATE.md`
* `.github/pull_request_template.md`
* `PULL_REQUEST_TEMPLATE.md`
* `docs/PULL_REQUEST_TEMPLATE.md`
* `.gitlab/merge_request_templates/Default.md`
* `pull_request_template.md`

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
