# Live Secret Detection

Our Live Secret Detection feature checks if exposed secrets are still active and assesses their potential risks. Based on the outcome, the issue's severity will be increased or lowered.

## Use Cases <a href="#use-cases" id="use-cases"></a>

* Identifying and flagging active secrets in code repositories.
* Reducing noise for secrets that are not active anymore (e.g., already rotated)
* Checking the scope of permissions granted by exposed secrets.
* Enhancing security by marking dangerous secrets for immediate action.

## How Live Secret Detection Works <a href="#how-live-secret-detection-works" id="how-live-secret-detection-works"></a>

**Identify and Verify**

Aikido sends the exposed secret to a secure endpoint to check whether it is still active. As a result, you may notice IPs in your logs that are coming from Aikido. Below the list of Aikido IPs:

* **52.214.244.18**
* **18.202.209.180**
* **52.50.198.227**
* **52.51.98.186**

**Assess permissions**

Aikido goes a step further and checks permissions of the active secrets. Based on that, we make an extra distinction in our severity upgrades. We check the following.

* **Expired Secrets**
* **Read-Only Scopes**
* **Write/Delete Scopes**

> We do different checks per provider (e.g., GitHub Access Tokens, Sendgrid tokens, etc). Contact us if there are secrets you want to have checked!

## False positives in Secrets Detection

While we strive to minimize false positives by filtering out test files, mock data, and example configurations, some legitimate-looking secrets may still be flagged incorrectly. This happens because certain patterns are impossible to distinguish programmatically.

### Common False Positive Scenarios

**Example secrets in documentation or comments** Sometimes developers include example API keys in code comments or documentation that match the exact pattern of real secrets:

* Real secret: `SENDGRID_API_KEY = “SG.actualR3alK3y_1234567890abcdef”`
* Example secret: `SENDGRID_API_KEY = “SG.exampleK3y_1234567890abcdef”`

Both follow the same format and pattern recognition cannot determine which is real without verification.

**Deactivated or rotated secrets** Our Live Secret Detection helps reduce these by checking if secrets are still active, but there may be edge cases where verification isn't possible.

**Test environment secrets** Secrets used exclusively in test environments that match production patterns but pose no real security risk.

### What You Can Do

If you encounter a false positive:

1. Review the context - Is this in a test file, documentation, or example code?
2. Verify if the secret is active using the relevant API or other methods
3. If confirmed as a false positive, you can ignore the issue to prevent future alerts

You can [ignore issues in the UI](/getting-started/core-functionalities/ignore-issues-to-remove-issues-from-main-feed.md) or [ignore secrets by using code comments](/code-scanning/scanning-practices/ignoring-secrets-via-code-comments.md).
