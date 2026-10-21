# GitHub Container Registry

You can integrate your GitHub Container Registry with Aikido to scan your containers for known vulnerabilities.

{% hint style="info" %}
When connecting, make sure the `Inherit access from source repository` setting is **enabled** for all the containers you want to scan. Otherwise Aikido won't be able to read those.
{% endhint %}

Follow the simple steps below to activate this feature:

**Step 1:** Log into your GitHub account to gather some data.

We'll have to gather your username (see screenshot)

[image: GitHub user profile overview with avatar, username, followers, and edit profile option.]

**Step 2:** Copy the organisation name where the container registry resides. This is visible in the github-url (see screenshot)

[image: URL input field displaying "https://github.com/Organization" for organization link entry.]

**Step 3:** Under profile settings, developer settings, Personal access tokens, Tokens (classic), generate a new classic token for Aikido

(direct link: <https://github.com/settings/tokens>).

the scope includes: `read:packages`

[image: GitHub token creation: No expiration, read package access scope selected.]

**Step 4:** Enter the collected data in Aikido (direct link: <https://app.aikido.dev/settings/container-image-registry/add/github>)

[image: Form to connect and authenticate a GitHub container image registry.]

**Step 5:** Aikido will now find all container repositories you can access and list them.

**Step 6:** Repositories can be linked to a code repository in order to perform better deduplication of findings. This step is optional!

**Step 7:** In the action menu next to the registry, click 'scan repos in registry' to get started. Results will appear in the Feed!

### &#x20;<a href="#set-up-github-container-registry-scanning" id="set-up-github-container-registry-scanning"></a>
