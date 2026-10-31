# GitHub Action Setup for Aikido CLI: Release Gating

The Aikido Security CI client allows you to integrate Aikido Security scans into CI pipelines. It helps ensure that security scans are part of your build process.

## 1. Get API token <a href="#id-1-get-api-token" id="id-1-get-api-token"></a>

1. Go to the [Continuous Integration Settings page](https://app.aikido.dev/settings/integrations/continuous-integration).
2. Generate a token and copy. Note that you will only be able to view this token once.
3. Save this token as *AIKIDO\_CLIENT\_API\_KEY* in your GitHub Secrets by going to Settings > Secrets and variables > Actions.

## 2. Create a new GitHub Action workflow <a href="#id-2-create-a-new-github-action-workflow" id="id-2-create-a-new-github-action-workflow"></a>

Create a workflow file using this example:

> Make sure that the local scanner is only triggered for your default branch. In the example below this is the 'main' branch. Adjust this if needed.

```
name: Aikido Release Gated Scan
on:
  push:
    branches:
      - main

jobs:
  aikido-scan:
    runs-on: ubuntu-latest

    steps:
      # Step 1: Set up Node.js
      - name: Set up Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '22'

      # Step 2: Install dependencies
      - name: Install Aikido CI API Client
        run: npm install --global @aikidosec/ci-api-client

      # Step 3: Run Aikido Scan Release
      - name: Run Aikido Scan
        run: aikido-api-client scan-release ${{ github.event.repository.name }} $GITHUB_SHA --apikey ${{ secrets.AIKIDO_CLIENT_API_KEY }} --fail-on-sast-scan --fail-on-iac-scan --fail-on-secrets-scan
```

Tweak the command if needed. Every flag is documented on the [`@aikidosec/ci-api-client`](https://www.npmjs.com/package/@aikidosec/ci-api-client) package.

## 3. Run your first scan <a href="#id-3-run-your-first-scan" id="id-3-run-your-first-scan"></a>

When a new change is pushed to your default branch, the workflow will be triggered. A scan will run and fail if new issues have been detected.
