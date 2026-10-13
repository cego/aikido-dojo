# Zen Firewall for Serverless Environments

Zen Firewall can be deployed directly into popular serverless platforms to protect your applications and event-driven workloads. We currently support:

* [AWS Lambda (NodeJS)](/zen-firewall/zen-installation-instructions/zen-firewall-for-nodejs-javascript-and-typescript/aws-lambda.md)
  * API Gateway
  * SQS
* [Google Cloud Functions (NodeJS)](/zen-firewall/zen-installation-instructions/zen-firewall-for-nodejs-javascript-and-typescript/google-cloud-functions.md)
* [Google Cloud Pub/Sub (NodeJS)](/zen-firewall/zen-installation-instructions/zen-firewall-for-nodejs-javascript-and-typescript/google-cloud-pubsub.md)

## How It Works

Zen integrates at the function or message level, inspecting requests and events before they reach your code. We’ve optimized our engine to ensure minimal latency and overhead, so your serverless applications remain fast and cost-efficient.

## Key Benefits

* **Lightweight and performant**: Optimized to run in short-lived serverless environments.
* **Automatic protection**: Detects and blocks critical attacks such as SQL injection, shell injection, path traversal, and SSRF.
* **Easy integration**: Deploy Zen without changing your application logic.

## Feature Limitations

Due to constraints of serverless environments, some firewall capabilities are not available:

* Tor traffic blocking
* Bot blocking
* IP blocking by known threat actors
* Blocking by country
* Manual user/IP blocking
