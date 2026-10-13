# Sprinto Integration

The Sprinto integration automatically pushes vulnerabilities to Sprinto for SOC2 and ISO27001:2022.

{% hint style="info" %}
Aikido syncs repositories and vulnerabilities to Sprinto every 24 hours.
{% endhint %}

### Set up the Sprinto integration <a href="#set-up-the-sprinto-integration" id="set-up-the-sprinto-integration"></a>

{% stepper %}
{% step %}
**Add the integration**

To activate the Sprinto integration, in Aikido: Go to [Settings > Integrations > Sprinto](https://app.aikido.dev/settings/integrations/compliance/sprinto) and click 'Add Sprinto integration'.
{% endstep %}

{% step %}
**Create an API key**

You'll need to create a Sprinto API Key. You can generate such key in Sprinto via [Settings > Developer API > Generate API Key](https://app.sprinto.com/app/admin/developerApi)
{% endstep %}

{% step %}
**Save the API key**

Next, copy your generated API Key.

Back in Aikido, paste the API Key, select your region and click 'Save'.

[image: Paste your Sprinto API key in Aikido, select the region, and save the integration.]
{% endstep %}

{% step %}
**Done**

Aikido will now sync repositories and vulnerabilities to Sprinto every 24 hours.

[image: Sprinto integration connected and ready to sync vulnerability data from Aikido.]%20\(4\).png?alt=media)
{% endstep %}
{% endstepper %}
