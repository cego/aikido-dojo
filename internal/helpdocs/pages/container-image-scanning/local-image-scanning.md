# Local Image Scanning

Local image scanning lets you analyze container images on your own machines without sending any code or images to Aikido. It is useful when you want to test builds during development, debug issues in isolated environments, or keep sensitive images fully local.

{% hint style="success" %}
For day to day workflows we recommend using a [cloud provider integration](/container-image-scanning/cloud-provider-registries.md) or the [standalone registries](/container-image-scanning/standalone-registries.md) since both give faster results and better coverage.
{% endhint %}

If you are looking for [source code scanning instead, check the code scanning docs](/code-scanning/local-code-scanning.md).

## General

{% content-ref url="/pages/UGmuTdNbSFWVbPodTQJA" %}
[Setting Up Image Scanning with Local Scanner](/container-image-scanning/local-image-scanning/setting-up-image-scanning-with-local-scanner.md)
{% endcontent-ref %}

{% content-ref url="/pages/Kxt0ghbRYyq1QVVqXTMb" %}
[PR Gating For Container Images Using Local Image Scanner](/container-image-scanning/local-image-scanning/pr-gating-for-container-images-using-local-image-scanner.md)
{% endcontent-ref %}

{% content-ref url="/pages/yUKARTnpQ5Xx9ZJOf00K" %}
[Release Gating For Container Images Using Local Image Scanner](/container-image-scanning/local-image-scanning/release-gating-for-container-images-using-local-image-scanner.md)
{% endcontent-ref %}

## CI/CD Integrations

{% content-ref url="/pages/JMJbsb8qfBp66ZTK7nIe" %}
[Azure DevOps Server Setup for Local Image Scanning](/container-image-scanning/local-image-scanning/azure-devops-server-setup-for-local-image-scanning.md)
{% endcontent-ref %}

{% content-ref url="/pages/ZC4RnYqy66X7mRQDjYHv" %}
[BitBucket Pipeline Setup for Local Image Scanning](/container-image-scanning/local-image-scanning/bitbucket-pipeline-setup-for-local-image-scanning.md)
{% endcontent-ref %}

{% content-ref url="/pages/Iz7HOx1fYMVH5JDlTMQ1" %}
[CircleCI Setup for Local Image Scanning](/container-image-scanning/local-image-scanning/circleci-setup-for-local-image-scanning.md)
{% endcontent-ref %}

{% content-ref url="/pages/Uu61wco68fGEBw4AUJI8" %}
[GitHub Action Setup for Local Image Scanning](/container-image-scanning/local-image-scanning/github-action-setup-for-local-image-scanning.md)
{% endcontent-ref %}

{% content-ref url="/pages/ZOZQ5MTytbyQvbBTBaVm" %}
[GitLab Setup for Local Image Scanning](/container-image-scanning/local-image-scanning/gitlab-setup-for-local-image-scanning.md)
{% endcontent-ref %}

{% content-ref url="/pages/K4TXf41zFabMa6z7EEoN" %}
[Jenkins Setup for Local Image Scanning](/container-image-scanning/local-image-scanning/jenkins-setup-for-local-image-scanning.md)
{% endcontent-ref %}

## Other

{% content-ref url="/pages/YPhPkWFWOMtz6C53n8Fo" %}
[Kaniko Image Scanning with Local Scanner](/container-image-scanning/local-image-scanning/setting-up-kaniko-image-scanning-with-local-scanner.md)
{% endcontent-ref %}
