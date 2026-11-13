# Running a Successful Pilot with Aikido

A pilot should answer two questions fast:

1. Does Aikido fit your environment and workflows?
2. Does Aikido improve signal, coverage, and speed for your team?

Most pilots run for 1 to 2 weeks. The goal is not to configure everything. The goal is to validate the parts that matter for your team.

### Suggested pilot plan

1. **Align on scope and success criteria.** Decide which teams, assets, and workflows you want to evaluate. Pick a clear owner for the pilot and a short list of must-have outcomes.
2. **Create your workspace and connect your first repositories.** Start by [creating your workspace](https://app.aikido.dev/login?regionRedirect=1). Then follow [Connect Your Repositories](https://help.aikido.dev/getting-started/setting-up-your-account/create-account-and-connect-your-repositories). If data handling is part of your review, also read [Aikido Never Stores Your Code](https://help.aikido.dev/getting-started/setting-up-your-account/aikido-never-stores-your-code).
3. **Validate the core product coverage.** Cover the main scanners and workflows you expect to use in production. Start with the core checklist below.
4. **Invite the right stakeholders.** Add security, platform, and developer leads early. Use [User Management](https://help.aikido.dev/getting-started/automated-user-management) to get the right people into the workspace.
5. **Review results against your buying criteria.** Measure setup time, issue quality, workflow fit, and expected ROI. Use the evaluation section below as a starting point.

{% hint style="info" %}
Need more time before the pilot starts? Your workspace can fall back to the free plan until you are ready. Your setup stays in place.
{% endhint %}

### Core coverage to validate

These are the most common areas to validate during a pilot:

1. **Code scanning.** Start with [Code Scanning Overview](https://help.aikido.dev/code-scanning/code-scanning-overview) and connect your main repositories first.
2. **Cloud scanning.** Review [Cloud Scanning Overview](https://help.aikido.dev/cloud-scanning/cloud-scanning-overview) if cloud posture management is in scope.
3. **Containers and virtual machines.** Validate [Container Image Scanning Overview](https://help.aikido.dev/container-image-scanning/container-image-scanning-overview) and/or [VM Scanning Overview](https://help.aikido.dev/virtual-machine-scanning/virtual-machines-overview) based on your runtime footprint.
4. **DAST and attack surface.** Review [DAST Overview](https://help.aikido.dev/dast-surface-monitoring/dast-surface-monitoring-overview) if you want live validation on apps or APIs.
5. **Task tracking and ownership.** Connect your workflows through [Task Management Tools](https://help.aikido.dev/getting-started/task-management-systems).
6. **Team access.** Make sure findings reach the right teams through [Inviting Users to Aikido](https://help.aikido.dev/getting-started/automated-user-management/inviting-users-to-aikido).

### Optional pilot tracks

If you want broader validation, add one or more of these tracks:

1. **Attack visibility and exploit paths.** Review [Reachability Analysis](https://help.aikido.dev/getting-started/reachability-analysis/introduction-to-reachability-analysis), [Attack Surface Monitoring](https://help.aikido.dev/dast-surface-monitoring/attack-surface-scanning), and [Threat Model](https://help.aikido.dev/pentests/coverage-and-findings/threat-model).
2. **Developer workflow coverage.** Test [IDE Plugins](https://help.aikido.dev/ai-and-dev-tools/ide-plugins-overview), [AI Coding Assistants (MCP)](https://help.aikido.dev/ai-and-dev-tools/aikido-mcp), and [PR Gating](https://help.aikido.dev/pr-and-release-gating/aikido-ci-gating-functionality).
3. **Malware and package protection.** Evaluate [Safe Chain](https://help.aikido.dev/code-scanning/aikido-malware-scanning), [Malware Detection in Open-Source Dependencies and Containers](https://help.aikido.dev/code-scanning/scanning-practices/malware-detection-in-open-source-dependencies), and [Aikido Device Protection](https://help.aikido.dev/aikido-device-protection/endpoint-protection).
4. **Advanced validation.** Run [Code Audit](https://help.aikido.dev/ai-code-analysis/ai-code-audit-overview) if you want deeper code-level reasoning. Run [Aikido Pentest](https://help.aikido.dev/pentests/aikido-pentest) if you want live target testing.
5. **Private and self-hosted environments.** If you scan internal apps or self-managed systems, review [Aikido Broker for Internal Applications](https://help.aikido.dev/miscellaneous-info/aikido-broker-for-internal-applications).
6. **Credit-based features.** If you plan to test credit-backed workflows, review [Wallet & Credits](https://help.aikido.dev/miscellaneous-info/wallet-and-credits) before you start.

### Suggested evaluation criteria

Use these questions to structure your pilot review:

#### Time to value

* How long did first setup take?
* How quickly did useful findings appear?
* Did onboarding feel clear for admins and developers?

#### Signal quality

* Did the findings feel relevant?
* Did noise reduction improve focus compared to other tools?
* Did prioritization help the team decide what to fix first?

#### Workflow fit

* Does Aikido fit your Git, CI, ticketing, and chat workflows?
* Can the right teams see and act on findings?
* Do the integrations support how your teams already work?

#### Breadth and depth

* Does Aikido cover the assets you care about most?
* Are the product areas deep enough for your requirements?
* Do advanced workflows like IDE scanning, PR gating, Code Audit, or Pentest add value for your team?

#### Buying confidence

* Does the product replace enough tools or manual work to justify the spend?
* Can you show expected ROI through better coverage, faster remediation, or less noise?
* Does the support model give your team confidence?

### Getting help during the pilot

You can work with the Aikido team in three ways:

1. **In-app chat.** Best for fast setup and troubleshooting questions.
2. **Slack or Microsoft Teams.** Good for shared pilot communication with your team.
3. **Documentation.** Use the docs to go deeper on setup, product behavior, and troubleshooting.
