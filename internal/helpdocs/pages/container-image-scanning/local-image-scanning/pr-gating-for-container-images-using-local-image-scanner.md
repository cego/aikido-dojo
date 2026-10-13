# PR Gating For Container Images Using Local Image Scanner

The Aikido Local Scanner can be used to enforce security gates in your CI pipeline.

PR gating ensures that new code meets your security standards before it is merged into the default branch. It scans only the changes introduced in the pull request. If the scan detects new issues that meet or exceed your configured severity threshold, the CI pipeline fails.

This prevents new vulnerabilities from being introduced while allowing existing findings to be addressed separately.

{% hint style="info" %}
We also support [release gating](/container-image-scanning/local-image-scanning/release-gating-for-container-images-using-local-image-scanner.md) for teams that want to enforce checks at release time.
{% endhint %}

To enable PR gating, add the `--fail-on <severity>` option to select your preferred severity level. Then, add the `--gating-mode pr` option to signify that you wish to perform PR gating. You must also specify the base (`--base-commit-id <commit-id>`) and head commit (`--head-commit-id <commit-id>`). If there a scan was previously performed on the base commit id, the scan results will be compared to those. If not, the results will be compared to the most recent scan on your image.

Example PR gating command:

```
./aikido-local-scanner image-scan your-image-name --apikey AIK_CI_xxx --fail-on critical --gating-mode pr --base-commit-id abc123 --head-commit-id def456 
```

### Examples <a href="#examples" id="examples"></a>

#### GitHub <a href="#github" id="github"></a>

For general information about setting up the Local Scanner in a GitHub environment, check out [this article](https://help.aikido.dev/code-scanning/local-code-scanning/github-action-setup-for-local-code-scanning).

Example `.github/workflows/aikido-scan.yml` for PR gating:

```yaml
name: Aikido Docker build and scan

on:
  pull_request:
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
            --fail-on critical \
            --gating-mode pr \
            --base-commit-id ${{ github.event.pull_request.base.sha }} \
            --head-commit-id ${{ github.event.pull_request.head.sha }} 
```

#### Azure DevOps Server <a href="#azure-devops-server" id="azure-devops-server"></a>

For general information about setting up the Local Scanner in an Azure DevOps environment, check out [this article](https://help.aikido.dev/container-image-scanning/local-image-scanning/azure-devops-server-setup-for-local-image-scanning).

Example for PR gating:

```yaml
trigger: none

pr:
  branches:
    include:
      - 'main'

pool:
  vmImage: ubuntu-latest

steps:
- script: docker build -t my-image .
  displayName: 'Build docker image'  
- script: | 
    docker run --rm -v /var/run/docker.sock:/var/run/docker.sock aikidosecurity/local-scanner image-scan my-image \
        --apikey $(AIKIDO_API_KEY) \
        --fail-on critical \
        --gating-mode pr \
        --base-commit-id abc123 \
        --head-commit-id def456
  displayName: 'Run Aikido scan'
```

Make sure to add this pipeline as a Build Validation to your main branch:

1. Go to Project Settings → Repositories → Policies → Branch Policies
2. Select your main branch
3. Add Build Validation
4. Choose this pipeline
5. Set:
   * Trigger: Automatic
