# JFrog Artifactory

Route your registry traffic through Aikido's registry proxy by configuring JFrog Artifactory as an intermediary. Artifactory pulls packages from Aikido instead of the public registries, so every package your team installs is checked for malware, age restrictions, and known vulnerabilities.

{% stepper %}
{% step %}

### Get your Aikido proxy URL

Go to [Settings > Registry Proxy](https://app.aikido.dev/issues/aikido-libraries/setup-instructions), select your ecosystem and copy the upstream proxy URL

{% endstep %}

{% step %}

### Create a Remote Repository in Artifactory

1. In Artifactory, go to **Administration > Repositories**.
2. Click **New Repository** and select **Remote**.
3. Set **Package Type**
4. Give the repository a key, for example `aikido`.
5. Set the **URL** field to your Aikido proxy URL from step 1.
6. Save the repository.

{% hint style="info" %}
If packages don't reflect the latest data from Aikido, set the **Metadata Cache Lifetime** to `0` in the remote repository settings. This forces Artifactory to always fetch fresh metadata from the Aikido proxy instead of serving a cached version.
{% endhint %}
{% endstep %}

{% step %}

### (Optional) Create a Virtual Repository

If you also have internal packages stored in a local Artifactory repository, create a Virtual Repository that aggregates both:

1. Go to **Administration > Repositories > New Repository > Virtual**.
2. Set **Package Type** to match your remote repository
3. Add your local repository and the remote repository as members.
4. Point your package manager at the virtual repository URL instead.

This gives your team a single registry endpoint that serves both internal packages and externally proxied packages through Aikido.
{% endstep %}

{% step %}

### Configure Developer Machines to use Artifactory

Point each developer's package manager at the Artifactory repository. See JFrog's documentation for the exact steps:

{% tabs %}
{% tab title="npm" %}
[npm repositories - connect your client](https://docs.jfrog.com/artifactory/docs/npm-repositories#connect-your-npm-client-to-artifactory)
{% endtab %}

{% tab title="PyPI" %}
[PyPI repositories - connect your client](https://docs.jfrog.com/artifactory/docs/pypi-repositories#connect-your-pypi-client-to-artifactory)
{% endtab %}

{% tab title="Maven" %}
[Maven repositories - connect your client](https://docs.jfrog.com/artifactory/docs/maven-repositories#connect-your-maven-client-to-artifactory)
{% endtab %}

{% tab title="NuGet" %}
[NuGet repositories - connect your client](https://docs.jfrog.com/artifactory/docs/nuget-repositories#connect-your-nuget-client-to-artifactory)
{% endtab %}
{% endtabs %}
{% endstep %}

{% step %}

### (Optional) Block the public registry with Device Protection

Once Artifactory is in place, developers can still bypass it by pointing their package manager directly at the public registry. Deploy Aikido Device Protection to enforce that all traffic goes through your private registry.

{% content-ref url="/pages/sa7N9YXthBxZbm6WmKBq" %}
[Block Public Package Registries](/aikido-device-protection/using-aikido-endpoint/block-public-package-registries.md)
{% endcontent-ref %}
{% endstep %}
{% endstepper %}
