# Custom Context for AutoFixes

Aikido lets you add extra custom context or instructions to SAST rules, enhancing AutoFix to your preferences.

Common use cases include:

* Introducing specific error handling or logging logic
* Adding custom comment styles
* Including unit tests in AutoFix
* Providing additional context not in the code repository
* Narrowing the scope of an AutoFix

## Custom Context on SAST Rule Level

Detailed documentation for the Custom Code Context for SAST Rules can be found [here](https://help.aikido.dev/code-scanning/scanning-practices/custom-code-context-for-sast-rules). Aikido uses the Custom Code Context for both the scanning (AutoTriage) and fixing of issues (AutoFix).

Custom Code Context can be configured for all code repositories or a set of repositories.

## Global Custom Context

**Custom Context per Repo**

If you want to provide custom context for all AutoFix rules in a code repository, add the following to the `.aikido` file in your code repository (more info on `.aikido` files [can be found here](/code-scanning/scanning-practices/ignore-via-code-with-aikido-files.md))

<pre><code>autofix:
<strong>  fix_context:
</strong><strong>    - path: aikido-autofix-all.md
</strong></code></pre>

For this example to work, the `aikido-autofix-all.md` needs to exist in the root of the respository and needs to contain the context.

**Custom Context per Language**

To scope the configuration down to a specific language use:

```
autofix:
  fix_context:
    - path: aikido-autofix-javascript.md
      language: javascript
```

The allowed values for language are:

* .net
* android
* c
* clojure
* dart
* elixir
* go
* ios
* java
* javascript
* kotlin
* php
* python
* ruby
* rust
* scala
* swift
* terraform
* typescript
* vb
* yaml
