# VM Patching

Aikido Images also covers your virtual machines (VMs). Aikido builds patched OS packages and binaries, and you install them on your VMs from the Aikido package repository.

## How it works

Aikido decides which patches to build from the CVEs detected in the images in the Aikido Images catalog. Each patch is available as a package in the Aikido package repository, and as a binary you can install on your VM.

* **SLA:** The same SLA as for all Aikido Images. Aikido patches each Critical and High CVE within 7 days.
* **Scope:** Patches cover OS packages. Kernel patches are not included and are not currently planned.

## Install patched packages

The install script is specific to your workspace, distribution, and version. Copy it from the Aikido app. The app shows it on the Aikido Libraries page, next to the install instructions for libraries.

1. Go to **AutoFix** → **Aikido Libraries** and click **Install instructions**.
2. Click **OS Packages**, then **Virtual Machine**.
3. Select your distribution and **Version**, copy the script, and run it as root on your VM.

After the script runs, install or upgrade packages as usual with your package manager.

{% hint style="warning" %}
The repository URL in the script contains a token specific to your workspace. Don't share the script or commit it to a public repository.
{% endhint %}

## Related

{% content-ref url="README.md" %}
[README.md](https://help.aikido.dev/autofix-and-remediation/aikido-images)
{% endcontent-ref %}

{% content-ref url="aikido-docker-token.md" %}
[aikido-docker-token.md](https://help.aikido.dev/autofix-and-remediation/aikido-images/aikido-docker-token)
{% endcontent-ref %}
