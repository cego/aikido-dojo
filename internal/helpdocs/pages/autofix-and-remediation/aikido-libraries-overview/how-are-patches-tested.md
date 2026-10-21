# How Are Patches Tested?

This applies to both **Aikido Images** and **Aikido Libraries**.

Every patch goes through multiple layers of automated and human testing, both before and after the patch is applied, to minimize the risk of introducing breaking changes, regressions, or additional vulnerabilities.

Automated testing includes:

* Unit and integration tests
* Functional tests
* Smoke tests
* Installation tests
* Fuzz tests
* Security tests
* PoC / exploit validation

Each patch is then independently reviewed and revalidated by human security engineers before release.
