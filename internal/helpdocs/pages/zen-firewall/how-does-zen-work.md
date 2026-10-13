# How Does Zen Work?

Zen Firewall by Aikido is a powerful Application Firewall that embeds directly into your code. This position allows it to see and understand what is really happening in your code, rather than just inspecting incoming requests from the outside like a WAF.

### How Zen differs from a traditional firewall

Zen instruments key packages and built-in modules in your runtime (e.g. `fs.readFile`, `http.request`, `fetch`, and `mysql.query` in Node.js). This feature enables Zen to trace potentially unsafe user inputs back to their source when they appear in risky contexts, such as when inputs break out of SQL query quotes, manipulate file paths, alter outbound HTTP connection hostnames, or create similar security vulnerabilities. When it detects this, Zen stops the operation directly in code by throwing an exception. It also captures a stack trace so developers can quickly see where and how the issue occurred.

**By working with the full context of your application, it can make more accurate decisions and reduce the need for constant tuning or rule updates compared to a traditional firewall or WAF.**

A traditional firewall or WAF, by contrast, sits at the network edge and decides what to block based only on the incoming HTTP request. It typically relies on pattern matching, for example spotting an SQL injection payload like `' OR 1=1 --`. This has limitations: it cannot see how that request is processed inside your application, so it may block harmless requests or miss cleverly disguised attacks.

### Why you might need both Zen and **a traditional firewall or WAF**

Zen cannot defend against distributed denial of service attacks (DDoS), because it is too far inside the network stack to stop massive floods of traffic before they reach your application. A traditional firewall is still needed to handle those situations. The strongest protection comes from using both: the firewall for large-scale network threats, and Zen for precise, code-level protection.

### More than blocking

Zen also includes features beyond what many traditional firewalls offer:

* **Bot and IP blocking** – While firewalls can block bots, setting up advanced bot filtering can be complicated. Zen makes it easy to block requests from bots, known threat actors, and Tor, as well as restrict access by country or continent when required for legal or business reasons.
* **User-aware rate limiting** – Firewalls offer a way to rate limit by IP address. Zen can rate limit based on the actual user and also provides a method to set the rate limiting group ID. For example, a B2B SaaS can rate limit endpoints by customer account rather than individual user or IP.
* **Automatic API documentation** – Zen samples live traffic to generate an OpenAPI specification, which can then be used in the Aikido platform for scanning and testing your API.
* **AI model tracking** – If you make calls to large language models, Zen can track input and output tokens to estimate usage and cost across different providers and models.

### Simple to install and maintain

Zen can be installed in seconds, often with a single command. The setup experience is similar to adding a tracing library such as OpenTelemetry. Because Zen works by instrumenting your runtime instead of relying on large sets of static rules, it does not require frequent updates or complex configuration.
