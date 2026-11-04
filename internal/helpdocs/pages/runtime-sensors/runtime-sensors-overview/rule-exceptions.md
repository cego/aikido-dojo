# Rule Exceptions

A rule exception tells Aikido to ignore matching activity for the rules you pick. Use exceptions to suppress known false positives, such as a deployment tool that legitimately runs a new binary in a container.

{% hint style="warning" %}
Rule exceptions apply to your whole workspace. They have no cluster or VM scope, so a condition such as a namespace or pod name suppresses matching activity on every cluster and VM that has the same value. Use specific conditions to avoid hiding activity you want to see.
{% endhint %}

## Prerequisites

* The admin role, to add, edit, or delete rule exceptions.

## View rule exceptions

Go to **Runtime Sensors** → **Rule exceptions** to see each exception with its name, its conditions, and the number of rules it applies to. Click the menu on a row to edit or delete the exception.

## Add a rule exception

{% stepper %}
{% step %}

### Open the form

Go to **Runtime Sensors** → **Rule exceptions** and click **Add Rule Exception**. You can also open a detection and click **Actions** → **Add Rule Exception**. Aikido then fills in the name and conditions from the detection, which you can change.
{% endstep %}

{% step %}

### Name the exception

Enter a name that explains why the exception exists.
{% endstep %}

{% step %}

### Define the conditions

Select a field, an operator, and a value. Click **Add Condition** to add more.

All conditions in one exception must match for the exception to apply. To match any of several conditions, create a separate exception for each one.
{% endstep %}

{% step %}

### Pick the rules

In **Rules to Apply to**, search for and select one or more rules. The exception applies to each selected rule.
{% endstep %}

{% step %}

### Save the exception

Click **Save**.
{% endstep %}
{% endstepper %}

## Fields

The field names match the names in the **Fields** list of a detection's event, so you can copy values from there.

| Field                        | Matches              |
| ---------------------------- | -------------------- |
| `proc.name`                  | Process name         |
| `user.name`                  | User name            |
| `k8s.ns.name`                | Kubernetes namespace |
| `k8s.pod.name`               | Kubernetes pod name  |
| `container.image.repository` | Container image      |
| `fd.name`                    | File name            |
| `fd.directory`               | File directory       |

## Operators

| Operator     | Matches when the field                                                      |
| ------------ | --------------------------------------------------------------------------- |
| `=`          | Equals the value.                                                           |
| `!=`         | Doesn't equal the value.                                                    |
| `startswith` | Starts with the value.                                                      |
| `contains`   | Contains the value.                                                         |
| `glob`       | Matches a pattern with a wildcard, e.g. `web-*`.                            |
| `in`         | Equals any value in a comma-separated list, e.g. `kube-system, monitoring`. |

## Limits

* A value can be at most 250 characters.
* A value can't contain parentheses or line breaks.
* An exception needs a name, at least one complete condition, and at least one rule.

{% hint style="info" %}
Pod names often end in a generated suffix. When you create an exception from a detection, Aikido uses `glob` with a pattern such as `web-*` for those pods, so the exception keeps matching after a restart.
{% endhint %}

## Troubleshooting

### A rule exception doesn't suppress a detection

Check that the exception lists the rule in **Rules to Apply to**, and that every condition matches the detection. Conditions are combined, so one condition that doesn't match is enough to stop the exception from applying. Open the detection and compare the namespace, pod, process, and file values with your conditions.

## Related

{% content-ref url="/pages/KZiqfv3WtFIfhwAFAdAq" %}
[Runtime Sensors Overview](/runtime-sensors/runtime-sensors-overview.md)
{% endcontent-ref %}

{% content-ref url="/pages/LeGgY6ZX0hQgT1470phu" %}
[Detection Rules](/runtime-sensors/runtime-sensors-overview/detection-rules.md)
{% endcontent-ref %}
