# Device Protection in GitLab CI/CD

GitLab CI/CD jobs run on Linux runners, so Device Protection installs the same way it does on any Linux machine. What is specific to GitLab is where the install goes: once on a shell executor host that stays around, or inside each job on the Docker and Kubernetes executors.

Read [Install Device Protection for Linux](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux.md) first for the token and the install flavors. The job-level snippets below use the container flavor, which [Device Protection in Containers](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/device-protection-in-containers.md) explains.

## Docker and Kubernetes executors

This covers GitLab.com hosted runners too. Add your Device Protection token as a masked CI/CD variable named `AIKIDO_TOKEN`, then start the agent in `before_script`:

{% code title=".gitlab-ci.yml" overflow="wrap" %}

```yaml
default:
  before_script:
    - curl -fsSL https://github.com/AikidoSec/safe-chain/raw/refs/heads/main/install-scripts/install-endpoint-linux.sh | sh -s -- --container
    - aikido-endpoint-protection-start --wait-for-healthy
    - for f in /etc/profile.d/aikido-*.sh; do . "$f"; done

install:
  image: node:20-bookworm
  script:
    - npm ci
```

{% endcode %}

* **`before_script` shares a shell with `script`**, so the sourced settings reach your commands. A job with its own `before_script` replaces the default one, so repeat the three lines there.
* **The job must run as root, and the image needs `curl`.**

## Shell executor

Install the agent on the runner host with `--headless`, then reboot. See [Install Device Protection for Linux](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux.md#choose-the-install-flavor).

## Verify it works

{% code title=".gitlab-ci.yml" overflow="wrap" %}

```yaml
verify-device-protection:
  image: node:20-bookworm
  script:
    - curl -sS -o /dev/null -w "%{http_code}\n" --max-time 15 https://registry.npmjs.org/safe-chain-test/-/safe-chain-test-0.0.1-security.tgz
```

{% endcode %}

`403` means the agent blocked the malware test package. `200` means the download skipped the agent.

## Troubleshooting

| Problem                                                    | Fix                                                                                         |
| ---------------------------------------------------------- | ------------------------------------------------------------------------------------------- |
| `Token is required`                                        | The job cannot see the variable. Protected variables only reach protected branches and tags |
| `Root privileges required`                                 | The image runs as a non-root user. Use an image that runs as root                           |
| Test package returns `200`, or `SELF_SIGNED_CERT_IN_CHAIN` | The job overrides `before_script`. Add the three lines to it                                |
