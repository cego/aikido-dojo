# Device Protection in GitHub Actions

GitHub Actions jobs run on Linux runners, so Device Protection installs the same way it does on any Linux machine. What is specific to GitHub Actions is where the install goes: once on a self-hosted runner that stays around, or inside each job on a hosted runner that is destroyed after the job.

Read [Install Device Protection for Linux](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux.md) first for the token and the install flavors. The job-level snippets below use the container flavor, which [Device Protection in Containers](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/device-protection-in-containers.md) explains.

## GitHub-hosted runners

Add your Device Protection token as a secret named `AIKIDO_TOKEN`, then start the agent before any step that installs packages:

{% code title=".github/workflows/build.yml" overflow="wrap" %}

```yaml
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Start Device Protection
        env:
          AIKIDO_TOKEN: ${{ secrets.AIKIDO_TOKEN }}
        run: |
          curl -fsSL https://github.com/AikidoSec/safe-chain/raw/refs/heads/main/install-scripts/install-endpoint-linux.sh \
            | sudo AIKIDO_TOKEN="$AIKIDO_TOKEN" sh -s -- --container
          sudo aikido-endpoint-protection-start --wait-for-healthy

      - uses: actions/checkout@v4

      - name: Install dependencies
        run: |
          for f in /etc/profile.d/aikido-*.sh; do . "$f"; done
          npm ci
```

{% endcode %}

* **Pass the token after `sudo`.** `sudo` drops the rest of your environment.
* **Source `/etc/profile.d/aikido-*.sh` in every step that installs packages.** Each step is a new shell, and without those settings downloads skip the agent.
* **The runner is a VM, but it takes `--container`.** A job cannot reboot the runner, and container mode needs no reboot. That is also why the agent is started explicitly. See [Choose the install flavor](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux.md#choose-the-install-flavor).

To skip the sourcing line, run every step in a login shell, which reads `/etc/profile.d/` on its own:

{% code title=".github/workflows/build.yml" overflow="wrap" %}

```yaml
defaults:
  run:
    shell: bash --login -eo pipefail {0}
```

{% endcode %}

## Job containers

Jobs with `container:` run as root, so use the same steps without `sudo`. The image needs `curl`. If the image already has the agent installed, skip the install and only start it. See [Job containers in GitHub Actions](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/device-protection-in-containers.md#job-containers-in-github-actions).

## Self-hosted runners

Install the agent on the runner host with `--headless`, then reboot. See [Install Device Protection for Linux](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux.md#choose-the-install-flavor).

## Verify it works

{% code title=".github/workflows/build.yml" overflow="wrap" %}

```yaml
      - name: Verify Device Protection
        run: |
          for f in /etc/profile.d/aikido-*.sh; do . "$f"; done
          curl -sS -o /dev/null -w "%{http_code}\n" --max-time 15 https://registry.npmjs.org/safe-chain-test/-/safe-chain-test-0.0.1-security.tgz
```

{% endcode %}

`403` means the agent blocked the malware test package. `200` means the download skipped the agent.

## Troubleshooting

| Problem                                                    | Fix                                                                                                                  |
| ---------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| `Token is required`                                        | The secret did not reach the install. Map it under `env` and pass it after `sudo`. Fork pull requests get no secrets |
| `Root privileges required`                                 | Prefix the install and the start command with `sudo`                                                                 |
| Test package returns `200`, or `SELF_SIGNED_CERT_IN_CHAIN` | The step did not source `/etc/profile.d/aikido-*.sh`                                                                 |
