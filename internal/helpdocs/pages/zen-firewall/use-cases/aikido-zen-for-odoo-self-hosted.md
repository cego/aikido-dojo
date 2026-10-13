# Aikido Zen for Odoo Self-Hosted

[Odoo](https://www.odoo.com/) is a full-featured, modular business management suite (ERP/CRM) that covers functions like CRM, e-commerce, accounting, inventory, manufacturing, project management, warehouse operations, accounting and more.

You can deploy Odoo on-premises or in the cloud. Aikido Zen does not support Odoo cloud.

Because it often handles critical business data, such as financials, inventory, purchase orders, customer data, invoices, production data, a breach or exploit of Odoo can have serious consequences.

Therefore Odoo is especially relevant for organizations that need to manage sensitive business data internally or want full control over their infrastructure.

## Why Aikido Zen matters for Odoo

Even if you secure your Odoo installation correctly, there remains a risk that a bug in Odoo itself, a third-party module, or a custom integration could lead to data leaks or unauthorized data access.

Aikido Zen acts as a runtime security layer on the server side. By embedding into the application environment, Zen intercepts suspicious actions. like malicious database queries, unsafe file access or external HTTP requests from the application before they execute. This way it helps mitigate exploitation of vulnerabilities, even in complex modular systems like Odoo.

This additional layer is particularly valuable for Odoo because:

* Odoo uses a modular architecture where core modules and many community or custom modules live side by side. That increases the attack surface.
* Business-critical operations, including inventory, accounting, orders, and manufacturing, often rely on database interaction, file storage, or external calls. Compromise of any of these could leak sensitive data or corrupt business logic.
* Many deployments choose on-premises hosting for full data ownership. In that context, adding runtime protection with Zen enhances security without moving data off-site.

## How Zen protects Odoo

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

Beta available soon, contact Aikido support for more information
