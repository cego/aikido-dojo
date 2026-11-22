# Protect Your App with Zen Firewall

Installing Zen Firewall is the first step. Work through this checklist next to get the most protection out of Zen. Each item links to the full guide.

## Turn on Blocking Mode

Zen starts in Detection Mode: it reports injection attacks but lets them through. Switch to Blocking Mode in the **Block Attacks** section of your app's **Firewall** tab so Zen stops them.

See [Blocking vs Detection Mode in Zen Firewall](https://help.aikido.dev/zen-firewall/zen-features/blocking-vs-detection-mode-in-zen-firewall).

## Block attackers with playbooks

Blocking Mode stops a single attack request. A playbook blocks the IP behind it, so the attacker can't keep trying. Add playbooks for **Attack wave detected** and **Injection attack detected** on the **Playbooks** tab.

See [Playbooks](https://help.aikido.dev/zen-firewall/zen-features/playbooks).

## Block known threat actors

Block IPs that threat intelligence lists flag as botnets, brute forcers, or exploit scanners. Click **Manage Threat Actors** on the **Firewall** tab.

See [Blocking Known Threat Actors with Zen Firewall](https://help.aikido.dev/zen-firewall/zen-features/blocking-known-threat-actors-with-zen-firewall).

## Monitor or block bots

Choose per bot category whether to ignore, monitor, or block it, for example scrapers or vulnerability scanners. Click **Manage Bots** on the **Firewall** tab.

See [Blocking Bot Traffic with Zen Firewall](https://help.aikido.dev/zen-firewall/zen-features/blocking-bot-traffic-with-zen-firewall).

## Block traffic by country

Block countries you don't serve, or allow only the ones you do. Click **Manage Countries** on the **Firewall** tab.

See [Blocking Traffic by Country with Zen Firewall](https://help.aikido.dev/zen-firewall/zen-features/blocking-traffic-by-country-with-zen-firewall).

## Identify your users

Call `setUser` in your app so Zen links attacks and requests to a user. You can then block a user or rate limit per user instead of per IP.

See [Tracking Users with Zen Firewall](https://help.aikido.dev/zen-firewall/zen-features/blocking-users-with-zen-firewall).

## Set up rate limiting

Limit how often a route can be called, for example login or password reset, to slow down brute force and abuse. Set it up per route on the **Routes** tab. Rate limiting isn't available in serverless environments.

See [Setting Up Rate Limiting for Routes](https://help.aikido.dev/zen-firewall/zen-features/setting-up-rate-limiting-for-routes).

## Review outbound domains

Check which external domains and IPs your app talks to, and turn on alerts for new ones. Open the **Outbound Domains** tab.

See [Monitor Outbound Domains](https://help.aikido.dev/zen-firewall/zen-features/monitor-outbound-domains).
