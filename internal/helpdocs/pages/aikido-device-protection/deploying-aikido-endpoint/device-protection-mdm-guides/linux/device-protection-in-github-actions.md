# Device Protection in GitHub Actions

GitHub Actions jobs run on Linux runners, so Device Protection installs the same way it does on any Linux machine. What is specific to GitHub Actions is where the install goes: once on a self-hosted runner that stays around, or inside each job on a hosted runner that is destroyed after the job.

Read [Install Device Protection for Linux](https://help.aikido.dev/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux) first for the token.

## GitHub-hosted runners

Add your Device Protection token as a secret named `AIKIDO_TOKEN`, then add the Aikido action as the first step of every job that installs packages:

{% code title=".github/workflows/build.yml" overflow="wrap" %}
```yaml
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: AikidoSec/safechain-internals/github-action-setup@action-v1
        with:
          token: ${{ secrets.AIKIDO_TOKEN }}

      - uses: actions/checkout@v4

      - name: Install dependencies
        run: npm ci
```
{% endcode %}

* **Add it to every job.** Each job runs on its own machine, and only the steps after the action are covered.
* **No extra setup in later steps.** The action installs and starts the agent, and exports the proxy and certificate settings to every later step, so you do not need to source `/etc/profile.d/aikido-*.sh`.
* **Fork pull requests get no secrets**, so the action fails there with an empty token.

| Input        | Required | Default | Description                                                                 |
| ------------ | -------- | ------- | --------------------------------------------------------------------------- |
| `token`      | Yes      |         | Your Device Protection token. Pass it from a secret.                        |
| `self-check` | No       | `false` | When `true`, the step fails unless a known test package is blocked afterwards. |

The action supports GitHub-hosted Linux runners (x64 and arm64) and uses passwordless `sudo`. It always installs the latest agent release.

## Job containers

Jobs with `container:` run as root. Use the same action step. The image must be Debian/Ubuntu or Enterprise Linux 9/10 and have `bash` and `curl`. If the image already has the agent installed, skip the action and only start the agent. See [Job containers in GitHub Actions](https://help.aikido.dev/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/device-protection-in-containers#job-containers-in-github-actions).

## Self-hosted runners

Install the agent on the runner host with `--headless`, then reboot. See [Install Device Protection for Linux](https://help.aikido.dev/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux#choose-the-install-flavor).

## Verify it works

Set `self-check: true` on the action to fail the job unless the agent blocks a known test package. You can also check manually:

{% code title=".github/workflows/build.yml" overflow="wrap" %}
```yaml
      - name: Verify Device Protection
        run: curl -sS -o /dev/null -w "%{http_code}\n" --max-time 15 https://registry.npmjs.org/safe-chain-test/-/safe-chain-test-0.0.1-security.tgz
```
{% endcode %}

`403` means the agent blocked the malware test package. `200` means the download skipped the agent.

## Troubleshooting

| Problem                                                    | Fix                                                                                              |
| ---------------------------------------------------------- | ------------------------------------------------------------------------------------------------ |
| Token error from the action                                | The secret did not reach the step. Check the secret name. Fork pull requests get no secrets      |
| Test package returns `200`, or `SELF_SIGNED_CERT_IN_CHAIN` | The package install ran before the action step, or in a different job. Add the action to that job first |
