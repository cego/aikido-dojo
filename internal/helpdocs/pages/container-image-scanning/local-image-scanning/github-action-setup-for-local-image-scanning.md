# GitHub Action Setup for Local Image Scanning

The Aikido Security Local Scanner is a tool that enables you to perform Aikido Security scans within your environment, ensuring your data never leaves your premises. The scans take place locally, and the results are then uploaded to the Aikido Security platform. This setup allows you to scan any images locally on your own machine or in your CI/CD pipeline before they get sent to the registry.

## How to set up Local Scanning for images <a href="#how-to-set-up-local-scanning-for-images" id="how-to-set-up-local-scanning-for-images"></a>

### 1. Get your authentication token <a href="#id-1-get-your-authentication-token" id="id-1-get-your-authentication-token"></a>

1. Go to the [Local Scanner Image setup page](https://app.aikido.dev/settings/container-image-registry/add/localscan)
2. Generate an authentication token and copy. Note that you will only be able to view this token once.
3. Add this token as argument `--apikey` when running the Local Scanner in your project.
4. Save this token in your GitHub Secrets by going to Settings > Secrets and variables > Actions.

[image: GitHub repository secrets interface showing an API key entry and an option to add more.]

### 2. Running the Local Scanner <a href="#id-2-running-the-local-scanner" id="id-2-running-the-local-scanner"></a>

Now all that is left to add a new YAML file to the `.github/workflows` folder and run the scanner on your image. You can find some examples below.

Building and scanning an image:

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
            --apikey ${{ secrets.AIKIDO_API_KEY }}
```

To run a release gated scan before pushing to your registry, you can run a workflow like:

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
            --apikey ${{ secrets.AIKIDO_API_KEY }} --fail-on critical
```

Adjust the severity specified in `--fail-on` to match your needs.

### 4. Check your scanning results <a href="#id-4-check-your-scanning-results" id="id-4-check-your-scanning-results"></a>

After your first scan is done, you can go to the Aikido Feed to check out your results. An image with the name you specified will have been created in the Aikido UI, containing all results from the scanning.
