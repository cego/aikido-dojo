# Deploy Device Protection with AWS Systems Manager

Deploy Aikido Device Protection to Windows and Linux machines managed by AWS Systems Manager with a State Manager association that runs the install command.

Use AWS Systems Manager (SSM) State Manager to install Aikido Device Protection on your Windows and Linux managed nodes. An association runs the install command on every node with a given tag, including nodes you tag later.

## Prerequisites

* **Managed nodes.** Each node runs SSM Agent and has permissions, either through an instance profile with the `AmazonSSMManagedInstanceCore` policy or through Default Host Management Configuration (needs IMDSv2 and SSM Agent 3.2.582.0 or later). Machines that aren't EC2 instances need a hybrid activation.
* **Permissions** to create State Manager associations and to run commands with Run Command.
* **The install command** for each operating system and [user group](/aikido-device-protection/deploying-aikido-endpoint/user-groups-for-aikido-endpoint.md) you roll out to. You copy it from the Aikido app in the first step.
* **A tag on the target nodes.** Use the key `AikidoDeviceProtection` with one value per operating system, user group, and Linux flavor, for example `windows-developers`, `linux-developers`, and `linux-servers`. A node carries one value only.
* **A supported operating system.** For Linux distributions, see [Supported ecosystems](/aikido-device-protection/miscellaneous-aikido-endpoint/supported-ecosystems.md).

## How it works

A State Manager association runs a script on every targeted node once when you create it, and then on a schedule. The script starts with a guard, so a node that already has the agent exits immediately. A node that gets the tag later runs the association right away.

Each install command carries the token of one user group. Create one association per tag value, so one per operating system, user group, and Linux flavor.

{% hint style="warning" %}
The install command contains the user group token. Anyone who can view the association, its command history, its S3 output, or the AWS CloudTrail log bucket can read the token. Restrict IAM access to State Manager, Run Command, and the CloudTrail bucket accordingly.
{% endhint %}

## Set up the association

{% stepper %}
{% step %}

### Copy the install command

