# Report a Problem / Device Protection Support

When contacting Aikido Support about Aikido Device Protection, please include the information below. This helps us troubleshoot issues faster and avoid unnecessary back and forth.

### Report a problem (recommended)

The easiest way to send diagnostics is directly from the app. Aikido will automatically attach agent diagnostics data (never personal data, deleted after debugging).

1. Open the Aikido Endpoint Protection app
2. Go to the **Help** tab
3. Click **Report A Problem**
4. Add any additional notes about the issue (optional)
5. Leave your email if you'd like us to follow up (optional)
6. Click **Submit Report**

### Run the diagnostics command (fallback)

If the app is not accessible, run this command and attach the diagnostics logs. We only keep this data for as long as we need it to fix the problem.

{% tabs %}
{% tab title="macOS" %}
{% code overflow="wrap" %}

```bash
sudo /Applications/Aikido\ Endpoint\ Protection.app/Contents/Resources/scripts/submit_diagnostics
```

{% endcode %}
{% endtab %}

{% tab title="Windows" %}
{% code overflow="wrap" %}

```powershell
& "C:\Program Files\AikidoSecurity\EndpointProtection\scripts\SubmitDiagnostics.ps1"
```

{% endcode %}
{% endtab %}

{% tab title="Linux" %}
{% code overflow="wrap" %}

```bash
sudo aikidoctl diagnostics
```

{% endcode %}
{% endtab %}
{% endtabs %}

### Find the diagnostics zip if the upload fails

If the report or the diagnostics command can't reach Aikido, the zip is still saved on the device. Attach it to your email or support ticket manually.

| Platform | Location                                                                                                |
| -------- | ------------------------------------------------------------------------------------------------------- |
| macOS    | `/Library/Logs/AikidoSecurity/EndpointProtection/aikido-endpoint-protection-logs-<timestamp>.zip`       |
| Windows  | `C:\ProgramData\AikidoSecurity\EndpointProtection\logs\aikido-endpoint-protection-logs-<timestamp>.zip` |
| Linux    | `/var/log/aikidosecurity/endpoint-protection/aikido-endpoint-protection-logs-<timestamp>.zip`           |
