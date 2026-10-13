# Cloudsmith Container Registry

You can integrate your Cloudsmith's Container Registry with Aikido, to scan your images for known vulnerabilities.

Follow the simple steps below to activate this feature:

**Step 1:** Log into your Cloudsmith account and navigate to API Settings.

See screenshot below.

We'll have to gather an API Key. Direct URL: <https://cloudsmith.io/user/settings/api/>

[image: User account dropdown menu displaying settings, documentation links, and logout option.]

You can copy the API key directly from this screen:

[image: API key management and restriction settings for Cloudsmith account integration and security.]

**Step 2:** Enter the collected data in Aikido (direct link: <https://app.aikido.dev/settings/container-image-registry/add/cloudsmith>)

[image: Form to connect and authenticate a Cloudsmith container image registry account.]

* **Username** can be found in your user account on top (see first screenshot above)
* **Organization Namespace** can be found here:

[image: Organization membership management dashboard showing user role and organization options.]

**Step 3:** Aikido will now find all container repositories you can access and list them.

**Step 4:** Repositories can be linked to a code repository in order to perform better deduplication of findings. This step is optional!

**Step 5:** In the action menu next to the registry, click 'scan repos in registry' to get started. Results will appear in the Feed!

***
