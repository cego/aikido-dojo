# Connect Bitbucket Cloud

## Connect Bitbucket Account to Aikido

Aikido requests read-only access to your Bitbucket organization to analyze your repositories. After analysis, your code is always wiped from the system.

### 1. Logging in using Bitbucket <a href="#id-1-logging-in-using-bitbucket" id="id-1-logging-in-using-bitbucket"></a>

To get started, navigate to <https://app.aikido.dev/> and log in with Bitbucket. This will look like the screenshot below. Here, Aikido only requests access to your identity on Bitbucket and the associated email address.

[image: Access request screen for account information by Aikido User Authentication, with grant or cancel options.]

### 2. Authorizing access to an organization <a href="#id-2-authorizing-access-to-an-organization" id="id-2-authorizing-access-to-an-organization"></a>

On the next screen, you can choose to connect a real organization or a sample workspace. If you choose a real organization you will be redirected back to Bitbucket. Once there, pick the organization you would like to authorize.

[image: Aikido requests account and repository access; user must approve or deny.]

### 3. Select your repos <a href="#id-3-select-your-repos" id="id-3-select-your-repos"></a>

You can optionally grant access to 1 or 2 repositories instead of all repositories as seen below:<br>

### 4. Checking results <a href="#id-4-checking-results" id="id-4-checking-results"></a>

After granting access and validating the repositories you want to scan, Aikido will automatically start scanning. After about 1 minute, you should see the first results come in!
