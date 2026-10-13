# Connect GitHub Organization

Aikido requests read-only access to your GitHub organization to analyze your repositories. We use the new GitHub App system, so we don't have to store any tokens in our database at all. After analysis, your code is always wiped from the system.

#### 1. Logging in using GitHub <a href="#id-1-logging-in-using-github" id="id-1-logging-in-using-github"></a>

To get started, navigate to <https://app.aikido.dev/> and log in with GitHub. This will look like the screenshot below. Here, Aikido only requests access to your identity on GitHub and the associated email address.

[image: Aikido Security requests GitHub authorization to access user identity and email addresses.]

#### 2. Authorizing access to an organization <a href="#id-2-authorizing-access-to-an-organization" id="id-2-authorizing-access-to-an-organization"></a>

On the next screen, you can choose to connect a real organization or a sample workspace. If you choose a real organization you will be redirected back to GitHub. Once there, pick the organization you would like to authorize. You can optionally grant access to 1 or 2 repositories instead of all repositories as seen below:

[image: Authorize Aikido Security to access all repositories and read email addresses on GitHub.]

#### 3. Checking results <a href="#id-3-checking-results" id="id-3-checking-results"></a>

After granting access and validating the repositories you want to scan, Aikido will automatically start scanning. After about 1 minute, you should see the first results come in!
