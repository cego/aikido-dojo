# Java (Spring MVC)

[Zen Firewall by Aikido](https://www.aikido.dev/zen) helps protect your application by blocking common attacks and unsafe behavior at runtime, with controls you can tune per app and environment. Use the guides below to install and set it up.

<https://github.com/AikidoSec/firewall-java>

## Requirements

* Java 17+ (tested up to Java 21; also supports Kotlin/Groovy with same agent).
* Spring MVC 3.x.
* [Aikido account](/getting-started/setting-up-your-account.md) & [Zen Firewall token](/zen-firewall/zen-installation-instructions/creating-an-aikido-zen-firewall-token.md)

## Installation & Configuration

{% stepper %}
{% step %}
**Install Zen Firewall by Aikido**

Download and extract the latest Java agent bundle:

{% tabs %}
{% tab title="zip" %}

```bash
curl -L https://github.com/AikidoSec/firewall-java/releases/latest/download/zen.zip -o zen.zip
unzip zen.zip
```

{% endtab %}

{% tab title="tar.gz" %}

```bash
curl -L https://github.com/AikidoSec/firewall-java/releases/latest/download/zen.tar.gz -o zen.tar.gz
tar -xzf zen.tar.gz
```

{% endtab %}
{% endtabs %}

Keep the extracted folder structure intact (for example `/opt/zen`) and make sure your process can read and write it.

Add the Java agent to your runtime:

```bash
-javaagent:/opt/zen/agent.jar
```

For rate limiting and user blocking features, add `agent_api.jar` to your project.

Gradle:

```gradle
dependencies {
    implementation files('/opt/zen/agent_api.jar')
}
```

Maven:

```xml
<dependency>
  <groupId>dev.aikido</groupId>
  <artifactId>agent_api</artifactId>
  <version>1.0</version>
  <systemPath>/opt/zen/agent_api.jar</systemPath>
</dependency>
```

{% endstep %}

{% step %}
**Enable Request Blocking and User Identification**

Use this middleware to enable rate limiting, user identification, and blocking features.

Zen Firewall does not require this middleware to block attacks. Core attack protection works without it. The middleware provides additional request context that Zen uses for protections such as [rate limiting](/zen-firewall/zen-features/setting-up-rate-limiting-for-routes.md), [user blocking](/zen-firewall/zen-features/blocking-users-with-zen-firewall.md), [bot blocking](/zen-firewall/zen-features/blocking-bot-traffic-with-zen-firewall.md), [country blocking](/zen-firewall/zen-features/blocking-traffic-by-country-with-zen-firewall.md), and [threat actor](/zen-firewall/zen-features/blocking-known-threat-actors-with-zen-firewall.md) blocking.

Adapt the example to fit how your application identifies users and handles requests.

{% @aikido-custom-code/code-highlight language="java" content="+\@Component
+\@Order(2)
+public class RateLimitingFilter implements Filter {

* @Override
* public void doFilter(ServletRequest request, ServletResponse response, FilterChain chain)
* ```
         throws IOException, ServletException {
  ```
* ```
     ShouldBlockRequest.ShouldBlockRequestResult result = ShouldBlockRequest.shouldBlockRequest();
  ```
*
* ```
     if (result.block()) {
  ```
* ```
         if (result.data().type().equals("ratelimited")) {
  ```
* ```
             String message = "You are rate limited by Zen.";
  ```
* ```
            if (result.data().trigger().equals("ip")) {
  ```
* ```
                message = message + " (Your IP: " + result.data().ip() + ")";
  ```
* ```
            }
  ```
* ```
            HttpServletResponse httpResponse = (HttpServletResponse) response;
  ```
* ```
            httpResponse.setStatus(429);
  ```
* ```
            httpResponse.getWriter().write(message);
  ```
* ```
            return;
  ```
* ```
        }
  ```
* ```
        if (result.data().type().equals("blocked")) {
  ```
* ```
            HttpServletResponse httpResponse = (HttpServletResponse) response;
  ```
* ```
            httpResponse.setStatus(403);
  ```
* ```
            httpResponse.getWriter().write("You are blocked by Zen.");
  ```
* ```
            return;
  ```
* ```
        }
  ```
* ```
    }
  ```
* ```
    chain.doFilter(request, response);
  ```
* }
  +}" %}

Set users in an earlier filter:

{% @aikido-custom-code/code-highlight language="java" content="+import dev.aikido.agent\_api.SetUser;
\+
+\@Component
+\@Order(0)
+public class SetUserFilter implements Filter {

* @Override
* public void doFilter(ServletRequest request, ServletResponse response, FilterChain chain)
* ```
        throws IOException, ServletException {
  ```
* ```
    SetUser.setUser(new SetUser.UserObject("123", "John Doe"));
  ```
* ```
    chain.doFilter(request, response);
  ```
* }
  +}" %}
  {% endstep %}

{% step %}
**Start Zen Firewall in dry / detection-only mode**

Start your app with Zen enabled and dry mode:

```bash
AIKIDO_TOKEN=AIK_RUNTIME_... AIKIDO_BLOCK=false java -javaagent:/opt/zen/agent.jar -jar build/myapp.jar
```

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
