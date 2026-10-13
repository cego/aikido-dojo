# DigitalOcean Container Registry

You can integrate your DigitalOcean Container Registry with Aikido to scan your containers for known vulnerabilities.

Follow the simple steps below to activate this feature:

**Step 1:** Log into your DigitalOcean account to gather some data.

We'll have to gather only a Personal Access Token (PAT). Direct URL: <https://cloud.digitalocean.com/account/api/tokens> (see screenshot below)

[image: Form to create a new personal access token with scope and expiration options.]

**Step 2:** Enter the collected data in Aikido (direct link: <https://app.aikido.dev/settings/container-image-registry/add/digitalocean>)

[image: Form to connect DigitalOcean container image registry with nickname and access token.]

**Step 3:** Aikido will now find all container repositories you can access and list them.

**Step 4:** Repositories can be linked to a code repository in order to perform better deduplication of findings. This step is optional!

**Step 5:** In the action menu next to the registry, click 'scan repos in registry' to get started. Results will appear in the Feed!

***
