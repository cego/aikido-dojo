# Aikido Zen for Legacy Apps

Aikido Zen is a runtime security layer that protects legacy applications while they are running. It is designed for systems that are hard to change but still need to be exposed to real users and real traffic.

Legacy apps often run on older frameworks, have large monolithic codebases, or are maintained with limited engineering time. Rewriting or deeply refactoring these systems is usually not realistic. Zen works with the application as it exists today.

Zen runs inside the application process and observes how requests interact with sensitive operations like database queries, file access, command execution, and outbound network calls.

## Legacy is hard to protect

Legacy applications were rarely built with modern attack techniques in mind. Input validation is often inconsistent, security libraries may be outdated or missing, and patching vulnerabilities can require risky changes to fragile code.

Traditional defenses do not solve this well. Scanners report issues but do not stop attacks. Network firewalls and classic web application firewalls operate outside the app and lack context about what the code is actually doing.

This leaves older systems exposed to common but high impact attacks such as injection, path traversal, and unsafe command execution. These vulnerabilities are especially dangerous because legacy apps often handle sensitive data and are still internet facing.

Aikido Zen reduces this risk immediately, without forcing code changes or architectural work.

## How Zen protects

**Mitigating abusive traffic**

Zen includes traffic control features that help slow down or block brute force style attempts and cut noise from scanning bots before they reach application logic.

**Blocking unsafe database operations**

Zen inspects database queries at runtime and can stop patterns that look like SQL injection or harmful query manipulation. This helps contain issues coming from custom code or third-party modules that interact with your data models.

**Preventing unauthorized file system access**

Many applications store configuration, attachments, and temporary files on disk. Zen monitors sensitive file operations and blocks path traversal or unexpected file reads and writes outside allowed directories.

**Protecting against unsafe outbound requests and SSRF**

Many applications integrate with external services via webhooks and APIs. Zen checks outbound HTTP calls and can block unsafe requests created from untrusted input. This prevents SSRF attempts and limits access to internal network targets.

**Monitoring and controlling outbound domains**

Zen shows every domain your server connects to. You can block or allow domains to control where data is sent. This is useful for detecting unexpected calls from custom modules, preventing data exfiltration and keeping a tight boundary around integrations your team approves.

**Local only inspection**

Zen makes security decisions locally on the server without sending request data to Aikido. This fits well with deployments where critical business data and privacy is a priority.

## Next step

To start protecting your legacy applications with Aikido Zen, create an Aikido account and follow the onboarding guide for your runtime in the Zen section of dashboard.

Depending on your setup, enabling Zen may require small, localized code changes or configuration updates. These changes are typically limited to application startup and do not affect business logic.

Once connected, Zen begins monitoring and blocking unsafe behavior at runtime.

Create your account at [aikido.dev](https://aikido.dev) to get started.
