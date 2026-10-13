# Track Progress and Findings

Follow an AI Code Analysis while it runs and read the findings it produces.

Once you've started an analysis, Aikido takes you to the assessment detail page. From there, you can watch progress and read findings as they land.

## Tracking progress

Open an assessment from the **AI Code Analysis** home to see its detail page. Every assessment has an **Overview** tab and an **Issues** tab. The **Issues** tab shows a badge when there are open issues.

* **Code Security Audit** and **Smart Contract Audit** also have a **Threat Model** tab. See [Threat Model](/pentests/coverage-and-findings/threat-model.md).
* **Data Exposure Audit** has a **Sensitive Data** tab instead of **Threat Model**.
* **Cloud Security Audit** has **Overview** and **Issues** only.

While the analysis is running, the **Overview** tab shows the current status and the header shows elapsed time and an estimated time remaining.

The **AI Code Analysis Information** block lists who requested the analysis and when it started. **Code & Documentation** lists the repositories linked to the analysis. Cloud Security Audits also list the cloud accounts.

### Analysis statuses

The status on the **Overview** tab reflects where the analysis is in its lifecycle:

* **Pending:** The analysis is queued and waiting for agent capacity.
* **Running:** Agents are actively reasoning over the selected code or cloud inventory. Findings can land in the **Issues** tab while this status is showing.
* **Completed:** The analysis finished. All findings are in the **Issues** tab and in the Feed.
* **Failed:** The analysis couldn't finish. Open the assessment for details, or [contact support](/miscellaneous-info/contact-support.md).
* **Cancelled:** The analysis was stopped before it completed.

## Reading findings

Findings appear in the **Issues** tab on the assessment and in the Feed alongside other Aikido findings. They follow the standard Aikido finding lifecycle. See [Manage Findings](/getting-started/core-functionalities.md) for ignore, snooze, severity, and task-tracker integrations.

Click **Attack Analysis** on any finding to open the standard Aikido finding modal, with sections for summary, description, risk, root cause, remediation, and more.

The reproduction block is titled **Evidence**. On a Code Security, Data Exposure, or Smart Contract audit, it points at the files and logic behind the finding. On a Cloud Security Audit, it points at the cloud resources, and at infrastructure as code files when you included repositories. Treat it as the justification for the finding, rather than a live HTTP exploit.

After you push a fix, click **Re-test** on the finding. Aikido checks a Code Security, Data Exposure, or Smart Contract finding against your latest code.

Cloud Security Audit findings still show **Re-test**. Retesting them is not supported. After you change the account or the infrastructure as code, run the audit again.

## Related

{% content-ref url="/pages/64GGpVtiCGGThzz1y8B7" %}
[AI Code Analysis](/ai-code-analysis/ai-code-analysis.md)
{% endcontent-ref %}
