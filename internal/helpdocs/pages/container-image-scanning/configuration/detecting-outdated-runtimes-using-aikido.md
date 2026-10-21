# Detecting Outdated Runtimes Using Aikido

### Runtimes we monitor <a href="#runtimes-we-monitor" id="runtimes-we-monitor"></a>

Once you [connected your docker registry](https://help.aikido.dev/container-image-scanning/container-image-scanning-overview), Aikido will start monitoring important runtimes.

Examples of monitored runtimes are:

* Python language runtime
* Node.js language runtime
* PHP language runtime
* Apache webserver
* Nginx webserver
* ..

We monitor these specific runtimes because they can have a big impact or are usually exposed to the web.

### Status overview of Runtimes <a href="#status-overview-of-runtimes" id="status-overview-of-runtimes"></a>

Our [Runtime Overview page](https://app.aikido.dev/containers/checks/runtimes) will give an overview of all runtimes. Different statuses are shown per monitored runtime.

* **Outdated**: This status indicates that the runtime has deprecated and is no longer maintained.
* **Almost Out-of-Date**: Assigned when a runtime is nearing the deprecation date, specifically when it will expire in 30 days.
* **Up-to-Date**: This status confirms that the runtime is current and fully updated.
* **Not Found**: This status is given when no runtime is found after Aikido scanning.

You can view the list of monitored runtimes in Aikido by clicking [here](https://app.aikido.dev/containers/checks/runtimes).

### Alerts in the feed <a href="#alerts-in-the-feed" id="alerts-in-the-feed"></a>

Once we detect a specific package version is no longer maintained by the vendor, Aikido will send an alert to your feed.

The first alert will come 2 months before the deprecation date. As the date grows closer, the severity of the issue will become higher. The following logic is applied.

* **Low Severity**: For runtimes expiring in the next 60 days.
* **Medium Severity** *(60)*: When the runtime expired within the last 30 days.
* **Lower End - High Severity -** *(70)*: If the runtime expired between 30 and 60 days ago.
* **High End - High Severity** *(89–99)*: For runtimes expired more than 60 days ago. The score keeps increasing the longer the runtime stays unmaintained, reaching 97 after 2 years and 99 after 3 years.

***
