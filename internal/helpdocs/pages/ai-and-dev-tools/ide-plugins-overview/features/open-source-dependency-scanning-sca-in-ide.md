# Open-Source Dependency Scanning (SCA) in IDE

Aikido’s IDE extension helps you find and fix vulnerabilities in your open-source dependencies without leaving your editor. It scans your project’s manifests and lockfiles to detect outdated or insecure packages, highlight affected versions, and suggest safe upgrades.

### How it works

* Aikido reads your dependency manifests and lockfiles to build an accurate list of packages and versions.
* Results include known CVEs, severity, affected versions, and safe upgrade ranges.
* After you run a manual SCA scan once, the extension watches your workspace for lockfile changes and refreshes results automatically.

### Run a manual SCA scan

#### VSCode

1. Open the Aikido sidebar in VS Code.
2. Go to Open-source dependencies.
3. Click Start scanning.
4. When results appear, select a package to view details, advisories, and fix guidance.
5. Each finding shows the minimum safe version or version ranges that resolve the issue.
6. For supported ecosystems, [AI AutoFix](/ai-and-dev-tools/ide-plugins-overview/features/aikido-ai-in-ide.md) can update the manifest or suggest a safe version bump that you can apply from the IDE.

#### Eclipse

1. Open "Aikido Workspace Scan"
2. Press "Scan All Code"
3. When results appear, select a package lockfile to view details, advisories, and fix guidance.
4. Each finding shows the minimum safe version or version ranges that resolve the issue.