1. Go to [Devices](https://app.aikido.dev/endpoint-protection/devices) in the Aikido app.
2. Click **Connect Device**.
3. Choose **Manual install**.
4. Select the [user group](/aikido-device-protection/deploying-aikido-endpoint/user-groups-for-aikido-endpoint.md) the machines should join.
5. Copy the **PowerShell install command** for Windows or the **Linux install command** for Linux.

For Linux servers and VMs, append `--headless` to the Linux install command. Leave it off for desktops. See [Install Device Protection for Linux](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux.md).
{% endstep %}

{% step %}

### Test on one node with Run Command

1. In the Systems Manager console, go to **Run Command** → **Run command**.
2. Select `AWS-RunPowerShellScript` for Windows or `AWS-RunShellScript` for Linux.
3. In **Command parameters**, paste the script from the next step.
4. In **Targets**, choose manual selection and pick one node.
5. Run the command and check the output for errors.
   {% endstep %}

{% step %}

### Create the association

In the Systems Manager console, go to **State Manager** → **Create association** and enter the values for your operating system.

{% tabs %}
{% tab title="Windows" %}

| Field                | Value                                                                               |
| -------------------- | ----------------------------------------------------------------------------------- |
| **Name**             | `AikidoDeviceProtection-windows-developers`                                         |
| **Document**         | `AWS-RunPowerShellScript`                                                           |
| **Parameters**       | **commands**: the script below                                                      |
| **Targets**          | **Specify instance tags**, key `AikidoDeviceProtection`, value `windows-developers` |
| **Specify schedule** | **On Schedule**, rate of every 1 day                                                |
| **Rate control**     | **Concurrency** `10%`, **Error threshold** `10%`                                    |

{% code title="commands" %}

```powershell
if (Get-Service -Name EndpointProtection -ErrorAction SilentlyContinue) { exit 0 }
# paste the install command from the Aikido app on the next line
```

{% endcode %}
{% endtab %}

{% tab title="Linux" %}

| Field                | Value                                                                             |
| -------------------- | --------------------------------------------------------------------------------- |
| **Name**             | `AikidoDeviceProtection-linux-developers`                                         |
| **Document**         | `AWS-RunShellScript`                                                              |
| **Parameters**       | **commands**: the script below                                                    |
| **Targets**          | **Specify instance tags**, key `AikidoDeviceProtection`, value `linux-developers` |
| **Specify schedule** | **On Schedule**, rate of every 1 day                                              |
| **Rate control**     | **Concurrency** `10%`, **Error threshold** `10%`                                  |

{% code title="commands" %}

```bash
if aikidoctl --simple >/dev/null 2>&1; then exit 0; fi

# paste the install command from the Aikido app on this line

# fail the run if the agent is missing or its service is not running
aikidoctl --simple || exit 1
systemctl is-active --quiet aikido-endpoint-protection || exit 1
```

{% endcode %}

Paste the install command on the empty line between the guard and the checks, not at the end. The last two lines make the association report **Failed** when the agent did not install or its service is not running.

For servers, create a second association with the value `linux-servers` and the `--headless` install command.
{% endtab %}
{% endtabs %}

Optionally, set **Output options** to write command output to S3. The output can contain the token, so protect the bucket.
{% endstep %}

{% step %}

### Block uninstall on Windows (optional)

Create a second association that stops users from launching the Aikido uninstaller with an AppLocker rule. It targets a separate tag, so it never runs before the install.

1. Create an association with the document `AWS-RunPowerShellScript` and the targets **Specify instance tags**, key `AikidoLockdown`, value `true`.
2. Set **commands** to the following script. It turns on the Application Identity service that AppLocker needs to enforce rules, then adds an AppLocker rule that blocks the uninstaller.
3. Add the tag `AikidoLockdown=true` to a node only after it has rebooted and shows **Active** in the [Aikido device list](https://app.aikido.dev/endpoint-protection/devices).

{% code title="aikido-applocker.ps1" %}

```powershell
$ErrorActionPreference = "Stop"
sc.exe config appidsvc start= auto
Start-Service -Name AppIDSvc
$tmp = "$env:TEMP\aikido-applocker.xml"
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/AikidoSec/safechain-internals/refs/heads/main/docs/aikido-applocker-deny.xml" -OutFile $tmp
Set-AppLockerPolicy -XmlPolicy $tmp -Merge
Remove-Item $tmp -Force
if ((Get-Service -Name AppIDSvc).Status -ne "Running") { exit 1 }
```

{% endcode %}

The first line stops the script at the first error and the last line fails the run if the Application Identity service is not running, so the association reports **Failed** instead of a silent success.

{% hint style="info" %}
If end users are local administrators, AppLocker rules can be bypassed. Combine this policy with App Control for Business (WDAC) and removal of local admin rights for stronger enforcement.
{% endhint %}
{% endstep %}

{% step %}

### Reboot or start a new login shell

The association doesn't reboot nodes. Protection becomes active at these points:

* **Windows:** the agent fully activates on the next boot.
* **Linux:** protection becomes active in a new login shell or after a reboot. See [When protection becomes active](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/linux-rollout-reference.md#when-protection-becomes-active).

Reboot nodes in a maintenance window.
{% endstep %}
{% endstepper %}

## Verify the rollout

1. In the Systems Manager console, go to **State Manager**, open the association, and check its execution history.
2. Check the agent on a node. You can run these with Run Command.

{% tabs %}
{% tab title="Windows" %}

```powershell
Get-Service -Name EndpointProtection
```

{% endtab %}

{% tab title="Linux" %}

```bash
aikidoctl --simple
systemctl is-active aikido-endpoint-protection
```

{% endtab %}
{% endtabs %}

3. Go to [Devices](https://app.aikido.dev/endpoint-protection/devices) in the Aikido app. Each device appears in the list with an **Active** status.

## Troubleshooting

| Problem                                     | What to check                                                                                                                                                                    |
| ------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| A node is missing from the targets          | Follow the AWS guide [Troubleshooting managed node availability](https://docs.aws.amazon.com/systems-manager/latest/userguide/fleet-manager-troubleshooting-managed-nodes.html). |
| The association status is **Failed**        | Open the command output for the node and read the error.                                                                                                                         |
| The device doesn't appear in the Aikido app | Reboot the node and check that it has internet access.                                                                                                                           |

### Windows

| Problem                                                      | Fix                                                                                                                                                                             |
| ------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Device doesn't appear in the dashboard after the script runs | Reboot the device; the agent fully registers on the next boot                                                                                                                   |
| Script exits with an error in the MDM console                | Run the script manually as SYSTEM on a test host to see the full error; confirm the device has internet access to download the installer                                        |
| Agent service is not running                                 | Reboot the device; if the service is still absent, re-run the install script                                                                                                    |
| Script is blocked by execution policy                        | Scripts run via MDM execute as SYSTEM and bypass user-level execution policy restrictions; confirm the MDM is not applying an additional policy that restricts script execution |
| Device shows as inactive after installation                  | The agent needs a reboot to fully activate; check back after rebooting                                                                                                          |

### Linux

| Problem                                                                                  | Fix                                                                                                                                                                                                                                        |
| ---------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| The package installs but the device never appears in the dashboard                       | `AIKIDO_TOKEN` was not set for the install command. Configuration management tools do not forward your local environment to the host, so set the variable on the task itself, then install the package again                               |
| `apt` or `dnf` looks for the package in your repositories instead of installing the file | Pass a path, not a name: `apt install ./EndpointProtection-amd64.deb`. The leading `./` is what makes the package manager treat it as a local file                                                                                         |
| The install fails on a Red Hat-family host                                               | Match the build to the major version: use the `el9` package on version 9 and the `el10` package on version 10                                                                                                                              |
| The service is not running                                                               | Run `systemctl status aikido-endpoint-protection`, then `sudo aikidoctl diagnostics` to send us the details                                                                                                                                |
| Node.js or uv still reject the Aikido certificate                                        | Installing sets `NODE_EXTRA_CA_CERTS` and `UV_SYSTEM_CERTS` system-wide, and running shells do not pick them up. Open a new login shell or reboot the device                                                                               |
| The token shows up in run output or logs                                                 | Use your tool's redaction: `no_log` in Ansible, `sensitive true` in Chef, `Sensitive()` in Puppet                                                                                                                                          |
| Several machines share one entry in the device list                                      | They booted with the same `/etc/machine-id`. Device identity on Linux follows that file, so clear it in the image you clone from and let systemd write a fresh one on first boot                                                           |
| No tray icon appears                                                                     | Expected on GNOME outside Ubuntu, and cosmetic. See [Tray Icon Support on Linux](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux/tray-icon-support-on-linux.md) |

## Related

{% content-ref url="/pages/CyKy9YG9FkQOjMYHbvLL" %}
[Linux Rollout Reference](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/linux-rollout-reference.md)
{% endcontent-ref %}

{% content-ref url="/pages/el1q3zjpdQmkekwTxga1" %}
[Deploy Device Protection with SCCM](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/windows/deploy-aikido-endpoint-with-sccm.md)
{% endcontent-ref %}

{% content-ref url="/pages/l6Jr0yFrUmsVezoR1yoV" %}
[User Groups for Device Protection](/aikido-device-protection/deploying-aikido-endpoint/user-groups-for-aikido-endpoint.md)
{% endcontent-ref %}
