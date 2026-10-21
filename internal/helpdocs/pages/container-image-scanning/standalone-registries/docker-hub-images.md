# Docker Hub images

You can now integrate your Docker Hub with Aikido to scan your containers for known vulnerabilities.

Follow the simple steps below to activate this feature:

**Step 1:** Log into your Docker Hub account and navigate to account settings, then to Personal access tokens

**Step 2:** Create an access token with scope 'Read-only'.

**Step 3:** Back in Aikido, go to settings, then [containers](https://app.aikido.dev/settings/container-image-registry). Click 'Connect registry' and pick Docker Hub. Enter your organization namespace and your username, in case you want to scan your personal repositories and the access token from the previous step.\
​\
You can find your username by clicking your avatar in the top navigation on the right, your username will then be shown underneath the avatar image. In our case it's `aikidotestaccount`.

**Step 4:** Aikido will now find all container repositories you can access and list them.

**Step 5:** Repositories can be linked to a code repository in order to perform better deduplication of findings. This step is optional!

**Step 6:** In the action menu next to the registry, click 'scan repos in registry' to get started. Results will appear in the Feed!
