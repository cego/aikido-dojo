# Cargo - Private Packages

For Aikido to update dependencies that include private packages, it needs access to your private registries so it can generate accurate lockfile updates. Follow the steps below to configure private Cargo packages.

1. Navigate to the [autofix settings](https://app.aikido.dev/issues/fix/settings) page and click on "*Connect Registry*"
2. Select *Cargo*

3. Cargo uses the `CARGO_REGISTRIES_<INDEX_NAME>_TOKEN` environment variables for token-based authentication to a registry named `<INDEX_NAME>`.

For example for the following `config.toml` file:

```
[registries.example]
index = "https://example.com/git/index"
```

You can use `CARGO_REGISTRIES_EXAMPLE_TOKEN` to set the token for the example registry.

4. Click on "*Save*"
