# Difference Between Aikido Cloud Scanning and IDE Scanning

Aikido offers security scanning in two places: inside the developer’s editor and in the cloud platform. Both scan code and dependencies, but with different depth and context.

### Cloud Scans

* Run SAST, SCA and IaC checks across full repositories and connected cloud environments.
* Analyse complete codebases, dependency graphs and configuration files rather than only what is open in the editor.
* Provide richer context, for example reachability, deployment paths, and whether an issue affects production.
* Use advanced auto triage to cut down false positives and highlight real risk.
* Suitable for full-scope audits, CI pipelines and organisation-wide posture.

### IDE Scans

* Run SAST and SCA directly in the editor for fast feedback while coding.
* Limited to the current project and the files or manifests present in the workspace.
* Shallower analysis because it needs to be fast and local.
* Less context about how the code is built or deployed, so more findings require manual validation.
* Helpful for preventing issues before commit, but not a replacement for full cloud-based scanning.
