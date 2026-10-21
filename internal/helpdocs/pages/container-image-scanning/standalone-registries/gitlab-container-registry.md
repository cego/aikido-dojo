# Gitlab Container Registry

You can now integrate your (cloud) GitLab Container Registry with Aikido to scan your containers for known vulnerabilities.

Follow the simple steps below to activate this feature:

**Step 1:** Log into your GitLab account to gather some data.

We'll have to gather your username (see screenshot)

[image: User profile dropdown displaying name, username, and status option in a web application.]

**Step 2:** Copy the group ID where the container registry resides (see screenshot)

[image: Organization profile showing group name, privacy status, plan type, and group ID.]

**Step 3:** Under personal preferences, Access tokens, create a new token for Aikido (direct link: <https://gitlab.com/-/user_settings/personal_access_tokens>).

The scopes included must be: read\_api, read\_registry

See screenshot:

[image: Form for creating a new personal access token with selectable API access scopes.]

**Step 4:** Enter the collected data in Aikido (direct link: <https://app.aikido.dev/settings/container-image-registry/add/gitlab>)

[image: GitLab container registry connection form: enter username, access token, and group ID.]

**Step 5:** Aikido will now find all container repositories you can access and list them.

**Step 6:** Repositories can be linked to a code repository in order to perform better deduplication of findings. This step is optional!

**Step 7:** In the action menu next to the registry, click 'scan repos in registry' to get started. Results will appear in the Feed!

***
