# Smoke testing

AutoShip performs a smoke test on a live environment to verify the changes

AutoShip validates the changes are not impacting the functionality of the application by running a smoke test on a live environment.

Smoke tests are configured per respository and environment, for each environment you specify:

* the domain AutoShip should connect to
* the credentials AutoShip should use to authenticate
* the instructions to follow to execute the smoke test

AutoShip will tell you wether the smoke test passed or failed and give you a report with screenshots on what it actually tested.

## IP addresses for AutoShip Smoke Testing

Aikido uses dedicated IP addresses to perform smoke testing in your environment. To prevent connectivity issues, rate limiting, or security blocks, add these IPs to your firewall’s allowlist or other security software:

* 34.252.102.184
* 52.48.122.82
* 52.49.182.62
* 52.209.168.11
* 52.210.210.125
* 54.76.103.212
* 54.194.175.200
* 54.217.255.121

### Private networks <a href="#testing-internal-applications" id="testing-internal-applications"></a>

If your application is hosted on a strictly private network, VPN, or intranet that cannot accept inbound traffic from the public internet, you can use [Aikido Broker](https://help.aikido.dev/pentests/configure-a-pentest/pentesting-with-the-broker). Install the broker inside your network, register the internal URLs, and add the **internal domains** in the smoke test configuration, AutoShip will automatically connect through the broker. Inbound allowlisting is not required.
