# Local Code Scanning

Local code scanning lets you analyze your source code directly on your own machines without sending any code to Aikido. It is useful when you want to test changes during development, work in isolated environments, or keep sensitive repositories fully local.

{% hint style="success" %}
For most workflows we still recommend using the [standard Aikido integrations](https://help.aikido.dev/code-scanning/connect-your-source-code) as they give faster results and better coverage.
{% endhint %}

Before you start, make sure you have an Aikido account/workspace that supports local scanning. Follow [Account Creation for Local Scanning](https://help.aikido.dev/code-scanning/local-code-scanning/account-creation-for-local-scanning-on-aikido).

To send Aikido a zip of a repository instead, follow [Upload a Repository as a Zip](https://help.aikido.dev/code-scanning/connect-your-source-code/upload-a-repository-as-a-zip). That upload sends the archive to Aikido, and Aikido stores it. Local scanning keeps the code on your machines.

If you are looking to scan container images instead, check the [container image scanning docs](https://help.aikido.dev/container-image-scanning/local-image-scanning).

## General

{% content-ref url="account-creation-for-local-scanning-on-aikido.md" %}
[account-creation-for-local-scanning-on-aikido.md](https://help.aikido.dev/code-scanning/local-code-scanning/account-creation-for-local-scanning-on-aikido)
{% endcontent-ref %}

{% content-ref url="pr-gating-for-code-using-local-scanner.md" %}
[pr-gating-for-code-using-local-scanner.md](https://help.aikido.dev/code-scanning/local-code-scanning/pr-gating-for-code-using-local-scanner)
{% endcontent-ref %}

{% content-ref url="release-gating-for-code-using-local-scanner.md" %}
[release-gating-for-code-using-local-scanner.md](https://help.aikido.dev/code-scanning/local-code-scanning/release-gating-for-code-using-local-scanner)
{% endcontent-ref %}

{% content-ref url="performing-nightly-scans-using-the-aikido-local-scanner.md" %}
[performing-nightly-scans-using-the-aikido-local-scanner.md](https://help.aikido.dev/code-scanning/local-code-scanning/performing-nightly-scans-using-the-aikido-local-scanner)
{% endcontent-ref %}

{% content-ref url="cli-options-for-local-scanner.md" %}
[cli-options-for-local-scanner.md](https://help.aikido.dev/code-scanning/local-code-scanning/cli-options-for-local-scanner)
{% endcontent-ref %}

## CI/CD Integrations

{% content-ref url="azure-devops-server-setup-for-local-code-scanning.md" %}
[azure-devops-server-setup-for-local-code-scanning.md](https://help.aikido.dev/code-scanning/local-code-scanning/azure-devops-server-setup-for-local-code-scanning)
{% endcontent-ref %}

{% content-ref url="bitbucket-pipeline-setup-for-local-code-scanning.md" %}
[bitbucket-pipeline-setup-for-local-code-scanning.md](https://help.aikido.dev/code-scanning/local-code-scanning/bitbucket-pipeline-setup-for-local-code-scanning)
{% endcontent-ref %}

{% content-ref url="bamboo-setup-for-local-code-scanning.md" %}
[bamboo-setup-for-local-code-scanning.md](https://help.aikido.dev/code-scanning/local-code-scanning/bamboo-setup-for-local-code-scanning)
{% endcontent-ref %}

{% content-ref url="circleci-setup-for-local-code-scanning.md" %}
[circleci-setup-for-local-code-scanning.md](https://help.aikido.dev/code-scanning/local-code-scanning/circleci-setup-for-local-code-scanning)
{% endcontent-ref %}

{% content-ref url="github-action-setup-for-local-code-scanning.md" %}
[github-action-setup-for-local-code-scanning.md](https://help.aikido.dev/code-scanning/local-code-scanning/github-action-setup-for-local-code-scanning)
{% endcontent-ref %}

{% content-ref url="gitlab-self-managed-setup-for-local-code-scanning.md" %}
[gitlab-self-managed-setup-for-local-code-scanning.md](https://help.aikido.dev/code-scanning/local-code-scanning/gitlab-self-managed-setup-for-local-code-scanning)
{% endcontent-ref %}

{% content-ref url="jenkins-setup-for-local-code-scanning.md" %}
[jenkins-setup-for-local-code-scanning.md](https://help.aikido.dev/code-scanning/local-code-scanning/jenkins-setup-for-local-code-scanning)
{% endcontent-ref %}

{% content-ref url="teamcity-pipeline-setup-for-local-code-scanning.md" %}
[teamcity-pipeline-setup-for-local-code-scanning.md](https://help.aikido.dev/code-scanning/local-code-scanning/teamcity-pipeline-setup-for-local-code-scanning)
{% endcontent-ref %}

## Operating Systems

{% content-ref url="linux-setup-for-local-code-scanning.md" %}
[linux-setup-for-local-code-scanning.md](https://help.aikido.dev/code-scanning/local-code-scanning/linux-setup-for-local-code-scanning)
{% endcontent-ref %}

{% content-ref url="mac-setup-for-local-code-scanning.md" %}
[mac-setup-for-local-code-scanning.md](https://help.aikido.dev/code-scanning/local-code-scanning/mac-setup-for-local-code-scanning)
{% endcontent-ref %}

{% content-ref url="windows-setup-for-local-code-scanning.md" %}
[windows-setup-for-local-code-scanning.md](https://help.aikido.dev/code-scanning/local-code-scanning/windows-setup-for-local-code-scanning)
{% endcontent-ref %}

## Other

{% content-ref url="local-scanning-in-existing-scm-integrated-workspaces.md" %}
[local-scanning-in-existing-scm-integrated-workspaces.md](https://help.aikido.dev/code-scanning/local-code-scanning/local-scanning-in-existing-scm-integrated-workspaces)
{% endcontent-ref %}
