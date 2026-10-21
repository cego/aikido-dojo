# Attack Wave Protection

The internet is constantly being scanned by automated tools. Security scanners like Nuclei, ZAP, Wapiti, and SQLMap probe applications for weaknesses, looking for sensitive files, testing hidden directories, and injecting malicious payloads. These reconnaissance attempts flood your servers with requests that waste CPU cycles and generate endless 404 errors.

## How It Works

Zen automatically detects these attack waves in real time by spotting probing patterns. For example:

* An IP making 15 suspicious requests within a minute.
* Attempts to access files like .env or wp-config.php.
* Requests for hidden directories like .git/config.
* Payloads designed to trigger SQL injection errors.<br>

When this behavior is detected, Zen flags it as an attack wave. This shifts your security posture from reactive to proactive, giving you instant visibility into who is systematically probing your applications.<br>

## What You’ll See

A live feed in your Aikido dashboard showing which IPs are targeting your applications.

### Try It Yourself

You can easily simulate an attack wave to confirm Zen is working. Just click Simulate Attack in your dashboard and watch the events appear.

## Block attack waves automatically

Use a [playbook](/zen-firewall/zen-features/playbooks.md) to block the IP behind an attack wave for a set period, without waiting for someone to act on the alert.
