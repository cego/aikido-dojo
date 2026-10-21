# Monitoring

AutoShip uses your monitoring systems to detect regressions during deployments.

Sentry and AWS CloudWatch are currently supported, after connecting these integrations in the global configuration, AutoShip requires additional information to know where to find the monitoring information for the specific repository and environment.

## Sentry

In the [global AutoShip configuration](https://app.aikido.dev/autoship/settings), add Sentry under monitoring and provide the slug for your Sentry organization as found in your Sentry URL (eg. `slug.sentry.io`).

Create an Auth Token on the [Sentry Auth Tokens](https://sentry.io/settings/account/api/auth-tokens.) page. Make sure to give the token the following permissions:

* `event:read`
* `org:read`
* `project:read`

Once the global configuration is done, enable Sentry monitoring in the AutoShip repository configuration. Provide the Sentry Project for the specific repo and environment.

## AWS CloudWatch

In the [global AutoShip configuration](https://app.aikido.dev/autoship/settings), use the `Create IAM Role & Policy` button to apply a CloudFormation template that creates a IAM role that allows AutoShip to read CloudWatch logs.

Once the global configuration is done, enable Cloudwatch in the AutoShip repository configuration. Provide the log group for the specific repo and environment.
