# JFrog Artifactory

You can integrate your JFrog Artifactory with Aikido to scan your images for known vulnerabilities.

Follow the simple steps below to activate this feature:

**Step 1:** In JFrog, we have to collect some data including your username and a JWT access token. The username is the value displayed in the Users table. To start, click 'User Management' on the top right.

[image: Platform configuration menu with options for projects, user management, and security settings.]

Then, click 'Access Tokens' on the left menu:

[image: Platform configurations sidebar showing user management and security access token options.]

In the top-right corner, click 'Generate token' and fill out the settings as below:

[image: Generate Token window for creating a scoped user access token with custom settings.]

**Step 2:** Enter the collected data in Aikido (direct link: <https://app.aikido.dev/settings/container-image-registry/add/artifactory>)

[image: Form for connecting to JFrog Artifactory container image registry using credentials and URL.]

**Step 3:** Aikido will now find all container repositories you can access and list them.

**Step 4:** Repositories can be linked to a code repository in order to perform better deduplication of findings. This step is optional!

**Step 5:** In the action menu next to the registry, click 'scan repos in registry' to get started. Results will appear in the Feed!

***
