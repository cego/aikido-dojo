# NPM & PNPM - Private Packages

When Aikido updates dependencies in repositories that use private packages, it needs access to those packages to correctly update your lockfiles. You can configure Aikido to authenticate with your private NPM registry to run these updates.

## Prerequisites <a href="#prerequisites" id="prerequisites"></a>

For repositories using NPM or PNPM, access to private packages is managed through an `.npmrc` file. The package manager uses the `.npmrc` file to know which registry to use and how to authenticate.

You can either:

* include an `.npmrc` file in your repository, referencing an [environment variable](https://help.aikido.dev/~/revisions/KY3QqEqZPV6ifxqjwoLa/aikido-autofix/connect-private-packages/custom-config-private-packages) for the token.
* provide the `.npmrc` configuration directly in the Aikido interface, [here](#configuration-in-aikido).

If your private package names look like this: `@pied-piper/***` and are hosted on GitHub's npm registry, your file should look something like this:

```
@pied-piper:registry=https://npm.pkg.github.com
//npm.pkg.github.com/:_authToken=<YOUR_GITHUB_TOKEN>
```

If a `.npmrc` file is set in Aikido, we will write the `.npmrc` file before updating the dependencies. If an `.npmrc` file is present in the repository, it will be overwritten by the `.npmrc` file set in Aikido.

You can find more info about `.npmrc` from NPM [here](https://docs.npmjs.com/cli/v8/configuring-npm/npmrc#auth-related-configuration).

## Configuration in Aikido <a href="#configuration-in-aikido" id="configuration-in-aikido"></a>

To allow repositories using private packages to be updated, provide your private registry configuration in the Aikido UI. Aikido stores the credentials encrypted and credentials cannot be retrieved through the Aikido UI or API.

1. Go to *Autofix* > *Settings* in Aikido, [here](https://app.aikido.dev/issues/fix/settings).
2. Click on "*Connect Registry*", and the selection modal will now be shown.

3. Select NPM to input your `.npmrc` file.

4. Fill in the contents of the `.npmrc` file. The example shows a configuration for a scoped package (`@pied-piper`) hosted on GitHub Packages.
5. Click "*Connect Registry*" to save the configuration.
