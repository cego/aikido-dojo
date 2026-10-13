# Data Sent by Aikido Zen

Aikido Zen periodically reports metadata to help analyze security activity and detect attacks in your environment. This data is aggregated, not sent in real-time, and does not include sensitive user data or payloads unless an attack is detected.

## Reporting frequency

* On startup – Zen sends a one-time snapshot of your environment.
* Every 10 minutes – Zen sends aggregated usage statistics.
* When an attack is detected – Zen immediately reports details about the attack and its metadata.

## Data collected

### On startup

* Packages – the list of installed dependencies.
* System information – OS, architecture, hostname, and machine IP.

### Every 10 minutes

* Routes – the list of API routes with their specifications and the total number of hits per route (no individual requests).
* Outbound hostnames – aggregated list of hostnames your app connects to and how often.
* AI usage – statistics of AI-related operations, this does not contain any prompts
* App statistics – counts of requests, operations, and general performance metrics.
* Users – last known IP address for each user, but only if user tracking is explicitly enabled via setUser.

### When an attack is detected

* Attack metadata – includes request method, endpoint, timestamp, attack vector type, and other diagnostic details required for analysis.

## Privacy and PII

Aikido Zen is designed to avoid collecting personally identifiable information (PII).

* All regular reports are aggregated and anonymized where possible.
* User IPs are only included if explicitly enabled.
* No raw data, headers, or payloads are sent during normal operation.
