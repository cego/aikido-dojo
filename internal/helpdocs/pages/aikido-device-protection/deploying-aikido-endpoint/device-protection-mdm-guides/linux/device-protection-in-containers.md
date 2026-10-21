# Device Protection in Containers

Device Protection can run inside a Linux container image, so package installs in CI jobs, build containers, and containerized dev environments get inspected the same way they do on a workstation. You install the agent while the image is built and let the agent's container entrypoint start it when the container runs.

Read [Install Device Protection for Linux](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux.md) first for the token and the packages that every method uses. The [Linux Rollout Reference](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/linux-rollout-reference.md) covers token handling and device identity across all of them.

{% hint style="info" %}
Container images always install with `--container`. Without it the agent waits for a desktop user to finish setup, its proxy never starts, and the container exits before your command runs. See [What `--container` changes](#what-container-changes).

`--container` is for the agent running *inside* a container, or on a hosted CI runner that is thrown away after the job. A Docker, Kubernetes, or Jenkins host that builds and runs those containers is a server, so install it there with `--headless` instead. See [Install Device Protection for Linux](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux.md#choose-the-install-flavor).
{% endhint %}

## Add Device Protection to an image

The build has to run as root, and the agent's own installer expects `curl` and `ca-certificates`, which Debian slim images leave out.

{% tabs %}
{% tab title="Install script" %}
The install script picks the package for the image's distribution and architecture, verifies its checksum, and installs it. It needs the token at build time.

{% code title="Dockerfile" overflow="wrap" %}

```dockerfile
FROM node:20-bookworm-slim

ARG AIKIDO_TOKEN

RUN apt-get update \
 && apt-get install -y --no-install-recommends curl ca-certificates \
 && curl -fsSL https://github.com/AikidoSec/safe-chain/raw/refs/heads/main/install-scripts/install-endpoint-linux.sh \
      | sh -s -- --token "$AIKIDO_TOKEN" --container \
 && rm -rf /var/lib/apt/lists/*

STOPSIGNAL SIGTERM
ENTRYPOINT ["aikido-endpoint-protection-container-entrypoint"]

# Your own command. This one downloads a package Aikido blocks, so it prints
# 403 while protection is working.
CMD ["curl", "-sS", "-o", "/dev/null", "-w", "%{http_code}\n", "--max-time", "15", \
     "https://registry.npmjs.org/safe-chain-test/-/safe-chain-test-0.0.1-security.tgz"]
```

{% endcode %}

Then build and run it:

{% code overflow="wrap" %}

```bash
docker build -t my-app --build-arg AIKIDO_TOKEN=<your-token> .
docker run --rm --init my-app
```

{% endcode %}

`--build-arg AIKIDO_TOKEN=<your-token>` is not optional: the build cannot see your shell environment, and the installer stops with an error when the token is empty. To pass a variable you already have, give Docker the bare name instead, `--build-arg AIKIDO_TOKEN`.
{% endtab %}

{% tab title="Package file" %}
Installing the package yourself lets you pin a version, install from an internal mirror, and build an image that has no token in it at all. Set the container-mode variables in the environment of the package manager, and set `AIKIDO_TOKEN` only if you want the token baked in.

{% code title="Dockerfile" overflow="wrap" %}

```dockerfile
FROM node:20-bookworm-slim

RUN apt-get update \
 && apt-get install -y --no-install-recommends curl ca-certificates \
 && curl -fLO https://github.com/AikidoSec/safechain-internals/releases/latest/download/EndpointProtection-amd64.deb \
 && AIKIDO_CONTAINER=1 AIKIDO_CI_CD=1 apt-get install -y --no-install-recommends ./EndpointProtection-amd64.deb \
 && rm -f EndpointProtection-amd64.deb \
 && rm -rf /var/lib/apt/lists/*

STOPSIGNAL SIGTERM
ENTRYPOINT ["aikido-endpoint-protection-container-entrypoint"]

CMD ["curl", "-sS", "-o", "/dev/null", "-w", "%{http_code}\n", "--max-time", "15", \
     "https://registry.npmjs.org/safe-chain-test/-/safe-chain-test-0.0.1-security.tgz"]
```

{% endcode %}

Then build the image and pass the token when you run it:

{% code overflow="wrap" %}

```bash
docker build -t my-app .
docker run --rm --init -e AIKIDO_TOKEN=<your-token> my-app
```

{% endcode %}

Set both variable names. `AIKIDO_CONTAINER` is the current one, and `AIKIDO_CI_CD` is the older name that packages up to 1.8.2 read, so setting both works whichever package you install.

Use the package that matches the image, not your laptop: `EndpointProtection-arm64.deb` for arm64 images, and the `el9` or `el10` RPM for Red Hat-family base images. [Install Device Protection for Linux](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux.md#manual-install) lists all of them.
{% endtab %}
{% endtabs %}

## Where the token comes from

The agent takes the token at install time or at run time, and the difference is where the token ends up.

| Passed as                                 | Ends up in                                                                           |
| ----------------------------------------- | ------------------------------------------------------------------------------------ |
| `--build-arg AIKIDO_TOKEN=<your-token>`   | The image, and its build history. Anyone who can pull the image can read it          |
| `docker run -e AIKIDO_TOKEN=<your-token>` | Only the running container. The entrypoint hands the token to the agent as it starts |

Prefer the run-time token for images you push to a registry, and keep the build-time token for images that never leave the machine that built them. BuildKit warns about the build-time form (`SecretsUsedInArgOrEnv`) for exactly this reason.

A token passed at run time only applies to an image that was built without one. Once a token is baked in, that is the one the agent uses.

## Point Node.js and uv at the Aikido CA

The agent inspects HTTPS with a local CA. Tools that read the system trust store pick it up on their own, but Node.js and uv carry their own, so the agent writes the variables that point them at the CA into `/etc/profile.d/`. The entrypoint sources only the proxy settings from there, so a `CMD` that runs `npm` fails with `SELF_SIGNED_CERT_IN_CHAIN` until you source the rest.

Source every file the agent drops before your command runs:

{% code title="Dockerfile" overflow="wrap" %}

```dockerfile
CMD ["sh", "-c", "for f in /etc/profile.d/aikido-*.sh; do . \"$f\"; done; exec npm ci"]
```

{% endcode %}

The agent writes those files while it starts, so read them at run time rather than copying the paths into `ENV` at build time.

## Inspect installs during the image build

The agent only runs while the container runs, so a `RUN npm ci` in your Dockerfile downloads packages uninspected. To cover the build too, start the agent inside the `RUN` step:

{% code title="Dockerfile" overflow="wrap" %}

```dockerfile
RUN aikido-endpoint-protection-start --wait-for-healthy \
 && for f in /etc/profile.d/aikido-*.sh; do . "$f"; done \
 && npm ci
```

{% endcode %}

The agent needs a token while it runs, so this only works in an image built with the token baked in.

## Job containers in GitHub Actions

`jobs.<id>.container` replaces the image's `ENTRYPOINT`, so the entrypoint never runs. Start the agent in the job's first step instead:

{% code title=".github/workflows/build.yml" overflow="wrap" %}

```yaml
jobs:
  build:
    runs-on: ubuntu-latest
    container:
      image: ghcr.io/your-org/my-app:latest
      env:
        AIKIDO_TOKEN: ${{ secrets.AIKIDO_TOKEN }}
    steps:
      - name: Start Device Protection
        run: aikido-endpoint-protection-start --wait-for-healthy

      - name: Install dependencies
        run: |
          for f in /etc/profile.d/aikido-*.sh; do . "$f"; done
          npm ci
```

{% endcode %}

Each step runs in its own shell, so source `/etc/profile.d/aikido-*.sh` in every step that installs packages.

## How container mode works

You don't need this to get an image working, but it explains what the flag and the entrypoint are doing.

### What `--container` changes

`--container` (or the `AIKIDO_CONTAINER=1` and `AIKIDO_CI_CD=1` variables when you install the package yourself) is what makes the agent work in a container:

* **The proxy keeps its CA in memory** instead of the OS keyring, which containers and CI runners do not have.
* **Traffic goes through an explicit proxy** on `127.0.0.1:7654` rather than the transparent eBPF datapath used on workstations. That means the container needs no `--privileged` and no added capabilities.

`--ci-cd` was the earlier name for this flag. The install script still accepts it, so existing Dockerfiles keep working, but use `--container` in anything you write now.

### What the entrypoint does

`aikido-endpoint-protection-container-entrypoint` is a small supervisor, not a wrapper that hands off and exits:

1. Stages the token from `AIKIDO_TOKEN` when one is set in the container's environment.
2. Starts the agent without systemd and waits up to 20 seconds for it to report healthy. If it never does, the container exits with an error instead of running your command unprotected.
3. Exports the proxy environment variables and runs your `CMD` as a child process.
4. On `SIGTERM` it stops the agent, then your command, and exits `143`. It also handles `SIGRTMIN+3`, which systemd-based images send on `docker stop`.
5. Exits with your command's exit code.

Because it stays PID 1 to receive that signal, it does not reap orphaned processes. `--init` (or `tini`) covers that, and is worth adding.

If your `CMD` is systemd itself, the entrypoint execs it and leaves the unit to start the agent.

## Verify it works

Run the image with the `CMD` above. A working container prints:

```
Aikido Endpoint Protection is healthy
403
```

The `403` is the agent refusing the malware test package. Swap the URL for a package that is not blocked and the same command prints `200`, which tells you the proxy is in the path rather than simply breaking downloads.

Inside a running container, `aikidoctl version` reports the installed version and the agent's state.

## Troubleshooting

| Problem                                                                                  | Fix                                                                                                                                                                                                           |
| ---------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `exec aikido-endpoint-protection-container-entrypoint failed: No such file or directory` | The image has an agent that predates container support. Install a current package                                                                                                                             |
| `Error: daemon started but never became healthy`                                         | The container has no token, or the agent was installed without `--container`. Pass a token with `-e AIKIDO_TOKEN=<your-token>`, and check that the image installed the agent in container mode                |
| `npm error SELF_SIGNED_CERT_IN_CHAIN`, or the same from `uv`                             | The Aikido CA variables never reached the process. Source `/etc/profile.d/aikido-*.sh` before the command that installs packages                                                                              |
| Packages are blocked as expected, but the container never appears in your device list    | The token was rejected. Malware blocking keeps working on a rejected token, so check the token itself, and that the container can reach `*.aikido.dev` over HTTPS                                             |
| Downloads are not inspected at all                                                       | Your command bypassed the entrypoint, for example a `docker run` that overrides it or a job container that replaces it. Start the agent explicitly with `aikido-endpoint-protection-start --wait-for-healthy` |
| `docker stop` takes the full grace period and kills the container                        | Something between the entrypoint and your process is swallowing `SIGTERM`. Run the image without a shell wrapper, or add `--init`                                                                             |

For anything that is not container-specific, such as picking the right package for a Red Hat-family base image, see the table on [Install Device Protection for Linux](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux.md#troubleshooting).
