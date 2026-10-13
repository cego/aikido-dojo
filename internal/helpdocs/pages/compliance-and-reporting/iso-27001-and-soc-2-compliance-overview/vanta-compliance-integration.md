# Vanta Integration

Vanta is a system to help your company attain SOC2 and ISO27001 compliance. Its system is integrated with 3rd party applications to feed it with evidence that it can monitor and bring together in 1 dashboard for your auditor.

Aikido can be found on Vanta's app marketplace. You can also start the installation via [this link](https://app.aikido.dev/settings/integrations/compliance/vanta).

[image: Vanta app authorization screen for granting Aikido Security account access.]

The integration moves data in a single direction. Aikido will **sync open dependency issue to Vanta on an hourly basis.** Aikido will copy the correct risk, so Vanta can assign an SLA to each issue. Snoozed issues and (auto-)ignored issues will not be synced to Vanta.

You may verify the integration is set up correctly by logging into Vanta, then clicking 'monitors' and searching for Aikido. Because Aikido is an aggregator of many risk types (dependencies, secrets, SAST,..) it is possible Aikido's findings may overlap with other security integrations. In such cases it is safe to deactivate monitoring for other tools in Vanta to avoid the duplication.

***
