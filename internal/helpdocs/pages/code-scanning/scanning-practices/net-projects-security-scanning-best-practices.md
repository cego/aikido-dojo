# NET Projects: Security Scanning Best Practices

Aikido can find known vulnerabilities (CVE) in your .NET dependencies as well as malware and dangerous licenses being used by those dependencies.

How does Aikido find those dependencies and their transitive subdependencies?

Out of the box, Aikido supports the following files for scanning:

* \*.csproj
* \*.deps.json
* packages.lock.json
* packages.config
* packages.props
* paket.lock

It should be noted that .csproj files often do not contain exact versions for subdependencies. This can cause Aikido to find some false positive CVEs in subdependencies.

## Why you should start using lockfiles <a href="#why-you-should-start-using-lockfiles" id="why-you-should-start-using-lockfiles"></a>

It's recommended to use lockfiles that contain a version and a hash for each NuGet dependency as well as the subdependencies.

There are other reasons to use lockfiles besides making security scanning easier for Aikido:

* Using a lockfile protects you against supply chain attacks via malicious packages. This kind of attack is becoming more popular
* Using a lockfile makes your build more predictable as everyone is using the exact same minor version of packages. Less chance of 'works on my machine'
* Faster build times: no need for dependency resolution anymore

### How to start using lockfiles in your .NET project? <a href="#how-to-start-using-lockfiles-in-your-net-project" id="how-to-start-using-lockfiles-in-your-net-project"></a>

To create a lock file, you need to add the following lines to your `.csproj` file:

```
<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <!-- Generate the lock file -->
    <RestorePackagesWithLockFile>true</RestorePackagesWithLockFile>

    <!-- Restore the exact packages as listed in the lock file -->
    <RestoreLockedMode Condition="'$(ContinuousIntegrationBuild)' == 'true'">true</RestoreLockedMode>
  </PropertyGroup>

  <ItemGroup>
    ...
  </ItemGroup>
</Project>
```

After these lines are added, run

```
dotnet.exe restore
```

This will generate a lockfile (packages.lock.json) that you can commit to the repo. You should never manually edit this file, similar to how you would use a package-lock.json file in NPM in NodeJS world.

To restore using an existing lockfile, run

```
dotnet.exe restore --locked-mode
```
