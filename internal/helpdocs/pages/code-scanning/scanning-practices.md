# Scanning Best Practices

How Aikido scans code and dependencies, and how to tune it for your repositories, languages, and build tools.

This section covers how Aikido's static application security testing (SAST), software composition analysis (SCA), and secret detection work under the hood, and how to configure them for your setup.

For automatic SAST noise reduction powered by [Aikido Agent](/aikido-agent/aikido-agent.md), see [SAST AutoTriage](/aikido-agent/sast-autotriage.md).

## SAST and custom rules

{% content-ref url="/pages/XrQNZ4vFYIDC1fgURgT8" %}
[SAST/IaC: Supported Languages and Security Focus](/code-scanning/scanning-practices/sast-by-aikido-supported-languages-and-security-focus.md)
{% endcontent-ref %}

{% content-ref url="/pages/77FNIm4RnhMsLqyf2rdi" %}
[Multi-File / Cross-File Vulnerability Tracing](/code-scanning/scanning-practices/multifile-vulnerability-tracing.md)
{% endcontent-ref %}

{% content-ref url="/pages/dAJk8tVbR5hRs1m9IGYK" %}
[Access Control Checks](/code-scanning/scanning-practices/enhancing-security-with-access-control-checks.md)
{% endcontent-ref %}

{% content-ref url="/pages/9WgfYSZ2aD4jhHWmStL6" %}
[Custom Rules Overview](/code-scanning/scanning-practices/custom-rules.md)
{% endcontent-ref %}

{% content-ref url="/pages/02CWVAixFxzD9llXWPbe" %}
[Add Custom SAST & IaC Rules](/code-scanning/scanning-practices/add-custom-sast-iac-rules.md)
{% endcontent-ref %}

{% content-ref url="/pages/ZS0remVii1WN8YqvvEeT" %}
[Custom Code Context for SAST Rules](/code-scanning/scanning-practices/custom-code-context-for-sast-rules.md)
{% endcontent-ref %}

{% content-ref url="/pages/1kOdieG9Z6lDtKIoLzWO" %}
[Excluding SAST Findings Using Comments](/code-scanning/scanning-practices/excluding-sast-findings-using-comments.md)
{% endcontent-ref %}

## Secret detection

{% content-ref url="/pages/6mEApLmmau4c7pjTdV4g" %}
[Live Secret Detection](/code-scanning/scanning-practices/live-secret-detection.md)
{% endcontent-ref %}

{% content-ref url="/pages/YQXdwysoAjwnXP4R9ddZ" %}
[Denoise via Secret Keyword Filter](/code-scanning/scanning-practices/denoise-via-secret-keyword-filter.md)
{% endcontent-ref %}

{% content-ref url="/pages/3ucgOGrylofIUiFnuU6m" %}
[Ignoring Secrets via Code Comments](/code-scanning/scanning-practices/ignoring-secrets-via-code-comments.md)
{% endcontent-ref %}

## Dependency scanning (SCA)

{% content-ref url="/pages/E5ENgfMdszxSZ3Wiv6vc" %}
[Support for Dependency Scanning by Language (SCA)](/code-scanning/scanning-practices/support-for-dependency-scanning-by-language.md)
{% endcontent-ref %}

{% content-ref url="/pages/Q5xBsk1QQM1s39nt9iP6" %}
[External Vulnerability Databases Used in Our SCA Engine](/code-scanning/scanning-practices/external-vulnerability-databases-used-in-our-sca-engine.md)
{% endcontent-ref %}

{% content-ref url="/pages/1jyYLPd2ZIz9ZQuSBzr0" %}
[Scanning Dev Dependencies for CVEs](/code-scanning/scanning-practices/scanning-dev-dependencies-for-cves.md)
{% endcontent-ref %}

{% content-ref url="/pages/WpRAD3vySlicQx8qVurl" %}
[Custom CVE Vulnerabilities](/code-scanning/scanning-practices/custom-vulnerabilities.md)
{% endcontent-ref %}

