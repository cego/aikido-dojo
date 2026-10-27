# Setting Up the Local VM Scanner on Linux

Aikido VM Scanner is a single-package that installs on your system, automatically scanning and identifying dependencies to provide a detailed view into your environment.

### Prerequisites

* Minimum system requirements: at least 1GB RAM.
* Preferred system requirements: at least 2GB RAM and 4 CPUs.
* Ensure you have sudo / admin privileges on your system
* Make sure to use the appropriate commands for your system or cloud provider
* If you need to run with root, place the AIKIDO\_TOKEN env var after sudo, like this:\
  `sudo AIKIDO_TOKEN=REPLACE_ME <install_command>`

### Installation and Upgrade

Make sure you run as sudo and replace `AIKIDO_TOKEN` with [valid token from Local VM scanning page in Aikido](https://app.aikido.dev/settings/integrations/vm-scan-agent). To specify the VM environment (that you will see later in Aikido), you can also set the `VM_TYPE` variable as one of: `production`, `staging` or `development`.

The VM Scanner Agent runs once a day, at a random time between 4:00 AM - 8:00 AM (machine time).

After install, **a first scan will start automatically**. If you want to run it on demand, you can manually execute:

```bash
/opt/aikido-vm-scanner-1.5.0/aikido-vm-scanner
```

#### **For Red Hat-based Systems (RHEL, CentOS, Fedora)**

**x86\_64**

```
AIKIDO_TOKEN=REPLACE_ME VM_TYPE=production dnf install -y https://aikido-vm-agent.s3.eu-west-1.amazonaws.com/v1.5.0/aikido-vm-scanner.x86_64.rpm
```

<table><thead><tr><th width="100.1015625">Hash</th><th>aikido-vm-scanner.x86_64.rpm</th></tr></thead><tbody><tr><td>MD5</td><td>227b56a1863b0b771bfeaa465c59f72e</td></tr><tr><td>SHA256</td><td>57ea588e9b0bc80bd9265b08101e16a833117e5420141c090be7e6a7df46c3a8</td></tr></tbody></table>

**aarch64**

```
AIKIDO_TOKEN=REPLACE_ME VM_TYPE=production dnf install -y https://aikido-vm-agent.s3.eu-west-1.amazonaws.com/v1.5.0/aikido-vm-scanner.aarch64.rpm
```

<table><thead><tr><th width="99.953125">Hash</th><th>aikido-vm-scanner.aarch64.rpm</th></tr></thead><tbody><tr><td>MD5</td><td>8222422733861e37061f5c9342d7db94</td></tr><tr><td>SHA256</td><td>3df217c38330b94c7c684c32d5ffeac0551ce1b45e5265b355a5c2a64942e7dd</td></tr></tbody></table>

#### **For Debian-based Systems (Debian, Ubuntu)**

{% hint style="warning" %}
We support Debian >= 10 and Ubuntu >= 20.04.
{% endhint %}

**x86\_64**

```
curl -L -O https://aikido-vm-agent.s3.eu-west-1.amazonaws.com/v1.5.0/aikido-vm-scanner.amd64.deb
AIKIDO_TOKEN=REPLACE_ME VM_TYPE=production apt-get install -y ./aikido-vm-scanner.amd64.deb
```

<table><thead><tr><th width="100.36328125">Hash</th><th>aikido-vm-scanner.amd64.deb</th></tr></thead><tbody><tr><td>MD5</td><td>818941bd3da7506ca5269872cd00ffbf</td></tr><tr><td>SHA256</td><td>222c1fa58bd841767707b720df130d1f40b84c00a02e12e8b7ddcd87f92b1286</td></tr></tbody></table>

**aarch64**

```
curl -L -O https://aikido-vm-agent.s3.eu-west-1.amazonaws.com/v1.5.0/aikido-vm-scanner.arm64.deb
AIKIDO_TOKEN=REPLACE_ME VM_TYPE=production apt-get install -y ./aikido-vm-scanner.arm64.deb
```

<table><thead><tr><th width="100.19921875">Hash</th><th>aikido-vm-scanner.arm64.deb</th></tr></thead><tbody><tr><td>MD5</td><td>9d26a64ddc625b11389ee398114abb4c</td></tr><tr><td>SHA256</td><td>f99cbf4acda7acc916d80a866c830c84ceae909804ecaaf7bb1a40c0edd87ae5</td></tr></tbody></table>

### Latest version

If you have an automated install process and you always want to be on the latest version as soon as we release it, you can replace the version in the install link with `latest`:

* <https://aikido-vm-agent.s3.eu-west-1.amazonaws.com/latest/aikido-vm-scanner.x86\\_64.rpm>
* <https://aikido-vm-agent.s3.eu-west-1.amazonaws.com/latest/aikido-vm-scanner.aarch64.rpm>
* <https://aikido-vm-agent.s3.eu-west-1.amazonaws.com/latest/aikido-vm-scanner.amd64.deb>
* <https://aikido-vm-agent.s3.eu-west-1.amazonaws.com/latest/aikido-vm-scanner.arm64.deb>

### Additional Configuration

#### Token setup

If for any reason you can't set the AIKIDO\_TOKEN at install time, you can set the token post install in one of two ways:

* Paste the token in `/opt/aikido-vm-scanner-1.5.0/.token`
* Change the contents of `/opt/aikido-vm-scanner-1.5.0/config.json` :

```
{
    "api_key": "AIKIDO_TOKEN_HERE",
    "exclude": []
}
```

#### Hostname change

By default, we automatically get the hostname for the scanned machine and submit that to Aikido, in order to be displayed in the Virtual Machines tab.

If you want to change the reported hostname, you can do that using the configuration file:

```
{
    "hostname": "YOUR_HOSTNAME_HERE",
    "exclude": []
}
```

#### Exclude files or paths

In the config.json you can exclude files and paths by adding additional items to the exclude list. You can find some examples below.

* Exclude a single file in a directory:

```
{
    "exclude": [
        ...
        "./home/testing/app/package.json",
    ]
}
```

* Exclude all releases subdirectories:

```
{
    "exclude": [
            ...
            "./home/testing/*/releases",
    ]
}
```

* Exclude all .json files in the out folder and all subdirectories:

```
{
    "exclude": [
        ...
        "./out/**/*.json",
    ]
}
```

* Exclude catalogers from scanning (eg: exclude Golang catalogers):

```
{
    "exclude": [
    ],
    "select_catalogers": "-go"
}
```

#### Output channel

If you want to control the output channel of the VM scanner, when installing you can specify the OUTPUT variable as `stdout`, `stderr` or `none`.

Example for rpm x86\_64:

```
AIKIDO_TOKEN=REPLACE_ME VM_TYPE=production OUTPUT=stderr dnf install -y https://aikido-vm-agent.s3.eu-west-1.amazonaws.com/v1.5.0/aikido-vm-scanner.x86_64.rpm
```

{% hint style="info" %}
This option does not affect the disk logs. These are created no matter what this setting is.
{% endhint %}

#### Disable initial scan

When a rpm/deb package is installed, it automatically triggers an initial scan.\
If you want to disable this feature, you can set the `INITIAL_SCAN` parameter to 0 at install time, like this:

```
AIKIDO_TOKEN=REPLACE_ME INITIAL_SCAN=0 dnf install -y https://aikido-vm-agent.s3.eu-west-1.amazonaws.com/v1.5.0/aikido-vm-scanner.x86_64.rpm
```

#### CLI parameters

If for any reason you need to start the scanning on demand via the command line, you can specify the following CLI parameters to be used for that scan:

```
  --api-key string
        API key for authentication
  --hostname string
        Hostname of the machine
  --select-catalogers string
        Comma-separated list of catalogers to use or remove
  --vm-type string
        Type of the virtual machine
```

These CLI parameters take precedence over those specified in config.json.

{% hint style="warning" %}
We trigger an automatic scan once per day and those scans will still use the configuration file and not the CLI parameters, so it is always preferred to place these in config.json.
{% endhint %}

* Example for setting a custom hostname via CLI:

```
/opt/aikido-vm-scanner-1.5.0/aikido-vm-scanner --hostname "My Custom VM"
```

* Example for excluding Golang catalogers via CLI:

```
/opt/aikido-vm-scanner-1.5.0/aikido-vm-scanner --select-catalogers "-go"
```

### Connecting through a proxy

If your machine reaches the internet through a corporate proxy, the VM Scanner picks it up from the standard `HTTPS_PROXY`, `HTTP_PROXY`, and `NO_PROXY` environment variables. There's no proxy option in `config.json` and nothing to enable, the scanner uses these variables automatically when they're set.

The scan itself runs entirely locally. The only outbound traffic is HTTPS (port 443) to `bg.aikido.dev`, or the regional equivalent such as `bg.us.aikido.dev`.

{% hint style="info" %}
Proxy support is available from version 1.4.5 onwards.
{% endhint %}

#### Setting the proxy for the daily scan

The daily scan runs from root's crontab, and cron jobs start with a minimal environment, so variables you export in a shell session never reach them. Declaring the variables in the crontab itself works on every distribution:

1. Run `sudo crontab -e`
2. Add the variables on their own lines, above the existing `aikido-vm-scanner` line

```
HTTPS_PROXY=http://proxy.example.com:8080
NO_PROXY=localhost,127.0.0.1
37 5 * * * /opt/aikido-vm-scanner-1.5.0/aikido-vm-scanner
```

If your proxy requires authentication, include the credentials in the URL: `http://user:password@proxy.example.com:8080`.

{% hint style="warning" %}
Variables declared at the top of a crontab apply to every job in that crontab, not just the scanner.
{% endhint %}

Keep the variables on separate lines rather than inlining them on the scanner line. Upgrading the package rewrites the scanner's cron line, but leaves your variable lines untouched.

To test the proxy immediately, run a scan on demand:

```bash
sudo HTTPS_PROXY=http://proxy.example.com:8080 /opt/aikido-vm-scanner-1.5.0/aikido-vm-scanner
```

#### Proxies that perform SSL inspection

When your proxy replaces certificates with ones signed by an internal CA, add that CA to the system trust store:

* **Red Hat-based systems**: copy the PEM file into `/etc/pki/ca-trust/source/anchors/` and run `update-ca-trust`
* **Debian-based systems**: copy the certificate as a `.crt` file into `/usr/local/share/ca-certificates/` and run `update-ca-certificates`

Alternatively, point the scanner at a PEM bundle by adding `SSL_CERT_FILE=/path/to/corporate-ca.pem` to the same crontab.

#### Verifying the proxy is used

Every run records which proxy it detected in the logs at `/var/log/aikido-vm-scanner-1.5.0/`:

```
Using proxy http://proxy.example.com:8080 for https://bg.aikido.dev/api/integrations/vm_scan/uploadScanResult
```

If the proxy isn't picked up, you'll see `No proxy detected for ..., connecting directly` instead.

### Uninstall

#### Manual uninstall

**For Red Hat-based Systems (RHEL, CentOS, Fedora)**

```
dnf remove -y aikido-vm-scanner
```

**For Debian-based Systems (Debian, Ubuntu)**

```
apt-get remove -y aikido-vm-scanner
```

### Logs

Logs are available here, along with the last generated SBOM: `/var/log/aikido-vm-scanner-1.5.0`
