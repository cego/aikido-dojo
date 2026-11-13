# Device Protection in TeamCity

TeamCity builds run on Linux build agents, so Device Protection installs the same way it does on any Linux machine. What is specific to TeamCity is where the install goes: once on a build agent that stays around, or inside each build on an agent that is thrown away afterwards.

Read [Install Device Protection for Linux](https://help.aikido.dev/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux) first for the token and the install flavors. The build-level snippets below use the container flavor, which [Device Protection in Containers](https://help.aikido.dev/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/device-protection-in-containers) explains.

## Permanent agents

Install the agent on each Linux build agent host with `--headless`, then reboot. See [Install Device Protection for Linux](https://help.aikido.dev/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux#choose-the-install-flavor). No build configuration changes needed.

For cloud agents started from a VM image, such as an Amazon EC2 cloud profile, bake the agent into the image instead. See [Deploy Device Protection in Golden Images](https://help.aikido.dev/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/deploy-device-protection-in-golden-images).

## JetBrains-hosted and other short-lived agents

JetBrains-hosted agents on TeamCity Cloud, and agents started per build on Kubernetes, have no host to prepare. Install the agent in the build and start it before anything installs packages.

{% stepper %}
{% step %}
### Add the token as a password parameter

In your project, go to **Project Settings** → **Parameters** and click **Add new parameter**. Name it `env.AIKIDO_TOKEN`, set the kind to **Environment variable**, and set the type to **Password**. TeamCity hides the value in build logs and exposes it to every build step as `AIKIDO_TOKEN`.
{% endstep %}

{% step %}
### Start Device Protection in the first build step

Add a **Command Line** build step before any step that installs packages, and set its custom script to:

{% code overflow="wrap" %}
```shell
curl -fsSL https://github.com/AikidoSec/safe-chain/raw/refs/heads/main/install-scripts/install-endpoint-linux.sh \
  | sudo AIKIDO_TOKEN="$AIKIDO_TOKEN" sh -s -- --container
sudo aikido-endpoint-protection-start --wait-for-healthy
```
{% endcode %}
{% endstep %}

{% step %}
### Source the Aikido settings in each step that installs packages

Start every later step that downloads packages with the sourcing line:

{% code overflow="wrap" %}
```shell
for f in /etc/profile.d/aikido-*.sh; do . "$f"; done
npm ci
```
{% endcode %}
{% endstep %}
{% endstepper %}

* **Pass the token after `sudo`.** `sudo` drops the rest of your environment. On an agent that runs builds as root, drop both `sudo` prefixes.
* **Source `/etc/profile.d/aikido-*.sh` in every step that installs packages.** Each build step is a new process, and without those settings downloads skip the agent.
* **The agent is a VM, but it takes `--container`.** A build cannot reboot the agent, and container mode needs no reboot. That is also why the agent is started explicitly. See [Choose the install flavor](https://help.aikido.dev/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux#choose-the-install-flavor).
* **Write `%%` for a literal `%`.** TeamCity reads `%name%` as a parameter reference in step scripts.

### Kotlin DSL

If you keep your build configurations in `.teamcity/settings.kts`, the same build looks like this:

{% code title=".teamcity/settings.kts" overflow="wrap" %}
```kotlin
object Build : BuildType({
    name = "Build"

    params {
        password("env.AIKIDO_TOKEN", "credentialsJSON:<token-id>")
    }

    steps {
        script {
            name = "Start Device Protection"
            scriptContent = """
                curl -fsSL https://github.com/AikidoSec/safe-chain/raw/refs/heads/main/install-scripts/install-endpoint-linux.sh \
                  | sudo AIKIDO_TOKEN="${'$'}AIKIDO_TOKEN" sh -s -- --container
                sudo aikido-endpoint-protection-start --wait-for-healthy
            """.trimIndent()
        }
        script {
            name = "Install dependencies"
            scriptContent = """
                for f in /etc/profile.d/aikido-*.sh; do . "${'$'}f"; done
                npm ci
            """.trimIndent()
        }
    }
})
```
{% endcode %}

* **Write `${'$'}` for every shell `$`.** Kotlin expands `$` inside `"""` strings, so `$f` breaks the build script.
* **Don't put the token in the DSL.** Add the password parameter in the TeamCity UI, then reference the `credentialsJSON:` token that TeamCity generates for it.

## Steps that run in a Docker container

With **Run step within Docker container** set on a step, TeamCity starts a fresh container for that step and replaces the image's `ENTRYPOINT`, so the agent's container entrypoint never runs. Start the agent in the same step that installs packages:

{% code overflow="wrap" %}
```shell
aikido-endpoint-protection-start --wait-for-healthy
for f in /etc/profile.d/aikido-*.sh; do . "$f"; done
npm ci
```
{% endcode %}

This assumes the image already has the agent installed. See [Add Device Protection to an image](https://help.aikido.dev/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/device-protection-in-containers#add-device-protection-to-an-image). For an image without it, put the install line from the previous section first in the step, without `sudo`, since the Docker wrapper runs the step as root. The image then needs `curl`.

TeamCity passes `env.AIKIDO_TOKEN` into the container. Every step is a new container, so repeat these lines in each step that installs packages.

## Verify it works

Add a **Command Line** step after Device Protection has started:

{% code overflow="wrap" %}
```shell
for f in /etc/profile.d/aikido-*.sh; do . "$f"; done
curl -sS -o /dev/null -w "%%{http_code}\n" --max-time 15 https://registry.npmjs.org/safe-chain-test/-/safe-chain-test-0.0.1-security.tgz
```
{% endcode %}

`403` means the agent blocked the malware test package. `200` means the download skipped the agent. The agent also appears in your [device list](https://app.aikido.dev/endpoint-protection/devices?regionRedirect=1). For more checks, see [How to test if Device Protection is working](https://help.aikido.dev/aikido-device-protection/miscellaneous-aikido-endpoint/how-to-test-device-protection).

## Troubleshooting

| Problem                                                    | Fix                                                                                                                                                     |
| ---------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `Token is required`                                        | The parameter is not named `env.AIKIDO_TOKEN`, or its kind is not **Environment variable**. Also check you pass the token after `sudo`                  |
| `Root privileges required`                                 | Prefix the install and the start command with `sudo`                                                                                                    |
| Test package returns `200`, or `SELF_SIGNED_CERT_IN_CHAIN` | The step did not source `/etc/profile.d/aikido-*.sh`                                                                                                    |
| Downloads in a Docker wrapper step are not inspected       | The wrapper replaced the image's entrypoint, so the agent never started. Run `aikido-endpoint-protection-start --wait-for-healthy` first in the step |
| The build fails with an unresolved parameter               | A `%` in the script was read as a parameter reference. Write `%%` instead                                                                               |

For anything that is not TeamCity-specific, such as certificate errors from `npm` or `uv`, see the troubleshooting table on [Device Protection in Containers](https://help.aikido.dev/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/device-protection-in-containers#troubleshooting).
