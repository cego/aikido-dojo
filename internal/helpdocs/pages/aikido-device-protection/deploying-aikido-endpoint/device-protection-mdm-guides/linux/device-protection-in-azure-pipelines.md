# Device Protection in Azure Pipelines

Azure Pipelines jobs run on Linux agents, so Device Protection installs the same way it does on any Linux machine. What is specific to Azure Pipelines is where the install goes: once on a self-hosted agent that stays around, or inside each job on a Microsoft-hosted agent.

Read [Install Device Protection for Linux](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux.md) first for the token and the install flavors. The job-level snippets below use the container flavor, which [Device Protection in Containers](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/device-protection-in-containers.md) explains.

## Microsoft-hosted agents

Add your Device Protection token as a secret variable named `AIKIDO_TOKEN`, then start the agent before any step that installs packages:

{% code title="azure-pipelines.yml" overflow="wrap" %}

```yaml
pool:
  vmImage: ubuntu-latest

steps:
  - script: |
      curl -fsSL https://github.com/AikidoSec/safe-chain/raw/refs/heads/main/install-scripts/install-endpoint-linux.sh \
        | sudo AIKIDO_TOKEN="$AIKIDO_TOKEN" sh -s -- --container
      sudo aikido-endpoint-protection-start --wait-for-healthy
    displayName: Start Device Protection
    env:
      AIKIDO_TOKEN: $(AIKIDO_TOKEN)

  - script: |
      for f in /etc/profile.d/aikido-*.sh; do . "$f"; done
      npm ci
    displayName: Install dependencies
```

{% endcode %}

* **Map the secret under `env`.** Azure Pipelines never exposes secret variables to a step on its own.
* **Pass the token after `sudo`.** `sudo` drops the rest of your environment.
* **Source `/etc/profile.d/aikido-*.sh` in every step that installs packages.** Each step is a new shell, and without those settings downloads skip the agent.
* **The runner is a VM, but it takes `--container`.** A job cannot reboot the runner, and container mode needs no reboot. That is also why the agent is started explicitly. See [Choose the install flavor](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux.md#choose-the-install-flavor).

## Self-hosted agents

Install the agent on the agent host with `--headless`, then reboot. See [Install Device Protection for Linux](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux.md#choose-the-install-flavor).

## Verify it works

{% code title="azure-pipelines.yml" overflow="wrap" %}

```yaml
  - script: |
      for f in /etc/profile.d/aikido-*.sh; do . "$f"; done
      curl -sS -o /dev/null -w "%{http_code}\n" --max-time 15 https://registry.npmjs.org/safe-chain-test/-/safe-chain-test-0.0.1-security.tgz
    displayName: Verify Device Protection
```

{% endcode %}

`403` means the agent blocked the malware test package. `200` means the download skipped the agent.

## Troubleshooting

| Problem                                                    | Fix                                                                                                |
| ---------------------------------------------------------- | -------------------------------------------------------------------------------------------------- |
| `Token is required`                                        | Map the secret under the step's `env`, and pass it after `sudo`. Fork pull requests get no secrets |
| `Root privileges required`                                 | Prefix the install and the start command with `sudo`                                               |
| Test package returns `200`, or `SELF_SIGNED_CERT_IN_CHAIN` | The step did not source `/etc/profile.d/aikido-*.sh`                                               |
