# Manually Adjust Issue Severity

### Introduction <a href="#introduction" id="introduction"></a>

When Aikido finds vulnerabilities in your code repos, cloud environments or public facing domains, our scoring engine gives the vulnerability a severity from 'low' to 'critical'. Our scoring engine takes into account a whole set of rules to assign this severity, but the most important one would be the urgency to fix.

If for some reason, you believe a vulnerability has been given a wrong severity, either to low or to high, Aikido gives you the opportunity to adjust the severity manually so it ends up higher or lower on your list of things to fix.

You can either adjust a **single issue's severity**, or the **severity of a whole group** of issues, in which case they will all get the same severity, regardless of their previous severity.

### Adjust severity of a single issue <a href="#adjust-severity-of-a-single-issue" id="adjust-severity-of-a-single-issue"></a>

A single issue's severity can be adjusted via the issue's action menu found in the sidebar, as shown in the image below.

{% hint style="warning" %}
Severity adjustments are not available for **pentest issues** for compliance reasons. If you believe the priority of a pentest finding is incorrect, please flag it to us so we can discuss it together.
{% endhint %}

### Adjust severity of a whole issue group <a href="#adjust-severity-of-a-whole-issue-group" id="adjust-severity-of-a-whole-issue-group"></a>

To adjust the severity of a whole issue group, you can click on "Adjust severity" on the issue group's action menu in a row in any table.

[image: Security vulnerabilities dashboard showing critical issues, severity, affected systems, and assigned team members.]

When adjusting the severity, you need to provide the new severity the vulnerability as well as a reason why you think the severity should be adjusted.

If you decide to lower or increase the severity of an issue group, Aikido's scoring engine will not apply that adjusted severity to any newly discovered issue in that group. This is because we believe that you should at that moment evaluate this new finding after which you can again adjust the severity.
