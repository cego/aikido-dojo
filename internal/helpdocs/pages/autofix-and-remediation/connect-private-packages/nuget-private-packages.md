# NuGet - Private Packages

For Aikido to update dependencies that include private packages, it needs access to your private NuGet registries so it can generate accurate lockfile updates.

Aikido supports 2 ways of configuring connections with your private NuGet registry:

1. Via a `nuget.config` file
2. Via credentials for the dotnet cli

## 1. Configuration via nuget.config file <a href="#id-1-configuration-via-nugetconfig-file" id="id-1-configuration-via-nugetconfig-file"></a>

### Prerequisites <a href="#prerequisites" id="prerequisites"></a>

For private NuGet packages, Aikido uses a `nuget.config` file to authenticate with the private registry. This file will overwrite the `nuget.config` in the root of the repository. It is possible to configure multiple private registries in this file.

Example `nuget.config` for accessing private packages on GitHub's NuGet registry:

```
<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <packageSources>
    <clear />
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" protocolVersion="3" />
    <add key="github" value="https://nuget.pkg.github.com/AikidoSec/index.json" />
  </packageSources>
  <packageSourceCredentials>
    <github>
      <add key="Username" value="AikidoSec" />
      <add key="ClearTextPassword" value="ghp_ABC123...XYZ" />
    </github>
  </packageSourceCredentials>
</configuration>
```

Take a look at the following docs for more information on authenticating with private NuGet registries.

* [Microsoft - Consuming packages from authenticated feeds](https://learn.microsoft.com/en-us/nuget/consume-packages/consuming-packages-authenticated-feeds)
* [GitHub - Working with the NuGet registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-nuget-registry)

## Configuration in Aikido <a href="#configuration-in-aikido" id="configuration-in-aikido"></a>

Once the prerequisites are fulfilled, you can configure Aikido to authenticate with your private NuGet registry when updating the dependencies by following the steps below:

1. Go to your account's settings page for the Autofix in Aikido, [here](https://app.aikido.dev/issues/fix/settings).
2. Click on "Connect Registry", the private registry modal will now be shown

<div data-with-frame="true"></div>

3. Select "Nuget" as package manager and select the `nuget.config` option

<div data-with-frame="true"></div>

4. Fill in your `nuget.config` with authentication information. Aikido securely encrypts your configuration file until just before they are used.
5. Click "Connect Registry" to save the configuration

## 2. Configuration via .NET's nuget CLI <a href="#id-2-configuration-via-nets-nuget-cli" id="id-2-configuration-via-nets-nuget-cli"></a>

For the CLI config, Aikido needs to have the registry's URL, username and password which can be used to authenticate with the private registry. Once you have this information, you can configure it following the steps below:

1. Go to your account's settings page for the Autofix in Aikido, [here](https://app.aikido.dev/issues/fix/settings).
2. Click on "Connect Registry", the configuration modal will now be shown
3. Select "Nuget" as package manager and then select the "NuGet CLI Registry" option

<div data-with-frame="true"></div>

4. Enter the details in the input fields and click on "Connect Registry" to save the configuration

## 3. Configuration via environment variables

For NuGet, you can also provide credentials through an environment variable in Autofix.

This requires a `nuget.config` file in your repository that defines the private feed. The feed name in the environment variable must match the feed name in `nuget.config`.

1. Go to your account's settings page for Autofix in Aikido, [here](https://app.aikido.dev/issues/fix/settings).
2. Add the following environment variable:
   1. Key: `NuGetPackageSourceCredentials_my-private-feed`
   2. Value: `Username=user;Password=secret-pat`

The PAT must have package read access.

<div data-with-frame="true"></div>
