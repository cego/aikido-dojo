# Bitbucket Pipeline Setup for Local Code Scanning

The Aikido Security Local Scanner is a tool that enables you to perform Aikido Security scans within your environment, ensuring your code never leaves your premises. The scans take place locally, and the results are then uploaded to the Aikido Security platform. This setup allows easy integration of the Local Scanner into your Bitbucket project for reporting purposes.

## How to set up Local Scanning <a href="#how-to-set-up-local-scanning" id="how-to-set-up-local-scanning"></a>

**Prerequisite**: make sure to have created an account that allows for Local Scanning. [More information on creating a Local Scanning Account](https://help.aikido.dev/en/articles/9070345-how-to-create-an-account-for-local-scanning-on-aikido).

## 1. Get your authentication token <a href="#id-1-get-your-authentication-token" id="id-1-get-your-authentication-token"></a>

1. Go to the [Local Scanner setup page](https://app.aikido.dev/settings/integrations/localscan).
2. Generate an authentication token and copy. Note that you will only be able to view this token once.
3. Add the authentication token to the Repository variables in Bitbucket to make it available for use in the pipeline. To do this, you need to go to your Repository Settings. Then navigate to **Pipelines > Repository variables.** There you can add the authentication token with `AIKIDO_API_KEY` as the key and the copied token contents as the value.

## 2. Running the Local Scanner <a href="#id-2-running-the-local-scanner" id="id-2-running-the-local-scanner"></a>

Make sure that the local scanner is only triggered for your default branch. By default, Aikido supports scanning one branch in your repository for dependency and code issues, typically the main or master branch. Therefore, we recommend running the local scanner exclusively on that branch to avoid mixing scan results on the Aikido platform.

### Using Docker <a href="#using-docker" id="using-docker"></a>

The easiest way to use our local scanner in your Bitbucket Pipelines is by using our Docker image.\
​\
Example `bitbucket-pipelines.yml`:

```yaml
image: aikidosecurity/local-scanner:latest

pipelines:
  branches:
    main:
      - step:
          name: Run Aikido scan
          script:
            - aikido-local-scanner scan . --repositoryname my-bitbucket-app --branchname main --apikey $AIKIDO_API_KEY
```

Specify your preferred branch using the `--branchname` option when executing the command.

If this is the first scan for this repository, Aikido will create a repository with the name you specified, containing all the scanning results. Subsequent scan results will be collected under this repository name in Aikido.

By default all scan types will be executed, if you'd like to run only a selection of scans (eg only a SAST scan), you can do so by supplying the scan types via the `--scan-types` option. More information on CLI options can be found [here](https://help.aikido.dev/en/articles/9027526-local-scanner-cli-options).

You can also run the scanner in release or PR gating mode. Release gating mode is helpful when scanning your repository prior to releasing, as it ensures there are no open issues before a potential release. When running in release gating mode, the scanner process will fail when there are any open issues of the chosen severity or higher after the scan is finished. PR gating mode can be used to scan for any potentially newly introduced issues in a PR.

More information about release or PR gating mode can be found in [this article](/code-scanning/local-code-scanning/pr-gating-for-code-using-local-scanner.md).

## 3. Check your scanning results <a href="#id-3-check-your-scanning-results" id="id-3-check-your-scanning-results"></a>

After your first scan is done, you can go to the Aikido Feed to check out your results. A repository with the name you specified will have been created, containing all results from the scanning.
