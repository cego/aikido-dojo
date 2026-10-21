# Manually Link Repository to Container

Linking your container images to their corresponding code repositories centralizes related vulnerabilities, so you can view, triage, and fix them in the right context. Once linked:

* Container issues will appear directly within the relevant repo.
* You'll unlock Container Autofix, which updates Dockerfiles automatically with available fixes.

#### Ways to Link Containers

You can link containers to repositories in two ways:

* **Manual Linking** – Link specific repos to containers one by one.
* **Smart Suggestions** – Automatically link multiple containers using intelligent repo suggestions. [Learn more here →](https://help.aikido.dev/container-image-scanning/configuration/automatically-link-repository-to-container-via-smart-suggestions)

## Manually link a repository to a container <a href="#how-to-link-a-repository-to-a-container" id="how-to-link-a-repository-to-a-container"></a>

**Step 1**: Navigate to [Settings > Containers](https://app.aikido.dev/settings/container-image-registry) via the top right corner of the Aikido interface.

**Step 2:** Click the 🔗 **Link** button next to the container you want to connect.\
Use the search bar to find and select the correct repository.

[image: Docker image "openscurity/mobsfscan:latest" prompts to link an unlinked code repository.]

[image: Dialog for linking a cloud image to the "Spoon-Knife" code repository.]

**Step 3**: To see all repositories linked to a container, go to in [Settings > Containers](https://app.aikido.dev/settings/container-image-registry). Use the `⋮` **three-dot menu** next to any linked repo to unlink it.