{% content-ref url="/pages/jBdYI8FmaGbI1tXOtWPk" %}
[Exploit Intelligence with KEV and EPSS](/code-scanning/miscellaneous/use-epss-values-to-further-reduce-noise.md)
{% endcontent-ref %}

{% content-ref url="/pages/uEAjBWRDm3hwdHyMsljo" %}
[Malware Detection in Open-Source Dependencies and Containers](/code-scanning/scanning-practices/malware-detection-in-open-source-dependencies.md)
{% endcontent-ref %}

{% content-ref url="/pages/Zlug5MrIdwTmOpWskaiT" %}
[Vulnerability Scanning on Private Packages - Maven](/code-scanning/scanning-practices/vulnerability-scanning-on-private-packages-maven.md)
{% endcontent-ref %}

{% content-ref url="/pages/hwlmIJvGdm8tWcVoNUla" %}
[C/C++ Lockfile-Less Scanning](/code-scanning/scanning-practices/c-c-lockfile-less-scanning.md)
{% endcontent-ref %}

{% content-ref url="/pages/31VJfNMHbKVUYXVEsQGB" %}
[Dependency Scanning for Bazel](/code-scanning/scanning-practices/dependency-scanning-for-bazel.md)
{% endcontent-ref %}

{% content-ref url="/pages/xLHbt2lWPuwCEi7tLqKy" %}
[Scala: Dependency Management and Scanning for build.sbt](/code-scanning/scanning-practices/scala-dependency-management-and-scanning-for-buildsbt.md)
{% endcontent-ref %}

{% content-ref url="/pages/qHjsTC9yEvcwfKJZDRkj" %}
[Java/Scala/Kotlin Projects Using Gradle: Security Scanning Best Practices](/code-scanning/scanning-practices/java-scala-kotlin-projects-using-gradle-security-scanning-best-practices.md)
{% endcontent-ref %}

{% content-ref url="/pages/Ks7xutHzOKnEni37okWP" %}
[NET Projects: Security Scanning Best Practices](/code-scanning/scanning-practices/net-projects-security-scanning-best-practices.md)
{% endcontent-ref %}

{% content-ref url="/pages/YKXDQ8IS2DlWbTSu7euP" %}
[Python Projects: Security Scanning Best Practices](/code-scanning/scanning-practices/python-projects-security-scanning-best-practices.md)
{% endcontent-ref %}

## Branches and repositories

{% content-ref url="/pages/oBnCZyQr3T8DWJCIZEeQ" %}
[Scanned Branches in Aikido](/code-scanning/miscellaneous/scanned-branches-in-aikido.md)
{% endcontent-ref %}

{% content-ref url="/pages/yM4cD4bReBfKrURSXvRf" %}
[Support for Multi-Branch Scanning](/code-scanning/miscellaneous/support-for-multi-branch-scanning.md)
{% endcontent-ref %}

{% content-ref url="/pages/fh16v1KvUA8ZKN20MLG9" %}
[Support for Specific Branch Scanning](/code-scanning/miscellaneous/support-for-specific-branch-scanning.md)
{% endcontent-ref %}

{% content-ref url="/pages/Mv7JkMAQhkkjZwhXzqtn" %}
[Ensuring Aikido Scans New Repositories](/code-scanning/miscellaneous/ensuring-aikido-scans-new-repositories.md)
{% endcontent-ref %}

{% content-ref url="/pages/BPTrSiOLyHwhITo8KIpd" %}
[Split Your Monorepo Per Path](/code-scanning/miscellaneous/split-your-monorepo-per-path.md)
{% endcontent-ref %}

## Other

{% content-ref url="/pages/zM9cpCYm3FltJJrzQKKM" %}
[Configure scanning with .aikido files](/code-scanning/scanning-practices/ignore-via-code-with-aikido-files.md)
{% endcontent-ref %}
