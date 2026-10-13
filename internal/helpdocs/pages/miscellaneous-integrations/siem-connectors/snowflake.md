# Snowflake

{% hint style="info" %}
This connector is currently available by request. [Reach out to support](https://www.aikido.dev/contact) to have it enabled for your workspace.
{% endhint %}

Send your Aikido issues to Snowflake, so your security findings are queryable alongside the rest of your data warehouse.

This is useful if your team wants to run custom analytics on Aikido findings, join them with other business data, or feed them into existing Snowflake-based dashboards and alerting.

### Requirements

Before connecting, get these from Snowflake:

* **Account identifier, Database, Warehouse, Schema, Role, and User** - for a dedicated service account with permission to create tables, stages, and file formats.
* **Table name** - the target table. If it doesn't exist yet, Aikido creates it automatically the first time it writes data.
* **Stage name** - used only when the table is first created.
* **Authentication** - either an RSA **Private Key** (recommended) or a **Password** for the service account.

### Connecting Snowflake

{% stepper %}
{% step %}

#### Set up a Snowflake service account

We recommend a dedicated service account rather than a personal login. Which kind of user you create depends on your authentication method.

**Key pair authentication (recommended)**

Follow [Snowflake's key pair authentication guide](https://docs.snowflake.com/en/user-guide/key-pair-auth) to generate an RSA key pair, then create a service user with the public key:

```sql
CREATE USER aikido_service_account
    COMMENT = 'Service account for Aikido Snowflake integration'
    TYPE = 'service'
    RSA_PUBLIC_KEY = 'Your RSA Public Key';
```

Keep the private key (raw PEM or base64) somewhere safe, you'll paste it into Aikido in a later step.

**Password authentication**

{% hint style="warning" %}
Snowflake service users (`TYPE = 'service'`) don't support password login. If you want to authenticate with a password, create a standard user instead of a service user.
{% endhint %}

```sql
CREATE USER aikido_service_account
    COMMENT = 'Service account for Aikido Snowflake integration'
    PASSWORD = 'a strong, generated password'
    MUST_CHANGE_PASSWORD = FALSE;
```

Key pair authentication is recommended over password authentication where possible.

**Create a role and grant permissions**

```sql
CREATE ROLE aikido_service_account_role;

GRANT ROLE aikido_service_account_role TO USER aikido_service_account;

GRANT USAGE ON WAREHOUSE {your_warehouse} TO ROLE aikido_service_account_role;
GRANT USAGE ON DATABASE {your_database} TO ROLE aikido_service_account_role;
GRANT USAGE ON SCHEMA {your_database}.{your_schema} TO ROLE aikido_service_account_role;
GRANT CREATE TABLE ON SCHEMA {your_database}.{your_schema} TO ROLE aikido_service_account_role;
GRANT CREATE STAGE ON SCHEMA {your_database}.{your_schema} TO ROLE aikido_service_account_role;
GRANT CREATE FILE FORMAT ON SCHEMA {your_database}.{your_schema} TO ROLE aikido_service_account_role;
```

This covers what the connector actually needs: query the warehouse, and create the table, stage, and file format on first write. You don't need `MONITOR` grants or `USAGE ON FUTURE SCHEMAS`, those give broader access than a single-table connector requires. Only add them if you separately want this role to have visibility into warehouse or database activity, or access to schemas created later.

If you're pointing Aikido at a table you created yourself, rather than letting Aikido create one, also grant write access to that specific table:

```sql
GRANT INSERT ON TABLE {your_database}.{your_schema}.{your_table} TO ROLE aikido_service_account_role;
```

{% endstep %}

{% step %}

#### Connect a SIEM

Go to **Settings > Integrations** in the Aikido app and add a new SIEM integration. In the **Connect a SIEM** dialog, select **Snowflake** and fill in:

**Settings**

* **Database** - the Snowflake database to connect to.
* **Warehouse** - the virtual warehouse used to execute queries and load data.
* **Schema** - the schema within the database where the target table lives.
* **User** - the username for the service account.
* **Account** - your Snowflake account identifier. Snowflake's current standard format is `organization_name-account_name` (run `SELECT CURRENT_ORGANIZATION_NAME(), CURRENT_ACCOUNT_NAME();` in Snowflake to find both). Older account-locator formats may also work depending on your account, see [Snowflake's account identifier guide](https://docs.snowflake.com/en/user-guide/admin-account-identifier) if you're unsure which applies.
* **Table** - the table data will be written to. If it doesn't exist, Aikido creates it automatically.
* **Case Insensitivity** - if checked, column names are converted to uppercase during table creation (Snowflake's default behavior). If unchecked, original casing is preserved. Data loading itself always matches columns case-insensitively.
* **Stage** - the Snowflake stage used to load data, only relevant when the table is first created.
* **Role** - the role the service account uses.
* **Authentication Type** - choose **private key** (recommended) or **password**.

**Options**

* **Batch Config** - controls how data is batched before writing:
  * **Batch Data Size** - maximum total size in bytes of a batch before it's sent, regardless of record count or time elapsed.
  * **Batch Record Count** - maximum number of records to accumulate before a batch is sent automatically.
  * **Maximum Flush Interval** - maximum time in seconds to hold records before sending, regardless of size or count.

**Authentication**

Fill in whichever field matches your chosen **Authentication Type**:

* **Private Key** - the RSA private key for the service account, in raw PEM or base64 format.
* **Password** - the service account's password, only if using password authentication (less secure than a private key).

Click **Test connection** to confirm Aikido can reach Snowflake, then **Save changes**.
{% endstep %}

{% step %}

#### Configure which issues sync

On the Snowflake integration page, under **Issue input settings**, configure which Aikido issues are synced:

* **Issue status** - for example, sync open issues only, or include resolved and ignored issues too.
* **Minimum severity** - the lowest severity level to sync.
* **Issue type** - limit to specific issue types, or sync all of them.
* **Language** - limit to specific languages, or sync all of them.
* **Sync frequency** - how often Aikido pushes matching issues to Snowflake (for example, every hour).

Click **Save Settings**.
{% endstep %}
{% endstepper %}

### How data is loaded

The first time data is written to a table that doesn't exist yet, Aikido stages the data and uses `COPY INTO` so Snowflake can infer a schema and create the table automatically. Once the table exists, later writes switch to direct `INSERT` statements, which are more cost-effective for ongoing small batches.

### Verifying events are arriving

In Snowflake, query your target table directly. Wait for your configured sync frequency to pass, or trigger a test issue in Aikido, and confirm new rows appear.

### Troubleshooting

* **Connection or timeout errors** - double check the account identifier format, and confirm your warehouse starts quickly (connections time out after 30 seconds).
* **Authentication errors** - for key pair auth, confirm the private key matches the public key registered on the service account, and try both PEM and base64 formats. For password auth, confirm the credentials are correct and password authentication is enabled.
* **Missing data** - verify the database, schema, table, and stage names are all correct.
* **Schema mismatches on insert** - once a table exists, incoming data must match its schema. If your data shape changes over time, inserts can start failing.
* **Table not created** - table creation requires stage access; check the role has `CREATE STAGE` and `CREATE FILE FORMAT` permissions on the schema.

### Limitations

* Only JSON data is supported.
* Explicit schema definition isn't supported when Aikido creates the table; it infers the schema from the data. You can also pre-create the table with your own schema and point Aikido at it, just make sure the service role has `INSERT` on that table (see the grants above).
* Direct inserts require the existing table schema to match incoming data, inconsistent data can cause insert failures.

### Managing the integration

From **Settings > Integrations > Snowflake** in Aikido, you can:

* **Manage Integration** - update connection details, batching, or authentication.
* **Remove Integration** - disconnect Snowflake. This stops all future syncs; it doesn't delete rows already written to Snowflake.
