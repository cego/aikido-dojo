# Aikido CI API

**Base URL** `https://app.aikido.dev`

The Aikido CI API allows you to start new scans for specific branches in a repo. You’ll be able to poll the state of the scans so you can build a gating mechanism inside your CI/CD platform.

Aikido also allows you to manage specific scans yourself. You can upload your own Checkov SARIF and custom SBOMs directly to Aikido to contribute to a feature branch scan. Custom CI systems such as Jenkins and CircleCI can also upload code coverage reports with the same API token. See [Upload code coverage](#upload-code-coverage).

## Authentication <a href="#authentication" id="authentication"></a>

1. Go to the [CI integration detail page](https://app.aikido.dev/settings/integrations/continuous-integration).
2. Generate an authentication token and copy. Note that you will only be able to view this token once.
3. Add the header `X-AIK-API-SECRET` to every request, with the authentication token as its value.

## Start a new feature branch scan <a href="#start-a-new-feature-branch-scan" id="start-a-new-feature-branch-scan"></a>

This async API kicks off a new feature branch scan on a specific commit.

`POST /api/integrations/continuous_integration/scan/repository`

**Body:**

```
{
    "version": "1.0.5",
    "branch_name": "branch-name",
    "repository_id": "R_kgDOJf6H8g", 
    "base_commit_id": "1c5a7dda074aa19ebc6aa9f25884b5e0b6bb3662",
    "head_commit_id": "b7deb06d8be6bc4c62277e0a28911b2284cab6da",
    "minimum_severity": "HIGH",
    "fail_on_sast_scan": true,
    "fail_on_dependency_scan": true,
    "fail_on_iac_scan": true,
    "fail_on_secrets_scan": true,
    "fail_on_malware_scan": false,
    "pull_request_metadata": {
        "url": "https://github.com/AikidoSec/some-repo/pull/12",
        "title": "checkov"
    },
    "self_managed_scanners": ["checkov", "json-sbom"],
    "expected_amount_json_sboms": 1,
}
```

**Required fields:**

* repository\_id: internal GitHub/Gitlab/Bitbucket/.. repository id
* base\_commit\_id: the base commit (the commit where you branched from for your PR - if unrelated to a PR, commit right before end\_commit\_id)
* head\_commit\_id: the latest commit we’re scanning (last commit on the branch)
* branch\_name: the branch name

**Optional fields:**

* pull\_request\_metadata.url: pull request URL
* pull\_request\_metadata.title: pull request title
* self\_managed\_scanners: array of custom scan types (supports `checkov` and `json-sbom`). Defaults to empty array \[].
* expected\_amount\_json\_sboms: integer of the amount of custom json-sbom uploads we’re waiting for.
* fail\_on\_dependency\_scan (default `true`): Determines whether Aikido should block on new dependency issues (CVEs)
* fail\_on\_sast\_scan (default `false`): Determines whether Aikido should block on new SAST issues
* fail\_on\_iac\_scan (default: `false`): Determines whether Aikido should block on new IaC scans
* fail\_on\_secrets\_scan (default `false`): Determines whether Aikido should block on new leaked secret issues
* fail\_on\_malware\_scan (default `false`): Determines whether Aikido should block on new malware issues.
  * **Note:** this feature is not available to all. Contact Aikido if you want to have it enabled.
* minimum\_severity (default `CRITICAL`): Determines on which (minimum) severity Aikido should respond with `FAILED`. This value can be one of `LOW`, `MEDIUM`, `HIGH`, `CRITICAL`.

**Response:**

```
{
    "scan_id": 22
}
```

## Poll for feature branch scan status <a href="#poll-for-feature-branch-scan-status" id="poll-for-feature-branch-scan-status"></a>

Fetches the current state of a scan. Use a polling mechanism (eg every 5 seconds) to implement gating for CI/CD.

`GET /api/integrations/continuous_integration/scan/repository?scan_id=1`

**Params:**

* scan\_id (int)

**Response:**

Example for feature branch *scan in progress*

```
{
    "all_scans_completed": false,
    "dependency_scan_completed": true,
    "sast_scan_completed": false,
    "iac_scan_completed": false,
    "secrets_scan_completed": false,
    "malware_scan_completed": false,
    "code_quality_scan_completed": false,
    "sbom_scan_completed": false // if 'json-sbom' in self_managed_scanners
}
```

Example for feature branch scan completed

```
{
    "all_scans_completed": true,
    "dependency_scan_completed": true,
    "sast_scan_completed": true,
    "iac_scan_completed": true,
    "secrets_scan_completed": true,
    "malware_scan_completed": true,
    "sbom_scan_completed": true, // if 'json-sbom' in self_managed_scanners
    "code_quality_scan_completed": false,
    "issues_outside_sla": 0,
    "new_issues_found": 1,
    "new_dependency_issues_found": 1,
    "new_sast_issues_found": 0,
    "new_iac_issues_found": 0,
    "new_leaked_secret_issues_found": 0,
    "new_malware_issues_found": 0,
    "new_license_issues_found": 0,
    "code_quality_issues_found": 0,
    "diff_url": "https://app.aikido.dev/featurebranch/scan/1?groupId=1",
    "gate_passed": false,
    "issue_breakdown": {
        "sast": { "critical": 0, "high": 0, "medium": 0, "low": 0 },
        "sca": { "critical": 1, "high": 0, "medium": 0, "low": 0 },
        "iac": { "critical": 0, "high": 0, "medium": 0, "low": 0 },
        "secrets": { "critical": 0, "high": 0, "medium": 0, "low": 0 },
        "malware": { "critical": 0, "high": 0, "medium": 0, "low": 0 },
        "license": { "critical": 0, "high": 0, "medium": 0, "low": 0 }
    }
}
```

Example for d*efault branch scan completed*

```
{
    "all_scans_completed": false,
    "dependency_scan_completed": true,
    "sast_scan_completed": false,
    "iac_scan_completed": false,
    "secrets_scan_completed": false,
    "malware_scan_completed": false,
    "open_issues_found": 2,
    "issue_links": [
        "https://app.aikido.dev/queue?sidebarIssue=1",
        "https://app.aikido.dev/queue?sidebarIssue=2",
    ],
    "gate_passed": false,
    "issue_breakdown": {
        "sast": { "critical": 0, "high": 0, "medium": 0, "low": 0 },
        "sca": { "critical": 1, "high": 0, "medium": 0, "low": 0 },
        "iac": { "critical": 0, "high": 0, "medium": 0, "low": 0 },
        "secrets": { "critical": 1, "high": 0, "medium": 0, "low": 0 },
        "malware": { "critical": 0, "high": 0, "medium": 0, "low": 0 },
        "license": { "critical": 0, "high": 0, "medium": 0, "low": 0 }
    }
}
```

```
{
    "all_scans_completed": false,
    "dependency_scan_completed": true,
    "sast_scan_completed": false,
    "iac_scan_completed": false,
    "secrets_scan_completed": false,
    "malware_scan_completed": false,
    "open_issues_found": 2,
    "issue_links": [
        "https://app.aikido.dev/queue?sidebarIssue=1",
        "https://app.aikido.dev/queue?sidebarIssue=2",
    ],
    "gate_passed": false,
    "issue_breakdown": {
        "sast": { "critical": 0, "high": 0, "medium": 0, "low": 0 },
        "sca": { "critical": 1, "high": 0, "medium": 0, "low": 0 },
        "iac": { "critical": 0, "high": 0, "medium": 0, "low": 0 },
        "secrets": { "critical": 1, "high": 0, "medium": 0, "low": 0 },
        "malware": { "critical": 0, "high": 0, "medium": 0, "low": 0 },
        "license": { "critical": 0, "high": 0, "medium": 0, "low": 0 }
    }
}
```

```
{
    "all_scans_completed": false,
    "dependency_scan_completed": true,
    "sast_scan_completed": false,
    "iac_scan_completed": false,
    "secrets_scan_completed": false,
    "malware_scan_completed": false,
    "open_issues_found": 2,
    "issue_links": [
        "https://app.aikido.dev/queue?sidebarIssue=1",
        "https://app.aikido.dev/queue?sidebarIssue=2",
    ],
    "gate_passed": false,
    "issue_breakdown": {
        "sast": { "critical": 0, "high": 0, "medium": 0, "low": 0 },
        "sca": { "critical": 1, "high": 0, "medium": 0, "low": 0 },
        "iac": { "critical": 0, "high": 0, "medium": 0, "low": 0 },
        "secrets": { "critical": 1, "high": 0, "medium": 0, "low": 0 },
        "malware": { "critical": 0, "high": 0, "medium": 0, "low": 0 },
        "license": { "critical": 0, "high": 0, "medium": 0, "low": 0 }
    }
}
```

## Upload custom scan results <a href="#upload-custom-scan-results" id="upload-custom-scan-results"></a>

This API can be used to upload your own Checkov SARIF or custom SBOMs generated with Syft. Use the JSON format when generating a Syft SBOM.

`POST /api/integrations/continuous_integration/scan/custom`

**Body:**

Example for *Checkov*

```
{
    "scan_id": 22, // optional
    "repository_id": 12,
    "payload_type": "checkov",
    "payload": <checkov_sarif_response>
}
```

Example for *Syft SBOM*

```
{
    "scan_id": 22, // optional
    "repository_id": 12,
    "container_image_name": "image-name",
    "payload_type": "json-sbom",
    "payload": <syft_json_sbom_response>
}
```

**Required fields:**

* container\_image\_name: required for json-sbom. This should not include the image tag, only a unique name for the image.
* payload\_type
* payload

**Optional fields:**

* scan\_id: (retrieved from *Start Scan*). Do not set this if you want the payload to go to the Aikido live feed.
* repository\_id: internal GitHub/Gitlab/Bitbucket/.. repository id (only required if scan\_id is not set)

**Response:**

```
{
    "success": 1
}
```

## Upload code coverage <a href="#upload-code-coverage" id="upload-code-coverage"></a>

Custom CI such as Jenkins and CircleCI can upload [LCOV](https://github.com/linux-test-project/lcov) or [Cobertura](https://cobertura.github.io/cobertura/) coverage reports with the same `X-AIK-API-SECRET` token.

Use the [Aikido CI API client](https://www.npmjs.com/package/@aikidosec/ci-api-client) `upload-coverage` command. Run it from a checked-out repository after tests write the report. The client reads the coverage files and source metadata, then sends them to Aikido. See [Code Coverage](/code-quality/code-coverage.md) for how uploads appear in Aikido.

```sh
$ aikido-api-client upload-coverage \
    --apikey <your-api-key> \
    --repo-name <owner/repo> \
    --commit-sha <commit_sha> \
    --branch-name <branch_name> \
    --repository-id <repository_id> \
    --file-paths coverage/lcov.info
```

`repository_id` is the Git provider repository id, the same id used by the other endpoints on this page. Pass more than one path to `--file-paths` to merge LCOV and Cobertura reports from separate packages or jobs.

Each upload replaces coverage for that repository and branch. In a monorepo, pass every app's report on each run. If you skip unchanged apps, include their latest reports in that same command. See [Monorepos That Test Only Changed Paths](/code-quality/code-coverage/monorepos-that-test-only-changed-paths.md).
