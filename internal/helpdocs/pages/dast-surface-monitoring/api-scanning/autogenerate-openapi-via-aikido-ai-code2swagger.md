# Autogenerate OpenAPI via Aikido AI (Code2Swagger)

Aikido can automatically generate an OpenAPI specification directly from your codebase, eliminating the need for manual documentation while ensuring accurate API scanning.

It uses LLM to parse your source code to find routes, parameters and body information without the need to build your own API specification.

Aikido will regularly rescan your code to keep the OpenAPI specification up to date.

{% hint style="success" %}
Aikido recommends using [Zen Firewall](/zen-firewall/getting-started-with-zen-firewall.md) whenever possible as it will give more accurate results because it's based on the actual traffic to your API.
{% endhint %}

### How to use

Follow the [Setting up REST API Scanning](/dast-surface-monitoring/api-scanning/rest-api-scanning.md#setting-up-rest-api-scanning) guide and choose the **Aikido AI** option.

### Supported Programming Languages

{% hint style="warning" %}
This feature supports back-end code analysis only; front-end code is not supported.
{% endhint %}

* C
* C++
* C# / Dotnet
* Dart
* Elixir
* Go
* Java
* Javascript
* Kotlin
* PHP
* Python
* Ruby
* Rust
* Swift
* Scala
* Typescript / TSX
* VB.NET (only route discovery)

<br>
