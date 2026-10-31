# Why is this IP blocked?

The **Why is this IP blocked?** tool lets you look up any IP address and see exactly which Zen Firewall rules are responsible for blocking it. It checks all active blocking features at once, so you can quickly diagnose unexpected blocks or verify that a rule is working as expected.

## How to use it

1. Go to the **Events** feed in a Zen Firewall service.
2. Click **Why is this IP blocked?** in the top right above the table.
3. Enter an IP address in the input field.
4. Click **Check**.

Zen will show you which blocking rules match that IP, such as country blocks, known threat actor lists, user blocks, or a [playbook](/zen-firewall/zen-features/playbooks.md).

{% hint style="info" %}
This tool checks the current state of your blocking rules. If a rule was recently changed, allow up to a minute for it to take effect before checking.
{% endhint %}
