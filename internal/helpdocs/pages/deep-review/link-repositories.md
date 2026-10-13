# Linked Repositories

Link related repositories when Deep Review needs context that lives outside the pull request's repository. This is useful when your services, shared libraries, or APIs are split across repositories.

{% stepper %}
{% step %}

### Open the repository

Open [**Repositories**](https://app.aikido.dev/repositories) and select the repository you want to configure.
{% endstep %}

{% step %}

### Open the Deep Review settings

Open **Repository Settings**, find the **Deep PR Review** card, and click **"Edit"** next to **Linked Repositories**.
{% endstep %}

{% step %}

### Add a repository

In the **Linked Repositories** dialog, click **"Link Repository"**.

<div data-with-frame="true">[image: Linked Repositories dialog listing one linked repository and a Link Repository button.]</div>
{% endstep %}

{% step %}

### Select the repository and add context

Search for the repository Deep Review should use as additional context. Optionally use the **Instruction** field to explain how the repository is related, for example that it holds the backend API your frontend calls. Then click **"Link Repository"**.

<div data-with-frame="true">[image: Link Repository dialog with a repository search field and an optional instruction field.]</div>
{% endstep %}
{% endstepper %}

Each repository can link the following number of repositories, depending on your Aikido plan:

| Plan       | Linked repositories per repository |
| ---------- | ---------------------------------- |
| Developer  | 0                                  |
| Basic      | 1                                  |
| Pro        | 1                                  |
| Advanced   | 5                                  |
| Enterprise | 10                                 |
