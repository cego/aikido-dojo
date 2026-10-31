# Runtime Sensors Overview

Runtime Sensors monitor your Kubernetes clusters and virtual machines (VMs) while they run. They use extended Berkeley Packet Filter (eBPF) to watch workload activity from the kernel and detect threats as they happen.

## How it works

Runtime Sensors run as part of the Aikido agent on each workload. The sensor observes activity on the host without changing your application code, and evaluates it against a set of detection rules. When activity matches a rule, Aikido records a detection.

Aikido ships a set of [detection rules](/runtime-sensors/runtime-sensors-overview/detection-rules.md) that cover behaviors such as container escapes and unexpected binaries running in a container. You can turn individual rules on or off for your workspace. When activity matches an enabled rule, Aikido creates a detection that you can investigate. It summarizes the severity, the affected workload, and the events behind it. If a rule flags activity you know is safe, add a [rule exception](/runtime-sensors/runtime-sensors-overview/rule-exceptions.md) to suppress it.

Threat detection is off by default for Kubernetes clusters. You turn it on per cluster. To get started, see [Enable Runtime Sensors](/runtime-sensors/runtime-sensors-overview/enable-runtime-sensors.md).

## Review detections

Go to **Runtime Sensors** → **Detections** to see activity that matched a detection rule. Aikido groups matching events by source, rule, and workload over a 24-hour window, so one detection can stand for many events.

Each detection shows:

* **Severity** of the rule that matched.
* **Rule** name.
* **Source**, the cluster or VM where the activity happened.
* **Workload**, for example the pod or the VM.
* **Last detected** time.

Use the search box to find a rule by name. Filter by cluster, VM, or severity. You can filter by a cluster or a VM, but not both at once.

Click a detection to investigate it. The details show:

* The rule name, the number of events found, and when the activity was last detected.
* The severity, plus the MITRE ATT\&CK technique and tactic the rule maps to (e.g. T1620 and defense evasion), with links to attack.mitre.org.
* A short summary of what the rule detects.
* Recent events, grouped by source and workload. Each event shows the fields that triggered it, such as the process command and user name, and a **View Details** link to the full event.

**View Details** opens the full event with the time it happened, the raw **Output** line the sensor produced, and a JSON list of **Fields**. The fields include the process, user, container, and Kubernetes namespace and pod. Use them to decide whether the activity is expected.

Use the arrows at the bottom of the details to move between detections without closing them.

From the details, click **Actions** → **Add Rule Exception** to suppress a false positive. See [Rule Exceptions](/runtime-sensors/runtime-sensors-overview/rule-exceptions.md).

## Related

{% content-ref url="/pages/Ob4FqWj79LARMNOkEWWM" %}
[Enable Runtime Sensors](/runtime-sensors/runtime-sensors-overview/enable-runtime-sensors.md)
{% endcontent-ref %}

{% content-ref url="/pages/LeGgY6ZX0hQgT1470phu" %}
[Detection Rules](/runtime-sensors/runtime-sensors-overview/detection-rules.md)
{% endcontent-ref %}

{% content-ref url="/pages/WqYU262B6J9GQcDOBAyW" %}
[Rule Exceptions](/runtime-sensors/runtime-sensors-overview/rule-exceptions.md)
{% endcontent-ref %}
