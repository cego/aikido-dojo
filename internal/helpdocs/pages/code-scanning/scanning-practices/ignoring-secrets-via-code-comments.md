# Ignoring Secrets via Code Comments

While you can use the Aikido UI to snooze or ignore secrets, you may also mark secrets as safe via code comments.

Aikido utilizes a modified version of Gitleaks under the hood, so you may mark secrets as safe by adding the string "gitleaks:allow" to the line with the secret.

An example for javascript:

`var a = "live_cdrBarsVQi4EGFRJi" // gitleaks:allow`

An example for python:

`a = "live_cdrBarsVQi4EGFRJi" # gitleaks:allow`
