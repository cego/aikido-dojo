# Yarn - Private Packages

When Aikido updates dependencies in repositories that use private packages, it needs access to those packages to correctly update your lockfiles. You can configure Aikido to authenticate with your private NPM registry to run these updates.

## Prerequisites <a href="#prerequisites" id="prerequisites"></a>

If you’re using **Yarn v1**, use the `.npmrc` option in Aikido as described in [the NPM & PNPM documentation](/autofix-and-remediation/connect-private-packages/npm-private-packages.md). For **Yarn v2 and newer**, use the `.yarnrc.yml` file to configure private package access.

### **Configuring `.yarnrc.yml`**

Private package access is defined in a `.yarnrc.yml` file in your project. This file declares `npmScopes` for your private packages and specifies which registry they use.

If a `.yarnrc.yml` file is set in Aikido, we write the `.yarnrc.yml` file next to your `yarn.lock` file before updating the dependencies. If the `.yarnrc.yml` file is pressent in the repository, it will be overwritten by the `.yarnrc.yml` file set in Aikido.

Example `.yarnrc.yml`

```yaml
npmScopes:
  piedPiper:
    npmAuthToken: "<YOUR_GITHUB_TOKEN>"
    npmRegistryServer: "https://npm.pkg.github.com"
```

You can find more information in the Yarn documentation [here](https://yarnpkg.com/configuration/yarnrc#npmScopes).

## Configuration in Aikido <a href="#configuration-in-aikido" id="configuration-in-aikido"></a>

To allow repositories using private packages to be updated, provide your private registry configuration in the Aikido UI. Aikido stores the credentials encrypted and credentials cannot be retrieved through the Aikido UI or API.

1. Go to *Autofix* > *Settings* in Aikido, [here](https://app.aikido.dev/issues/fix/settings).
2. Click on "*Connect Registry*", and the selection modal will now be shown.

3. Select Yarn to input your `.yarnrc.yml` file. (Note that for Yarn V1, you will need to select NPM.)

4. Fill in the contents of the `.yarnrc.yml` file. The example shows a configuration for GitHub Packages in the piedPiper GitHub organisation.
5. Click "*Connect Registry*" to save the configuration.
