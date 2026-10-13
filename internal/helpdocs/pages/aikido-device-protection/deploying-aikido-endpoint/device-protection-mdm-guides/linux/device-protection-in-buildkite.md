# Device Protection in Buildkite

Buildkite jobs run on Linux, so Device Protection installs the same way it does on any Linux machine. What is specific to Buildkite is where the install goes: once on an agent host that stays around, or inside each job on an agent that is destroyed after every build.

Read [Install Device Protection for Linux](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux.md) first for the token and the packages, and [Device Protection in Containers](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/device-protection-in-containers.md) for the container install the job-level snippets below use.

## Agents that stay around

Install the agent on the host, as you would on any build server, and every job it runs is inspected without a single change to your pipelines. Use [Deploying on Linux](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux.md) to roll it out across a fleet of agents with Ansible, cloud-init, or an image you bake yourself.

## Agents that are destroyed after every build

Ephemeral, containerized, and Buildkite hosted agents have no host to prepare, so install the agent in the job and start it before anything installs packages:

{% code title=".buildkite/pipeline.yml" overflow="wrap" %}

```yaml
agents:
  queue: linux-small

steps:
  - label: ":shield: install dependencies"
    command: |
      export AIKIDO_TOKEN=$$(buildkite-agent secret get AIKIDO_TOKEN)

      curl -fsSLO https://github.com/AikidoSec/safechain-internals/releases/latest/download/EndpointProtection-amd64.deb
      AIKIDO_CONTAINER=1 AIKIDO_CI_CD=1 apt-get install -y --no-install-recommends ./EndpointProtection-amd64.deb

      aikido-endpoint-protection-start --wait-for-healthy
      for f in /etc/profile.d/aikido-*.sh; do . "$$f"; done

      npm ci
```

{% endcode %}

* **Export the token before the install.** The package reads `AIKIDO_TOKEN` while it installs.
* **The install itself is the standard container install.** The mode variables, the start command, and the sourcing are explained in [Device Protection in Containers](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/device-protection-in-containers.md), along with the package to use on arm64 and Red Hat-family agents.

The install needs root, which Docker-based agents already have. Where the job runs as the `buildkite-agent` user, prefix the install with `sudo`.

{% hint style="info" %}
`buildkite-agent secret get` needs agent 3.106.0 or later. On that version you can also skip the `$$` escaping and let Buildkite set the variable for you:

```yaml
steps:
  - command: .buildkite/install-deps.sh
    secrets:
      AIKIDO_TOKEN: aikido_token
```

{% endhint %}

### Every job, without editing pipelines

A step that forgets the block installs packages uninspected. Move it into a `pre-command` hook instead, in the repository at `.buildkite/hooks/pre-command` or in the agent's hooks directory:

{% code title=".buildkite/hooks/pre-command" overflow="wrap" %}

```bash
#!/bin/bash

export AIKIDO_TOKEN="$(buildkite-agent secret get AIKIDO_TOKEN)"

if ! command -v aikidoctl >/dev/null 2>&1; then
  curl -fsSLO https://github.com/AikidoSec/safechain-internals/releases/latest/download/EndpointProtection-amd64.deb
  AIKIDO_CONTAINER=1 AIKIDO_CI_CD=1 apt-get install -y --no-install-recommends ./EndpointProtection-amd64.deb
  rm -f EndpointProtection-amd64.deb
fi

aikido-endpoint-protection-start --wait-for-healthy

for f in /etc/profile.d/aikido-*.sh; do
  [ -e "$f" ] && . "$f"
done
```

{% endcode %}

Hooks are scripts on disk rather than pipeline YAML, so write `$` here and not `$$`. Buildkite sources hooks, so what the hook exports reaches your command.

## Jobs that run in containers

If your steps run through the Docker or Docker Compose plugin, install the agent in the image rather than in the step. See [Device Protection in Containers](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/device-protection-in-containers.md).

{% code title=".buildkite/pipeline.yml" overflow="wrap" %}

```yaml
steps:
  - label: ":shield: install dependencies"
    command: |
      for f in /etc/profile.d/aikido-*.sh; do . "$$f"; done
      npm ci
    plugins:
      - docker#v5.14.0:
          image: "ghcr.io/your-org/my-app:latest"
          environment:
            - AIKIDO_TOKEN
```

{% endcode %}

The plugin keeps the image's `ENTRYPOINT`, so the agent starts on its own. Set the plugin's `entrypoint` or `command` option and that entrypoint is gone, so run `aikido-endpoint-protection-start --wait-for-healthy` first in your command. The Kubernetes stack replaces the container's command, so it needs the same.

## Verify it works

Add a step that downloads a package Aikido blocks:

{% code title=".buildkite/pipeline.yml" overflow="wrap" %}

```yaml
steps:
  - label: ":shield: verify Device Protection"
    command: |
      aikidoctl version
      curl -sS -o /dev/null -w "%{http_code}\n" --max-time 15 https://registry.npmjs.org/safe-chain-test/-/safe-chain-test-0.0.1-security.tgz
```

{% endcode %}

Run it in a step where the agent is already running. The `curl` prints `403`, the agent refusing the malware test package, and the agent appears in your [device list](https://app.aikido.dev/endpoint-protection/devices). For more checks, see [How to test if Device Protection is working](/aikido-device-protection/miscellaneous-aikido-endpoint/how-to-test-device-protection.md).

## Troubleshooting

| Problem                                                            | Fix                                                                                                                                                                                       |
| ------------------------------------------------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| The token is empty, and the agent never registers                  | The `$` was not escaped, so Buildkite resolved `$(buildkite-agent secret get ...)` at upload time. Use `$$` in pipeline YAML, and a single `$` in hook scripts                            |
| `buildkite-agent: command not found`, or the secret cannot be read | The agent is older than 3.106.0, the secret lives in another cluster, or the job runs in a container without the agent binary. Enable `mount-buildkite-agent` on the Docker plugin        |
| Downloads in a job are not inspected at all                        | The plugin's `entrypoint` or `command` option replaced the image's entrypoint, so the agent never started. Run `aikido-endpoint-protection-start --wait-for-healthy` first in the command |

For anything that is not Buildkite-specific, such as certificate errors from `npm` or `uv`, see the troubleshooting table on [Device Protection in Containers](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/device-protection-in-containers.md#troubleshooting).
