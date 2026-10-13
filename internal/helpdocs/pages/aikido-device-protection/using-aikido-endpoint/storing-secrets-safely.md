# Storing Secrets Safely

[Secret Scanning](/aikido-device-protection/using-aikido-endpoint/secret-scanning.md) shows you where credentials are sitting in plaintext on your device. Cleaning up the file is half the fix. This page is the other half: keeping the secret out of those files in the first place, by storing it in a secrets manager and fetching it when you need it.

## How secrets end up in plaintext

There are three routes, and Aikido finds all three:

* **Credential files** - You wrote the value into a config file once and forgot about it, like `~/.aws/credentials` or `.npmrc`.
* **The command you typed** - It lands in `.zsh_history` or `.bash_history`, and while it runs the value is visible to every other process on the machine through `ps`.
* **AI agent sessions** - If an AI coding agent ran the command, the value is in that agent's session transcript too.

Here's a real pattern Aikido catches:

```bash
AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCY aws-vault add default
```

aws-vault puts that key in your keychain, which is exactly where it belongs. The problem is the trip. By the time the key reached the keychain it was already in your shell history, in the process list, and in the agent transcript if an agent ran it. Storing a secret well doesn't help if you leak it on the way in.

## Do you need a long-lived credential at all?

The best secret is the one that doesn't exist. Short-lived credentials expire on their own, so there's nothing on disk to leak and nothing to rotate. Check whether your provider offers them before you reach for a vault:

* **AWS** - Use IAM Identity Center. Run [`aws configure sso`](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-sso.html) once, then `aws sso login` when your session expires.
* **Google Cloud** - Run `gcloud auth application-default login` instead of downloading a service account key.
* **GitHub** - Run `gh auth login`, which stores the token in your OS keychain, instead of pasting a personal access token into `.npmrc` or `.git-credentials`.

## Storing a secret in a secrets manager

When you do need a long-lived credential, put it in a manager and read it back at run time. Each tab below covers the two commands you need. Follow the linked guide for the rest.

{% hint style="warning" %}
Fetch the value into an environment variable or pipe it to stdin. Passing it as an argument, like `mytool --key "$(op read ...)"`, puts it straight back into your shell history and the process list.
{% endhint %}

{% tabs %}
{% tab title="macOS Keychain" %}
Built into macOS, no extra tooling. Leaving the value off the `-w` flag makes `security` prompt for it, so the secret never appears in your command.

```bash
security add-generic-password -a "$USER" -s stripe-api-key -w
```

Read it back when you need it:

```bash
export STRIPE_API_KEY=$(security find-generic-password -a "$USER" -s stripe-api-key -w)
```

See [Apple's Keychain Access guide](https://support.apple.com/guide/keychain-access/welcome/mac) for managing items in the UI.
{% endtab %}

{% tab title="1Password" %}
Create the item in the 1Password app rather than through CLI arguments. Each field then has a secret reference in the form `op://vault/item/field`, which contains no secret value itself.

```bash
export STRIPE_API_KEY=$(op read "op://Private/Stripe/credential")
```

For a whole set of variables, put the references in an env file and let `op run` resolve them. The values exist only inside that subprocess, for as long as the command runs:

```bash
op run --env-file=./aws.env -- terraform plan
```

For AWS specifically, `op plugin init aws` lets the AWS CLI authenticate through 1Password with biometrics, so nothing is stored in your home directory at all.

See [Secret reference syntax](https://www.1password.dev/cli/secret-reference-syntax/), [Load secrets into the environment](https://www.1password.dev/cli/secrets-environment-variables/), and [the AWS shell plugin](https://www.1password.dev/cli/shell-plugins/aws/).
{% endtab %}

{% tab title="Bitwarden" %}
Add the item in the Bitwarden app, then unlock your vault and read the value back:

```bash
export BW_SESSION=$(bw unlock --raw)
export STRIPE_API_KEY=$(bw get password "Stripe API key")
```

For secrets your team or your machines share, Bitwarden Secrets Manager injects them directly into a command with `bws run -- terraform plan`.

{% hint style="warning" %}
`BWS_ACCESS_TOKEN` is itself a secret. Keep it in your OS keychain and read it at run time, don't paste it into `.zshrc`.
{% endhint %}

See the [Password Manager CLI](https://bitwarden.com/help/cli/) and [Secrets Manager CLI](https://bitwarden.com/help/secrets-manager-cli/) guides.
{% endtab %}

{% tab title="aws-vault" %}
`aws-vault` stores AWS keys in your OS keychain and hands out temporary credentials instead. Run it with no value and it prompts for both parts, so nothing enters your shell history:

```bash
aws-vault add prod
```

Then run anything through it:

```bash
aws-vault exec prod -- aws s3 ls
```

To stop typing `aws-vault exec` every time, point the AWS CLI at it in `~/.aws/config` and leave `~/.aws/credentials` empty:

```ini
[profile prod]
credential_process = aws-vault export --format=json prod
```

{% hint style="warning" %}
Once `credential_process` is set for a profile, don't also run `aws-vault exec` by hand for it. The AWS SDK invokes aws-vault for you.
{% endhint %}

See the [aws-vault usage guide](https://github.com/99designs/aws-vault/blob/master/USAGE.md).
{% endtab %}
{% endtabs %}

## After you move a secret

Moving the credential into a manager doesn't undo the exposure. Finish the job:

1. Rotate the old value at its source. Treat it as compromised, because you don't know who read the file or the history.
2. Remove the plaintext copy. Use the **Secrets** tab in the desktop client for shell history and AI agent sessions, and edit credential files yourself.
3. Confirm the finding is gone from the **Secrets** tab.

{% content-ref url="/pages/Uk7VvDDNYtvQt1fdvQYn" %}
[Secret Scanning](/aikido-device-protection/using-aikido-endpoint/secret-scanning.md)
{% endcontent-ref %}

For secrets Aikido found in your code rather than on a device, see [Resolving Malware and Leaked Secret Issues](/getting-started/core-functionalities/resolving-malware-and-leaked-secret-issues.md).
