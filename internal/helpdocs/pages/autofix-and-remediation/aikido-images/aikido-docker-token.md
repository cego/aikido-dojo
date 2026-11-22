# Aikido Docker Token

The Aikido Docker token is a read-only token for pulling Aikido Images. It is part of the image path: `docker.aikido.io/<token>/<image>:<tag>`.

By default, Container AutoFix writes the token into the Dockerfile. This means you don't need to configure anything on each developer's machine.

## Keep the token out of a public repository

If your Dockerfile is in a public repository, keep the token out of the repository. Use a build arg instead.

Declare the arg before the first `FROM`, and use it in the image path:

{% code title="Dockerfile" %}
```dockerfile
ARG AIKIDO_TOKEN
FROM docker.aikido.io/${AIKIDO_TOKEN}/golang:1.24-bookworm
```
{% endcode %}

Aikido passes the `AIKIDO_TOKEN` build arg automatically when it builds your image for Container AutoFix. You don't need to configure it. In your own pipeline, pass the same `AIKIDO_TOKEN` build arg from your CI secret store.

To find your token, go to [**AutoFix** → **Aikido Images**](https://app.aikido.dev/issues/aikido-images?regionRedirect=1) and click an image. The token is the part of the image path between `docker.aikido.io/` and the image name.

## Related

{% content-ref url="autofix-for-containers-using-aikido-images.md" %}
[autofix-for-containers-using-aikido-images.md](https://help.aikido.dev/autofix-and-remediation/aikido-images/autofix-for-containers-using-aikido-images)
{% endcontent-ref %}

{% content-ref url="../configure/setting-up-container-autofix.md" %}
[setting-up-container-autofix.md](https://help.aikido.dev/autofix-and-remediation/configure/setting-up-container-autofix)
{% endcontent-ref %}
