# SIEM Field Reference

This page describes the fields included in each Aikido issue when it's sent to your connected SIEM.

For how to connect a specific SIEM, see [SIEM Connectors](/miscellaneous-integrations/siem-connectors.md).

### Payload format

Each sync sends a JSON array, with one object per issue that matches your [issue input settings](/miscellaneous-integrations/siem-connectors.md). The fields are the same across every SIEM connector.

### Core fields

| Field                  | Type           | Description                                                                                                                                                                                                         |
| ---------------------- | -------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `id`                   | integer        | The individual issue's ID.                                                                                                                                                                                          |
| `group_id`             | integer        | The ID of the issue group this issue belongs to. Aikido groups related occurrences of the same underlying issue together.                                                                                           |
| `attack_surface`       | string         | One of `frontend`, `backend`, `docker_container`, `cloud`, `cloud_instance`.                                                                                                                                        |
| `status`               | string         | One of `open`, `ignored`, `snoozed`, `closed`.                                                                                                                                                                      |
| `severity`             | string         | One of `critical`, `high`, `medium`, `low`.                                                                                                                                                                         |
| `severity_score`       | integer        | Aikido's priority score for the issue, from 1 to 100.                                                                                                                                                               |
| `type`                 | string         | One of `open_source`, `leaked_secret`, `cloud`, `cloud_instance`, `docker_container`, `iac`, `sast`, `surface_monitoring`, `malware`, `eol`, `scm_security`, `ai_pentest`, `ai_code_analysis`, `mobile`, `license`. |
| `rule`                 | string or null | The rule title, for SAST, IaC, cloud, and surface monitoring issues. Null for other types.                                                                                                                          |
| `programming_language` | string or null | Detected from the affected file path, for example `JS`, `PHP`, `PY`, `Java`.                                                                                                                                        |
| `how_to_fix`           | string or null | Aikido's remediation guidance for the issue, if available.                                                                                                                                                          |

### Vulnerability fields

| Field              | Type           | Description                                                                                         |
| ------------------ | -------------- | --------------------------------------------------------------------------------------------------- |
| `cve_id`           | string or null | The CVE ID, or an internal Aikido CVE ID (for example `AIKIDO-2020-123`) when no public CVE exists. |
| `cve_id_alias`     | string or null | The related public CVE, when `cve_id` is an internal Aikido CVE ID.                                 |
| `affected_package` | string or null | The affected package name, if any.                                                                  |
| `affected_file`    | string or null | The affected file or target path, if any.                                                           |

### Location fields

Only the field group matching where the issue was found is populated; the rest are null.

| Field                                        | Type                            | Description                                     |
| -------------------------------------------- | ------------------------------- | ----------------------------------------------- |
| `code_repo_id`, `code_repo_name`             | integer or null, string or null | The code repository, for source-code issues.    |
| `container_repo_id`, `container_repo_name`   | integer or null, string or null | The container repository, for container issues. |
| `cloud_id`, `cloud_name`                     | integer or null, string or null | The cloud account, for cloud issues.            |
| `domain_id`, `domain_name`                   | integer or null, string or null | The domain, for surface monitoring issues.      |
| `virtual_machine_id`, `virtual_machine_name` | integer or null, string or null | The VM, for VM issues.                          |

### Status and timing fields

| Field               | Type                   | Description                                                                                                                                                                |
| ------------------- | ---------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `first_detected_at` | unix timestamp         | When the issue was first created.                                                                                                                                          |
| `ignored_at`        | unix timestamp or null | When the issue was ignored, if it has been.                                                                                                                                |
| `closed_at`         | unix timestamp or null | When the issue was closed, if it has been.                                                                                                                                 |
| `sla_remediate_by`  | unix timestamp or null | The SLA deadline for remediation, if one applies.                                                                                                                          |
| `snooze_until`      | unix timestamp or null | When a snooze expires, if the issue is snoozed.                                                                                                                            |
| `snooze_reason`     | string                 | Why the issue is snoozed. Empty string if it isn't.                                                                                                                        |
| `ignore_reasons`    | array                  | Why the issue is ignored. Empty array if it isn't. Each entry has a `kind` (`manual_ignore`, `auto_ignore`, `rule_ignore`, or `api_ignore`) and a human-readable `reason`. |

### Issue type metadata

`issue_type_metadata` is an object that's only included for some issue types, with at most one of the following keys set, matching the issue's `type`:

| Key                  | Fields                                                                                               |
| -------------------- | ---------------------------------------------------------------------------------------------------- |
| `sast`               | `start_line`, `end_line`                                                                             |
| `cloud`              | `resource_id`, `resource_name`, `resource_kind`                                                      |
| `surface_monitoring` | `cve_id`, `package_name`, `package_version`, `evidence`, `evidence_description`, `domain_target_url` |
| `open_source`        | `installed_version`, `patched_versions` (array)                                                      |
| `eol`                | `eol_since` (unix timestamp or null), `affected_version`                                             |
| `leaked_secret`      | `commit`, `line`                                                                                     |

### Example

```json
[
  {
    "id": 1,
    "group_id": 1,
    "attack_surface": "backend",
    "status": "open",
    "severity": "critical",
    "severity_score": 90,
    "type": "open_source",
    "rule": null,
    "cve_id": "AIKIDO-2020-123",
    "cve_id_alias": "CVE-2020-7598",
    "affected_package": "minimist",
    "affected_file": "package-lock.json",
    "first_detected_at": 1700489005,
    "code_repo_id": 1,
    "code_repo_name": "test-service",
    "container_repo_id": null,
    "container_repo_name": null,
    "cloud_id": null,
    "cloud_name": null,
    "domain_id": null,
    "domain_name": null,
    "virtual_machine_id": null,
    "virtual_machine_name": null,
    "ignored_at": null,
    "closed_at": null,
    "sla_remediate_by": 1700924603,
    "snooze_until": null,
    "snooze_reason": "",
    "ignore_reasons": [],
    "issue_type_metadata": {
      "open_source": {
        "installed_version": "4.2.0",
        "patched_versions": ["4.2.6"]
      }
    },
    "how_to_fix": "Upgrade package minimist to 4.2.6 or later.",
    "programming_language": "JS"
  }
]
```
