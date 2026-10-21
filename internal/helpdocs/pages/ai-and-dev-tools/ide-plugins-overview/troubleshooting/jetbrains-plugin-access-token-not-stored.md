# JetBrains Plugin - Access Token Not Stored

If you are experiencing issues with storing the Aikido access token in your JetBrains IDE, navigate to **Settings > Appearance & Behavior > System > Settings > Passwords**. This section shows which mechanism your IDE uses to store passwords and tokens.

Make sure the option "**Do not save, forget passwords after restart**" is **NOT** enabled.

If "In native Keychain" or "In KeePass" is selected, verify that no mechanisms are interfering with JetBrains accessing your system's native keychain or KeePass on your machine.

*Example configuration*

More info can be found in [JetBrains documentation.](https://www.jetbrains.com/help/idea/reference-ide-settings-password-safe.html#0)
