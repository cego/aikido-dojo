# Automatic Creation of Pull Requests

Aikido allows for automatic Pull Request creation on a **daily basis**, handling mass dependency updates, fixing SAST and AI Code Analysis issues, updating container base images, and fixing Code Quality issues. Pentest issues are fixed as soon as they're reported, see [AutoFix for Pentest](https://help.aikido.dev/autofix-and-remediation/scope/autofix-for-pentest). You can fine tune the pull request title, description, labels and other settings with "PR Configuration".

### Setup automatic creation <a href="#setup-automatic-creation" id="setup-automatic-creation"></a>

**Step 1:** Go to [Autofix Settings](https://app.aikido.dev/issues/fix/settings?regionRedirect=1) page

**Step 2:** Scroll down to Automated PRs. Enable the toggle for 'Autocreation of PRs'

**Step 3:** The configuration modal opens. You can enable autocreation for dependencies in code repos, automated fixes for Docker containers, or automatically fix SAST/AI Code Analysis & IaC, pentest, and Code Quality issues. Make sure to select code repositories and containers for which you want PRs created.

**Step 4:** Test whether new PRs are created correctly by clicking the **Test Autocreation** button.

### Code Dependency Automatic PR creation limits <a href="#pull-request-configuration" id="pull-request-configuration"></a>

AutoFix creates a maximum of **10 PRs per repo** during the **daily** run.

### SAST/AI Code Analysis and IaC automatic PR creation limits

When you perform code scanning on a repository for the first time, you will typically discover a large amount of issues. To keep the number of pull requests manageable have the following limits when creating PRs for SAST and AI Code Analysis issues:

* PRs are only created for issues where AutoFix has **medium** or **high confidence**.
* During the **daily** run, AutoFix will create a **maximum** of **5 PRs per repo**. A SAST PR contains a maximum of **5 changed files**. An AI Code Analysis PR fixes up to **5 findings** in the same file.
* AutoFix limits the amount of **open SAST and AI Code Analysis AutoFix PRs** to **10 per repo**. Once that limit is hit, AutoFix doesn't automatically create new PRs. When you close or merge the PRs, you will receive new PRs during the next daily run.

SAST and AI Code Analysis PRs share these limits.

If you need higher limits, reach out to us in the Aikido UI.

### Pull Request Configuration <a href="#pull-request-configuration" id="pull-request-configuration"></a>

You can configure pull request properties (such as the title, labels, branch name, and commit message) on the AutoFix PR Configuration page. More docs can be found [here](https://help.aikido.dev/autofix-and-remediation/configure/autofix-pr-configuration).
