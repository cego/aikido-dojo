# Cloudsmith

Route your registry traffic through Aikido's registry proxy by adding it as an upstream on your Cloudsmith repository. Cloudsmith pulls packages from Aikido instead of the public registries, so every package your team installs is checked for malware, age restrictions, and known vulnerabilities.

{% stepper %}
{% step %}

### Get your Aikido proxy URL

Go to [Settings > Registry Proxy](https://app.aikido.dev/issues/aikido-libraries/setup-instructions), select your ecosystem and copy the upstream proxy URL

{% endstep %}

{% step %}

### Add an Upstream in Cloudsmith

1. In Cloudsmith, open the repository you want to route through Aikido.
2. Go to the **Sources** tab and click **+ Configure upstream**.
3. Choose **Manual setup** and select the format that matches your ecosystem (npm, Python, Maven, or NuGet).
4. Give the upstream a **Name**, for example `aikido`.
5. Set the **Proxy URL** to your Aikido proxy URL from step 1.
6. Set the **Mode** to **Cache and Proxy** so Cloudsmith serves packages that Aikido has approved.
7. Leave **SSL Verification** enabled and set a **Priority** (higher priority upstreams are resolved first).
8. Save the upstream.
   {% endstep %}

{% step %}

### Configure Developer Machines to use Cloudsmith

Point each developer's package manager at your Cloudsmith repository. Cloudsmith provides copy-paste setup snippets, pre-configured with your namespace and repository, under the repository's **Setup** instructions. See Cloudsmith's documentation for the exact steps, including authentication:

{% tabs %}
{% tab title="npm" %}
[npm Registry](https://docs.cloudsmith.com/formats/npm-registry)
{% endtab %}

{% tab title="Python" %}
[Python Repository](https://docs.cloudsmith.com/formats/python-repository)
{% endtab %}

{% tab title="Maven" %}
[Maven Repository](https://docs.cloudsmith.com/formats/maven-repository)
{% endtab %}

{% tab title="NuGet" %}
[NuGet Feed](https://docs.cloudsmith.com/formats/nuget-feed)
{% endtab %}
{% endtabs %}
{% endstep %}

{% step %}

### (Optional) Block the public registry with Device Protection

Once Cloudsmith is in place, developers can still bypass it by pointing their package manager directly at the public registry. Deploy Aikido Device Protection to enforce that all traffic goes through your private registry.

{% content-ref url="/pages/sa7N9YXthBxZbm6WmKBq" %}
[Block Public Package Registries](/aikido-device-protection/using-aikido-endpoint/block-public-package-registries.md)
{% endcontent-ref %}
{% endstep %}
{% endstepper %}
