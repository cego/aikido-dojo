# Send Alerts to Multiple Slack Channels

## Introduction <a href="#introduction" id="introduction"></a>

Streamline your security response with Aikido by setting up tailored alerts to be sent to designated Slack channels. This guide will help you ensure that your team receives timely and relevant security notifications.

## Setting Up Slack Alerts <a href="#setting-up-slack-alerts" id="setting-up-slack-alerts"></a>

_Prerequisites: Ensure your Slack integration is active. You can find how to set up the Slack Integration_ [_here_](https://help.aikido.dev/getting-started/chat-and-alerts/slack-notifications)_._

#### Step 1: Create an Alert <a href="#step-1-create-an-alert" id="step-1-create-an-alert"></a>

* Navigate to the [Slack integration page](https://app.aikido.dev/settings/integrations/notifications/slack?regionRedirect=1) and click **Create alert**.

#### Step 2: Opt for Team-Based Notifications <a href="#step-2-opt-for-team-based-notifications" id="step-2-opt-for-team-based-notifications"></a>

* Under **Team to notify**, select **All Teams** or one team. Link repositories, containers, and other resources to that team first, so the alert only covers what that team owns.

[image: Configure Slack alerts for all users or specific teams with custom options.]

#### Step 3: Map Channels to Teams <a href="#step-3-map-channels-to-teams" id="step-3-map-channels-to-teams"></a>

Create a separate alert for each notification. Instant alerts and the weekly digest can be limited to one team. SLA alerts apply to the whole workspace.

See [Slack Notifications](https://help.aikido.dev/getting-started/chat-and-alerts/slack-notifications) for every option on each alert type, including issue type, severity, and SLA timing.

#### Step 4: Save Your Configuration <a href="#step-4-save-your-configuration" id="step-4-save-your-configuration"></a>

* Confirm your choices and click 'Save Alert' to activate the alerting system.

## Implementation Tips <a href="#implementation-tips" id="implementation-tips"></a>

* Regularly review and correct team settings in Aikido to maintain alert accuracy.
* Confirm that the Slack channels correspond to the correct teams to ensure proper communication

***
