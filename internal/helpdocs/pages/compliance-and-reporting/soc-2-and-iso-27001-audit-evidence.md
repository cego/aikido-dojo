# SOC 2 and ISO 27001 Audit Evidence

Aikido evidences the technical controls of SOC 2 and ISO 27001:2022. For each covered control, Aikido runs named checks against your repositories, cloud accounts, containers, and domains. It reports each control as 'complying', 'failing', or 'unknown'.

Organizational controls such as the control environment, user access reviews, physical security, and incident response fall outside Aikido's scope and are listed under [What Aikido does not cover](#what-aikido-does-not-cover) so your auditor knows where other evidence is needed.

Go to **Reports** → **Compliance** in the Aikido app to open the reports described here.

## The SOC 2 compliance report

The mapping is the same for every Aikido workspace, but the results can differ based on the clouds you connect, your teams, and the features you enable. The checks named below are summaries. The report lists each individual check against each resource.

| Criterion | Control objective                                                                                                                                                                                                             | What Aikido checks                                                                                                                                                                                                                                                                                                                                       |
| --------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| CC3.2     | Estimate significance of risks identified                                                                                                                                                                                     | Identity of cloud users, threat detection, severe surface monitoring and open source dependency issues, monitoring of repositories, container images, and domains                                                                                                                                                                                        |
| CC3.3     | Consider the potential for fraud                                                                                                                                                                                              | Least privilege for cloud resources, multi-factor authentication (MFA) for cloud users, security logging, threat detection                                                                                                                                                                                                                               |
| CC5.2     | Selects and develops general control activities over technology to support the achievement of objectives                                                                                                                      | Least privilege, deletion protection for cloud resources, identity of cloud users                                                                                                                                                                                                                                                                        |
| CC6.1     | Restricts logical access. Consider network segmentation. Restrict access to information assets. Manages credentials for infrastructure and software. Use encryption to protect data.                                          | MFA, least privilege, encryption at rest and in transit, exposure of sensitive data, protection against SQL injection, server-side request forgery (SSRF), command injection, and cross-site scripting (XSS), public file storage, SSH access, separate production and test environments, load balancer access points, up to date cryptography libraries |
| CC6.6     | Restrict access. Require additional authentication or credentials. Implement boundary protection system.                                                                                                                      | Public access to cloud resources, enforced MFA, DNS zone transfer lock, least privilege                                                                                                                                                                                                                                                                  |
| CC6.7     | Use encryption technologies or secure communication channels to protect data                                                                                                                                                  | Latest Transport Layer Security (TLS) version, up to date cryptography libraries                                                                                                                                                                                                                                                                         |
| CC6.8     | Restrict application and software installation. Detect unauthorized changes to software and configuration parameters. Use anti-virus and anti-malware software.                                                               | Runtime access protection, security logging, threat detection, Aikido malware scanner                                                                                                                                                                                                                                                                    |
| CC7.1     | Monitor infrastructure and software. Implement change detection mechanism. Detect unknown or unauthorized components. Conduct vulnerability scans. Implement filters to analyze anomalies. Restore the affected environments. | Security logging, service level agreements (SLAs) configured and met, connected repositories, cloud and public domains, alerting, security contact, risky licenses, lockfiles, threat detection, up to date runtimes, critical open source issues                                                                                                        |
| CC8.1     | Protect confidential information. Track system changes.                                                                                                                                                                       | Exposed secrets, issue tracker integration                                                                                                                                                                                                                                                                                                               |

## ISO 27001:2022 controls Aikido evidences

| Control | Title                                                       | What Aikido checks                                                                                                                   |
| ------- | ----------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| A.5.15  | Access control                                              | Least privilege for cloud users and resources, exposure of sensitive data                                                            |
| A.5.16  | Identity management                                         | Identity of cloud users                                                                                                              |
| A.5.28  | Collection of evidence                                      | Budget alerts, security logging, threat detection                                                                                    |
| A.5.33  | Protection of records                                       | Public access to cloud resources, encryption at rest, high availability                                                              |
| A.8.2   | Privileged access rights                                    | Permission enforcement checks                                                                                                        |
| A.8.3   | Information access restriction                              | Public access, access controls for cloud resources                                                                                   |
| A.8.5   | Secure authentication                                       | MFA for cloud users, encryption in transit                                                                                           |
| A.8.6   | Capacity management                                         | Budget alerts, load balancer configuration                                                                                           |
| A.8.7   | Protection against malware                                  | Unwanted filesystem writes, threat detection, lockfiles                                                                              |
| A.8.8   | Management of technical vulnerabilities                     | No issues outside their SLA                                                                                                          |
| A.8.9   | Configuration management                                    | Up to date cloud instance versions, lockfiles                                                                                        |
| A.8.12  | Data leakage prevention                                     | Secure file storage, encryption at rest and in transit, protection against remote code execution (RCE), SSRF, SQL injection, and XSS |
| A.8.13  | Backups                                                     | Backups for stateful resources, cross account backups                                                                                |
| A.8.15  | Logging                                                     | Security logging for cloud instances                                                                                                 |
| A.8.16  | Monitoring activities                                       | Connected cloud, real time security alerts, threat detection                                                                         |
| A.8.18  | Use of privileged utility programs                          | Secure access, MFA, exposure of sensitive data, access controls                                                                      |
| A.8.20  | Network security                                            | Unauthorized network access, SSH from anywhere                                                                                       |
| A.8.24  | Use of cryptography                                         | Safe SSL protocols, secure cookies, up to date cryptographic libraries                                                               |
| A.8.25  | Secure development lifecycle                                | Continuous integration (CI) integration, connected repository and cloud                                                              |
| A.8.28  | Secure coding                                               | No issues outside their SLA                                                                                                          |
| A.8.31  | Separation of development, test and production environments | Separate production and test environments                                                                                            |
| A.8.32  | Change management                                           | Issue tracker integration                                                                                                            |

