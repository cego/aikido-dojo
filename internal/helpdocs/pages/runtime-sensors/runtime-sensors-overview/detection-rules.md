# Detection Rules

A detection rule describes one kind of suspicious behavior, for example a container escape or running a binary that isn't part of the container image. Aikido provides the rules. You can't create, edit, or delete them, but you can turn each one on or off.

## Prerequisites

* The admin role, to turn rules on or off. Other roles can view the list.

## View detection rules

Go to **Runtime Sensors** → **Detection Rules** to see every rule with:

| Column         | Description                                                       |
| -------------- | ----------------------------------------------------------------- |
| **Severity**   | Critical, High, Medium, or Low.                                   |
| **Rule**       | The rule name and a short description of the behavior it detects. |
| **Exceptions** | The number of rule exceptions that apply to the rule.             |
| **Enabled**    | A switch that turns the rule on or off.                           |

Use the search box to find a rule by name or description. Filter by severity, status (enabled or disabled), or MITRE ATT\&CK tactic. Click a rule to see its full description and the condition it evaluates.

## Turn a rule on or off

Click the switch in the **Enabled** column. The change applies to your whole workspace, across all clusters and VMs. While a rule is off, Aikido doesn't store new detections for it. Detections already recorded stay in **Detections**.

{% hint style="info" %}
To silence a rule for specific activity only, for example one namespace, add a [rule exception](/runtime-sensors/runtime-sensors-overview/rule-exceptions.md) instead of turning the rule off. Exceptions apply across your whole workspace.
{% endhint %}

## Related

{% content-ref url="/pages/KZiqfv3WtFIfhwAFAdAq" %}
[Runtime Sensors Overview](/runtime-sensors/runtime-sensors-overview.md)
{% endcontent-ref %}

{% content-ref url="/pages/WqYU262B6J9GQcDOBAyW" %}
[Rule Exceptions](/runtime-sensors/runtime-sensors-overview/rule-exceptions.md)
{% endcontent-ref %}
