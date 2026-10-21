# GCP Scanner for GCP Artifact Registry

**Table of contents:**

* [Setup](#setup)
  * [Step 1: Set Up Aikido Integration with GCP](#step-1-set-up-aikido-integration-with-gcp)
  * [Step 2: Enable Artifact Registry Scanning in GCP](#step-2-enable-artifact-registry-scanning-in-gcp)
  * [Step 3: Push the latest version of your images to artifact registry](#step-3-push-the-latest-version-of-your-images-to-artifact-registry)
  * [Step 4: Start a scan in Aikido to process the results](#step-4-start-a-scan-in-aikido-to-process-the-results)
  * [Step 5. Link Images](#step-5-link-images)

## GCP Scanner for GCP Artifact Registry

Google Cloud Platform (GCP) provides an efficient way to store container images through the Artifact Registry. Leveraging the power of Aikido in conjunction with GCP's Artifact Registry ensures a robust security framework. Let's walk through the process of enabling container analysis for images stored in Artifact Registry.

> Note: if you use the deprecated 'Container Registry' instead of Artifact Registry, see [GCP Legacy Container Registry](https://help.aikido.dev/container-image-scanning/cloud-provider-registries/image-scanning-for-gcp-legacy-container-registry) for setup steps.

### Setup <a href="#setup" id="setup"></a>

Aikido will use the findings reported by GCP Container Analysis and run them through the same deduplication and de-noising engine you are familiar with. Let's dive into the details of this new functionality and how to enable it.

#### Step 1: Set Up Aikido Integration with GCP <a href="#step-1-set-up-aikido-integration-with-gcp" id="step-1-set-up-aikido-integration-with-gcp"></a>

Before you begin, make sure your GCP cloud environment has been linked with Aikido. If you have not done this already, navigate to the[ cloud overview](https://app.aikido.dev/clouds) in Aikido. Click on "Connect cloud" and follow the steps to get set up. More information can be found in [this article](https://help.aikido.dev/cloud-scanning/connect-your-cloud/gcp/connect-gcp-account-to-aikido).

#### Step 2: Enable Artifact Registry Scanning in GCP <a href="#step-2-enable-artifact-registry-scanning-in-gcp" id="step-2-enable-artifact-registry-scanning-in-gcp"></a>

After connecting your GCP environment, navigate to the [Artifact Registry](https://console.cloud.google.com/artifacts) service. Click on "Settings" on the left-hand side, where you can turn on vulnerability scanning for your container images.

[image: Vulnerability scanning settings with option to enable scanning for container security issues.]

#### Step 3: Push the latest version of your images to artifact registry <a href="#step-3-push-the-latest-version-of-your-images-to-artifact-registry" id="step-3-push-the-latest-version-of-your-images-to-artifact-registry"></a>

GCP will only scan newly pushed images for vulnerabilities. So for the analysis to start on your images, you should push the latest version of your images again to the artifact registry.

#### Step 4: Start a scan in Aikido to process the results <a href="#step-4-start-a-scan-in-aikido-to-process-the-results" id="step-4-start-a-scan-in-aikido-to-process-the-results"></a>

After enabling the scanner in GCP, you can go back to Aikido and start a scan for your GCP cloud environment.

#### Step 5. Link Images <a href="#step-5-link-images" id="step-5-link-images"></a>

The last step is to link a cloud image to a code repository. During the scan in the previous step, Aikido will look for any repositories you host on the Artifact Registry in GCP. Go to the 'Images' tab on the cloud detail page and link the images to the code repository where the source code is hosted.

[image: Navigation menu with "Images" tab highlighted and selected.]

Once the cloud images are linked to the code repositories, Aikido will assess and score the findings from Container Analysis and link them to the related cloud environment and code repository.

Aikido also supports the scanning of your containers hosted in GCP's Container Registry. You can find those instructions [here](https://help.aikido.dev/container-image-scanning/cloud-provider-registries/image-scanning-for-gcp-legacy-container-registry).

***
