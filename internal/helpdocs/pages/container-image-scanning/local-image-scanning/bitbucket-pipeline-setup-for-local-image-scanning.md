# BitBucket Pipeline Setup for Local Image Scanning

The Aikido Security Local Scanner is a tool that enables you to perform Aikido Security scans within your environment, ensuring your data never leaves your premises. The scans take place locally, and the results are then uploaded to the Aikido Security platform. This setup allows you to scan any images locally on your own machine or in your CI/CD pipeline before they get sent to the registry.

## How to set up Local Scanning for images <a href="#how-to-set-up-local-scanning-for-images" id="how-to-set-up-local-scanning-for-images"></a>

### 1. Get your authentication token <a href="#id-1-get-your-authentication-token" id="id-1-get-your-authentication-token"></a>

1. Go to the [Local Scanner Image setup page](https://app.aikido.dev/settings/container-image-registry/add/localscan)
2. Generate an authentication token and copy. Note that you will only be able to view this token once.
3. Add the authentication token to the Repository variables in Bitbucket to make it available for use in the pipeline. To do this, you need to go to your Repository Settings. Then navigate to **Pipelines > Repository variables.** There you can add the authentication token with `AIKIDO_API_KEY` as the key and the copied token contents as the value.

### 2. Running the Local Scanner <a href="#id-2-running-the-local-scanner" id="id-2-running-the-local-scanner"></a>

Now all that is left to run the scanner on your image. You can use this example file as a guide:

```yaml
image: docker:25.0.3

pipelines:
  default:
    - step:
        name: Build Docker Image
        image: docker:25.0.3
        services:
          - docker
        script:
          - export PATH=/usr/bin:$PATH
          - docker build -t my-image .
          - docker save --output my-image.tar my-image
        artifacts:
          - my-image.tar
    - step:
        name: Scan Docker Image
        image: aikidosecurity/local-scanner
        script:
          - aikido-local-scanner image-scan my-image.tar --image-name my-image --apikey $AIKIDO_API_KEY
```

### 3. Check your scanning results <a href="#id-3-check-your-scanning-results" id="id-3-check-your-scanning-results"></a>

After your first scan is done, you can go to the Aikido Feed to check out your results. An image with the name you specified will have been created in the Aikido UI, containing all results from the scanning.
