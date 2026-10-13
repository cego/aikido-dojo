# CI/CD

AutoShip checks the CI/CD results and fixes issues if needed.

AutoShip currently support GitHub Actions and CircleCI.

## GitHub Actions

Just select GitHub Actions as your CI/CD in the [global AutoShip settings](https://app.aikido.dev/autoship/settings), no additional configuration is required. The GitHub Action results are pulled through the Aikido AutoShip GitHub App.

Select GitHub actions as the CI/CD system for the repository during [the repository setup](https://app.aikido.dev/autoship/repositories).

## CircleCI

Enable CircleCI in the [global AutoShip settings](https://app.aikido.dev/autoship/settings). A Peronsal API Token is required for AutoShip to connect to CircleCI. The PAT can be set up here: <https://app.circleci.com/settings/user/tokens>.

Afer connecting AutoShip to CircleCI through the PAT, no additional configuration is required. AutoShip will use the name of the repository to locate the pipelines for that repo.

Select CircleCI as the CI/CD system for the repository during [the repository setup](https://app.aikido.dev/autoship/repositories).
