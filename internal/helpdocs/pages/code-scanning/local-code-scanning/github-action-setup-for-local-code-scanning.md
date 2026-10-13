# GitHub Action Setup for Local On-Prem Code Scanning

## Introduction

The Aikido Security Local Scanner enables you to run scans directly within your GitHub Action and/or if you're using a **GitHub Enterprise Server (GHES)** environment. Your code never leaves your premises: scans are executed locally and only the results are uploaded to the Aikido Security platform. This setup ensures maximum privacy while still letting you manage vulnerabilities through the Aikido dashboard.

{% hint style="danger" %}
This is **NOT** applicable for normal cloud based scanner. This GitHub Action setup is only intended for on-prem scanning with local binaries.
{% endhint %}

### Prerequisites <a href="#how-to-set-up-local-scanning" id="how-to-set-up-local-scanning"></a>

* A GitHub Enterprise Server environment (self-hosted or enterprise-managed).
* A CI/CD runner - GitHub Actions
* An Aikido Local Scanning account.
* System requirements:
  * 2–4 CPU cores
  * 8–16 GB RAM
  * Outbound HTTPS (443) access enabled

### How to set up Local On-Prem Scanning <a href="#how-to-set-up-local-scanning" id="how-to-set-up-local-scanning"></a>

**Prerequisite**: make sure to have created an account that allows for Local Scanning. [More information on creating a Local Scanning Account](https://help.aikido.dev/en/articles/9070345-how-to-create-an-account-for-local-scanning-on-aikido).

#### 1. Generate an Authentication Token <a href="#id-1-get-your-authentication-token" id="id-1-get-your-authentication-token"></a>

1. Go to the [Local Scanner setup page](https://app.aikido.dev/settings/integrations/localscan).
2. Generate an authentication token and copy. **Note** that you will only be able to view this token once.
3. Add this token as argument `--apikey` when running the Local Scanner in your project.
4. Save this token in your GitHub Secrets by going to Settings > Secrets and variables > Actions.

   [image: Repository secrets management interface showing the stored AIKIDO\_API\_KEY secret.]

#### 2. Running the Local Scanner <a href="#id-2-running-the-local-scanner" id="id-2-running-the-local-scanner"></a>

Now all that is left to run the scanner on your repository.\
​\
Make sure that the local scanner is only triggered for your default branch. By default, Aikido supports scanning one branch in your repository for dependency and code issues, typically the main or master branch. Therefore, we recommend running the local scanner exclusively on that branch to avoid mixing scan results on the Aikido platform. You can specify this in the 'on' section of your workflow file.

#### Using Docker <a href="#using-docker" id="using-docker"></a>

The easiest way to use our local scanner in your GitHub Actions is by using our Docker image.\
​\
Example `.github/workflows/aikido-scan.yml`:

```yaml
​on:
  push:
    branches:
      - main

name: Aikido Scan
jobs:
  aikido-local-scan-repo:
    runs-on: ubuntu-latest
    container:
      image: aikidosecurity/local-scanner:latest
    steps: 
      - uses: actions/checkout@v4 
        with: 
          token: ${{ secrets.GITHUB_TOKEN }} 
          path: my-repo 
      - name: Run scan
        run: aikido-local-scanner scan my-repo --apikey ${{ secrets.AIKIDO_API_KEY }} --repositoryname MyRepo --branchname main
```

Specify your preferred branch using the `--branchname` option when executing the command.

If this is the first scan for this repository, Aikido will create a repository with the name you specified, containing all the scanning results. Subsequent scan results will be collected under this repository name in Aikido.

By default all scan types will be executed, if you'd like to run only a selection of scans (eg only a SAST scan), you can do so by supplying the scan types via the `--scan-types` option. More information on CLI options can be found [here](https://help.aikido.dev/en/articles/9027526-local-scanner-cli-options).

You can also run the scanner in release or PR gating mode. Release gating mode is helpful when scanning your repository prior to releasing, as it ensures there are no open issues before a potential release. When running in release gating mode, the scanner process will fail when there are any open issues of the chosen severity or higher after the scan is finished. PR gating mode can be used to scan for any potentially newly introduced issues in a PR.

More information about release or PR gating mode can be found in [this article](/code-scanning/local-code-scanning/pr-gating-for-code-using-local-scanner.md).

#### 3. Check your scanning results <a href="#id-3-check-your-scanning-results" id="id-3-check-your-scanning-results"></a>

After the first scan:

* A repository will appear in Aikido with the name you provided in `-repositoryname`.
* Subsequent scans will update the same repository with new findings.
* You can configure scans to run in **release gating** or **PR gating** mode if you want to block deployments with unresolved issues.

### Notes & Limitations

* Local Scanner accounts **does not include AutoFix** in the UI. AutoFix will be available via [IDE integrations](/ai-and-dev-tools/ide-plugins-overview.md)
* By default, all scan types are enabled. You can limit scans (e.g., SAST only) using the `-scan-types` flag.
  * Read "[CLI option for Local Scanner](/code-scanning/local-code-scanning/cli-options-for-local-scanner.md)" to view all options that can be passed when running the local scanner.
* Only one branch (usually `main` or `master`) should be scanned to avoid mixing results.
