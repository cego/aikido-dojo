# Connect GitLab Self-Managed Server (Broker Set-Up)

Aikido lets you connect a self-managed GitLab instance to scan and secure your code. Follow the steps below to connect your GitLab server to Aikido.

### Before you start

* The person setting this up needs access to both the GitLab instance **and** the GitLab group you want to connect.

{% hint style="warning" %}
An Aikido workspace always maps to a single GitLab group.

We recommend connecting Aikido to a top-level (root) group that contains all subgroups. If you don’t have a root group, create one workspace per GitLab group. You can do this after creating your first workspace via the top-left dropdown: “Add another workspace”.
{% endhint %}

## Configuration

{% stepper %}
{% step %}
**Create an Aikido account**

To connect your GitLab server, first sign up or log in to Aikido using Google or Microsoft. On the [signup screen](https://app.aikido.dev/login), click **Google / Microsoft**.
{% endstep %}

{% step %}
**Start the GitLab setup**

Once you’re authenticated, create a new workspace by clicking **Self-Managed** in the GitLab section.

{% endstep %}

{% step %}
**Enable Broker access**

Enable the **Aikido broker** by clicking the toggle in the "Advanced Connection Options" section.

Click **Create Group** to continue.

{% endstep %}

{% step %}
**Configure the Broker**

Configure the Broker using the guide below. After you deploy it, return here to finish the GitLab setup.

{% content-ref url="/pages/h0JtJS4MeX5E8efcCvRF" %}
[Aikido Broker for Internal Applications](/miscellaneous-info/aikido-broker-for-internal-applications.md)
{% endcontent-ref %}

{% endstep %}

{% step %}
**Complete GitLab server setup**

After you create the broker and deploy it in your environment, you can complete the GitLab setup.

First, copy the **Broker URL**. It looks similar to this:

```
https://55673-355dsfwea68cb.aikidobroker.com
```

Then click **Complete GitLab Server Setup**.

{% endstep %}

{% step %}
**Enter the Broker URL**

Paste the Broker resource URL from the previous step in the input field for the "server URL".

{% endstep %}

{% step %}
**Create a GitLab personal access token (PAT)**

Next, create a personal access token (PAT). We recommend using a dedicated service account

* Log in to your GitLab server
* Go to the admin area (`/admin`)
* Go to **Settings** → **Service accounts**
* Click **Add service account**, enter a name, then save

{% hint style="warning" %}
Add the service account to the GitLab group you want to connect, like any other user. See GitLab’s docs on [group members](https://docs.gitlab.com/ee/user/group/members/).
{% endhint %}

Now that the service account is created, you can create a PAT for it by clicking the three dots and select "**Manage access tokens**"

* Click on "**Add new token**"
* Enter a name for the token, for example: `Aikido Security Access Token`
* Set an expiration date that matches your internal policy. Rotate the token before it expires.
* We need the following scopes to be selected:
  * **read\_user**
  * **read\_api**
  * **read\_repository**
* Click the **Create token** button at the bottom of the form.

* Copy the token being shown on the screen and enter it in the input field.

**Important:** You won’t be able to see the token again after you leave this screen. Copy it before you continue.

Aikido will now check the connection to your GitLab server. If it fails, double-check the server URL and the token.
{% endstep %}

{% step %}
**Complete the installation**

After you click **Next, Connect Group**, select the group you want to start with. You can always connect more groups later.

In the final step, select the repositories you want Aikido to monitor.
{% endstep %}
{% endstepper %}

{% hint style="info" %}
Also scanning your GitLab container registry through the Broker? The registry usually runs on its own host that needs a separate Broker resource. See [GitLab Self Managed Container Registry](/container-image-scanning/standalone-registries/gitlab-self-managed-container-registry.md#connect-a-registry-on-an-internal-network-broker).
{% endhint %}
