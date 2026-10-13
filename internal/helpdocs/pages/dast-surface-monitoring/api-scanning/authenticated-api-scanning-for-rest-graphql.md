# Authenticated API Scanning for REST/GraphQL

This guide will walk you through the steps to set up authenticated domain scanning in Aikido, ensuring thorough and secure assessments.

> API Scanning is only available on Pro and Advanced Plans

## Use Cases <a href="#use-cases" id="use-cases"></a>

* Ensure comprehensive security assessments for protected areas of your website.
* Identify vulnerabilities in authenticated sections of your APIs.

## Setting up authentication on a domain <a href="#setting-up-authentication-on-a-domain" id="setting-up-authentication-on-a-domain"></a>

**Step 1:** Go to the [**Domains Overview**](https://app.aikido.dev/settings/domains) and open the action menu for a REST/ GraphQL API domain of your choice by clicking the triple dots. Select **Authenticate Domain.**

[image: Domain action menu with options to scan, configure, authenticate, or delete a domain.]

**Step 2:** Select your preferred option to authenticate.

[image: Authentication setup screen with multiple login method options and credential input fields.]

> All types of scan credentials are securely stored using PKCS1 encryption

## Authentication Options <a href="#authentication-options" id="authentication-options"></a>

**Login via Form**

Fill in the URL and email/password for the domain authentication. Click **Test** to let Aikido check whether it can access the domain with those credentials.

You can also enter the one time password information if two factor authentication is enabled. For [more information about the OTP URL and how it works, see the documentation](/dast-surface-monitoring/miscellaneous/using-2fa-in-front-end-and-api-scans.md).

{% hint style="success" %}

* Microsoft / Google SSO is currently **not** supported. As a workaround, you can manually authenticate and pass a valid session using the Cookie header via custom headers.
* Is your case not supported? Let us know via the chat and we will look into it!
  {% endhint %}

[image: Domain authentication setup screen for form-based login credentials configuration.]

Aikido will attempt to submit the form using the following rules:

1. Find a visible button or input field with `type=submit` while ignoring popular OAuth options like Google and Facebook
2. Find button based on the label or text. Looks for text equal or similar to: login, log in, submit, sign in, .. It does so in multiple languages.
3. Find button based on set of HTML ID's, for example `id=form-submit`
4. Find first visible button on page

**Custom Headers**

If your endpoints accepts a fixed key, cookie or token which should not change after creation, you can add it as a custom header via this option.

Use-cases:

* Cookies: Set the `Cookie` header.

  ```
  Cookie: sessionId=38afes7a8
  ```
* JWT Bearer token: Set the `Authorization` header

  ```
  Authorization: Bearer <token>
  ```

**OAuth Client Credentials**

This option can be used when you want to bypass MFA. Aikido makes a request to the provided login URL which follows the OAuth spec for a Client Credentials flow. This means that we'll make a POST request to the configured login url, with `grant_type` set to `client_credentials` and a basic authorization header containing the client\_id and client\_secret as the username and password respectively.

**OAuth Password Grant**

Choose this option if your API requires direct user credentials (username and password) using the OAuth Password Grant flow. Aikido sends a POST request to your login URL with `grant_type` set to `password`, including the username and password you provide. Depending on the API, a Client ID and Secret might also be required.

**Basic Auth**

Select this option if your API uses standard Basic Authentication. Provide your username and password, and Aikido will automatically include the necessary `Authorization: Basic ...` header with every request made to your API.

**Login via AI Agent (beta)**

The AI Agent uses an LLM to control a real browser session. It follows the instructions you provide to complete the authentication process, just like a human would. This works well for:

* Form based logins
* Static OTP or one time codes
* Multi step authentication flows
* Custom or non standard login screens

Provide clear step by step instructions, including credentials and any required actions. The agent can handle redirects, dynamic pages, and intermediate steps without needing detailed configuration. In most cases, short and simple instructions are enough.
