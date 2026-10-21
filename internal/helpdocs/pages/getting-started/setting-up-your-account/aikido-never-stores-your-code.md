# Aikido Never Stores Your Code

{% hint style="success" %}
In short: Aikido does not store your code after analysis has taken place. Some of the analysis jobs such as SAST or Secrets Detection require a git clone operation. Below we talk about the technical measures we take to ensure your code is protected:
{% endhint %}

* We perform different actions such as git clones in a fresh docker container for each repository. After analysis, the data is wiped and the docker container is terminated.
* For GitHub, no refresh or access tokens are ever stored in our database. We use the new GitHub Apps which do not require this. Even a database breach of Aikido itself would not result in your GitHub code being downloadable.
* By default, our integrations require a very minimal read-only scope. Only if you enable special features such as Autofix Pull Requests, Aikido will request write accesses.
* If you want to keep your code completely on-premise, without ever leaving your environment, you can use our [Local Scanner](https://help.aikido.dev/category/aikido-local-scan-setup/sg4xF4OsJciW). The results will seamlessly populate on the Aikido platform.
* Aikido has SOC2 Type 2 & ISO27001:2022 certification. A report is available [upon request](http://trustcenter.aikido.dev/). That means we adhere to several organizational and technical policies by default.
* Aikido's servers are located in AWS, with data stored in-region for EU, US, AU and ME. Meaning customer data is stored for [app.aikido.dev](http://app.aikido.dev) in the EU region, for [app.us.aikido.dev](http://app.us.aikido.dev) in the US region, for [app.me.aikido.dev](http://app.me.aikido.dev) in the ME region and for [app.au.aikido.dev](http://app.au.aikido.dev) in the AU region.

The process we use to ensure code security:

[image: Secure repository scanning workflow: select, clone, scan, encrypt findings, destroy containers.]

**Disclaimer.**

Aikido has some features where certain parts of your code are stored. This is in the case for the following functionalities:

* AutoFix: Aikido stores the diffs (original and AutoFixed code) - only files that are part of the AutoFix
* Aikido stores the calltree for each AutoTriaged SAST finding for up to 2 weeks

All code that is stored is ran through Gitleaks. If there are any obvious secrets in the code, we make sure to definitely not store these.
