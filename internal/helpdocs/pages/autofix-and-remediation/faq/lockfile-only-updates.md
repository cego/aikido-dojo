# Lockfile-only updates

AutoFix makes the smallest changes required to update dependencies. This means we perform lockfile-only updates whenever applicable for the following package managers:

* NPM
* Yarn
* PNPM
* Bun
* Poetry
* UV

**The example below uses NPM, but the same logic applies for the other package mangers mentioned above.**

Updates to `package.json` are required in most scenarios, however for transitive dependencies this is not required if the parent dependency includes the transitive dependency with a version range that includes the fixed version.

For example, you project includes axios 1.15.0:

{% code title="package.json" %}

```json
{
	"dependencies": {
		"axios": "1.15.0"
	}
}
```

{% endcode %}

Axios itself depends on `follow-redirects` :

{% code title="package.json (axios)" %}

```json
{
    "dependencies": {
    "follow-redirects": "^1.15.11",
    "form-data": "^4.0.5",
    "proxy-from-env": "^2.1.0"
  }
}
```

{% endcode %}

If at the time of installing the dependency the latest version of `follow-redirects` was `1.15.11` , your lockfile will refer to `follow-redirects` version `1.15.11` , this version contains **CVE-2026-40895**.

This CVE is fixed in version `1.16.0` of `follow-redirects`. follow-redirects is **not present** in the package.json but axios allows version `1.16.0` as defined by version selector `^1.15.11`. This means the issue can be fixed in the lockfile without changing the package.json, just like when you perform a clean npm install.

For this case AutoFix will only update the lockfile and not make changes to the package.json.

## Will this be reverted on a fresh install?

No, since axios is specifying the version range as `^1.15.11` , npm/yarn/pnpm/bun will automatically choose the latest minor version, which is `1.16.0` at the time of writing.

Same behavior applies for Poetry and UV.
