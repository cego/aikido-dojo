# Device Protection in Dev Containers

Add Device Protection to a Dev Container as a [Feature](https://containers.dev/implementors/features/), so any `devcontainer.json` can add it to any base image without a custom Dockerfile.

Read [Device Protection in Containers](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/device-protection-in-containers.md) first. The Feature installs the same package in the same `--container` mode, so the entrypoint, the Aikido CA variables, and the token rules all work the same way.

## Add the Feature

Create this alongside your project, for example at `.devcontainer/features/device-protection/`:

{% code title=".devcontainer/features/device-protection/devcontainer-feature.json" overflow="wrap" %}

```json
{
  "id": "device-protection",
  "version": "1.0.0",
  "name": "Aikido Device Protection",
  "description": "Installs Aikido Device Protection in --container mode so package manager traffic inside the container is inspected for malware.",
  "entrypoint": "aikido-endpoint-protection-container-entrypoint"
}
```

{% endcode %}

`entrypoint` runs the agent's supervisor alongside the image's own entrypoint, instead of replacing it.

{% code title=".devcontainer/features/device-protection/install.sh" overflow="wrap" %}

```bash
#!/usr/bin/env bash
set -euo pipefail

case "$(uname -m)" in
  x86_64) ARCH=amd64 ;;
  aarch64) ARCH=arm64 ;;
  *) echo "Aikido Device Protection: unsupported architecture $(uname -m)" >&2; exit 1 ;;
esac

if command -v apt-get >/dev/null 2>&1; then
  apt-get update
  apt-get install -y --no-install-recommends curl ca-certificates
  curl -fLO "https://github.com/AikidoSec/safechain-internals/releases/latest/download/EndpointProtection-${ARCH}.deb"
  AIKIDO_CONTAINER=1 AIKIDO_CI_CD=1 apt-get install -y --no-install-recommends "./EndpointProtection-${ARCH}.deb"
  rm -f "EndpointProtection-${ARCH}.deb"
  rm -rf /var/lib/apt/lists/*
elif command -v dnf >/dev/null 2>&1; then
  . /etc/os-release
  EL="el${VERSION_ID%%.*}"
  dnf install -y --setopt=install_weak_deps=False curl ca-certificates
  curl -fLO "https://github.com/AikidoSec/safechain-internals/releases/latest/download/EndpointProtection-${ARCH}.${EL}.rpm"
  AIKIDO_CONTAINER=1 AIKIDO_CI_CD=1 dnf install -y --setopt=install_weak_deps=False "./EndpointProtection-${ARCH}.${EL}.rpm"
  rm -f "EndpointProtection-${ARCH}.${EL}.rpm"
else
  echo "Aikido Device Protection: no apt-get or dnf found, unsupported base image" >&2
  exit 1
fi
```

{% endcode %}

This is the package-file flavor from Device Protection in Containers, so no token is baked into the image. It picks the `.deb` or `.rpm` that matches the base image's package manager and architecture, the same files the manual Linux install uses.

{% hint style="warning" %}
Don't add `token` as a Feature option. Options are written into a `devcontainer-features.env` file that install scripts source during the image build, which means a token passed that way can end up in a build layer. Pass `AIKIDO_TOKEN` at container start instead, through `containerEnv` below.
{% endhint %}

## Wire it up in devcontainer.json

{% code title=".devcontainer/devcontainer.json" overflow="wrap" %}

```json
{
  "name": "my-app",
  "image": "mcr.microsoft.com/devcontainers/javascript-node:20",
  "features": {
    "./features/device-protection": {}
  },
  "containerEnv": {
    "AIKIDO_TOKEN": "${localEnv:AIKIDO_TOKEN}"
  },
  "postCreateCommand": "echo 'for f in /etc/profile.d/aikido-*.sh; do [ -e \"$f\" ] && . \"$f\"; done' >> ~/.bashrc"
}
```

{% endcode %}

* `features` references the Feature by its local path, relative to `devcontainer.json`.
* `containerEnv` passes the token when the container starts. Set `AIKIDO_TOKEN` in your shell (or a Codespaces secret) before opening the Dev Container.
* `postCreateCommand` appends the CA sourcing to `~/.bashrc`, so terminals opened later inside the container (VS Code's integrated terminal, Codespaces) pick up the Aikido CA instead of failing installs with `SELF_SIGNED_CERT_IN_CHAIN`.

If your Dev Container's default shell is zsh instead of bash, append to `~/.zshrc` instead:

{% tabs %}
{% tab title="bash" %}

```json
"postCreateCommand": "echo 'for f in /etc/profile.d/aikido-*.sh; do [ -e \"$f\" ] && . \"$f\"; done' >> ~/.bashrc"
```

{% endtab %}

{% tab title="zsh" %}

```json
"postCreateCommand": "echo 'for f in /etc/profile.d/aikido-*.sh; do [ -e \"$f\" ] && . \"$f\"; done' >> ~/.zshrc"
```

{% endtab %}
{% endtabs %}

If you already build a custom `Dockerfile` for other reasons, `features` still applies on top of it. You no longer need to add the agent's install or `ENTRYPOINT` lines to that `Dockerfile` yourself, the Feature does it.

To reuse the Feature across repositories instead of vendoring the folder into each one, publish it once to a registry your team controls with `devcontainer features publish`, then reference it as `ghcr.io/<your-org>/devcontainer-features/device-protection:1` instead of by path. See the [Dev Containers CLI docs](https://containers.dev/implementors/features-distribution/) for the publish command.

## GitHub Codespaces secrets

For Codespaces, add `AIKIDO_TOKEN` as a repository or organization Codespaces secret rather than committing it. `${localEnv:AIKIDO_TOKEN}` picks up a secret exposed as an environment variable the same way it picks up a local shell variable.

## Verify it works

Open a new terminal in the Dev Container (a fresh one, so it picks up the `postCreateCommand` change) and run:

{% code overflow="wrap" %}

```bash
curl -sS -o /dev/null -w "%{http_code}\n" --max-time 15 https://registry.npmjs.org/safe-chain-test/-/safe-chain-test-0.0.1-security.tgz
aikidoctl version
```

{% endcode %}

The `curl` prints `403`, the agent refusing the malware test package, and `aikidoctl version` reports the installed version and the agent's state.

## Troubleshooting

The entrypoint and install steps are the same as any Linux container, so start with the troubleshooting table on Device Protection in Containers. Two additions specific to Dev Containers:

| Problem                                                                                                                          | Fix                                                                                                                                                                                                              |
| -------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Package installs fail with `SELF_SIGNED_CERT_IN_CHAIN` in the integrated terminal, but work when run through `postCreateCommand` | The terminal was open before `postCreateCommand` finished, or it's a different shell than the one it appended to. Open a new terminal, and check `~/.bashrc` (or `~/.zshrc`) for the `aikido-*.sh` sourcing line |
| `no apt-get or dnf found, unsupported base image`                                                                                | The base image uses a package manager the Feature doesn't handle yet, for example Alpine's `apk`. Add a branch for it or install the agent manually as shown on Device Protection in Containers                  |
