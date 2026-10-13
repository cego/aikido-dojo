# Sonatype Nexus Repository

Route your registry traffic through Aikido's registry proxy by configuring Nexus Repository as an intermediary. Nexus pulls packages from Aikido instead of the public registries, so every package your team installs is checked for malware, age restrictions, and known vulnerabilities.

{% stepper %}
{% step %}

### Get your Aikido proxy URL

Go to [Settings > Registry Proxy](https://app.aikido.dev/issues/aikido-libraries/setup-instructions), select your ecosystem and copy the upstream proxy URL

{% endstep %}

{% step %}

### Create a Proxy Repository in Nexus

1. In Nexus, go to **Settings > Repository > Repositories**.
2. Click **Create Repository** and select your ecosystem.
3. Give the repository a name, for example `aikido`.
4. Under **Remote Storage**, paste your Aikido proxy URL from step 1.
5. Save the repository.

{% hint style="warning" %}
Do not change the Remote Storage URL after the repository is created. Doing so can cause 404 errors as Nexus tries to reconcile cached metadata. If you need a different upstream URL, create a new proxy repository instead.
{% endhint %}
{% endstep %}

{% step %}

### (Optional) Create a Group Repository

If you also have internal packages in a hosted Nexus repository, create a Group Repository to expose both under a single URL:

1. Go to **Settings > Repository > Repositories > Create Repository** and select **aikido**
2. Add your hosted repository and the proxy repository as members.
3. Point your package manager at the group repository URL instead.

{% hint style="info" %}
Group repositories are read-only in Nexus Community Edition. To publish packages, always push directly to the hosted repository.
{% endhint %}
{% endstep %}

{% step %}

### Configure Developer Machines to use Nexus

Point each developer's package manager at the Nexus repository. See Sonatype's documentation for the exact steps, including authentication setup:

{% tabs %}
{% tab title="npm" %}
[Configuring npm](https://help.sonatype.com/en/configuring-npm.html)
{% endtab %}

{% tab title="PyPI" %}
[Configuring PyPI](https://help.sonatype.com/en/configure-pypi-with-nexus.html)
{% endtab %}

{% tab title="Maven" %}
[Configuring Maven](https://help.sonatype.com/en/maven-repositories.html)
{% endtab %}

{% tab title="NuGet" %}
[Configuring NuGet](https://help.sonatype.com/en/configure-nuget-with-nexus.html)
{% endtab %}
{% endtabs %}
{% endstep %}

{% step %}

### (Optional) Block the public registry with Device Protection

Once Nexus is in place, developers can still bypass it by pointing their package manager directly at the public registry. Deploy Aikido Device Protection to enforce that all traffic goes through your private registry.

{% content-ref url="/pages/sa7N9YXthBxZbm6WmKBq" %}
[Block Public Package Registries](/aikido-device-protection/using-aikido-endpoint/block-public-package-registries.md)
{% endcontent-ref %}
{% endstep %}
{% endstepper %}
