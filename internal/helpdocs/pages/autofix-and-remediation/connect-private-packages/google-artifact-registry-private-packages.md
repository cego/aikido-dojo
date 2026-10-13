# Google Artifact Registry - Private Packages

Private packages stored in a Google Artifact Registry can be retrieved by AutoFix for **Poetry, UV, NPM, PNPM, Yarn & Bun AutoFixes**. Authentication to Google Artifact Registry can be done through [Workload Identity Federation](#using-workload-identity-federation) or a [Service Account](#using-a-service-account).

## Using Workload Identity Federation

Workload Identity Federation allows you to grant the Aikido read role access to your Google Artifact Registry without sharing credentials with Aikido. This is the **recommended approach** for setting up a private packages hosted in Google Artifact Registry.

### Setup in the Google Cloud Console

#### Step 1: Create the Workload Identity Pool and AWS Provider

1. Open the Google Cloud Console and navigate to [**IAM & Admin** > **Workload Identity Federation**](https://console.cloud.google.com/iam-admin/workload-identity-pools)**.**
2. Click **+ Get Started** (or **+ Create Pool** if you already have pools setup).
3. **Name your pool:** Enter `aikido-identity-pool` as the name and click **Continue**.
4. **Add a provider to your pool:**
   * Select **AWS** from the *Select a provider* dropdown list
   * Enter the **Provider name:** `aws-provider`
   * Enter Aikido's **AWS Account ID:** `881830977366`
   * Click **Continue**.
5. **Configure provider attributes:**
   * Under **Attribute mapping**, Google automatically maps `google.subject` to `aws.arn`. Leave this as default or modify it if you want specific attribute extractions.
   * Click **Save** to create the pool and provider.

#### Step 2: Download the federated identity JSON

1. Go to the [aikido-identity-pool](http://console.cloud.google.com/iam-admin/workload-identity-pools/pool/aikido-identity-pool) details
2. Click **+ Grant access**
3. Use Grant access using federated identities
4. Click `Download config` to retrieve the configuration JSON

#### Step 3: Grant access on Artifact Registry

1. Go to [Artifact Registry](https://console.cloud.google.com/artifacts)
2. Select the repository you want to provide Aikido access to
3. Click **Permissions** in the top menu (this only appears when you have selected a repository)
4. Click **Add principal**
5. In the **New principals** textbox enter the following principal & replace `<project-id>` with your 12-digit project id:

```
principalSet://iam.googleapis.com/projects/<project-id>/locations/global/workloadIdentityPools/aikido-identity-pool/attribute.aws_role/arn:aws:sts::881830977366:assumed-role/autofix-codeartifact-readonly-role
```

6. **Assign** the `Artifact Registry Reader` role
7. Click **Save**

### **Setup in the Aikido UI**

Google Artifact Registry is currently only supported for Poetry, UV, NPM, PNPM, Yarn and Bun AutoFixes. You can configure the GCP credentials in the Google Artifact Registry section under private package configuration.

Once the steps in the GCP console are done, you can configure AutoFix to authenticate with your private registry by following the steps below:

1. Go to the [AutoFix > Settings](https://app.aikido.dev/issues/fix/settings) page in the Aikido UI
2. Click on **+ Connect Registry**, the configuration modal will now be shown

   
3. Select Google Artifact Registry:\
   
4. Paste your saved **JSON config** content from Step 3 above in the **GCP credentials** field
5. Click **Save** to save the configuration.

## Using a Service Account

### Setup in the Google Cloud Console

#### Step 1: Create a Service Account

First, create a service account in your Google Cloud project:

1. Go to the Google Cloud Console.
2. Navigate to [**IAM & Admin** > **Service Accounts**](https://console.cloud.google.com/iam-admin/serviceaccounts).
3. Click **Create Service Account**.
4. Fill in a **Service account name** such as `Aikido Artifact Registry Reader` and click **Create And Continue**.
5. Grant the service account with the **Artifact Registry Reader** role.
6. Click **Continue** and **Done**.

#### Step 2: Create a Key for the Service Account <a href="#id-2-create-a-key-for-the-service-account" id="id-2-create-a-key-for-the-service-account"></a>

1. On the **Service Accounts** page, find the service account you just created.
2. Click on the three dots on the right and select **Manage Keys**.
3. Click **Add Key** > **Create New Key**.
4. Choose **JSON** and click **Create**.
5. Save the **JSON key file** to a secure location.

### **Setup in the Aikido UI**

Google Artifact Registry is currently only supported for Poetry, UV, NPM, PNPM, Yarn and Bun AutoFixes. You can configure the GCP credentials in the Google Artifact Registry section under private package configuration.

Once the steps in the GCP console are done, you can configure AutoFix to authenticate with your private registry by following the steps below:

1. Go to the [AutoFix > Settings](https://app.aikido.dev/issues/fix/settings) page in the Aikido UI
2. Click on **+ Connect Registry**, the configuration modal will now be shown

   
3. Select Google Artifact Registry:\
   
4. Paste your saved **Service Acount Key** content from Step 2 above in the **GCP credentials** field
5. Click **Save** to save the configuration.
