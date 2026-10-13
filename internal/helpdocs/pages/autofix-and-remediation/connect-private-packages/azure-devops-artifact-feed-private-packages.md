# Azure DevOps Artifact Feed - Private Packages

For node project, the connection with an Azure DevOps Artifact Feed can be configured by providing an `.npmrc` file in the Aikido Autofix configuration

1. Go to your account's settings page for the autofixer in Aikido, [here](https://app.aikido.dev/issues/fix/settings).
2. Click on "Connect Registry", an select "NPM"

3. Fill in the `.nmprc` file like described below

4. Click "Connect Registry" to save the configuration.

### Obtaining the .npmrc for Azure DevOps

1. In Azure Devops, navigate to the "Artifacts" page of the project, click on "Connect to Feed" button and select "NPM" from the list
2. You will now land on page describing how to setup the project, select "Other" From the project setup tabs, you should see a page with setup instructions similar to this. Follow the steps and paste the `.npmrc` content into Aikido.

[image: Azure npm registry setup instructions with authentication and access token configuration steps.]
