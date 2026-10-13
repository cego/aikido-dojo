# Setting Up the Local VM Scanner on Windows

Aikido VM Scanner is a single-package that installs on your system, automatically scanning and identifying dependencies to provide a detailed view into your environment.

### Prerequisites

* Minimum system requirements: at least 1GB RAM.
* Preferred system requirements: at least 2GB RAM and 4 CPUs.
* Ensure you have admin privileges on your system
* Make sure to use the appropriate commands for your system or cloud provider
* [On Windows, ensure that PowerShell is installed and enabled](https://learn.microsoft.com/en-us/powershell/scripting/install/installing-powershell-on-windows?view=powershell-7.5)

### Installation and Upgrade

Start PowerShell in admin mode, replace `AIKIDO_TOKEN` with [valid token from Local VM scanning page in Aikido](https://app.aikido.dev/settings/integrations/vm-scan-agent). You can also specify the following optional parameters:

* `VM_TYPE` as one of: `production`, `staging` or `development`.
* `OUTPUT` as one of: `stdout`, `stderr` or `none`.
* `INITIAL_SCAN` as one of: `1` or `0`.

```
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
Invoke-WebRequest -Uri "https://aikido-vm-agent.s3.eu-west-1.amazonaws.com/v1.4.12/AikidoVmScanner.msi" -OutFile "AikidoVmScanner.msi"
msiexec /i AikidoVmScanner.msi /qn /norestart AIKIDO_TOKEN=REPLACE_ME VM_TYPE=production OUTPUT=stdout INITIAL_SCAN=1
```

<table><thead><tr><th width="99.64453125">Hash</th><th>AikidoVmScanner.msi</th></tr></thead><tbody><tr><td>MD5</td><td>4f57b9af4b8117a464e8252ae100804d</td></tr><tr><td>SHA256</td><td>5b00e28dadd50bc0269b88f17b3521efb906ba732648f0457c28787368e58985</td></tr></tbody></table>

If you have an automated process for installation and you always want to be on the latest version, you can use the following link:

* <https://aikido-vm-agent.s3.eu-west-1.amazonaws.com/latest/AikidoVmScanner.msi>

The VM Scanner Agent runs once a day, at a random time between **4:00 AM - 8:00 AM** (machine time).

After install, **a first scan will start automatically**. If you want to run it on demand, you can manually execute it with the following command in a PowerShell with admin privileges.

```powershell
& 'C:\Program Files\AikidoVmScanner\AikidoVmScanner.exe'
```

### Connecting through a proxy

If your machine reaches the internet through a corporate proxy, the VM Scanner picks it up from the standard `HTTPS_PROXY`, `HTTP_PROXY`, and `NO_PROXY` environment variables. When none of these are set, it falls back to the Windows system proxy settings, including WinHTTP defaults and PAC scripts. There's no proxy option in `config.json` and nothing to enable.

The scan itself runs entirely locally. The only outbound traffic is HTTPS (port 443) to `bg.aikido.dev`, or the regional equivalent such as `bg.us.aikido.dev`.

{% hint style="info" %}
Proxy support is available from version 1.4.5 onwards.
{% endhint %}

#### Setting the proxy for the daily scan

The daily scan runs as a scheduled task under the SYSTEM account, so the variables have to be set at machine level. Proxy settings configured for your own user account aren't visible to SYSTEM. In an admin PowerShell:

```powershell
[Environment]::SetEnvironmentVariable("HTTPS_PROXY", "http://proxy.example.com:8080", "Machine")
[Environment]::SetEnvironmentVariable("NO_PROXY", "localhost,127.0.0.1", "Machine")
```

If your proxy requires authentication, include the credentials in the URL: `http://user:password@proxy.example.com:8080`.

Instead of environment variables, you can configure the machine-wide WinHTTP proxy, which the scanner also honors:

```powershell
netsh winhttp set proxy proxy-server="http=proxy.example.com:8080;https=proxy.example.com:8080" bypass-list="localhost;127.0.0.1"
```

#### Proxies that perform SSL inspection

When your proxy replaces certificates with ones signed by an internal CA, import that CA into the machine's trusted root store. In an admin PowerShell:

```powershell
Import-Certificate -FilePath "C:\path\to\corporate-ca.crt" -CertStoreLocation Cert:\LocalMachine\Root
```

#### Verifying the proxy is used

Every run records which proxy it detected in the logs at `C:\ProgramData\AikidoVmScanner\Logs`:

```
Using proxy http://proxy.example.com:8080 for https://bg.aikido.dev/api/integrations/vm_scan/uploadScanResult
```

If the proxy isn't picked up, you'll see `No proxy detected for ..., connecting directly` instead.

### Uninstall

Start PowerShell in admin mode and run:

```powershell
msiexec.exe /x (Get-WmiObject -Class Win32_Product | Where-Object { $_.Name -like "AikidoVmScanner*" }).IdentifyingNumber /qn /norestart
```

### Logs

Logs are available here, along with the last generated SBOM: `C:\ProgramData\AikidoVmScanner\Logs`
