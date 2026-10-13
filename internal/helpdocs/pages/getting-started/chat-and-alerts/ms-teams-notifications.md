# Microsoft Teams Notifications

## Connecting Aikido to Microsoft Teams <a href="#connecting-aikido-to-ms-teams" id="connecting-aikido-to-ms-teams"></a>

{% hint style="warning" %}
Please note that only **Aikido admin users** can link a new Microsoft Team to your Aikido workspace.

If you are using the **OLD version** of the Microsoft Teams Integration, **please delete all alerts** before setting up the new version. You can setup the new version by removing the integration and then following the steps below.\
\
App requires application-level permissions; access is limited to the Microsoft it's installed in, not the whole org.
{% endhint %}

Follow these steps to connect Aikido to Microsoft Teams:

#### Step 0. Uninstall the MS Teams Integration v1 <a href="#id-1-click-add-integration" id="id-1-click-add-integration"></a>

In case you are still running the old MS Teams integration, first **uninstall** this one from the Aikido UI!

#### Step 1. Install the Aikido app in Microsoft Teams <a href="#id-1-click-add-integration" id="id-1-click-add-integration"></a>

* Open the [Aikido app](https://teams.microsoft.com/l/app/c2baec07-db8a-49de-b066-c0ddc19cc9c0?source=store-copy-link) in the Microsoft Teams app store and click "Add to a Team".
* Select a channel for initial installation. You can configure alerts for all public channels of the selected Team later.

[image: Select the team to which the messages should be sent.]

You can find more information on how to install an app on Microsoft Teams in [this Microsoft support article](https://support.microsoft.com/en-us/office/add-an-app-to-microsoft-teams-b2217706-f7ed-4e64-8e96-c413afd02f77).

#### Step 2. Link the Microsoft Team to an Aikido workspace <a href="#id-2-select-which-ms-teams-team-youd-like-to-connect" id="id-2-select-which-ms-teams-team-youd-like-to-connect"></a>

* The Aikido bot will send a message to the selected Microsoft Team channel. Click "Open Aikido" and log in.
* Complete the app setup by confirming the installation in the modal.

[image: Select a Microsoft Teams group to connect with Aikido integration.]

#### Step 3. Configure the alerts you'd like to receive <a href="#id-3-configure-the-alerts-youd-like-to-receive" id="id-3-configure-the-alerts-youd-like-to-receive"></a>

{% hint style="warning" %}
Due to limitation on Microsoft side, some public channels may not be discovered automatically. You can add them manually by selecting "Public channel missing..." in the channel selection options.
{% endhint %}

Once the connection to Microsoft Teams was successful, you can start adding alerts.

* Choose whether you want to send global or team-based notifications
* Choose for which severity types
* Choose for which issue types

***

## Unlinking Aikido and Microsoft Teams

1. Open the Microsoft Teams app and right click on the team where the Aikido bot is installed.
2. Click **Manage team** and select the **Apps** tab.
3. Search for the Aikido app in the list of installed apps, click on the three dots next to it, and select **Remove**.
4. Confirm the removal in the pop-up dialog. All settings for this team will be deleted automatically.

***

## FAQ

<details>

<summary>Can I link multiple Teams to one Aikido Workspace?</summary>

Yes, this is possible. Just follow the same steps again.

</details>

<details>

<summary>Is it possible to link multiple Aikido Workspaces to one Microsoft Team?</summary>

To link additional workspaces to the same Microsoft Team, click the button "Link another workspace" in the confirmation message sent by Aikido to your selected Teams channel. You will find theses message inside the Microsoft channel where you first installed the Aikido App. Please **do not** try to install the app multiple times.

</details>

<details>

<summary>How can I receive alerts in a private channel?</summary>

This is not easily possible due to restrictions imposed by Microsoft. We recommend creating a private Team with a public channel instead. But we describe a [workaround on this page](/getting-started/chat-and-alerts/send-alerts-to-private-microsoft-teams-channels.md).

</details>

<details>

<summary>I don't receive the "Connect to Aikido" chat message</summary>

Please ensure that no moderation settings are configured for the selected channel that prevent bots from sending messages. Uninstall the app and try the process again using a different channel. If this does not help, please contact the Aikido support team.

</details>

<details>

<summary>Can I use the app in the Government Community Cloud (GCC)?</summary>

This is currently not possible. Please contact us to let us know that you would like to use our Teams app in a GCC environment.

</details>

***
