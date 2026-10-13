# Dependency Scanning for Bazel

With Bazel, you can import dependencies from a variety of sources, such as Maven repositories, Git repositories, or internal artifact registries. This is typically achieved using dependency rules like `http_jar`, `maven_install` .

Aikido supports scanning the targets you define within Bazel for security vulnerabilities. A build is not required, which speeds up the analysis process.

There are a limited amount of Bazel rules we currently support, but we're extending support along the way. Contact our team in case you want to suggest new Bazel rules to support.

#### Supported Bazel rules

Currently, we scan for dependencies imported via the following Bazel rules:

<table><thead><tr><th width="184.8897705078125">Bazel Rule</th><th>Note</th><th>Version</th></tr></thead><tbody><tr><td><code>java_library</code></td><td><ul><li><p>Analyzes the dependency-related attributes to identify dependencies brought in via the following dependency rules:</p><ul><li><code>maven_install</code></li></ul></li></ul></td><td><ul><li>Bazel >= 5</li></ul></td></tr><tr><td><code>rules_jvm_external</code></td><td><ul><li>Parses the lock file referenced in the `lock_file` attribute of `maven.install(...)` definitions. View details <a href="https://github.com/bazel-contrib/rules_jvm_external?tab=readme-ov-file#pinning-artifacts-and-integration-with-bazels-downloader"><strong>here</strong></a>.</li></ul></td><td><ul><li>Bazel >= 5</li></ul></td></tr></tbody></table>
