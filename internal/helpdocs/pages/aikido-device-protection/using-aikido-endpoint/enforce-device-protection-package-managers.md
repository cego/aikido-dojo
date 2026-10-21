# Enforce Device Protection During Package Installation

Block package installations on machines where Aikido Device Protection is not active. This lets you enforce protection as part of your development workflow, without relying on policy alone.

## Why use this

Supply chain attacks often arrive through dependencies. Device Protection intercepts malicious packages at install time but only if it's running. Adding a pre-install check ensures developers cannot accidentally install packages on an unprotected machine.

## How it works

Most package managers support lifecycle hooks or wrapper scripts that run before installation. If the check exits with a non-zero code, the install is aborted and the developer sees a clear message.

{% tabs %}
{% tab title="macOS" %}
{% code overflow="wrap" %}

```bash
/Applications/Aikido\ Endpoint\ Protection.app/Contents/Resources/scripts/healthy
```

{% endcode %}
{% endtab %}

{% tab title="Windows" %}
{% code overflow="wrap" %}

```powershell
& "C:\Program Files\AikidoSecurity\EndpointProtection\scripts\Healthy.ps1"
```

{% endcode %}
{% endtab %}
{% endtabs %}

## Package manager examples

The examples below use the macOS path. On Windows, swap in the PowerShell command above.

{% tabs %}
{% tab title="npm" %}
Add a `preinstall` script to `package.json`:

```json
{
  "scripts": {
    "preinstall": "'/Applications/Aikido Endpoint Protection.app/Contents/Resources/scripts/healthy'"
  }
}
```

{% endtab %}

{% tab title="yarn / pnpm" %}
Both yarn and pnpm respect the `preinstall` script in `package.json`:

```json
{
  "scripts": {
    "preinstall": "'/Applications/Aikido Endpoint Protection.app/Contents/Resources/scripts/healthy'"
  }
}
```

{% endtab %}

{% tab title="pip / uv" %}
pip and uv have no native pre-install hook. Use a `Makefile` target that wraps the install:

```makefile
install: check-aikido
	pip install -r requirements.txt

check-aikido:
	/Applications/Aikido\ Endpoint\ Protection.app/Contents/Resources/scripts/healthy
```

Run `make install` instead of `pip install` directly. The same pattern works for `uv sync` or `poetry install`.
{% endtab %}

{% tab title="Maven" %}
Add the `exec-maven-plugin` to the `validate` phase in `pom.xml`. This runs before any dependencies are resolved:

```xml
<plugin>
  <groupId>org.codehaus.mojo</groupId>
  <artifactId>exec-maven-plugin</artifactId>
  <version>3.1.0</version>
  <executions>
    <execution>
      <id>check-aikido</id>
      <phase>validate</phase>
      <goals><goal>exec</goal></goals>
      <configuration>
        <executable>bash</executable>
        <arguments>
          <argument>-c</argument>
          <argument>/Applications/Aikido\ Endpoint\ Protection.app/Contents/Resources/scripts/healthy</argument>
        </arguments>
      </configuration>
    </execution>
  </executions>
</plugin>
```

{% endtab %}

{% tab title="NuGet" %}
Add a `.targets` file to your project (e.g. `CheckAikido.targets`) and import it in your `.csproj`:

```xml
<!-- CheckAikido.targets -->
<Project>
  <Target Name="CheckAikido" BeforeTargets="Restore">
    <Exec Command="bash -c &quot;/Applications/Aikido\ Endpoint\ Protection.app/Contents/Resources/scripts/healthy&quot;" />
  </Target>
</Project>
```

```xml
<!-- YourProject.csproj -->
<Import Project="CheckAikido.targets" />
```

{% endtab %}
{% endtabs %}

## Git hooks

Git hooks fire before commits or checkouts, not package installs, but they are a useful second line of enforcement for teams that want to block all git activity on unprotected machines.

Add the check to `.git/hooks/pre-commit`:

```bash
#!/bin/bash
/Applications/Aikido\ Endpoint\ Protection.app/Contents/Resources/scripts/healthy
```

Make it executable:

```bash
chmod +x .git/hooks/pre-commit
```

{% hint style="info" %}
`.git/hooks` is not tracked by git. To share hooks with your team, use a tool like [Lefthook](https://github.com/evilmartians/lefthook) or [Husky](https://typicode.github.io/husky/) to manage and commit hooks as part of your repository.
{% endhint %}
