# Scanning Best Practices

This section covers how Aikido's static application security testing (SAST), software composition analysis (SCA), and secret detection work under the hood, and how to configure them for your setup.

For automatic SAST noise reduction powered by [Aikido Agent](https://help.aikido.dev/aikido-agent/aikido-agent), see [SAST AutoTriage](https://help.aikido.dev/aikido-agent/sast-autotriage).

## SAST and custom rules

{% content-ref url="sast-by-aikido-supported-languages-and-security-focus.md" %}
[sast-by-aikido-supported-languages-and-security-focus.md](https://help.aikido.dev/code-scanning/scanning-practices/sast-by-aikido-supported-languages-and-security-focus)
{% endcontent-ref %}

{% content-ref url="multifile-vulnerability-tracing.md" %}
[multifile-vulnerability-tracing.md](https://help.aikido.dev/code-scanning/scanning-practices/multifile-vulnerability-tracing)
{% endcontent-ref %}

{% content-ref url="enhancing-security-with-access-control-checks.md" %}
[enhancing-security-with-access-control-checks.md](https://help.aikido.dev/code-scanning/scanning-practices/enhancing-security-with-access-control-checks)
{% endcontent-ref %}

{% content-ref url="custom-rules.md" %}
[custom-rules.md](https://help.aikido.dev/code-scanning/scanning-practices/custom-rules)
{% endcontent-ref %}

{% content-ref url="add-custom-sast-iac-rules.md" %}
[add-custom-sast-iac-rules.md](https://help.aikido.dev/code-scanning/scanning-practices/add-custom-sast-iac-rules)
{% endcontent-ref %}

{% content-ref url="custom-code-context-for-sast-rules.md" %}
[custom-code-context-for-sast-rules.md](https://help.aikido.dev/code-scanning/scanning-practices/custom-code-context-for-sast-rules)
{% endcontent-ref %}

{% content-ref url="excluding-sast-findings-using-comments.md" %}
[excluding-sast-findings-using-comments.md](https://help.aikido.dev/code-scanning/scanning-practices/excluding-sast-findings-using-comments)
{% endcontent-ref %}

## Secret detection

{% content-ref url="live-secret-detection.md" %}
[live-secret-detection.md](https://help.aikido.dev/code-scanning/scanning-practices/live-secret-detection)
{% endcontent-ref %}

{% content-ref url="denoise-via-secret-keyword-filter.md" %}
[denoise-via-secret-keyword-filter.md](https://help.aikido.dev/code-scanning/scanning-practices/denoise-via-secret-keyword-filter)
{% endcontent-ref %}

{% content-ref url="ignoring-secrets-via-code-comments.md" %}
[ignoring-secrets-via-code-comments.md](https://help.aikido.dev/code-scanning/scanning-practices/ignoring-secrets-via-code-comments)
{% endcontent-ref %}

## Dependency scanning (SCA)

{% content-ref url="support-for-dependency-scanning-by-language.md" %}
[support-for-dependency-scanning-by-language.md](https://help.aikido.dev/code-scanning/scanning-practices/support-for-dependency-scanning-by-language)
{% endcontent-ref %}

{% content-ref url="external-vulnerability-databases-used-in-our-sca-engine.md" %}
[external-vulnerability-databases-used-in-our-sca-engine.md](https://help.aikido.dev/code-scanning/scanning-practices/external-vulnerability-databases-used-in-our-sca-engine)
{% endcontent-ref %}

{% content-ref url="scanning-dev-dependencies-for-cves.md" %}
[scanning-dev-dependencies-for-cves.md](https://help.aikido.dev/code-scanning/scanning-practices/scanning-dev-dependencies-for-cves)
{% endcontent-ref %}

{% content-ref url="custom-vulnerabilities.md" %}
[custom-vulnerabilities.md](https://help.aikido.dev/code-scanning/scanning-practices/custom-vulnerabilities)
{% endcontent-ref %}

{% content-ref url="../miscellaneous/use-epss-values-to-further-reduce-noise.md" %}
[use-epss-values-to-further-reduce-noise.md](https://help.aikido.dev/code-scanning/miscellaneous/use-epss-values-to-further-reduce-noise)
{% endcontent-ref %}

{% content-ref url="malware-detection-in-open-source-dependencies.md" %}
[malware-detection-in-open-source-dependencies.md](https://help.aikido.dev/code-scanning/scanning-practices/malware-detection-in-open-source-dependencies)
{% endcontent-ref %}

{% content-ref url="vulnerability-scanning-on-private-packages-maven.md" %}
[vulnerability-scanning-on-private-packages-maven.md](https://help.aikido.dev/code-scanning/scanning-practices/vulnerability-scanning-on-private-packages-maven)
{% endcontent-ref %}

{% content-ref url="c-c-lockfile-less-scanning.md" %}
[c-c-lockfile-less-scanning.md](https://help.aikido.dev/code-scanning/scanning-practices/c-c-lockfile-less-scanning)
{% endcontent-ref %}

{% content-ref url="dependency-scanning-for-bazel.md" %}
[dependency-scanning-for-bazel.md](https://help.aikido.dev/code-scanning/scanning-practices/dependency-scanning-for-bazel)
{% endcontent-ref %}

{% content-ref url="scala-dependency-management-and-scanning-for-buildsbt.md" %}
[scala-dependency-management-and-scanning-for-buildsbt.md](https://help.aikido.dev/code-scanning/scanning-practices/scala-dependency-management-and-scanning-for-buildsbt)
{% endcontent-ref %}

{% content-ref url="java-scala-kotlin-projects-using-gradle-security-scanning-best-practices.md" %}
[java-scala-kotlin-projects-using-gradle-security-scanning-best-practices.md](https://help.aikido.dev/code-scanning/scanning-practices/java-scala-kotlin-projects-using-gradle-security-scanning-best-practices)
{% endcontent-ref %}

{% content-ref url="net-projects-security-scanning-best-practices.md" %}
[net-projects-security-scanning-best-practices.md](https://help.aikido.dev/code-scanning/scanning-practices/net-projects-security-scanning-best-practices)
{% endcontent-ref %}

{% content-ref url="python-projects-security-scanning-best-practices.md" %}
[python-projects-security-scanning-best-practices.md](https://help.aikido.dev/code-scanning/scanning-practices/python-projects-security-scanning-best-practices)
{% endcontent-ref %}

## Branches and repositories

{% content-ref url="../miscellaneous/scanned-branches-in-aikido.md" %}
[scanned-branches-in-aikido.md](https://help.aikido.dev/code-scanning/miscellaneous/scanned-branches-in-aikido)
{% endcontent-ref %}

{% content-ref url="../miscellaneous/support-for-multi-branch-scanning.md" %}
[support-for-multi-branch-scanning.md](https://help.aikido.dev/code-scanning/miscellaneous/support-for-multi-branch-scanning)
{% endcontent-ref %}

{% content-ref url="../miscellaneous/support-for-specific-branch-scanning.md" %}
[support-for-specific-branch-scanning.md](https://help.aikido.dev/code-scanning/miscellaneous/support-for-specific-branch-scanning)
{% endcontent-ref %}

{% content-ref url="../miscellaneous/ensuring-aikido-scans-new-repositories.md" %}
[ensuring-aikido-scans-new-repositories.md](https://help.aikido.dev/code-scanning/miscellaneous/ensuring-aikido-scans-new-repositories)
{% endcontent-ref %}

{% content-ref url="../miscellaneous/split-your-monorepo-per-path.md" %}
[split-your-monorepo-per-path.md](https://help.aikido.dev/code-scanning/miscellaneous/split-your-monorepo-per-path)
{% endcontent-ref %}

## Other

{% content-ref url="ignore-via-code-with-aikido-files.md" %}
[ignore-via-code-with-aikido-files.md](https://help.aikido.dev/code-scanning/scanning-practices/ignore-via-code-with-aikido-files)
{% endcontent-ref %}
