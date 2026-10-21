# Zen Performance & Reliability

## Performance Impact <a href="#performance-impact" id="performance-impact"></a>

At Zen, we prioritize your application's performance. Our agents are designed to have minimal impact on your systems, and we continuously monitor this through testing, automated benchmarks and using Zen ourselves.

## Reliability & Failsafes <a href="#reliability--failsafes" id="reliability--failsafes"></a>

Zen is designed to maintain security and performance even when facing connectivity issues with the Aikido platform. Here's how Zen handles different scenarios:

### When Zen Starts Without Server Connection <a href="#when-zen-starts-without-server-connection" id="when-zen-starts-without-server-connection"></a>

* ✅ Your application continues running normally
* ✅ Basic attack blocking (sql, path injection, ..) remain active
* ⚠️ Some advanced security features (bot detection, threat actors, ..) may be limited until connection is restored

### During Normal Operation If Zen Loses Connection <a href="#during-normal-operation-if-zen-loses-connection" id="during-normal-operation-if-zen-loses-connection"></a>

* ✅ All existing security configurations remain active
* ✅ Your application continues running without interruption
* ✅ Previously downloaded security rules stay in effect
* ⚠️ New security updates will queue until connection resumes

### Request Processing <a href="#request-processing" id="request-processing"></a>

* ✅ All security checks happen locally on your machine
* ✅ No dependency on Zen server response times. Zen talks to Aikido servers in a background process and not for each request.

### Configuration Updates <a href="#configuration-updates" id="configuration-updates"></a>

Zen maintains strict validation of all configuration updates to ensure:

* ✅ Only valid configurations are applied
* ✅ Invalid configurations are safely rejected
* ✅ Your security posture remains stable

## Best Practices <a href="#best-practices" id="best-practices"></a>

* Keep your Zen agent updated to the latest version
* Monitor your application logs for any Zen-related messages
* Contact support if you notice any unexpected behavior
