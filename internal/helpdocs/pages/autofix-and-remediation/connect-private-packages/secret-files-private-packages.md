# Secret Files - Private packages

For Aikido to update dependencies that include private packages, it needs **access to your private registries** so it can generate accurate lockfile updates.

Many container-based systems (like Kubernetes and Docker) use **Secret Files** that allow you to mount a secret on the filesystem. In Aikido you can provide secrets files, which will be stored encrypted and injected into your workflows.

During Autofix, the Secret File will be written to a temporary file on the filesystem. The path of the temporary file will be stored in the environment variable you provide.

This means when our system detects vulnerabilities in dependencies, it can seamlessly authenticate with private registries, automatically patch the affected packages, and update the lockfiles, all while keeping your credentials safe.

## Configuration in Aikido

Once write access to the repos is set up, you can configure Aikido to authenticate with your private registry by following the steps below:

1. Go to your account's settings page for AutoFix, [here](https://app.aikido.dev/issues/fix/settings).
2. Click on "*Connect registry*" to see the modal below

3. When you select *"Secret Files"* you will be able to enter the content of the Secret File and the name of the environment variable that will be used to read the path to the secret file.

4. Fill in the "Secret Content" and "Environment Variable" and add as many secret files as needed.
