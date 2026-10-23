# Device Protection in Docker Containers

Package installs inside a Docker container run in their own Linux environment, outside the reach of the agent on the host. With **Protect Containers** turned on, Device Protection on the host finds the running containers and installs the Linux agent into them, so `npm`, `pip`, and other package managers in the container get inspected the same way they do on the host.

## Prerequisites

* Device Protection on a macOS, Windows, or Linux host, signed in to by a user whose `docker` command can reach the Docker engine. Docker Desktop is the tested setup. Other runtimes, such as Podman or Rancher Desktop, are covered only when they provide a working `docker` command.
* Containers based on a Debian-family image (Debian, Ubuntu) or a Red Hat-family image (RHEL, CentOS, Fedora, Rocky Linux, AlmaLinux). Alpine and other distributions are skipped.
* Outbound HTTPS from each container to `github.com` to download the installer.

## Turn on container protection

1. Go to [**Device Protection** → **Settings**](https://app.aikido.dev/endpoint-protection/settings).
2. Under **Device Settings**, turn on **Protect Containers**.

[image: Device Settings card with the Protect Containers and Protect WSL2 toggles]

The setting applies to every device in your workspace. Devices pick it up on their next check-in.

## How it works

The host agent checks the running containers every 30 seconds and installs into new containers as soon as they start. It only does this while protection is on, not paused, and the host is enrolled with a valid token.

* **Only running containers are touched.** Aikido never starts a stopped container.
* **The Linux agent is installed when it is missing.** Aikido runs the standard Linux install script as root inside the container with `--container`, using the host's token, then starts the agent. If the image has neither `curl` nor `wget`, Aikido installs `curl` with the image's package manager first. A container that already has the agent is only started.
* **2 containers at a time.** An install can take a few minutes, so containers that start meanwhile wait their turn.
* **Each container is its own device.** A protected container reports to Aikido as a separate Linux device, linked to the host that protected it. When the container stops, the host unlinks it.
* **The agent lives in the container, not the image.** Removing the container removes the agent. A new container from the same image gets a fresh install when it starts.

These containers are skipped:

| Container                                                           | Why                                                                                    |
| ------------------------------------------------------------------- | -------------------------------------------------------------------------------------- |
| Images that are not Debian-family or Red Hat-family, e.g. Alpine    | The Linux agent doesn't support them                                                   |
| Containers with the label `com.aikido.endpoint.protection=disabled` | Opted out. See [Opt a container out](#opt-a-container-out)                             |
| Containers whose name starts with `k8s_`                            | Kubernetes system containers, such as the control plane of Docker Desktop's Kubernetes |

Turning **Protect Containers** off, or pausing protection on the host, stops new installs. It doesn't remove the agent from containers that already have it.

## Opt a container out

Add the opt-out label when you start the container:

```shell
docker run --label com.aikido.endpoint.protection=disabled <image>
```

In Docker Compose, set it under `labels` on the service. To opt out every container from an image, add `LABEL com.aikido.endpoint.protection=disabled` to its Dockerfile.

## Use the agent in a container

The agent routes package downloads through a proxy inside the container. Processes that were already running when the agent started don't use it. Commands you run afterwards need the agent's proxy and certificate settings, which it writes to `/etc/profile.d/`. Source them before you install packages:

```shell
docker exec -it <container> sh -c 'for f in /etc/profile.d/aikido-*.sh; do . "$f"; done; exec bash'
```

To protect containers from the moment they start, including the command they run on startup, build Device Protection into the image instead. See [Device Protection in Containers](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/device-protection-in-containers.md).

## Verify it works

From the host, run:

```shell
docker exec -u 0 <container> aikidoctl version
```

It prints the installed version and the agent's state. If the command is not found, the install hasn't happened yet: check that the container is running and isn't on the skipped list above.

## Troubleshooting

| Problem                                   | Fix                                                                                                                                                                                                                                                                     |
| ----------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| A container is never protected            | Check that it is running, that its image is Debian-family or Red Hat-family, and that it has no opt-out label. Then check that `docker ps` works for the user signed in at the host                                                                                     |
| Packages in the container are not blocked | The shell or process doesn't use the agent's proxy. Source `/etc/profile.d/aikido-*.sh` as shown in [Use the agent in a container](#use-the-agent-in-a-container)                                                                                                       |
| The install keeps failing                 | Aikido retries after 1, 4, then 16 minutes, and gives up after 4 attempts until the host agent restarts. The host agent log has the installer output on lines starting with `Docker agent:`. Send that log to [Contact Support](/miscellaneous-info/contact-support.md) |
| Nothing happens while nobody is signed in | Expected. Aikido runs `docker` as the user signed in at the host                                                                                                                                                                                                        |

The host agent log is at:

| Host    | Log file                                                                      |
| ------- | ----------------------------------------------------------------------------- |
| macOS   | `/Library/Logs/AikidoSecurity/EndpointProtection/endpoint-protection.log`     |
| Windows | `%ProgramData%\AikidoSecurity\EndpointProtection\logs\EndpointProtection.log` |
| Linux   | `/var/log/aikidosecurity/endpoint-protection/endpoint-protection.log`         |

## Related

{% content-ref url="/pages/U1zOPz2xS3bI8ZXDYgin" %}
[Device Protection in Containers](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/device-protection-in-containers.md)
{% endcontent-ref %}

{% content-ref url="/pages/XPg332v7R3oyny1XQC2p" %}
[Device Protection in WSL2](/aikido-device-protection/deploying-aikido-endpoint/device-protection-in-wsl2.md)
{% endcontent-ref %}
