# Real-time SAST, Secrets and IaC scanning in IDE

The Aikido IDE plugin helps you catch security issues the moment they’re introduced. Every time you open or save a file, Aikido runs a quick scan in the background and highlights problems directly in your editor.

#### What It Scans

Aikido performs two types of scans in real time:

* SAST (Static Application Security Testing): Detects insecure coding patterns, potential injections, unsafe deserialization, and other code-level vulnerabilities.
* Secrets: Finds exposed credentials such as API keys, passwords, or tokens.
* IAC: (Infrastructure as Code): Detects cloud and infrastructure misconfigurations in Terraform, CloudFormation, Dockerfiles, and similar files.

#### How It Works

When you open or save a file, Aikido scans the code using the same analysis engine as the Aikido platform.

Detected issues appear:

* Inline, underlined or highlighted in the editor.
* In the Aikido sidebar, grouped by severity and category.
* In the Problems panel, for quick navigation.

Hover over any finding to see context and remediation details. For supported findings, you can [analyse using AI Autotriage or apply an AI AutoFix](/ai-and-dev-tools/ide-plugins-overview/features/aikido-ai-in-ide.md) to safely patch the issue without leaving your IDE.
