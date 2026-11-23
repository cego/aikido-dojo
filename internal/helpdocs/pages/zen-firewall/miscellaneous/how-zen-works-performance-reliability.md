# How Zen Works: Performance & Reliability

## Performance Impact {#performance-impact}

At Zen, we prioritize your application's performance. Our agents are designed to have minimal impact on your systems, and we continuously monitor this through testing, automated benchmarks and using Zen ourselves.

## Reliability & Failsafes {#reliability--failsafes}

Zen is designed to maintain security and performance even when facing connectivity issues with the Aikido platform. Here's how Zen handles different scenarios:

### When Zen Starts Without Server Connection  {#when-zen-starts-without-server-connection-}

- ✅ Your application continues running normally
- ✅ Basic attack blocking (sql, path injection, ..) remain active
- ⚠️ Some advanced security features (bot detection, threat actors, ..) may be limited until connection is restored

### During Normal Operation If Zen Loses Connection {#during-normal-operation-if-zen-loses-connection}

- ✅ All existing security configurations remain active
- ✅ Your application continues running without interruption
- ✅ Previously downloaded security rules stay in effect
- ⚠️ New security updates will queue until connection resumes

### Request Processing  {#request-processing-}

- ✅ All security checks happen locally on your machine
- ✅ No dependency on Zen server response times. Zen talks to Aikido servers in a background process and not for each request. 

### Configuration Updates  {#configuration-updates-}

Zen maintains strict validation of all configuration updates to ensure:

- ✅ Only valid configurations are applied
- ✅ Invalid configurations are safely rejected
- ✅ Your security posture remains stable

## Best Practices {#best-practices}

- Keep your Zen agent updated to the latest version
- Monitor your application logs for any Zen-related messages
- Contact support if you notice any unexpected behavior
