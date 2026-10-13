# ASP.NET Framework (including OWIN)

This guide will walk you through installing and setting up Zen Firewall by Aikido for your application. Follow the steps below to protect your application.

If you encounter any issues or problems, don't hesitate reach out on support chat or Github issues

<https://github.com/AikidoSec/firewall-dotnet>

## Requirements

* C#/.NET Framework 4.6+ (tested up to 4.8.x).
* ASP.NET Framework (OWIN supported).
* [Aikido account](/getting-started/setting-up-your-account.md) & [Zen Firewall token](/zen-firewall/zen-installation-instructions/creating-an-aikido-zen-firewall-token.md)

## Installation & Configuration

{% stepper %}
{% step %}
**Install Zen Firewall by Aikido**

Ensure your project runs on .NET Framework 4.6 or higher.

Install the package from NuGet:

```bash
dotnet add package Aikido.Zen.DotNetFramework
```

Or, in the Package Manager Console:

```powershell
Install-Package Aikido.Zen.DotNetFramework
```

Add your token to `Web.config`:

```xml
<appSettings>
  <add key="Aikido:AikidoToken" value="your-api-key" />
</appSettings>
```

Start Zen when the app starts:

{% @aikido-custom-code/code-highlight language="csharp" content="! Global.asax.cs
+using Aikido.Zen.DotNetFramework;

protected void Application\_Start()
{
// other code

* Zen.Start();
  }" %}

If you use OWIN:

{% @aikido-custom-code/code-highlight language="csharp" content="! Startup.cs
+using Aikido.Zen.DotNetFramework;

public void Configuration(IAppBuilder app)
{
// other code

* Zen.Start();
  }" %}
  {% endstep %}

{% step %}
**Enable User Identification & Blocking**

Use this middleware to identify the current user and enable user based features such as user blocking.

Zen Firewall does not require this middleware to block attacks. Core attack protection works without it. Add this middleware when you want Zen to recognize which user is making the request.

Adapt the example to match how your application identifies users.

{% @aikido-custom-code/code-highlight language="csharp" content="! Global.asax.cs
public void Application\_Start()
{

* Zen.SetUser(context => new User(context.User.Identity.Name, context.User.Identity.Name));
  Zen.Start();
  }" %}

If you use OWIN:

{% @aikido-custom-code/code-highlight language="csharp" content="! Startup.cs
public void Configuration(IAppBuilder app)
{

* Zen.SetUser(context => new User(context.User.Identity.Name, context.User.Identity.Name));
  Zen.Start();
  }" %}
  {% endstep %}

{% step %}
**Start Zen Firewall in dry / detection-only mode**

Set the token as an environment variable so the Aikido Zen agent can pick it up. If you don't have a token yet, follow [instructions here](/zen-firewall/zen-installation-instructions/creating-an-aikido-zen-firewall-token.md).

```bash
AIKIDO_TOKEN=AIK_RUNTIME_
```

We recommend to start your app in dry mode to ensure it works as expected without blocking any requests. We advise running Zen Firewall in staging for two weeks to avoid false positives.

```bash
AIKIDO_BLOCK=false
```

{% hint style="info" %}
You can use `AIKIDO_DEBUG=true` to enable debug mode for more detailed information about what the agent is doing. For more information about your environment variables: [Configuration via Environment Variables](/zen-firewall/zen-installation-instructions/configuration-via-environment-variables.md)
{% endhint %}
{% endstep %}

{% step %}
**Test your app**

Browse to your application and perform a couple of actions or open a couple of pages. Zen will automatically discover the routes in your application.

{% hint style="info" %}
Zen sends data back to Aikido every 10 minutes
{% endhint %}

You can verify a working agent by looking at the following pages of your Zen application:

* **Events**: Should show an "Application started" event.
* **Routes**: After some time your application routes will start showing here with the method, route and requests.
* **Instances**: Should show the number of active instances for your application where Zen is installed.

{% endstep %}

{% step %}
**Setup rate limiting in the dashboard**

When you've added the Zen Firewall middleware you can test protecting a route from brute force attacks, you do this by setting up rate limit in the Aikido Dashboard:

1. Click on the created app.
2. Go to the **Routes** tab.
3. Find the route you would like to limit and click **Setup rate limiting**.
4. Follow the instructions to configure the rate limit (e.g., 5 requests per minute).

[image: API route management interface showing authentication routes with protection and rate limiting options.]

[image: Set rate limiting for POST /auth/login to 5 requests per minute.]

**Verify Rate Limiting**

Start your app and try to access the route you've rate limited 5 times within a minute. After the fifth attempt, you should receive a rate limit error:

```
You are rate limited by Aikido firewall. (Your IP: 1.2.3.4)
```

{% endstep %}

{% step %}
**Next steps**

Congrats you've successfully installed Zen Firewall. If you encountered any problems, have concerns or feature requests, don't hesitate to reach out to support.

You can now go and explore the many features that Zen Firewall provides:

* [Blocking Bot traffic with Zen Firewall](/zen-firewall/zen-features/blocking-bot-traffic-with-zen-firewall.md)
* [Blocking or Monitoring Tor traffic with Zen Firewall](/zen-firewall/zen-features/blocking-tor-traffic-with-zen-firewall.md)
* [Tracking Users with Zen Firewall](/zen-firewall/zen-features/blocking-users-with-zen-firewall.md)
* [Blocking Known Threat Actors with Zen Firewall](/zen-firewall/zen-features/blocking-known-threat-actors-with-zen-firewall.md)
* [Blocking Traffic by Country with Zen Firewall](/zen-firewall/zen-features/blocking-traffic-by-country-with-zen-firewall.md)
* [Setting Up Rate Limiting for Routes](/zen-firewall/zen-features/setting-up-rate-limiting-for-routes.md)
* [Monitor Outbound Domains](/zen-firewall/zen-features/monitor-outbound-domains.md)

Additional information:

* [Zen Performance & Reliability](/zen-firewall/miscellaneous/how-zen-works-performance-reliability.md)
* [Blocking vs Detection Mode in Zen Firewall](/zen-firewall/zen-features/blocking-vs-detection-mode-in-zen-firewall.md)
* [Understanding Your Zen Statistics](/zen-firewall/zen-features/understanding-your-zen-statistics.md)
  {% endstep %}
  {% endstepper %}
