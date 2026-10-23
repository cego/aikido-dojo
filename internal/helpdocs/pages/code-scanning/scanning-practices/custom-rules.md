# Custom Rules Overview

You can define your own rules in Aikido across code security, cloud security, and code quality, so checks align with your your environment and standards.

### **Where you can add custom rules**

* [**Custom SAST & IaC Rules**](https://help.aikido.dev/general-information/add-custom-sast-iac-rules)

  Define rules to detect application or infrastructure patterns that are specific to your stack, such as enforcing internal security APIs or common misconfigurations specific to your organization.
* [**Custom CSPM Rules**](https://help.aikido.dev/cloud-scanning/custom-cspm-rules)

  Add your own cloud misconfiguration checks on top of Aikido’s managed rules. Custom CSPM rules can be mapped to compliance benchmarks via tags, and they show up in compliance reports like any other cloud rule (for example, mapping a backup rule to ISO 27001 sections).
* [**Custom Code Quality Rules**](https://help.aikido.dev/code-quality/add-custom-code-rules)

  Describe code patterns or conventions you want to enforce in natural language. Aikido turns these into AI-powered checks that run on pull requests and help you enforce team-specific coding guidelines. Add an [AutoFix](/autofix-and-remediation/overview-aikido-autofix.md) prompt on the rule so Aikido can fix violations.

### Typical use cases

* Encode internal security guidelines, such required wrappers or safe helpers
* Flag cloud configurations that are risky in your particular architecture and tying them to your compliance framework
* Enforce local coding standards and architectural rules
