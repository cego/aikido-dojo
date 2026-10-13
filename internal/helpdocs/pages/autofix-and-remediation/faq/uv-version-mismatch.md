# UV - version mismatch

AutoFix uses UV version **0.11.19** when `required-version` is not set in your `pyproject.toml`.

If your project needs a different UV version, set it in `pyproject.toml`. For example, add the following to make AutoFix use UV version **0.10.3**:

{% code title="pyproject.toml" %}

```toml

[tool.uv]
required-version = "==0.10.3"

```

{% endcode %}

For more information see <https://docs.astral.sh/uv/reference/settings/#required-version>
