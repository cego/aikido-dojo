# Understanding Your Zen Statistics

Zen gather statistics to offer a clear, at-a-glance view of your application's traffic and Zen's protective actions.

### Stats tab <a href="#stats-tab" id="stats-tab"></a>

At the top of the page, you'll find two key graphs. The **Requests** graph shows the "Amount requests seen" by your application over time, helping you understand typical traffic patterns and spot unusual spikes. Below it, the **Attacks** graph displays the number of "Attacks blocked" and "Attacks detected" by Zen, visually confirming active protection and showing when attacks occur.

[image: App statistics dashboard showing request volume and detected attacks over time.]

### Invocation Details: Understanding Attack Paths (Sinks) <a href="#invocation-details-understanding-attack-paths-sinks" id="invocation-details-understanding-attack-paths-sinks"></a>

The **Invocation Details** section, also known as "Sinks," provides a detailed breakdown of the different attack vectors Zen monitors within your application. Use this to identify frequently targeted or potentially vulnerable parts of your application. These differ per language.

[image: Invocation statistics for HTTP, SQL, shell, and I/O operations over the past week.]

### Route-Specific Statistics <a href="#route-specific-statistics" id="route-specific-statistics"></a>

For granular insights, the **Routes** table details traffic and protection for each of your application's paths.

[image: API routes dashboard showing request rates, rate limits, and detection status for each endpoint.]

### Outbound-Specific Statistics <a href="#outbound-specific-statistics" id="outbound-specific-statistics"></a>

Finally, the **Outbound** section summarises all outbound calls & services Zen is currently monitoring. This view helps you quickly verify that all intended services are monitored and check their recent activity.

[image: Outbound network activity summary for My App, showing requests by hostname and last seen times.]
