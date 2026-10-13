# Setting Up Image Scanning with Local Scanner

The Aikido Security Local Scanner is a tool that enables you to perform Aikido Security scans within your environment, ensuring your data never leaves your premises. The scans take place locally, and the results are then uploaded to the Aikido Security platform. This setup allows you to scan any images locally on your own machine or in your CI/CD pipeline before they get sent to the registry.

## How to set up Local Scanning for images <a href="#how-to-set-up-local-scanning-for-images" id="how-to-set-up-local-scanning-for-images"></a>

### 1. Get your authentication token <a href="#id-1-get-your-authentication-token" id="id-1-get-your-authentication-token"></a>

1. Go to the [Local Scanner Image setup page](https://app.aikido.dev/settings/container-image-registry/add/localscan)
2. Generate an authentication token and copy. Note that you will only be able to view this token once.
3. Add this token as argument `--apikey` when running the Local Scanner in your project .

   [image: Local image scanning setup instructions with binary downloads, token generation, and setup guide.]

### 2. Adding the Local Scanner to your project <a href="#id-2-adding-the-local-scanner-to-your-project" id="id-2-adding-the-local-scanner-to-your-project"></a>

Download the local scanner binary from the Aikido UI.

### 3. Running the Local Scanner <a href="#id-3-running-the-local-scanner" id="id-3-running-the-local-scanner"></a>

Now all that is left to run the scanner on your image.

```
aikido-local-scanner image-scan your-image-name --apikey AIK_CI_xxx
```

Or if you are using our docker image (**Prerequisite:** You need to have a running Docker daemon)

```
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock aikidosecurity/local-scanner image-scan your-image-name --apikey AIK_CI_xxx
```

If you are running this in your CI, be careful to not run scan on images that are based on feature branches, to prevent production and feature results being mixed together in Aikido.

You can also run the scanner in release or PR gating mode, check out [this article](/container-image-scanning/local-image-scanning/pr-gating-for-container-images-using-local-image-scanner.md) for more info.

### 4. Check your scanning results <a href="#id-4-check-your-scanning-results" id="id-4-check-your-scanning-results"></a>

After your first scan is done, you can go to the Aikido Feed to check out your results. An image with the name you specified will have been created in the Aikido UI, containing all results from the scanning.

***
