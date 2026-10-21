# Custom CSPM Rules

Aikido’s Custom CSPM Rules let you define cloud misconfiguration rules in natural language, generating new issues in the feed based on your cloud asset searches. All these issues follow the same functionality as other issues in the feed, e.g, **ignore, snooze, alert notifications**, **task creation**, **etc.**

### **Use cases**

* Extending the out-of-the-box Aikido cloud checks.
* Organization-specific requirements, such as data residency or tagging policies.
* Responding to zero-day cloud misconfigurations.

### **Create a Custom CSPM Rule**

**Step 1.** Go to the [**Cloud Assets**](https://app.aikido.dev/clouds/assets) page in Aikido.

**Step 2.** Click on **Custom CSPM Rules**

**Step 3.** Click **Create custom rule** on the Custom Rules Page.

**Step 4. Fill in the necessary details to create the rule**

* Add a **related search query** to define the condition that triggers the issue (e.g. `buckets outside eu`).
* Issue title
* TL;DR
* How to fix
* Score (this impacts severity)

[image: Custom cloud rule example]

**Step 5**. Your rule will be added to the page. You can always edit or delete by clicking the action dropdown menu.

**Step 6.** Trigger a new cloud scan manually to have matching results appear in the feed.

### **Important Notes**

* Custom CSPM rules are evaluated after they are created and with each cloud scan.
* These rules generate **cloud misconfiguration issues**, which support all standard issue features (severity adjustment, Slack integration, task creation, etc.).
* Custom CSPM rules are applied **across all connected cloud accounts**.

### Mapping to Compliance Reports

You can map custom CSPM rules to compliance reports by adding compliance tags when creating or editing a rule. Aikido uses these tags to automatically include the rule in all relevant compliance sections of the supported benchmarks. The custom CSPM rules will appear in the compliance reports right after you add the tags, and will behave like any other cloud rule/check.

For example, if you wish to map a rule to "ISO 27001:2022, A.8.13 - Backups", you add the `Backup` tag. Please reach out to us if you need support or aren't sure which tag to use.
