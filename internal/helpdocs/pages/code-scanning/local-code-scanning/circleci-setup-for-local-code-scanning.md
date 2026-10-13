# CircleCI Setup for Local Code Scanning

The Aikido Security Local Scanner is a tool that enables you to perform Aikido Security scans within your environment, ensuring your code never leaves your premises. The scans take place locally, and the results are then uploaded to the Aikido Security platform. This setup allows easy integration of the Local Scanner into your CircleCI pipelines for reporting purposes.

## How to set up Local Scanning <a href="#how-to-set-up-local-scanning" id="how-to-set-up-local-scanning"></a>

**Prerequisite**: make sure to have created an account that allows for Local Scanning. [More information on creating a Local Scanning Account](https://help.aikido.dev/en/articles/9070345-how-to-create-an-account-for-local-scanning-on-aikido).

### 1. Get your authentication token <a href="#id-1-get-your-authentication-token" id="id-1-get-your-authentication-token"></a>

1. Go to the [Local Scanner setup page](https://app.aikido.dev/settings/integrations/localscan).
2. Generate an authentication token and copy. Note that you will only be able to view this token once.
3. Add this token as argument `--apikey` when running the Local Scanner in your project.
4. Save this token as an environment variable by going to **Project Settings** > **Environment Variables** > **Add Environment Variable**.

   [image: Environment variable management interface showing an API key entry and options to add or import variables.]

## 2. Running the Local Scanner <a href="#id-2-running-the-local-scanner" id="id-2-running-the-local-scanner"></a>

Now all that is left to run the scanner on your repository.\
​\
Make sure that the local scanner is only triggered for your default branch. By default, Aikido supports scanning one branch in your repository for dependency and code issues, typically the main or master branch. Therefore, we recommend running the local scanner exclusively on that branch to avoid mixing scan results on the Aikido platform. You can specify this using branch filters.

## Using Docker <a href="#using-docker" id="using-docker"></a>

The easiest way to use our local scanner in your CircleCI Pipeline is by using our Docker image.

Example `.circleci/config.yml` :

```yaml
version: 2.1

jobs:
  run-aikio-scan:
    docker:
      - image: aikidosecurity/local-scanner:latest

    steps:
      - checkout:
          path: my-app
      - run:
          name: "Run Aikido scan"
          command: aikido-local-scanner scan my-app --repositoryname my-app --branchname main --apikey "$AIKIDO_API_KEY"
workflows:
  my-workflow:
    jobs:
      - run-aikio-scan:
          filters:
            branches:
              only:
                - main
```

Specify your preferred branch using the `--branchname` option when executing the command.

If this is the first scan for this repository, Aikido will create a repository with the name you specified, containing all the scanning results. Subsequent scan results will be collected under this repository name in Aikido.

By default all scan types will be executed, if you'd like to run only a selection of scans (eg only a SAST scan), you can do so by supplying the scan types via the `--scan-types` option. More information on CLI options can be found [here](https://help.aikido.dev/en/articles/9027526-local-scanner-cli-options).

You can also run the scanner in release or PR gating mode. Release gating mode is helpful when scanning your repository prior to releasing, as it ensures there are no open issues before a potential release. When running in release gating mode, the scanner process will fail when there are any open issues of the chosen severity or higher after the scan is finished. PR gating mode can be used to scan for any potentially newly introduced issues in a PR.

More information about release or PR gating mode can be found in [this article](/code-scanning/local-code-scanning/pr-gating-for-code-using-local-scanner.md).

## 3. Check your scanning results <a href="#id-3-check-your-scanning-results" id="id-3-check-your-scanning-results"></a>

After your first scan is done, you can go to the Aikido Feed to check out your results. A repository with the name you specified will have been created, containing all results from the scanning.
