# Scanning Fly.io with Aikido

Aikido fully supports protecting workloads on Fly through specific integrations. A native integration may be added in the future, but you can already achieve full coverage by combining:<br>

* [Code Scanning](/code-scanning/code-scanning-overview.md) for application code, dependencies and Dockerfiles
* [Local Scanner](/code-scanning/local-code-scanning.md) for Fly container images

## Features

### Container image scanning

Fly’s registry doesn’t support listing all stored images/tags. Fly only exposes images once they’re referenced by a release, so we can’t enumerate the full registry contents through an API.

As a workaround, we scan the images that are currently deployed (running) by retrieving their image references, pulling them locally, and scanning them on a daily schedule with the [Aikido Local Scanner](/code-scanning/local-code-scanning.md).

```bash
export AIKIDO_API_KEY="AIK_CI_..."
fly releases --image \
  | awk 'NR > 1 && $2 == "running" { print $NF }' \
  | sort -u \
  | xargs -r -I {} sh -c 'docker pull "{}" && aikido-local-scanner image-scan --apikey "$AIKIDO_API_KEY" --force-create-image-for-tag "{}"'
```

The above will need to be done for each application on Fly.io
