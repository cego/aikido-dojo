# Windows Setup for Local Code Scanning

**Table of contents:**

* [How to set up Local Scanning](#how-to-set-up-local-scanning)
  * [1. Get your authentication token](#1-get-your-authentication-token)
  * [2. Adding the Local Scanner to your project](#2-adding-the-local-scanner-to-your-project)
  * [3. Running the Local Scanner](#3-running-the-local-scanner)
  * [4. Check your scanning results](#4-check-your-scanning-results)

## Windows Setup for Local Code Scanning

The Aikido Security Local Scanner is a tool that enables you to perform Aikido Security scans within your environment, ensuring your code never leaves your premises. The scans take place locally, and the results are then uploaded to the Aikido Security platform. This setup allows you to scan any repos locally own your own machine.

### How to set up Local Scanning <a href="#how-to-set-up-local-scanning" id="how-to-set-up-local-scanning"></a>

**Prerequisite**: make sure to have created an account that allows for Local Scanning. [More information on creating a Local Scanning Account](https://help.aikido.dev/en/articles/9070345-how-to-create-an-account-for-local-scanning-on-aikido).

#### 1. Get your authentication token <a href="#id-1-get-your-authentication-token" id="id-1-get-your-authentication-token"></a>

1. Go to the [Local Scanner setup page](https://app.aikido.dev/settings/integrations/localscan).
2. Generate an authentication token and copy. Note that you will only be able to view this token once.
3. Add this token as argument `--apikey` when running the Local Scanner in your project .

#### 2. Adding the Local Scanner to your project <a href="#id-2-adding-the-local-scanner-to-your-project" id="id-2-adding-the-local-scanner-to-your-project"></a>

Download the local scanner binary from the Aikido UI.

#### 3. Running the Local Scanner <a href="#id-3-running-the-local-scanner" id="id-3-running-the-local-scanner"></a>

Now all that is left to spin up a container and scan your repository. Navigate to the root of your repository and run the following command; filling in your key, repository name and branch name.

```shellscript
aikido-local-scanner.exe scan . --apikey AIK_CI_xxx --repositoryname DemoApp --branchname main
```

Alternatively, you can use our [Docker image](https://hub.docker.com/r/aikidosecurity/local-scanner) to scan your repository. Ensure that you have [Docker installed](https://docs.docker.com/desktop/) before proceeding. Navigate into the directory you want to scan and run a scan using a command like:

If you are using the command prompt:

```shellscript
docker run --rm -v "%cd%:/my-app" aikidosecurity/local-scanner scan /my-app --apikey AIK_CI_xxx --repositoryname RepoName --branchname main
```

If you are using Powershell:

```shellscript
docker run --rm -v "${PWD}:/my-app" aikidosecurity/local-scanner scan /my-app --apikey AIK_CI_xxx --repositoryname RepoName --branchname main
```

By default all scanners will be executed, if you'd like to run only a selection of scanners, you can do so by supplying the scanner names `--scanners` option. More information on CLI options can be found [here](https://help.aikido.dev/en/articles/9027526-local-scanner-cli-options).

#### 4. Check your scanning results <a href="#id-4-check-your-scanning-results" id="id-4-check-your-scanning-results"></a>

After your first scan is done, you can go to the Aikido Feed to check out your results. A repository with the name you specified will have been created, containing all results from the scanning.
