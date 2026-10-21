# Add Custom SAST & IaC Rules

## Introduction <a href="#introduction" id="introduction"></a>

With these custom rules you can make Aikido scan for specific risks in your codebase, especially those risks that are particularly relevant for your environment. This way you can detect vulnerabilities that broader SAST or IaC rules might overlook.

## Step-by-Step Guide <a href="#step-by-step-guide" id="step-by-step-guide"></a>

{% hint style="info" %}
This functionality is not by default enabled in your workspace. Contact Aikido to enable.
{% endhint %}

**Step 1:** Go to the [repositories checks](https://app.aikido.dev/repositories/checks) page.

**Step 2:** Click on "[Create Custom Rule](https://app.aikido.dev/repositories/sast/custom/add)" in the SAST section\
​

[image: Repository checks dashboard with options for creating custom security and compliance rules.]

​**Step 3:** Enter the following details for your rule:

* **Opengrep rule:** Define the rule Aikido will search for. **Tip:** Use the [Opengrep playground](https://github.com/opengrep/opengrep-playground?tab=readme-ov-file#installation)[ ](https://semgrep.dev/playground)to test your rule's effectiveness before saving.
* **Title:** Name your rule for easy identification.
* **TL;DR:** Provide a concise description of the issue. This will show up in the sidebar.
* **How to fix it:** Let your team know the best way to fix this issue.
* **Language:** Specify the programming language.
* **Aikido Score:** Set the priority level for issue reporting in the main Feed.

**Step 4:** Once you're satisfied with the rule's configuration, click "Save" to add it to your Aikido SAST checks. Your custom rule is now active and will be automatically applied in future scans.

### Extra Info <a href="#extra-info" id="extra-info"></a>

* Overall, the language attribute in the Opengrep rule will always prevail. This can be helpful when you are looking to implement a custom rule that needs to be applied to all languages and files at once.
* If you want to create IaC rules, you can do this by setting the language to yaml/terraform/...

## Examples <a href="#examples" id="examples"></a>

* **SAST Rule:** Looking for use of the weak MD5 hashing algorithm in javascript.

  ```
  rules:
    - id: CUSTOM-RULE-XXXX
      message: It looks like MD5 is used 
      languages:
        - javascript
      paths: 
        include
          - "*.js"
      severity: WARNING
      pattern-either:
        - pattern: $CRYPTO.createHash("md5")
        - pattern: CryptoJS.MD5(...)
  ```
* **IaC Rule:** A custom rule for detecting lambda functions that might be dangerous.

  ```
  rules:
     - id: CUSTOM-RULE-530
       languages:
         - hcl
       severity: WARNING
       message: >
         A Lambda function was found with the "type:monitored" tag, but without a "service" tag.
       patterns:
         - pattern: |-
             resource "aws_lambda_function" $ANYTHING {
               ...
               tags = {..., type = "monitored", ...}
             }
         - pattern-not: |-
             resource "aws_lambda_function" $ANYTHING {
               ...
               tags = {..., service= "...", ...}
             }
  ```
