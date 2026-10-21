# Release Gating For Container Images Using Local Image Scanner

The Aikido Local Scanner can enforce security gates at different stages of your CI pipeline, including pull requests and releases.

Release gating scans the full build artifact, such as a container image, before it is published or deployed. Unlike [PR gating](/container-image-scanning/local-image-scanning/pr-gating-for-container-images-using-local-image-scanner.md), which scans only the changes introduced in a pull request, release gating evaluates the entire build. If the scan finds issues that meet or exceed your configured severity threshold, the CI pipeline fails and the release is blocked.

This ensures that the final artifact meets your security standards before it reaches production.

You can enable release gating mode by using the `--fail-on <severity>` flag. This feature is helpful when scanning your image prior to publishing it to an image library, as it ensures there are no open issues before release. When running in release gating mode, the scanner process will fail when there are any open issues of the chosen severity or higher after the scan is finished.

Example release gating command:

```
./aikido-local-scanner image-scan your-image-name --apikey AIK_CI_xxx --fail-on critical
```

### Examples <a href="#examples" id="examples"></a>

#### GitHub <a href="#github" id="github"></a>

For general information about setting up the Local Scanner in a GitHub environment, check out [this article](https://help.aikido.dev/code-scanning/local-code-scanning/github-action-setup-for-local-code-scanning). Example `.github/workflows/aikido-scan.yml` for release gating:

```yaml
name: Aikido Docker build and scan

on:
  push:
    branches:
      - main

jobs:
  build-and-scan:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Build Docker image
        run: docker build -t your-local-image-name .

      - name: Run Aikido image scan
        run: |
          docker run --rm \
            -v /var/run/docker.sock:/var/run/docker.sock \
            aikidosecurity/local-scanner \
            image-scan your-local-image-name \
            --apikey ${{ secrets.AIKIDO_API_KEY }} \
            --fail-on critical
```

#### Azure DevOps Server <a href="#azure-devops-server" id="azure-devops-server"></a>

For general information about setting up the Local Scanner in an Azure DevOps environment, check out [this article](https://help.aikido.dev/container-image-scanning/local-image-scanning/azure-devops-server-setup-for-local-image-scanning). Example for release gating:

```yaml
trigger:
- main

pool:
  vmImage: ubuntu-latest

steps:
- script: docker build -t my-image .
  displayName: 'Build docker image'  
- script: |
    docker run --rm -v /var/run/docker.sock:/var/run/docker.sock aikidosecurity/local-scanner image-scan my-image \
    --apikey $(AIKIDO_API_KEY) \
    --fail-on critical
  displayName: 'Run Aikido scan'
```
