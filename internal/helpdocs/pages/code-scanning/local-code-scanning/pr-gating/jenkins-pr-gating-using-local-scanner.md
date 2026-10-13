# Jenkins PR gating using Local Scanner

The Aikido Local Scanner can be used to enforce security gates in your CI pipeline.

{% hint style="info" %}
You can also use PR gating with your cloud-connected workspaces (GitHub, BitBucket, etc). Check out our CI [Integration section](https://help.aikido.dev/section/ci-integrations/sg3q6UrIf4qE).
{% endhint %}

PR gating ensures that new code meets your security standards before it is merged into the default branch. It scans only the changes introduced in the pull request. If the scan detects new issues that meet or exceed your configured severity threshold, the CI pipeline fails.

This prevents new vulnerabilities from being introduced while allowing existing findings to be addressed separately.

{% hint style="info" %}
We also support [release gating](/container-image-scanning/local-image-scanning/release-gating-for-container-images-using-local-image-scanner.md) for teams that want to enforce checks at release time.
{% endhint %}

To enable PR gating, add the `--fail-on <severity>` option to select your preferred severity level. Then, add the `--gating-mode pr` option to signify that you wish to perform PR gating. You must also specify the base (`--base-commit-id <commit-id>`) and head commit (`--head-commit-id <commit-id>`). If there a scan was previously performed on the base commit id, the scan results will be compared to those. If not, the results will be compared to the most recent scan on your default branch.

For general information about setting up the Local Scanner in Jenkins environment, check out [this article](/code-scanning/local-code-scanning/github-action-setup-for-local-code-scanning.md).

#### How it works

There are three moving parts:

1. **Integration between your SCM and Jenkins.**
2. **A Multibranch Pipeline job** with the *Discover pull requests* behavior. This makes Jenkins build each open pull request. On a PR build, Jenkins checks out the source branch and exposes the standard change-request environment variables (`CHANGE_ID`, `CHANGE_TARGET`, `CHANGE_BRANCH`).
3. **The Aikido Local Scanner**, run from your `Jenkinsfile` in `--gating-mode pr`.

```
PR opened/updated  ──►  SCM sends webhook  ──►  Jenkins Multibranch PR build
                                                      │
                                                      ▼
                                        aikido-local-scanner --gating-mode pr
                                                      │
                                  new issues ≥ --fail-on?  ── yes ──► build FAILS
                                                      │                    │
                                                      no                   ▼
                                                      ▼            build status = FAILED
                                                build PASSES               │
                                                      │                    ▼
                                                      ▼          SCM blocks merge
                                              SCM checks pass
```

#### Multibranch Pipeline setup

Pull request builds require a **Multibranch Pipeline** job.

1. In Jenkins, select **New Item → Multibranch Pipeline**, name it, and select **OK**.
2. Under **Branch Sources**, select **Add source → Select your SCM**.
3. Choose your **repository**.
4. Under **Behaviors**, select **Add → Discover pull requests** and disable **Discover branches.**
5. Under **Scan Multibranch Pipeline Triggers**, enable the **webhook trigger → Pull request opened or source branch updated** option.
6. Leave **Build Configuration** as **by Jenkinsfile**, with **Script Path** = `Jenkinsfile`.
7. Select **Save**. Jenkins scans the repository and starts a build for each open pull request.

#### Add Aikido PR gating stage to your Jenkinsfile

Add (or extend) a `Jenkinsfile` at the root of your repository. The stage below runs **only on pull request builds**, derives the base and head commits from Git, and runs the scanner in PR gating mode. The scanner's non-zero exit on a failing gate is what fails the build.

```groovy
pipeline {
    agent any

    environment {
        // Picks up credentials as setup in https://help.aikido.dev/code-scanning/local-code-scanning/jenkins-setup-for-local-code-scanning#id-1-get-your-authentication-token
        AIKIDO_API_KEY = credentials('aikido-local-scanner-api-key')
        AIKIDO_REPO_NAME = 'my-repo'
    }

    stages {
        stage('Aikido PR Gating') {
            // CHANGE_ID is only set on pull request builds
            when { expression { env.CHANGE_ID != null } }
            steps {
                script {
                    // Ensure the target (destination) branch is available locally
                    sh "git fetch origin ${env.CHANGE_TARGET}"

                    // Head = tip of the PR source branch (what's being merged in)
                    def headCommit = sh(
                        script: 'git rev-parse HEAD',
                        returnStdout: true
                    ).trim()

                    // Base = tip of the destination branch (what it's merging into)
                    def baseCommit = sh(
                        script: "git rev-parse origin/${env.CHANGE_TARGET}",
                        returnStdout: true
                    ).trim()

                    echo "PR #${env.CHANGE_ID}: ${env.CHANGE_BRANCH} -> ${env.CHANGE_TARGET}"
                    echo "Base commit: ${baseCommit}"
                    echo "Head commit: ${headCommit}"

                    // Docker variant. The scanner exits non-zero when new issues
                    // at or above --fail-on are found, which fails the build.
                    sh """
                        docker run --rm \
                          -v "\$(pwd):/scan-target" \
                          aikidosecurity/local-scanner \
                          scan /scan-target \
                          --apikey "\$AIKIDO_API_KEY" \
                          --repositoryname "\$AIKIDO_REPO_NAME" \
                          --branchname "${env.CHANGE_BRANCH}" \
                          --gating-mode pr \
                          --fail-on critical \
                          --base-commit-id ${baseCommit} \
                          --head-commit-id ${headCommit}
                    """
                }
            }
        }
    }
}
```

**If the scanner binary is installed on the agent** instead of Docker, replace the `docker run …` block with:

```groovy
sh """
    aikido-local-scanner scan ./ \
      --apikey "\$AIKIDO_API_KEY" \
      --repositoryname "\$AIKIDO_REPO_NAME" \
      --branchname "${env.CHANGE_BRANCH}" \
      --gating-mode pr \
      --fail-on critical \
      --base-commit-id ${baseCommit} \
      --head-commit-id ${headCommit}
"""
```

* **`--fail-on critical`** sets the severity threshold. Change it to `high`, `medium`, or `low` to gate more strictly.

Setup the gate enforcement in your SCM to match your needs.
