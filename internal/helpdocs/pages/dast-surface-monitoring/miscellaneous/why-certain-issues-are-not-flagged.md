# Why We Flag Certain Issues the Way We Do for Domains & API's

Aikido's Domain and API scans are intentionally opinionated. We don't report every deviation from a best-practice checklist, and where we do report something, the severity reflects how urgently we think you should fix it rather than how the finding is named.

That means our results won't always line up with other scanners. This page explains the reasoning behind the decisions customers ask about most: issues we deliberately don't flag, and issues we rate higher than you might expect.

If our rating doesn't fit your environment, you can [manually adjust the severity](/getting-started/core-functionalities/manually-adjust-issue-severity.md) of a single issue or a whole issue group.

## Issues we don't flag

### TLS 1.2

We don't flag the use of TLS 1.2 on its own. TLS 1.2 is still considered secure as long as it's configured with strong cipher suites. Alerting on TLS 1.2 itself would create noise on domains that are perfectly safe.

Instead, we flag the insecure ciphers in use. If a domain runs TLS 1.2 with weak ciphers, you'll see a finding for those ciphers. If the ciphers are strong, there's nothing to fix, and no alert.

This is consistent with our CSPM and IaC scanning, where the baseline is a minimum of TLS 1.2 and only weak ciphers are flagged.

### Missing Permissions-Policy header

We don't flag a missing `Permissions-Policy` header. This header lets a page declare which native browser features its JavaScript is allowed to request, such as the microphone, camera or geolocation.

Its absence is low impact. Even without the header, the browser still prompts the user to accept or deny access to those native features. The header restricts what your own scripts can ask for, it doesn't open a hole an attacker can use. Flagging it would add noise without pointing to a real risk.

## Severity decisions

### Missing Content-Security-Policy header flagged as critical

A Content Security Policy (CSP) is a header that tells the browser which sources of scripts and other content it's allowed to load. Without one, any script that reaches your page runs with full access to it, including your users' input and their session.

We rate a missing CSP as high because of how much protection it buys. It's the strongest defense available against cross-site scripting (XSS, where an attacker gets their own JavaScript to run in your users' browsers), and it's a standard expectation on any public-facing app. There's no good reason to run without one.

We don't wait for a confirmed XSS before flagging it. Our severity reflects [urgency to fix](/getting-started/core-functionalities/how-is-severity-score-calculated.md), not proof of an exploit today.
