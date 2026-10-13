# Blocking Bot traffic with Zen Firewall

Zen Firewall by Aikido helps you identify and block various types of automated traffic accessing your application. This feature enhances your security measures by providing control over unwanted bot activity.

> **Important**: Bot blocking operates independently of the global "[Blocking/Detection Mode](/zen-firewall/zen-features/blocking-vs-detection-mode-in-zen-firewall.md)" setting. When you enable bot blocking, it will always be enforced, even if Zen is in Detection Mode.

## Use Cases <a href="#use-cases" id="use-cases"></a>

* 🛡️ **Protect Server Resources**: Prevent bots from overwhelming your infrastructure with excessive requests
* 💰 **Reduce Costs**: Lower bandwidth and computing expenses by blocking unnecessary bot traffic
* 🔒 **Secure Content**: Protect your intellectual property from automated scraping and copying
* ⚡ **Improve Performance**: Enhance site speed by reducing bot-related server load

## Bot Categories We Detect <a href="#bot-categories-we-detect" id="bot-categories-we-detect"></a>

Aikido groups automated traffic into clear categories so you can decide what to block, allow, or monitor. A full list of specific user agents is available on the [dedicated bot list page](/zen-firewall/miscellaneous/bot-protection-details.md).

Examples of the types of categories available:

* **Search & Discovery Bots**\
  Includes traditional search engines, advertisement crawlers, SEO tools, social media bots, and messenger preview crawlers. These bots index, analyze, or preview your content for search, ads, or sharing.
* **AI Crawlers & Assistants**\
  Covers AI search crawlers, AI data scrapers, and AI assistants that browse or extract content to power generative AI systems or respond to user queries.
* **Archiving & Monitoring Services**\
  Includes web archivers and uptime or performance monitoring tools that periodically access your pages.
* **Security & Exploitation Tools**\
  Vulnerability scanners, data harvesters, headless browsers, and other automation tools commonly used for probing, scraping, or abuse.

{% hint style="info" %}
**Tip**: Consider carefully which bot categories to block, as some legitimate services (like Google's search crawler) might be necessary for your site's functionality.
{% endhint %}

## How to block bot traffic <a href="#how-to-block-bot-traffic" id="how-to-block-bot-traffic"></a>

Select a specific app and continue to the **Firewall** tab. Click the "**Manage Bots**" next to "**Manage Bots**" to configure Bot blocking.

[image: Firewall settings interface displaying attack, Tor traffic, and bot blocking options for a demo app.]

Use the **Bots** view to select the bot types you want to ignore, monitor or block and click on "**Update bots**" when you're finished.

[image: Interface for selecting and blocking specific bots from accessing your app.]

{% hint style="info" %}
Note that bot blocking is not immediate; it takes up to a minute for the block to take effect.
{% endhint %}
