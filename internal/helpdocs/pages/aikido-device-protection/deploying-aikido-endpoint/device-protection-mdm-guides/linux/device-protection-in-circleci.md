# Device Protection in CircleCI

CircleCI jobs run on Linux executors, so Device Protection installs the same way it does on any Linux machine. What is specific to CircleCI is where the install goes: once on a self-hosted runner that stays around, or inside each job on the Docker and machine executors.

Read [Install Device Protection for Linux](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux.md) first for the token and the install flavors. The job-level snippets below use the container flavor, which [Device Protection in Containers](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/device-protection-in-containers.md) explains.

## Docker and machine executors

Add your Device Protection token as a project environment variable, or to a context, named `AIKIDO_TOKEN`. Then define a reusable command and run it first in each job:

{% code title=".circleci/config.yml" overflow="wrap" %}

```yaml
version: 2.1

commands:
  start-device-protection:
    steps:
      - run:
          name: Start Device Protection
          command: |
            curl -fsSL https://github.com/AikidoSec/safe-chain/raw/refs/heads/main/install-scripts/install-endpoint-linux.sh \
              | sudo AIKIDO_TOKEN="$AIKIDO_TOKEN" sh -s -- --container
            sudo aikido-endpoint-protection-start --wait-for-healthy
            echo 'for f in /etc/profile.d/aikido-*.sh; do . "$f"; done' >> "$BASH_ENV"

jobs:
  build:
    docker:
      - image: cimg/node:lts
    steps:
      - checkout
      - start-device-protection
      - run: npm ci

workflows:
  build:
    jobs:
      - build
```

{% endcode %}

* **Pass the token after `sudo`.** `sudo` drops the rest of your environment. On an image that runs as root without `sudo`, drop both `sudo` prefixes.
* **The `$BASH_ENV` line covers every later step.** CircleCI sources that file at the start of each `run` step.
* **The runner is a VM, but it takes `--container`.** A job cannot reboot the runner, and container mode needs no reboot. That is also why the agent is started explicitly. See [Choose the install flavor](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux.md#choose-the-install-flavor).

## Self-hosted runners

Install the agent on the runner host with `--headless`, then reboot. See [Install Device Protection for Linux](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux.md#choose-the-install-flavor).

## Verify it works

{% code title=".circleci/config.yml" overflow="wrap" %}

```yaml
      - run:
          name: Verify Device Protection
          command: curl -sS -o /dev/null -w "%{http_code}\n" --max-time 15 https://registry.npmjs.org/safe-chain-test/-/safe-chain-test-0.0.1-security.tgz
```

{% endcode %}

`403` means the agent blocked the malware test package. `200` means the download skipped the agent.

## Troubleshooting

| Problem                                                    | Fix                                                                                                                   |
| ---------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------- |
| `Token is required`                                        | The job cannot see the variable. Check the job uses the right context. Fork pull requests get no variables by default |
| `Root privileges required`                                 | Prefix the install and the start command with `sudo`                                                                  |
| Test package returns `200`, or `SELF_SIGNED_CERT_IN_CHAIN` | The step does not run in bash, so it skips `$BASH_ENV`. Source `/etc/profile.d/aikido-*.sh` in that step              |
