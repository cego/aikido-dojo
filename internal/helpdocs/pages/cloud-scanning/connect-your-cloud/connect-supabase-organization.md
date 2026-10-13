# Connect Supabase Organization

### Why connect my Supabase organization?

Securing your cloud infrastructure is crucial to protecting your user data. You can leverage Aikido's security checks to detect and address any misconfigurations in your Supabase projects.

#### **Main use cases**

* Aikido surfaces critical cloud misconfigurations that could expose your Supabase projects to attackers - things like overly permissive policies, weak authentication settings, or misconfigured storage buckets. All configuration checks can be found [here.](https://app.aikido.dev/clouds/checks)
* Aikido continuously monitors your Supabase organization for new risks as your setup evolves.

Aikido performs daily compliance scans on the above.

### Getting started <a href="#getting-started" id="getting-started"></a>

To get started, head to the [cloud overview page](https://app.aikido.dev/clouds) on Aikido and click **"Connect Cloud"**. Select **Supabase** from the list and follow the step-by-step setup wizard.

{% stepper %}
{% step %}
**Authorize Aikido in Supabase**

Aikido will redirect you to Supabase to authorize API access. Review the permissions, select the organization you want to connect from the dropdown, and click **"Authorize Aikido Security"**.\
\

{% endstep %}

{% step %}
**Name your cloud configuration**

Give your connected project a name in Aikido and specify the environment it operates in (production, staging, etc.). This helps Aikido prioritize findings based on severity and business impact. Click **"Save"**.\
\

{% endstep %}
{% endstepper %}

Within 1-2 minutes after connecting your account, Aikido will report misconfigurations that could pose a threat.
