# GCP Legacy Container Registry

You can now integrate your GCP Container Registry with Aikido to scan your containers for known vulnerabilities.

Follow the simple steps below to activate this feature:

1. Create a service account

   First, you need to create a service account in GCP's 'IAM & Admin' console. The service account needs to have been assigned the `Container Registry Service Agent` role.
2. Create a key for the service account

   Once the service account is created, you need to create a key for this account which Aikido can use to access the containers. Navigate to the service account's details from the overview by clicking on the name, and then selecting the `KEYS` tab.\
   ​\
   On this tab, click on `ADD KEY` and create a new JSON key. When the key is created, it will automatically download a JSON file with the key contents. You need to provide the key contents in the next step.
3. Connect the registry in Aikido

   Now that you created the service account and obtained the key, you can connect the registry in Aikido. Start by going to [Aikido's container overview page](https://app.aikido.dev/containers). Click on `Connect registry` and select `GCP Container Registry` from the list of registries.\
   ​\
   On this page, you need to enter the project ID from GCP where the Container Registry is hosted and then you can upload the JSON key file you obtained in the previous step. Aikido will store this key encrypted in a secure place and use it to scan the images.\
   ​\
   Once the registry is connected, we'll look for the repositories in your Container Registry, after which you can start scanning the ones you want.

[image: Connect GCP Container Registry by entering project ID and uploading a JSON key file.]

Aikido also supports the scanning of your containers hosted in GCP's Artifact Registry. You can find those instructions [here](https://help.aikido.dev/container-image-scanning/cloud-provider-registries/aikido-scanner-for-gcp-artifact-registry).

### &#x20;<a href="#set-up-gcp-container-registry" id="set-up-gcp-container-registry"></a>
