# NPM, PNPM, Yarn, Bun - corrupt lockfile

A corrupt lockfile is most commonly caused by a mismatch in node, npm, pnpm, yarn or bun versions between your project and AutoFix. Make sure to configure the node and package manager version in your `package.json`.

AutoFix uses the node & npm/pnpm/yarn/bun version specified in the [engines section](https://docs.npmjs.com/cli/v11/configuring-npm/package-json#engines) of your `package.json` file. Make sure the node version and npm version in your `package.json` match the version you're using locally and in your CI/CD. By algining the versions, you will avoid changes to the `package-lock.json` files by AutoFix.

For example:

```
{
  "engines": {
    "node": "22",
    "npm": "11.14.1",
    "pnpm": "11.7.0",
    "yarn": "^4.17.11",
    "bun": "1.3.0"
  }
}
```
