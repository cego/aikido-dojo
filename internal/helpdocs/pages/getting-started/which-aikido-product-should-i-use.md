# Which Aikido product for what?

These are different solutions, not different names for the same thing. You do not need all of them on day one.

Writing a book is a useful way to see how the code-side products differ.

<table><thead><tr><th width="156.2734375">Product</th><th width="221.4453125">Book analogy</th><th width="367.28515625">Question it answers</th></tr></thead><tbody><tr><td><a href="https://help.aikido.dev/code-scanning/scanning-practices/sast-by-aikido-supported-languages-and-security-focus">SAST</a></td><td>Spell checker</td><td>Does this code contain a dangerous pattern?</td></tr><tr><td><a href="https://help.aikido.dev/code-quality/code-quality-overview">Code Quality</a></td><td>Style guide</td><td>Does this code meet your team's quality standards?</td></tr><tr><td><a href="https://help.aikido.dev/deep-review/how-deep-review-works">Deep Review</a></td><td>Review a new part in context of the full book</td><td>Does this PR/MR introduce any new issues?</td></tr><tr><td><a href="https://help.aikido.dev/ai-code-analysis/ai-code-analysis">AI Code Analysis</a></td><td>Proofread full version of the book (digitally)</td><td>Are there any existing issues from the past we need to fix?</td></tr><tr><td><a href="https://help.aikido.dev/pentests/aikido-pentest">Aikido Pentest</a></td><td>Proofread full version of the book (printed)</td><td>Are there any issues in our live environment we need to fix? Required for compliance.</td></tr></tbody></table>

SAST traces untrusted user input across files. It runs on connected repositories, and on pull requests through PR gating. The spell-checker line is the question it answers (known patterns), not a claim that the engine is simple. Code Quality is the same idea for quality: checks on new code as you add it.

Deep Review is not "SAST in CI" and not a second Code Quality scan. PR gating runs SAST (and Code Quality, if you enable it) on the diff. Deep Review reads the change with context beyond those lines, and it covers **both security and quality**.

## Internet-facing apps

These two sit outside the book analogy.

* [**Attack Surface Monitoring**](https://help.aikido.dev/dast-surface-monitoring/attack-surface-scanning) finds internet-facing assets you may not know about, such as a forgotten subdomain or dangling DNS.
* [**DAST**](https://help.aikido.dev/dast-surface-monitoring/dast-surface-monitoring-overview) ([Quick Scan](https://help.aikido.dev/dast-surface-monitoring/front-end-scanning/scan-front-end-app-domains-with-aikido)): continuously checks a known domain for exposure and config, such as security headers, cookies, and TLS. It is not the tool for deep injection or for proving that users can see each other's data. That is Pentest.

## What to pick

**A search box can be tricked into querying the database** (SQL injection). Start with SAST. You do not need a pentest to catch a known pattern.

**A pull request accidentally lets users see each other's orders.** The app still runs. Use Deep Review. SAST may see no dangerous pattern. Deep Review can also flag quality problems in that same change.

**The app has been live a long time and you are not sure customers are isolated from each other.** Use AI Code Analysis on the source.

**You need a pentest report for SOC 2, ISO, or a customer security review.** Use Pentest. Run it on a testing environment.

**An old login page is still on the internet, or a proxy change exposed an admin API.** Use Attack Surface Monitoring to discover what is reachable. Use Quick Scan for ongoing exposure on domains you already know.

## Start here

You do not have to run all of these at once.

1. **SAST and Code Quality:** [Connect repos](https://help.aikido.dev/getting-started/setting-up-your-account/create-account-and-connect-your-repositories) so SAST can scan connected repositories. Then enable SAST and Code Quality in [PR Checks](https://help.aikido.dev/pr-and-release-gating/aikido-ci-gating-functionality) for the repositories you care about, so new code is checked as you add it.
2. **On every pull request:** Enable Deep Review in PR Checks if you want a security and quality read of the change (beyond the SAST and Code Quality checks PR gating already runs).
3. **On purpose:** AI Code Analysis or Pentest when you need depth or proof, such as a release, an audit, or a high-risk app. Code Quality [repository scans](https://help.aikido.dev/code-quality/repository-scans) cover existing code.

For internet-facing apps, add Attack Surface Monitoring on your root domain and Quick Scan on domains you already know.

If you are in a heavily regulated industry such as finance or health tech, or you have a complex setup with connected systems, [Continuous Pentesting](https://help.aikido.dev/pentests/continuous-pentesting) retests the running app after deploys. That is not a substitute for Deep Review, which reads the pull request in source.

Credits are explained in [Wallet & Credits](https://help.aikido.dev/miscellaneous-info/wallet-and-credits).
