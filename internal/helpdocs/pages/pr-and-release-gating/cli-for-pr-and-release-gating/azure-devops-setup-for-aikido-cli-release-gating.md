# Azure DevOps Setup for Aikido CLI: Release Gating

The Aikido Security CI client allows you to integrate Aikido Security scans into CI pipelines. It helps ensure that security scans are part of your build process.

## 1. Get API token <a href="#id-1-get-api-token" id="id-1-get-api-token"></a>

1. Go to the [Continuous Integration Settings page](https://app.aikido.dev/settings/integrations/continuous-integration?regionRedirect=1).
2. Generate a token and copy. Note that you will only be able to view this token once.
3. Save this token as _AIKIDO\_CLIENT\_API\_KEY_ as a secret variable in your Azure DevOps pipeline by clicking _Variables_ when editing the pipeline, adding a new variable and checking _Keep this value secret_.

## 2. Create a new Azure DevOps pipeline <a href="#id-2-create-a-new-azure-devops-pipeline" id="id-2-create-a-new-azure-devops-pipeline"></a>

Navigate to Pipelines in the left hand navigation, click _Pipelines_ and select _New Pipeline_ or _Create Pipeline_. In the Configure step, select _Starter Pipeline_ and replace its contents with this example:

> Make sure that the scan is only triggered for your default branch. In the example below this is the 'main' branch. Adjust this if needed.

```yaml
trigger:
  branches:
    include:
      - main

pool:
  vmImage: 'ubuntu-latest'

steps:
- task: UseNode@1
  displayName: 'Set up Node.js'
  inputs:
    version: '22.x'
    checkLatest: true

- bash: |
    npm install --global @aikidosec/ci-api-client
    aikido-api-client scan-release \
      "$(Build.Repository.Name)" \
      "$(Build.SourceVersion)" \
      --apikey "$AIKIDO_CLIENT_API_KEY" \
      --fail-on-sast-scan \
      --fail-on-iac-scan \
      --fail-on-secrets-scan
  workingDirectory: '$(Build.SourcesDirectory)'
  displayName: 'Run Aikido Scan'
  env:
    AIKIDO_CLIENT_API_KEY: $(AIKIDO_CLIENT_API_KEY)
```

Tweak the command if needed. Every flag is documented on the [`@aikidosec/ci-api-client`](https://www.npmjs.com/package/@aikidosec/ci-api-client) package.

## 3. Run your first scan <a href="#id-3-run-your-first-scan" id="id-3-run-your-first-scan"></a>

When a new change is pushed to your default branch, the pipeline will be triggered. A scan will run and fail if new issues have been detected.
