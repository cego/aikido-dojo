# Monorepos That Test Only Changed Paths

Run and upload coverage for every app in a monorepo. Skip unchanged apps only when you still upload their latest reports.

Run coverage for every app and upload every report on each run. This is the simple setup. One test run and one upload stay correct as apps are added, renamed, or removed, without a store of old reports to maintain.

{% hint style="warning" %}
Each upload replaces coverage for that repository and branch. Aikido does not merge it with the previous upload. If you upload only the app that changed, every other app is removed from Code Coverage until a later upload includes it again.
{% endhint %}

Skipping tests for apps that did not change is an optimization for a faster pipeline. The upload still has to contain every app. Keep the latest report for each app that did not run, and send it with the reports from the apps that did.

## Skip unchanged apps

When you use that optimization, keep one artifact per app and upload all of them on every run:

1. When an app's path changes, run its tests and publish that report as its own artifact, such as `coverage-web`.
2. When an app's path does not change, download that app's artifact from the last successful run and publish the same report again. Republishing refreshes how long the artifact is kept, so an app that rarely changes is not dropped when the previous artifact expires.
3. In a final job, download every app's artifact from this run and pass every report path in one upload. A second upload replaces the first.

{% hint style="warning" %}
Run coverage for every app once on the main branch before you start skipping unchanged apps. The upload step expects a report for every app, so the job fails until those artifacts exist on main.
{% endhint %}

## GitHub Actions

{% code title=".github/workflows/coverage.yml" overflow="wrap" %}

```yaml
name: Coverage

on:
  push:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        app: [web, api]
    permissions:
      contents: read
      actions: read
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Did this app change?
        id: changes
        run: |
          if git diff --name-only "${{ github.event.before }}" "${{ github.sha }}" | grep -q "^apps/${{ matrix.app }}/"; then
            echo "changed=true" >> "$GITHUB_OUTPUT"
          else
            echo "changed=false" >> "$GITHUB_OUTPUT"
          fi

      - name: Test the app that changed
        if: steps.changes.outputs.changed == 'true'
        run: npm test --workspace=apps/${{ matrix.app }} -- --coverage

      - name: Reuse the latest available report when nothing changed
        if: steps.changes.outputs.changed == 'false'
        env:
          GH_TOKEN: ${{ github.token }}
        run: |
          mkdir -p "apps/${{ matrix.app }}/coverage"
          run_id=$(gh run list --workflow coverage.yml --branch main --status success --limit 1 --json databaseId --jq '.[0].databaseId')
          gh run download "$run_id" --name "coverage-${{ matrix.app }}" --dir "apps/${{ matrix.app }}/coverage"

      - uses: actions/upload-artifact@v4
        with:
          name: coverage-${{ matrix.app }}
          path: apps/${{ matrix.app }}/coverage/lcov.info

  upload-coverage:
    needs: test
    runs-on: ubuntu-latest
    permissions:
      id-token: write
      contents: read
      actions: read
    steps:
      - uses: actions/checkout@v4
      - uses: actions/download-artifact@v4
        with:
          pattern: coverage-*
          path: coverage-reports

      - uses: AikidoSec/code-coverage-github-action@latest
        with:
          file-paths: |
            coverage-reports/coverage-web/lcov.info
            coverage-reports/coverage-api/lcov.info
```

{% endcode %}

Each test job either creates a new `lcov.info` or downloads the one saved on main. The upload job sends both. `actions/download-artifact` puts each artifact in a folder of the same name.

Name the workflow file `coverage.yml`. On a pull request, download the saved reports from `main`. The pull request branch has no earlier run to restore from.

## Other CI platforms

GitLab CI, Bitbucket Pipelines, and Azure DevOps can use the same optimization. Publish each app's report, restore the reports this pipeline did not regenerate, and upload every report together. [Upload a coverage report](/code-quality/code-coverage.md#upload-a-coverage-report) shows the upload step for each platform.

## Related

{% content-ref url="/pages/tlSdjHfgkHkVdfbFV5fb" %}
[Code Coverage](/code-quality/code-coverage.md)
{% endcontent-ref %}

{% content-ref url="/pages/ejEAOR4kBk0qpdGBxobg" %}
[PR Gating on Code Coverage](/code-quality/code-coverage/pr-gating-on-code-coverage.md)
{% endcontent-ref %}
