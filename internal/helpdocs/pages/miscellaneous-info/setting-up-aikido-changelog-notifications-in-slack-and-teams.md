# Setting Up Aikido Changelog Notifications in Slack & Teams

Stay informed about the latest product updates, bug fixes, and new features by connecting our changelog’s RSS feed to your team’s communication tools. This guide will show you how to configure automatic notifications in both Slack and Microsoft Teams, ensuring your team never misses important updates that could enhance your security posture or workflow efficiency.

## Slack

1. **Add the RSS app to your Slack workspace**:
   * Visit the [Slack App Directory](https://slack.com/apps/A0F81R7U7-rss)
   * Click “Add to Slack” and select the workspace where you want to receive notifications
   * Authorize the app when prompted
2. **Configure the RSS feed**:
   * In your Slack workspace and the channel where you want to receive the notifications, type\
     `/feed subscribe https://aikido.help/changelog.xml`
3. **Verify the connection**:

* Once configured, you should see a confirmation message in your Slack channel
* The next changelog update will automatically appear in your selected channel

## Microsoft Teams

1. **Add the RSS connector to your Teams channel**:
   * Navigate to the Teams channel where you want to receive updates
   * Click the “…” (more options) button next to the channel name
   * Select “Connectors” from the dropdown menu
   * Search for “RSS” and click “Configure” on the RSS connector
2. **Set up the RSS feed**:
   * In the configuration window, enter a name (e.g., “Aikido Product Changelog”)
   * Paste our changelog RSS feed URL: `https://aikido.help/changelog.xml`
   * Select how often you want updates to be delivered
   * Click “Save”
3. **Verify the connection**:
   * Once configured, you should see a confirmation message in your Teams channel
   * The next changelog update will automatically appear in your selected channel
