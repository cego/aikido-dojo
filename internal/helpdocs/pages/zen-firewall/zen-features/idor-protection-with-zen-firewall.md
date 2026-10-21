# IDOR Protection with Zen Firewall

IDOR stands for Insecure Direct Object Reference. It is a class of access control bugs where one account can access another account’s data because a database query is not properly scoped.

In multi tenant SaaS applications, data is usually separated with a column such as tenant\_id, account\_id, or organization\_id. If a query forgets to filter on that column, or uses the wrong value, data can leak across accounts.

IDOR Protection on Zen Firewall analyzes SQL queries at runtime and enforces tenant scoping automatically. If a query is unsafe, Zen throws an error immediately.

This helps you catch multi tenant isolation bugs during development and testing, not after a production incident.

For a deeper explanation of the vulnerability itself, see the blog post: [IDOR Vulnerabilities Explained](https://www.aikido.dev/blog/idor-vulnerability-explained)

## Supported databases

{% tabs %}
{% tab title="NodeJS" %}
Currently supported drivers:

* MySQL via mysql and mysql2
* PostgreSQL via pg
* SQLite (via better-sqlite3 and node:sqlite packages)

Any ORM or query builder that uses these drivers under the hood is supported, for example Drizzle, Knex, Sequelize, and TypeORM.

ORMs that use their own query engine, such as Prisma, are not supported unless configured to use a supported driver adapter.

{% hint style="info" %}
If you use ESM, review the [ESM caveats](https://github.com/AikidoSec/firewall-node/blob/main/docs/esm.md). Queries inside uninstrumented ESM sub dependencies cannot be checked by Zen.
{% endhint %}
{% endtab %}
{% endtabs %}

## Limitations

* [NodeJS](https://github.com/AikidoSec/firewall-node/blob/main/docs/idor-protection.md#limitations)

## What Zen protects against

When enabled, Zen validates that every relevant query:

* Filters on the correct tenant column for SELECT, UPDATE, and DELETE
* Uses the correct tenant ID value
* Includes the tenant column in every INSERT
* Sets the tenant column to the correct tenant ID in INSERT<br>

Examples of issues Zen blocks:

* A SELECT without a tenant filter, allowing one account to read another account’s data
* An UPDATE or DELETE without tenant scoping
* An INSERT that omits the tenant column
* An INSERT that sets the wrong tenant ID

{% hint style="warning" %}
IDOR protection always throws an Error on violations, regardless of block or detect mode. A missing tenant filter is a developer bug, not an external attack.
{% endhint %}

## How it works

Zen instruments supported database drivers and inspects SQL queries at runtime. For each request, you set the active tenant ID. Zen then verifies that all relevant queries:

1. Reference the configured tenant column.
2. Use the same tenant ID that was set for the request.

If a query does not meet these conditions, Zen throws an error and prevents it from executing.

## Setup

{% stepper %}
{% step %}
**Enable IDOR protection at startup**

Enable IDOR protection when your application starts

{% tabs %}
{% tab title="NodeJS" %}

```js
import Zen from "@aikidosec/firewall";

Zen.enableIdorProtection({
  tenantColumnName: "tenant_id",
  excludedTables: ["users"],
});
```

Configuration options:

* `tenantColumnName`

  The column that identifies the tenant in your tables, for example tenant\_id, account\_id, or organization\_id.
* `excludedTables`

  Tables that should not be checked because they are not scoped to a single tenant. For example, a shared users table that contains users across all tenants.
  {% endtab %}
  {% endtabs %}
  {% endstep %}

{% step %}
**Set the tenant ID per request**

Every request must have a tenant ID when IDOR protection is enabled.

Call setTenantId early in your request lifecycle, typically in middleware after authentication.

{% tabs %}
{% tab title="NodeJS" %}

```js
import Zen from "@aikidosec/firewall";

app.use((req, res, next) => {
  // Retrieve the tenant ID from your authentication layer
  Zen.setTenantId(req.user.organizationId);

  next();
});
```

{% hint style="warning" %}
If setTenantId is not called for a request, Zen will throw an Error when a SQL query is executed.
{% endhint %}
{% endtab %}
{% endtabs %}
{% endstep %}

{% step %}
**Bypass IDOR checks for specific queries**

Some queries are intentionally global, for example cross tenant analytics or internal admin reports.

To skip IDOR checks for a specific block of code, use withoutIdorProtection:

```js
import Zen from "@aikidosec/firewall";

const result = await Zen.withoutIdorProtection(async () => {
  return await db.query(
    "SELECT count(*) FROM agents WHERE status = 'running'"
  );
});
```

Only queries executed inside the callback are excluded from IDOR validation.

Use this carefully and only where cross tenant access is explicitly required.
{% endstep %}
{% endstepper %}

## Troubleshooting

### Missing tenant filter

```
Zen IDOR protection: query on table 'orders' is missing a filter on column 'tenant_id'
```

Your query does not filter on the configured tenant column. This applies to SELECT, UPDATE, and DELETE.

Example:

```
SELECT * FROM orders WHERE status = 'active';
```

Add a condition that filters on tenant\_id.

### Wrong tenant ID value

```
Zen IDOR protection: query on table 'orders' filters 'tenant_id' with value '456' but tenant ID is '123'
```

The query filters on the tenant column, but the value does not match the tenant ID set via `setTenantId`.<br>

Ensure that all queries use the tenant ID from the current request context.

### Missing tenant column in INSERT

```
Zen IDOR protection: INSERT on table 'orders' is missing column 'tenant_id'
```

Every INSERT must include the tenant column.

#### Wrong tenant ID in INSERT

```
Zen IDOR protection: INSERT on table 'orders' sets 'tenant_id' to '456' but tenant ID is '123'
```

The INSERT includes the tenant column, but the value does not match the active tenant ID.

### Missing \`setTenantId\` call

```
Zen IDOR protection: setTenantId() was not called for this request. Every request must have a tenant ID when IDOR protection is enabled.
```

Ensure setTenantId is called for every request before any database queries run.

## Statements that are always allowed

Zen only checks statements that read or modify row data:

* SELECT
* INSERT
* UPDATE
* DELETE

The following are recognized and never trigger IDOR errors:

* DDL statements such as CREATE TABLE, ALTER TABLE, DROP TABLE
* Session commands such as SET, SHOW
* Transaction statements such as BEGIN, COMMIT, ROLLBACK
