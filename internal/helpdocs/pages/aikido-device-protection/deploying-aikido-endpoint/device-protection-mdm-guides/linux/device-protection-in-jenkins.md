# Device Protection in Jenkins

Jenkins builds run on Linux agents, so Device Protection installs the same way it does on any Linux machine. What is specific to Jenkins is where the install goes: once on a permanent agent, or inside each build on a Docker agent.

Read [Install Device Protection for Linux](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux.md) first for the token and the install flavors. The job-level snippets below use the container flavor, which [Device Protection in Containers](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/device-protection-in-containers.md) explains.

## Permanent agents

Install the agent on each Linux agent host with `--headless`, then reboot. See [Install Device Protection for Linux](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux.md#choose-the-install-flavor). No pipeline changes needed.

## Docker agents

Add your Device Protection token as a **Secret text** credential with the ID `aikido-token`, then start the agent before any stage that installs packages:

{% code title="Jenkinsfile" overflow="wrap" %}

```groovy
pipeline {
  agent {
    docker {
      image 'node:20-bookworm'
      args '-u root'
    }
  }
  environment {
    AIKIDO_TOKEN = credentials('aikido-token')
  }
  stages {
    stage('Start Device Protection') {
      steps {
        sh '''
          curl -fsSL https://github.com/AikidoSec/safe-chain/raw/refs/heads/main/install-scripts/install-endpoint-linux.sh | sh -s -- --container
          aikido-endpoint-protection-start --wait-for-healthy
        '''
      }
    }
    stage('Install dependencies') {
      steps {
        sh '''
          for f in /etc/profile.d/aikido-*.sh; do . "$f"; done
          npm ci
        '''
      }
    }
  }
}
```

{% endcode %}

* **`args '-u root'` is required.** Jenkins runs the container as the agent's user, and the install needs root.
* **Source `/etc/profile.d/aikido-*.sh` in every `sh` step that installs packages.** Each `sh` step is a new shell, and without those settings downloads skip the agent.
* **Use `'''`, not `"""`.** Groovy expands `$` inside double quotes and breaks the loop.

## Verify it works

{% code title="Jenkinsfile" overflow="wrap" %}

```groovy
stage('Verify Device Protection') {
  steps {
    sh '''
      for f in /etc/profile.d/aikido-*.sh; do . "$f"; done
      curl -sS -o /dev/null -w "%{http_code}\\n" --max-time 15 https://registry.npmjs.org/safe-chain-test/-/safe-chain-test-0.0.1-security.tgz
    '''
  }
}
```

{% endcode %}

`403` means the agent blocked the malware test package. `200` means the download skipped the agent.

## Troubleshooting

| Problem                                                    | Fix                                                                 |
| ---------------------------------------------------------- | ------------------------------------------------------------------- |
| `Root privileges required`                                 | Add `args '-u root'` to the `docker` agent                          |
| `Token is required`                                        | Check the credential ID, and that the credential is **Secret text** |
| Test package returns `200`, or `SELF_SIGNED_CERT_IN_CHAIN` | The `sh` step did not source `/etc/profile.d/aikido-*.sh`           |
