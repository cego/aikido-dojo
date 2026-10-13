# Scaleway Container Registry

You can integrate your Scaleway Container Registry with Aikido to scan your containers for known vulnerabilities.

Follow the simple steps below to activate this feature:

**Step 1:** Log into your Scaleway account to gather some data.

We'll have to create an API key. Direct URL: <https://console.scaleway.com/iam/api-keys> (see screenshots below)

[image: User profile dropdown menu showing organization and account management options.]

[image: IAM dashboard for managing users, API keys, and access permissions in an organization.]

[image: API key generation form with user options and expiration settings.]

**Step 2:** Enter the collected data in Aikido (direct link: <https://app.aikido.dev/settings/container-image-registry/add/scaleway>)

[image: Form to connect and configure a Scaleway container image registry.]

**Step 3:** Aikido will now find all container repositories you can access and list them.

**Step 4:** Repositories can be linked to a code repository in order to perform better deduplication of findings. This step is optional!

**Step 5:** In the action menu next to the registry, click 'scan repos in registry' to get started. Results will appear in the Feed!

***
