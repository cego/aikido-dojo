# Jenkins Setup for Local Image Scanning

The Aikido Security Local Scanner is a tool that enables you to perform Aikido Security scans within your environment, ensuring your data never leaves your premises. The scans take place locally, and the results are then uploaded to the Aikido Security platform. This setup allows you to scan any images locally on your own machine or in your CI/CD pipeline before they get sent to the registry.

## How to set up Local Scanning for images <a href="#how-to-set-up-local-scanning-for-images" id="how-to-set-up-local-scanning-for-images"></a>

### 1. Get your authentication token <a href="#id-1-get-your-authentication-token" id="id-1-get-your-authentication-token"></a>

1. Go to the [Local Scanner Image setup page](https://app.aikido.dev/settings/container-image-registry/add/localscan)
2. Generate an authentication token and copy. Note that you will only be able to view this token once.
3. Add this token as argument `--apikey` when running the Local Scanner in your project .
4. Store this token in your Jenkins. To do this, navigate to the "Credentials" page in the Jenkins settings.

[image: Jenkins interface for adding a new secret text credential for global use.]

### 2. Adding the Local Scanner to your project <a href="#id-2-adding-the-local-scanner-to-your-project" id="id-2-adding-the-local-scanner-to-your-project"></a>

Download the local scanner binary from the Aikido UI. In Jenkins, add it on your build node and add it to /usr/local/bin to access it by name. Make sure the binary is executable (example: `chmod +x -R aikido-local-scanner`). If you are using a runner that has Docker installed, it is possible to skip this step and use the Local Scanner docker image to perform the scan (see below).

### 3. Running the Local Scanner <a href="#id-3-running-the-local-scanner" id="id-3-running-the-local-scanner"></a>

Now all that is left to run the scanner on your repository.\
​\
Make sure that the local scanner is only triggered for your default branch. By default, Aikido supports scanning one branch in your repository for dependency and code issues, typically the main or master branch. Therefore, we recommend running the local scanner exclusively on that branch to avoid mixing scan results on the Aikido platform. If you use git, you can specify this in the *Branches to build* section.

[image: Configuring Jenkins to build only the 'main' branch.]

Now, make the API key available to use in the project by adding it to your build environment.

[image: Configuring secret API key binding for Aikido in the build environment.]

Add an 'Execute shell' step to your project with the following content. You can either use name of your image or use the path of a tar artifact of your image.

```shellscript
aikido-local-scanner image-scan your-image --apikey $AIKIDO_API_KEY
```

Or if you are using our docker image

```shellscript
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock aikidosecurity/local-scanner image-scan your-image --apikey $AIKIDO_API_KEY
```

To run a release gated scan before pushing to your registry, you can add `--fail-on critical` to the commands above. Adjust the severity to meet your specific gating needs.

### 4. Check your scanning results <a href="#id-4-check-your-scanning-results" id="id-4-check-your-scanning-results"></a>

After your first scan is done, you can go to the Aikido Feed to check out your results. An image with the name you specified will have been created in the Aikido UI, containing all results from the scanning.
