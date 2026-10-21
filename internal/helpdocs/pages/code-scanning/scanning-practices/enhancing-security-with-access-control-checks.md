# Access Control Checks

Aikido's checks on Access Controls offers robust security by informing about critical access control practices. This way, you can ensure that only authorized and verified changes are made to your codebase. Some examples of checks are multi-factor authentication, restricting default access rights, and requiring mandatory code reviews.

> All Access Controls checks [can be found here](https://app.aikido.dev/repositories/access_control).

### Prerequisite <a href="#prerequisite" id="prerequisite"></a>

* Only available for GitHub Cloud & GitLab Cloud connected workspaces.

### Access Control Setup for GitHub <a href="#access-control-setup-for-github" id="access-control-setup-for-github"></a>

> For GitLab, no extra authorisation steps need to be taken.

**Step 1.** In the Main Feed, filter on **Access Controls.** Click **Authorise on GitHub** in order to allow Aikido scan for configurations related the access controls.

[image: Access controls filter active; Aikido requests extra GitHub permissions for analysis.]

**Step 2.** In GitHub, grant permissions to **install** the Aikido GitHub Config Scanner. It is recommended to select **All Repositories**.

[image: Installing Aikido GitHub Config Scanner with repository and organization read access.]

**Step 3.** After connecting, Aikido will do a scan for [checks mentioned here. ](https://app.aikido.dev/repositories/access_control)After a couple of minutes, you will be able to view them in the Aikido feed. The sidebar will give more information about which repos need configuration adjustments.

[image: Access control security issues listed by severity and estimated fix time.]
