# VS Code - Extension Keeps Disconnecting

This usually means your key isn’t persisting, check whether VSCode’s storage location is retained.

You can run VS Code in verbose mode `code --verbose` and look for log lines mentioning SecretStorage or globalState for errors and other information.

See additional details below for your OS.

* Linux
  * Encrypted mode: Uses VSCode SecretStorage → KeyStorageLinux → OS-level library (e.g. libsecret), typically stored in the system keyring.
    * For libsecret you can use `secret-tool search application code`
  * Unencrypted mode: Uses globalState, stored at \~/.config/Code/User/globalStorage/state.vscdb. If keys aren’t persisting, verify this path is retained across sessions.
    * You read the database to view existing secrets:\
      `sqlite3 ~/.config/Code/User/globalStorage/state.vscdb "SELECT * FROM ItemTable;"`
