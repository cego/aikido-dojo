# Python Projects: Security Scanning Best Practices

Aikido can find known vulnerabilities (CVEs) in your Python dependencies and flag risky or unwanted licenses used by those dependencies.

To do this accurately, Aikido needs access to the exact versions of all your dependencies, including the ones that those packages depend on.

## How does Aikido find dependencies and subdependencies?

Out of the box, Aikido supports the following files for dependency scanning:

* `uv.lock`
* `Pipfile.lock`
* `poetry.lock`
* `pdm.lock`
* `requirements.txt` *(fallback only when no lockfile is present)*
* `requirements.yml` *(fallback only when no lockfile is present)*

If your project only includes a `requirements.txt`, Aikido might not detect the full dependency tree.

That’s because pip does not record the exact versions of subdependencies, so results may be incomplete or outdated.

For complete and reliable scans, we recommend using uv, a modern Python package manager that generates lockfiles.

## Why use lockfiles?

* **Accurate security scans**: Lockfiles includes every resolved dependency with an exact version, giving Aikido a full view of your risk surface.
* **Protection against supply chain attacks**: Locking prevents unwanted version changes that could introduce malicious packages.
* **Predictable builds**: Every environment installs the same versions. Fewer “works on my machine” issues.
* **Faster installs**: Once dependencies are locked, lockfiles skip heavy resolver steps.

## How to start using lockfiles in your project

### UV

To get started, [install uv](mailto:undefined) and [set up your project with a pyproject.toml](https://docs.astral.sh/uv/guides/projects/#creating-a-new-project) if it doesn’t already have one.

From there, you can [define your dependencies and let uv generate a uv.lock file](mailto:undefined) that pins exact versions for every dependency and subdependency.

Both pyproject.toml and uv.lock should be committed to your repository, and you should avoid editing the lockfile manually.

If your project currently uses a requirements.txt, [uv can import those dependencies directly and create a lockfile for you](mailto:undefined). This ensures Aikido can detect your full dependency tree with complete version accuracy.

You can find detailed steps for project setup, dependency management, and migration in the [uv documentation](https://help.aikido.dev/).

### Poetry

To get started, [install Poetry and set up your project with a pyproject.toml ](https://python-poetry.org/docs/)if it doesn’t already have one.

From there, you can [define your dependencies and let Poetry generate a poetry.lock file](https://python-poetry.org/docs/managing-dependencies/) that pins exact versions for every dependency and the packages they rely on.

Both pyproject.toml and poetry.lock should be committed to your repository, and you should avoid editing the lockfile manually.

You can find detailed steps for project setup, dependency management, and migration in the [Poetry documentation](https://python-poetry.org/).
