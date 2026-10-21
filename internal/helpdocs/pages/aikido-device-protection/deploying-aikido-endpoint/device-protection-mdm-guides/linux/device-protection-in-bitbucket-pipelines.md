# Device Protection in Bitbucket Pipelines

Bitbucket Pipelines steps run on Linux, so Device Protection installs the same way it does on any Linux machine. What is specific to Bitbucket is where the install goes: once on a Linux Shell runner that stays around, or inside each step on Bitbucket Cloud and Docker runners.

Read [Install Device Protection for Linux](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux.md) first for the token and the install flavors. The job-level snippets below use the container flavor, which [Device Protection in Containers](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/device-protection-in-containers.md) explains.

## Bitbucket Cloud and Docker runners

Add your Device Protection token as a secured repository variable named `AIKIDO_TOKEN`. Then add this script to your repository:

{% code title="ci/start-device-protection.sh" overflow="wrap" %}

```bash
curl -fsSL https://github.com/AikidoSec/safe-chain/raw/refs/heads/main/install-scripts/install-endpoint-linux.sh | sh -s -- --container
aikido-endpoint-protection-start --wait-for-healthy
for f in /etc/profile.d/aikido-*.sh; do . "$f"; done
```

{% endcode %}

Source it at the start of each step that installs packages:

{% code title="bitbucket-pipelines.yml" overflow="wrap" %}

```yaml
image: node:20-bookworm

pipelines:
  default:
    - step:
        name: Install dependencies
        script:
          - . ci/start-device-protection.sh
          - npm ci
```

{% endcode %}

* **Source it with `.`, don't run it with `sh`.** Sourcing keeps the settings for the rest of the step.
* **Every step needs it.** Each step is a new container.
* **The step must run as root, and the image needs `curl`.**

## Linux Shell runners

Install the agent on the runner host with `--headless`, then reboot. See [Install Device Protection for Linux](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux.md#choose-the-install-flavor).

## Verify it works

{% code title="bitbucket-pipelines.yml" overflow="wrap" %}

```yaml
    - step:
        name: Verify Device Protection
        script:
          - . ci/start-device-protection.sh
          - curl -sS -o /dev/null -w "%{http_code}\n" --max-time 15 https://registry.npmjs.org/safe-chain-test/-/safe-chain-test-0.0.1-security.tgz
```

{% endcode %}

`403` means the agent blocked the malware test package. `200` means the download skipped the agent.

## Troubleshooting

| Problem                                                    | Fix                                                          |
| ---------------------------------------------------------- | ------------------------------------------------------------ |
| `Token is required`                                        | The step cannot see the variable. Check its name             |
| `Root privileges required`                                 | The image runs as a non-root user. Use one that runs as root |
| Test package returns `200`, or `SELF_SIGNED_CERT_IN_CHAIN` | The script ran with `sh` or `bash`. Source it with `.`       |