## Who approves a finding decision?

The approval and remediation path is the same for every control.

* **Approval owner.** Users with ignore or snooze permission dismiss a finding directly. There is no mandatory second approver. Admins can turn on ignore requests under [Triage settings](https://app.aikido.dev/settings/advanced), after which users without the permission submit a request with a reason, and a reviewer approves or rejects it in the [Approval Inbox](/getting-started/core-functionalities/approval-inbox-for-ignored-issues.md). Each request records who asked, why, who decided, and when. See [Setting Roles and Permissions](/getting-started/automated-user-management/setting-roles-and-permissions.md) for which roles hold these permissions.
* **Remediation deadline.** [SLAs](/getting-started/core-functionalities/enable-slas-in-aikido.md) are set per severity in calendar days and start when Aikido first detects an issue. The check "Does not have any issues outside of their SLA" evidences CC7.1, A.8.8, and A.8.28.
* **Remediation action.** Findings flow into your issue tracker, which evidences CC8.1 and A.8.32. AutoFix pull requests are recorded in the activity log.
* **Snooze.** A finding can be snoozed for a day, a week, a month, a year, or a custom period, with a reason.

## Trace a finding across rescans

Use the numeric issue ID from the API together with the repository or asset and the path. An identifier of the form `AIKIDO-2026-922647` names the vulnerability, not the finding, so the same identifier can appear in several repositories, paths, or assets.

On every rescan, Aikido deduplicates by grouped issue, repository, target path, and vulnerability identifier, and reuses the existing issue record. The numeric issue ID stays stable.

## What Aikido does not cover

Aikido can only evidence technical controls. The following SOC 2 criteria are not technical controls and require evidence from other sources:

* **Availability:** A1.1 capacity management, A1.2 environmental and recovery infrastructure.
* **Confidentiality, Processing Integrity, and Privacy:** All criteria.
* **Security (Common Criteria):** CC1 control environment, CC2 communication and information, CC3.1 and CC3.4 risk identification, CC4 monitoring activities, CC5.1 and CC5.3 control activities, CC6.2 user registration and access, CC6.3 role based access management, CC6.4 physical access, CC6.5 decommissioning of assets, CC7.2 to CC7.5 anomaly monitoring, security event evaluation, incident response and recovery, and CC9 risk mitigation.

## FAQ

### Does Aikido provide audit-ready evidence for SOC 2 and ISO 27001?

Yes, for the technical controls: 23 SOC 2 control objectives across 9 criteria and 22 ISO 27001:2022 Annex A controls.

### Can Aikido help satisfy customer security reviews?

Yes. The [Security Audit Report](/compliance-and-reporting/reports/security-audit-report.md) is built for customers, auditors, and leadership. You choose its sections, for example SLA compliance, scan history, and compliance frameworks.

### Does Aikido keep an audit trail of accepted risks and suppressed findings?

Yes. Every ignore, snooze, re-open, close, and severity change is logged with the user, date, comment, and issue, exportable as CSV or JSON.

### Does Aikido replace a compliance platform like Vanta or Drata?

No. Aikido evidences the technical controls and pushes its results to compliance platforms such as [Drata](/compliance-and-reporting/iso-27001-and-soc-2-compliance-overview/drata-compliance-integration.md), [Sprinto](/compliance-and-reporting/iso-27001-and-soc-2-compliance-overview/sprinto-compliance-integration.md), and [Thoropass](/compliance-and-reporting/iso-27001-and-soc-2-compliance-overview/thoropass-compliance-integration.md).

### Which compliance frameworks does Aikido report on?

ISO 27001:2022, SOC 2, CIS v8.1, CIS AWS Benchmark, NIS2, NIST 800-53, PCI DSS v4.0, HIPAA, DORA, HITRUST, GDPR, ENS, and UK Cyber Essentials. See [Compliance Reporting](/compliance-and-reporting/reports/compliance-reporting.md) for the full list of report pages.
