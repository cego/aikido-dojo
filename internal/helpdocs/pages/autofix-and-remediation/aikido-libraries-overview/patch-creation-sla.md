# Patch Creation SLA

For every new CVE that is detected in an existing **Aikido Image**, or in an **Aikido Library**, the following SLA apply:

| Severity        | SLA                     |
| --------------- | ----------------------- |
| CISA KEV        | 48 hours                |
| Critical / High | 7 days                  |
| Medium / Low    | Commercially reasonable |

SLA timelines start once a CVE is published and a valid fix candidate is available in the upstream ecosystem.

### A note on pre-existing CVEs

**Aikido Libraries** - the SLA covers new CVEs that arise after you enroll your repository. Pre-existing CVEs will still be fixed, but they are not covered under the SLA commitment. This is because repositories can carry a large backlog of open CVEs at the time of enrollment, and working through that debt is handled separately.
