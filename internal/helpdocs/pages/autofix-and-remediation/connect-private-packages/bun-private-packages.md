# Bun - Private Packages

When Aikido updates dependencies in repositories that use private packages, it needs access to those packages to correctly update your lockfiles. You can configure Aikido to authenticate with your private NPM registry to run these updates.

## Prerequisites <a href="#prerequisites" id="prerequisites"></a>

Bun can use either `.npmrc` or `bunfig.toml`. See the [NPM documentation](/autofix-and-remediation/connect-private-packages/npm-private-packages.md) for more information about `.npmrc`.

### **Configuring `bunfig.toml`**

Private package access is defined in a `bunfig.toml` file in your project. See the Bun documentation [here](https://bun.com/docs/pm/scopes-registries).

If a `bunfig.toml` file is set in Aikido, we write the `bunfig.toml` file next to your `bun.lock` file before updating the dependencies. If the `bunfig.toml` file is pressent in the repository, it will be overwritten by the `.bunfig.toml` file set in Aikido.

## Configuration in Aikido <a href="#configuration-in-aikido" id="configuration-in-aikido"></a>

To allow repositories using private packages to be updated, provide your private registry configuration in the Aikido UI. Aikido stores the credentials encrypted and credentials cannot be retrieved through the Aikido UI or API.

1. Go to *Autofix* > *Settings* in Aikido, [here](https://app.aikido.dev/issues/fix/settings).
2. Click on "*Connect Registry*", and the selection modal will now be shown.

3. Select Bun to input your `bunfig.toml` file. (Note that for `.npmrc`, you will need to select NPM.)

4. Fill in the contents of the `bunfig.toml` file. The example shows a configuration for a scoped package (`@piedpiper`) hosted on GitHub Packages.
5. Click "*Connect Registry*" to save the configuration.
