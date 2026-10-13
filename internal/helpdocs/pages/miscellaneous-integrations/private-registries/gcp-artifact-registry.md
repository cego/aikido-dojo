# Google Cloud Artifact Registry

Route your registry traffic through Aikido's registry proxy by configuring Google Cloud Artifact Registry as an intermediary. Artifact Registry pulls packages from Aikido instead of the public registries, so every package your team installs is checked for malware, age restrictions, and known vulnerabilities.

{% hint style="warning" %}
Artifact Registry does not support NuGet, so NuGet packages cannot be routed through Aikido this way.
{% endhint %}

{% stepper %}
{% step %}

### Get your Aikido proxy URL

Go to [Settings > Registry Proxy](https://app.aikido.dev/issues/aikido-libraries/setup-instructions), select your ecosystem and copy the upstream proxy URL

{% endstep %}

{% step %}

### Create a Remote Repository in Artifact Registry

1. In the Google Cloud Console, go to **Artifact Registry > Repositories**.
2. Click **Create Repository**.
3. Give the repository a name, for example `aikido`.
4. Set **Format** to match your ecosystem (npm, Python, or Maven).
5. Set **Repository mode** to **Remote**.
6. Under **Remote repository source**, select **Custom**.
7. In the **Custom repository** field, paste your Aikido proxy URL from step 1.
8. Set **Remote repository authentication mode** to **Unauthenticated**.
9. Choose a **Region** and click **Create**.

{% hint style="warning" %}
The **Custom repository URL** is immutable after creation. If you need to change it, delete the repository and create a new one.
{% endhint %}
{% endstep %}

{% step %}

### (Optional) Create a Virtual Repository

If you also publish internal packages to a Standard repository in Artifact Registry, create a Virtual Repository that aggregates both:

1. Go to **Artifact Registry > Repositories > Create Repository**.
2. Set **Format** to match your remote repository and **Repository mode** to **Virtual**.
3. Add your Standard repository and the Remote repository from step 2 as upstream repositories.
4. Assign a priority to each upstream: higher priority is preferred when a package exists in multiple upstreams.
5. Point your package manager at the Virtual Repository URL instead.

This gives your team a single registry endpoint that serves both internal packages and externally proxied packages through Aikido.
{% endstep %}

{% step %}

### Configure Developer Machines to use Artifact Registry

Point each developer's package manager at the Artifact Registry repository (the Virtual Repository, if you created one). See Google Cloud's documentation for the exact steps:

{% tabs %}
{% tab title="npm" %}
[Configure authentication for npm](https://cloud.google.com/artifact-registry/docs/nodejs/authentication)
{% endtab %}

{% tab title="Python" %}
[Configure authentication for Python](https://cloud.google.com/artifact-registry/docs/python/authentication)
{% endtab %}

{% tab title="Maven" %}
[Configure authentication for Maven](https://cloud.google.com/artifact-registry/docs/java/authentication)
{% endtab %}
{% endtabs %}
{% endstep %}

{% step %}

### (Optional) Block the public registry with Device Protection

Once Artifact Registry is in place, developers can still bypass it by pointing their package manager directly at the public registry. Deploy Aikido Device Protection to enforce that all traffic goes through your private registry.

{% content-ref url="/pages/sa7N9YXthBxZbm6WmKBq" %}
[Block Public Package Registries](/aikido-device-protection/using-aikido-endpoint/block-public-package-registries.md)
{% endcontent-ref %}
{% endstep %}
{% endstepper %}
