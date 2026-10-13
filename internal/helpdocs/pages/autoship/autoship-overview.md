# AutoShip Overview

AutoShip ships security fixes so you can focus on shipping features.

Although AutoFix optimizes to create non-breaking pull requests, it is still required to test the changes in a live environment to ensure there are no unwanted runtime effects. Shipping pull requests created by AutoFix thus can take significant time. You need to:

1. Check the CI/CD results
2. Deploy to a test environment
3. Ensure there are no new errors / regressions in the test environment
4. Perform a test to see nothing is broken

Once these tests are done, the pull request can safely be shipped to production.

AutoShip integrates with all your runtime systems in order to automate this process. AutoShip starts based on a label on the pull request. Once the label is added, AutoShip starts the process of shipping the pull request.

## Getting started

Go to [AutoShip](https://app.aikido.dev/autoship) in the Aikido UI and first configure the global settings: provide write access, the label to trigger on and your environments to deploy to. After connecting the CI/CD, Monitoring and Alerting integrations, you can start enabling AutoShip on your repositories. For each repository you need to configure which CI/CD they use and what environments AutoShip should deploy to. For each environment you can specify where the monitoring information can be found to detect regressions and how to smoke test that environment.

### Global Configuration

On the [AutoShip > Settings page](https://app.aikido.dev/autoship/settings) you can setup AutoShip:

#### Write Access for AutoShip

Currently we only support **GitHub**. Write access is provided through the `Aikido AutoShip` GitHub app. This app is separate from the `Aikdio AutoFix` GitHub app. You can configure independently on which repos AutoFix is enabled and on which repos AutoShip is enabled.

#### Trigger label

AutoShip triggers on pull requests that have a designated label. This label can be configured and will be used as the trigger across all repos on which AutoShip is enabled.

#### Environments

Specify the different environments AutoShip can deploy to. Typically the first environment is a staging environment, a safe space to test the changes. Once the changes are tested and verified on the first environment, it will move on to the next environment. Typicaly the second environment is the production environment, AutoShip will monitor the deploy and perform a smoke test to verify the changes.

Environments can be reused across multiple repositories. Specify the name of the environment and the branch in the global configuration to easily reuse them across your repositories.

The name of the branch is used by AutoShip as a safe-guard, AutoShip can only merge to this branch. For example if you deploy to a staging environment by merging to the `staging` branch, ou need to specify that branch name, otherwise AutoShip will not be able to merge. This prevents AutoShip from merging to sensitive branches (eg. main or production) in your repository.

#### CI/CD

AutoShip connects to your CI/CD system to check whether the build is passing. When build problems arise, AutoShip will fix them and push the fix on the pull request branch.

In the global configuration you connect your CI/CD system. In most cases this requires an API key, see the [CI/CD](/autoship/ci-cd.md) page for details per system. We currently support Github Actions and CircleCI.

Additional configuration per repository (eg. project) can be provided when configuring the repository.

#### Monitoring

Once the deploy to an environment succeeded, AutoShip will look at the data from your monitoring system to verify there are no new errors. If a regression is detected, AutoShip will roll back the changes.

In the global configuration you connect your monitoring systems. In most case this requires an API key, see the [Monitoring](/autoship/monitoring.md) page for details per system. We currently support Sentry and CloudWatch.

Additional configuration per repository (eg. project) can be provided when configuring the repository.

#### Alerting

AutoShip keeps you informed about its progress. This can be done through Slack. AutoShip creates a thread per deployment and posts progress messages:

### Per repository configuration

Once the global configuration is finished and the integrations are connected, it's only a matter of enabling them for your repositories.

#### CI/CD

Choose the CI/CD system that is used for this particular repository. For GitHub Action and CircleCI it's just a matter for selecting them, no further configuration required.

#### Environments

Select the environments you want to enable for this repository. You can choose to deploy some repositories to a testing environment without shipping all the way to production to perform a final manual verification.

In the environment configuration you will specify:

* **How AutoShip should deploy to this environment.** These are a free-text instructions, eg. `merge the PR branch into the staging branch`.
* **How AutoShip should test the functionality.** Here you configure which domain should be tested, how to authenticate to the domain and provide free-text instructions on how to test, eg. `add a product to the basket and check it is present on the checkout page`. [More information about smoke testing.](/autoship/smoke-testing.md)
* **How AutoShip should detect regressions.** Point AutoShip to the monitoring system used for this particular repository & environment. Eg. where to find the Cloudwatch logs for the service in staging.

#### Alerting

Specify on which Slack channel you want to be notified about the AutoShip progress.

### Try it!

Add the label you specified (eg. autoship) on a pull request.

* The Aikido AutoShip app will show up in the GitHub checks
* A Slack message will be posted
