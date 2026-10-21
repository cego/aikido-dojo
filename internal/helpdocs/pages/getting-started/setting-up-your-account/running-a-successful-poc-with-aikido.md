# Running a Successful Pilot with Aikido

A pilot should answer two questions fast:

1. Does Aikido fit your environment and workflows?
2. Does Aikido improve signal, coverage, and speed for your team?

Most pilots run for 1 to 2 weeks. The goal is not to configure everything. The goal is to validate the parts that matter for your team.

### Suggested pilot plan

1. **Align on scope and success criteria.** Decide which teams, assets, and workflows you want to evaluate. Pick a clear owner for the pilot and a short list of must-have outcomes.
2. **Create your workspace and connect your first repositories.** Start by [creating your workspace](https://app.aikido.dev/login). Then follow [Connect Your Repositories](/getting-started/setting-up-your-account/create-account-and-connect-your-repositories.md). If data handling is part of your review, also read [Aikido Never Stores Your Code](/getting-started/setting-up-your-account/aikido-never-stores-your-code.md).
3. **Validate the core product coverage.** Cover the main scanners and workflows you expect to use in production. Start with the core checklist below.
4. **Invite the right stakeholders.** Add security, platform, and developer leads early. Use [User Management](/getting-started/automated-user-management.md) to get the right people into the workspace.
5. **Review results against your buying criteria.** Measure setup time, issue quality, workflow fit, and expected ROI. Use the evaluation section below as a starting point.

{% hint style="info" %}
Need more time before the pilot starts? Your workspace can fall back to the free plan until you are ready. Your setup stays in place.
{% endhint %}

### Core coverage to validate

These are the most common areas to validate during a pilot:

1. **Code scanning.** Start with [Code Scanning Overview](/code-scanning/code-scanning-overview.md) and connect your main repositories first.
2. **Cloud scanning.** Review [Cloud Scanning Overview](/cloud-scanning/cloud-scanning-overview.md) if cloud posture management is in scope.
3. **Containers and virtual machines.** Validate [Container Image Scanning Overview](/container-image-scanning/container-image-scanning-overview.md) and/or [VM Scanning Overview](/virtual-machine-scanning/virtual-machines-overview.md) based on your runtime footprint.
4. **DAST and attack surface.** Review [DAST Overview](/dast-surface-monitoring/dast-surface-monitoring-overview.md) if you want live validation on apps or APIs.
5. **Task tracking and ownership.** Connect your workflows through [Task Management Tools](/getting-started/task-management-systems.md).
6. **Team access.** Make sure findings reach the right teams through [Inviting Users to Aikido](/getting-started/automated-user-management/inviting-users-to-aikido.md).

### Optional pilot tracks

If you want broader validation, add one or more of these tracks:

1. **Attack visibility and exploit paths.** Review [Reachability Analysis](/getting-started/reachability-analysis/introduction-to-reachability-analysis.md), [Attack Surface Monitoring](/dast-surface-monitoring/attack-surface-scanning.md), and [Threat Model](/pentests/coverage-and-findings/threat-model.md).
2. **Developer workflow coverage.** Test [IDE Plugins](/ai-and-dev-tools/ide-plugins-overview.md), [AI Coding Assistants (MCP)](/ai-and-dev-tools/aikido-mcp.md), and [PR Gating](/pr-and-release-gating/aikido-ci-gating-functionality.md).
3. **Malware and package protection.** Evaluate [Safe Chain](/code-scanning/aikido-malware-scanning.md), [Malware Detection in Open-Source Dependencies and Containers](/code-scanning/scanning-practices/malware-detection-in-open-source-dependencies.md), and [Aikido Device Protection](/aikido-device-protection/endpoint-protection.md).
4. **Advanced validation.** Run [Code Audit](/ai-code-analysis/ai-code-audit-overview.md) if you want deeper code-level reasoning. Run [Aikido Pentest](/pentests/aikido-pentest.md) if you want live target testing.
5. **Private and self-hosted environments.** If you scan internal apps or self-managed systems, review [Aikido Broker for Internal Applications](/miscellaneous-info/aikido-broker-for-internal-applications.md).
6. **Credit-based features.** If you plan to test credit-backed workflows, review [Wallet & Credits](/miscellaneous-info/wallet-and-credits.md) before you start.

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
