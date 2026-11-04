# GitLab Pipeline Setup for Aikido CLI: Release Gating

The Aikido Security CI client allows you to integrate Aikido Security scans into CI pipelines. It helps ensure that security scans are part of your build process.

### 1. Get API key <a href="#id-1-get-api-token" id="id-1-get-api-token"></a>

1. Go to the [Continuous Integration Settings page](https://app.aikido.dev/settings/integrations/continuous-integration).
2. Generate a token and copy. Note that you will only be able to view this token once.
3. Save this token as *AIKIDO\_CLIENT\_API\_KEY* in your GitLab CI/CD variables as a secret by going to Settings > CI/CD and Variables > Actions. This can either be done on group or project level in GitLab

> When adding the API key, make sure the variable is available on all branches (uncheck protect variable) and it is masked in any logs (check the box to "Mask variable"). You can leave "Expand variable reference" checked.

### 2. Create a new GitLab Pipeline job

Create a pipeline job using this this example:

> Make sure that the local scanner is only triggered for your default branch. In the example below this is the 'main' branch. Adjust this if needed.

{% tabs %}
{% tab title="Yaml" %}

```yaml
stages:
  - security_check

aikido_release_gate:
  stage: security_check
  image: node:22
  script:
    - npm install --global @aikidosec/ci-api-client
    - |
      aikido-api-client scan-release \
        "$CI_PROJECT_ID" \
        "$CI_COMMIT_SHA" \
        --apikey "$AIKIDO_CLIENT_API_KEY" \
        --fail-on-sast-scan \
        --fail-on-iac-scan \
        --fail-on-secrets-scan
  only:
    - main
```

{% endtab %}

{% tab title="" %}

{% endtab %}
{% endtabs %}

Tweak the command if needed. Every flag is documented on the [`@aikidosec/ci-api-client`](https://www.npmjs.com/package/@aikidosec/ci-api-client) package.

### 3. Run your first scan <a href="#id-3-run-your-first-scan" id="id-3-run-your-first-scan"></a>

When a new change is pushed to your default branch, the workflow will be triggered. A scan will run and fail if new issues have been detected.
