# Configuring Registry Scan Scope

Configure wildcard patterns to control which container images are scanned in your connected registries.

Use wildcard patterns to control which container images Aikido scans in your connected registries. This helps you focus on relevant containers and avoid clutter when working with registries containing hundreds of images across multiple projects.

#### How Scan Scope Works

You can configure scan scope in two ways:

* **Include containers matching**: Aikido will exclusively scan containers that match your specified patterns. All other containers will be ignored.
* **Exclude containers matching**: Aikido will scan all containers except those matching your patterns.

{% hint style="info" %}
**Wildcard patterns match against container names only.** Tags, image IDs, and other metadata are not considered when evaluating scope rules.
{% endhint %}

#### Setting Up Scan Scope

1. Navigate to [**Settings** > **Containers**](https://app.aikido.dev/settings/container-image-registry)
2. Click the triple dots on the registry you want to configure, and click '**Edit Registry Configuration**'<br>

   
3. In the **Registry Settings** modal, scroll to the **Registry Scan Scope** section<br>

   
4. Select your scan mode:
   * **Include containers matching** to create an allowlist
   * **Exclude containers matching** to create a blocklist
5. Enter your wildcard pattern (e.g., `backend-*`)
6. Click **Add Pattern** to add additional patterns
7. Click **Save**

After saving, Aikido applies your scan scope rules **during the next scheduled scan**. Only containers matching your criteria will be added to your Aikido workspace and scanned.

#### Wildcard Pattern Examples

Wildcard patterns use the `*` character to match multiple container names:

* `frontend-*` - Matches all containers starting with "frontend-" (e.g., `frontend-web`, `frontend-api`)
* `*-prod` - Matches all containers ending with "-prod"
* `production/*` - Matches all containers in the production project or namespace
* `*test*` - Matches any container with "test" in the name
