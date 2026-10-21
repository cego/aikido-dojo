# Azure DevOps Server Setup for Local Code Scanning

The Aikido Security Local Scanner is a tool that enables you to perform Aikido Security scans within your environment, ensuring your code never leaves your premises. The scans take place locally, and the results are then uploaded to the Aikido Security platform. This setup allows easy integration of the Local Scanner into your Azure DevOps Server project for reporting purposes.

## How to set up Local Scanning <a href="#how-to-set-up-local-scanning" id="how-to-set-up-local-scanning"></a>

**Prerequisite**: make sure to have created an account that allows for Local Scanning. [More information on creating a Local Scanning Account](https://help.aikido.dev/en/articles/9070345-how-to-create-an-account-for-local-scanning-on-aikido).

## 1. Setting up a Local Scanner pipeline <a href="#id-1-setting-up-a-local-scanner-pipeline" id="id-1-setting-up-a-local-scanner-pipeline"></a>

Navigate to Pipelines in the left hand navigation, click *Pipelines* and select *New Pipeline* or *Create Pipeline*. In the Configure step, select *Starter Pipeline*. This example pipeline below will trigger a scan on changes on the main branch. It's important to run you first scan on your default branch, as this will setup the repository in Aikido correctly.

```yaml
trigger:
  branches:
    include:
    - main

pool:
  vmImage: 'ubuntu-latest'

container:
  image: aikidosecurity/local-scanner:latest
  options: --entrypoint=""

steps:
- bash: |
    aikido-local-scanner scan . \
      --apikey "$(AIKIDO_API_KEY)" \
      --repositoryname "$(Build.Repository.Name)" \
      --branchname main
  workingDirectory: $(Build.SourcesDirectory)
  displayName: 'Run Aikido Scan'
```

If this is the first scan for this repository, Aikido will create a repository with the name you specified, containing all the scanning results. Subsequent scan results will be collected under this repository name in Aikido.

By default all scan types will be executed, if you'd like to run only a selection of scans (eg only a SAST scan), you can do so by supplying the scan types via the `--scan-types` option. More information on CLI options can be found [here](https://help.aikido.dev/en/articles/9027526-local-scanner-cli-options).

You can also run the scanner in release or PR gating mode. Release gating mode is helpful when scanning your repository prior to releasing, as it ensures there are no open issues before a potential release. When running in release gating mode, the scanner process will fail when there are any open issues of the chosen severity or higher after the scan is finished. PR gating mode can be used to scan for any potentially newly introduced issues in a PR.

More information about release or PR gating mode can be found in [this article](/code-scanning/local-code-scanning/pr-gating-for-code-using-local-scanner.md).

Now it is time to add the api key as a secret variable. Click variables in the right hand corner and add a new variable. Fill in AIKIDO\_API\_KEY as the name. Check the *Keep this value secret* checkbox. To find the secret value, move on to step 2.

## 2. Get your authentication token <a href="#id-2-get-your-authentication-token" id="id-2-get-your-authentication-token"></a>

1. Go to the [Local Scanner setup page](https://app.aikido.dev/settings/integrations/localscan)
2. Generate an authentication token and copy. Note that you will only be able to view this token once.
3. Paste this token as the value for the AIKIDO\_API\_KEY variable in the pipeline setup.

Now you can create and run your pipeline for the first time.

## 3. Check your scanning results <a href="#id-3-check-your-scanning-results" id="id-3-check-your-scanning-results"></a>

After your first scan is done, you can go to the Aikido Feed to check out your results. A repository will have been created, containing all results from the scanning.
