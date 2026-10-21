# Uninstall Device Protection

Use the commands below in your terminal to uninstall Aikido Device Protection from your system.

{% tabs %}
{% tab title="macOS" %}
{% code overflow="wrap" %}

```bash
sudo "/Applications/Aikido Endpoint Protection.app/Contents/Resources/scripts/uninstall"
```

{% endcode %}
{% endtab %}

{% tab title="Windows" %}
{% code overflow="wrap" %}

```powershell
& powershell.exe -NoProfile -ExecutionPolicy Bypass -File "C:\Program Files\AikidoSecurity\EndpointProtection\scripts\Uninstall.ps1"
```

{% endcode %}

Alternatively, open **Control Panel** > **Uninstall a program**, select **Aikido Endpoint Protection**, and click **Uninstall**.
{% endtab %}

{% tab title="Linux" %}
On Debian-family distributions:

{% code overflow="wrap" %}

```bash
sudo apt remove aikido-endpoint-protection
```

{% endcode %}

On Red Hat-family distributions:

{% code overflow="wrap" %}

```bash
sudo dnf remove aikido-endpoint-protection
```

{% endcode %}

Removing the package stops the service, restores the system trust store, and deletes the agent's state and logs.
{% endtab %}
{% endtabs %}

{% hint style="danger" %}
After uninstalling, reboot your device to avoid connection errors and ensure it's fully removed.
{% endhint %}
