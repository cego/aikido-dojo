# OSCAL DAST Report

The OSCAL FedRAMP DAST Report gives you your dynamic scan results in the format FedRAMP reviewers expect. OSCAL (Open Security Controls Assessment Language) is a machine-readable format defined by NIST, so the reviewer's tooling can ingest the file directly instead of someone reading a PDF.

Aikido builds the file from the findings of your domain and API scans, covering the DAST (Dynamic Application Security Testing) part of your evidence. No configuration is needed.

Use this report for FedRAMP authorization packages, monthly continuous monitoring evidence, or handing scan evidence to a 3PAO (third-party assessment organization). For human readers, use the [Security Audit Report](/compliance-and-reporting/reports/security-audit-report.md) or a [pentest report](/pentests/coverage-and-findings/pentest-reports.md) instead.

## Download the report

1. Go to **Reports** in Aikido.
2. Open **OSCAL DAST Report**.
3. Click **Download** to get the OSCAL JSON file.

The export reflects the current state of your scans, so generate it after your latest scan has finished.

{% hint style="info" %}
Validate the file against the FedRAMP OSCAL tooling your reviewer uses, as requirements change between baseline revisions. Your package may need other artifacts that Aikido does not generate, such as a POA\&M (Plan of Action and Milestones).
{% endhint %}
