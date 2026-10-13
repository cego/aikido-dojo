# Transferring Ownership of the Aikido App in Github Enterprise (Cloud & Server)

When you create the Aikido GitHub App inside your GitHub Enterprise environment, it is set up under the account of the user performing the installation. This follows the default GitHub flow when no prior ownership choice is known.

You can transfer the app to any other user or organization in your GitHub Enterprise Server environment at any time. This has no impact on how the app operates.

To transfer ownership, follow the steps below.

{% stepper %}
{% step %}
**Go to the Github App's developer Settings page**

In your GitHub Server environment, navigate to Settings > Developer Settings > GitHub Apps.\
Click on "Edit".

{% endstep %}

{% step %}
**Go to Advanced settings**

In the left-hand navigation, go to the "Advanced" settings and scroll all the way down on the page until you see this section.\
\
Click on "Transfer ownership"<br>

{% endstep %}

{% step %}
**Transfer ownership**

Once you've clicked the button, you need to confirm the name of the App, and then you can select the organization or user to which you want to transfer the ownership.

{% endstep %}

{% step %}
**Success**

The App is now transferred to the new owner.
{% endstep %}
{% endstepper %}
