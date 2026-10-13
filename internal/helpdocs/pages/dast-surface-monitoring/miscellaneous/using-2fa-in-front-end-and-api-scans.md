# Using 2FA in Front-End and API Scans

If your application login requires a TOTP code, Aikido can log in during scans as long as you provide the otpauth URL for that 2FA token.

This allows Aikido to generate the same 6 digit codes your authenticator app would.

## What is an otpauth URL

An otpauth URL contains the TOTP secret and configuration needed to produce time based codes.

It looks like this:

```
otpauth://totp/MyApp:user@example.com?secret=JBSWY3DPEHPK3PXP&issuer=MyApp
```

## How to get your otpauth URL

You can get the URL in one of these ways:

### 1. From the 2FA setup screen

Look for:

* Show secret
* Manual setup
* Enter key instead

Some providers show the otpauth URL directly. Others show only the Base32 secret.

If you only have the Base32 secret, you can build the URL yourself:

```
otpauth://totp/<Label>?secret=<BASE32SECRET>&issuer=<YourApp>
```

### 2. By scanning the QR code

If the setup shows a QR code:

* Use any QR reader that reveals the text content
* The scanned result will be the full otpauth URL

### 3. From your IdP or admin panel

Some IdPs expose the TOTP seed for service accounts.

Look for fields like:

* OTP URI
* TOTP configuration
* MFA seed

If your login uses SMS, push notifications, or hardware keys, you will not get an otpauth URL. Those are currently not supported by Aikido.
