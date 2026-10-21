# Bamboo Setup for Local Code Scanning

The Aikido Security Local Scanner is a tool that enables you to perform Aikido Security scans within your environment, ensuring your code never leaves your premises. The scans take place locally, and the results are then uploaded to the Aikido Security platform. This setup allows easy integration of the Local Scanner into Bamboo for reporting purposes.

## How to set up Local Scanning <a href="#how-to-set-up-local-scanning" id="how-to-set-up-local-scanning"></a>

**Prerequisite**: make sure to have created an account that allows for Local Scanning. [More information on creating a Local Scanning Account](https://help.aikido.dev/en/articles/9070345-how-to-create-an-account-for-local-scanning-on-aikido).

## 1. Get your authentication token <a href="#id-1-get-your-authentication-token" id="id-1-get-your-authentication-token"></a>

1. Go to the [Local Scanner setup page](https://app.aikido.dev/settings/integrations/localscan).
2. Generate an authentication token and copy. Note that you will only be able to view this token once.
3. Add the authentication token to the Plan variables in Bitbucket to make it available for use in the plan. To do this, you need to navigate to the plan configuration > Variables. There you can add the authentication token with `AIKIDO_API_KEY` as the key and the copied token contents as the value. You can also choose to define the variable on Project or Global level (see [here](https://confluence.atlassian.com/bamboo/bamboo-variables-289277087.html)).

## 2. Running the Local Scanner <a href="#id-2-running-the-local-scanner" id="id-2-running-the-local-scanner"></a>

Make sure that the local scanner is only triggered for your default branch. By default, Aikido supports scanning one branch in your repository for dependency and code issues, typically the main or master branch. Therefore, we recommend running the local scanner exclusively on that branch to avoid mixing scan results on the Aikido platform.

## Using Docker <a href="#using-docker" id="using-docker"></a>

If your agent has the capability to run Docker, you can use the Aikido Local Scanner Docker image to run the scans. First, make sure there is 'Source Code Checkout' task to ensure the code is available to be scanned. After that add a 'Docker' task like this:

[image: Search for Docker task type in Atlassian integration task selection window.]

Select 'Run a Docker container' as the command. Our image can be found on [Docker Hub](https://hub.docker.com/r/aikidosecurity/local-scanner) or [AWS Elastic Container Registry](https://gallery.ecr.aws/z1o3v2w5/aikidosecurity/local-scanner).

Fill in the Container Command field like this: `scan . --apikey ${bamboo.AIKIDO_API_KEY} --repositoryname ${bamboo.planRepository.1.name} --branchname ${bamboo.planRepository.1.branch}`

## Using the binary <a href="#using-the-binary" id="using-the-binary"></a>

If your agent does not have the capability to run Docker, you can use the Aikido Local Scanner binary directly. To do this, download the correct binary from the [Local Scanner Setup page](https://app.aikido.dev/settings/integrations/localscan).\
Now you need to add this a capability to your agent.

* **Navigate to the Desired Agent:**
  * In Bamboo, go to **Build Resources** > **Agents**.
  * Select the agent where you want to add the capability.
* **Add the Executable Capability:**
  * In the **Agent-specific capabilities** section, click **Add capability**.
  * Choose **Capability type** > **Executable**.
  * In the **Executable label**, enter a name for the executable, such as Aikido Local Scanner.
  * In the **Path** field, specify the path to the aikido-local-scanner binary on the agent machine.

Now, you can use the Aikido Local scanner as a Task in your Plan. First, make sure there is 'Source Code Checkout' task to ensure the code is available to be scanned. After that add a 'Command' task like this:

[image: Task types menu showing command execution options in a build automation tool.]

As the executable, choose 'Aikido Local Scanner' and add the command parameters in the Argument field:

```shellscript
scan . --apikey ${bamboo.AIKIDO_API_KEY} --repositoryname ${bamboo.planRepository.1.name} --branchname ${bamboo.planRepository.1.branch}
```

If this is the first scan for this repository, Aikido will create a repository with the name you specified, containing all the scanning results. Subsequent scan results will be collected under this repository name in Aikido.

By default all scan types will be executed, if you'd like to run only a selection of scans (eg only a SAST scan), you can do so by supplying the scan types via the `--scan-types` option. More information on CLI options can be found [here](https://help.aikido.dev/en/articles/9027526-local-scanner-cli-options).

You can also run the scanner in release or PR gating mode. Release gating mode is helpful when scanning your repository prior to releasing, as it ensures there are no open issues before a potential release. When running in release gating mode, the scanner process will fail when there are any open issues of the chosen severity or higher after the scan is finished. PR gating mode can be used to scan for any potentially newly introduced issues in a PR.

More information about release or PR gating mode can be found in [this article](/code-scanning/local-code-scanning/pr-gating-for-code-using-local-scanner.md).

### 3. Check your scanning results <a href="#id-3-check-your-scanning-results" id="id-3-check-your-scanning-results"></a>

After your first scan is done, you can go to the Aikido Feed to check out your results. A repository with the name you specified will have been created, containing all results from the scanning.
