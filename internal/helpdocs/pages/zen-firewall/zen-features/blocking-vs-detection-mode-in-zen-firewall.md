# Blocking vs Detection Mode in Zen Firewall

Zen Firewall by Aikido offers two operational modes for handling attack patterns: **Detection Mode** and **Blocking Mode**. These modes specifically apply to application-level attacks like SQL injection, shell injection, path traversal, and SSRF.

Other security features like [user blocking](/zen-firewall/zen-features/blocking-users-with-zen-firewall.md), [bot protection](/zen-firewall/zen-features/blocking-bot-traffic-with-zen-firewall.md), [tor blocking](/zen-firewall/zen-features/blocking-tor-traffic-with-zen-firewall.md) and [country blocking](/zen-firewall/zen-features/blocking-traffic-by-country-with-zen-firewall.md) are managed separately through their respective settings.

## Understanding the Modes <a href="#understanding-the-modes" id="understanding-the-modes"></a>

### Detection Mode (Default) <a href="#detection-mode-default" id="detection-mode-default"></a>

* 📝 Logs potential attacks without blocking them
* 🔬 Allows you to monitor threats without impacting traffic
* 📊 Helps evaluate security patterns before enabling blocks

### Blocking Mode <a href="#blocking-mode" id="blocking-mode"></a>

* 🛡️ Actively blocks detected attack patterns
* 🚫 Blocks malicious requests before they can execute
* ⚡ Provides real-time protection against threats

## Configuration Options <a href="#configuration-options" id="configuration-options"></a>

### Using Environment Variables <a href="#using-environment-variables" id="using-environment-variables"></a>

```
# Enable Blocking Mode 
AIKIDO_BLOCK="true" 

# Enable Detection Mode (default) 
AIKIDO_BLOCK="false"
```

### Using the Dashboard <a href="#using-the-dashboard" id="using-the-dashboard"></a>

1. Navigate to your app's firewall settings
2. Find the "**Block Attacks**" section
3. Toggle between "Detection-Only" and "Blocking-Mode"

[image: Firewall settings for blocking attacks and Tor traffic in Demo app.]

{% hint style="info" %}
Note: The environment variable sets the initial mode, but the dashboard configuration takes precedence once loaded.
{% endhint %}

## What Gets Blocked? <a href="#what-gets-blocked" id="what-gets-blocked"></a>

Blocking mode only affects these attack types:

* SQL Injection attempts
* NoSQL Injection attempts
* Shell Injection attacks
* Path Traversal exploits
* SSRF attempts
