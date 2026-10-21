# Configuration via Environment Variables

## Token <a href="#token" id="token"></a>

Use the `AIKIDO_TOKEN` env variable to set the token you generated in Aikido:

```bash
AIKIDO_TOKEN="token"
```

## Blocking Mode <a href="#blocking-mode" id="blocking-mode"></a>

You can use the `AIKIDO_BLOCK` env variable to enable blocking mode.

```bash
AIKIDO_BLOCK="true"
```

> You can also toggle between blocking mode and detection only mode in the dashboard. The environment variable will be used as starting mode. When the config is retrieved, Zen will switch to the configured mode from the dashboard.

## Disable the Agent <a href="#disable-the-agent" id="disable-the-agent"></a>

You can use the `AIKIDO_DISABLE` env variable to disable Zen entirely. In the unlikely event that Zen causes problems, you can disable Zen without removing any code or environment variables (e.g. `AIKIDO_TOKEN` etc).

```bash
AIKIDO_DISABLE="true"
```

## Debugging <a href="#debugging" id="debugging"></a>

You can use the `AIKIDO_DEBUG` env variable to output more information:

* Whether certain packages are supported or not
* Whether a token was detected
* Whether Zen is running in detection only mode

```bash
AIKIDO_DEBUG="true"
```

## Proxy Settings <a href="#proxy-settings" id="proxy-settings"></a>

You can use the `AIKIDO_TRUST_PROXY` env variable to [trust proxies](/zen-firewall/zen-installation-instructions/proxy-load-balancer-settings.md):

> By default, trust proxy is **enabled**.

```bash
AIKIDO_TRUST_PROXY="false"
```

## Instance Name Override <a href="#limit-api-discovery-samples" id="limit-api-discovery-samples"></a>

`AIKIDO_INSTANCE_NAME` sets a custom name for your application instance, useful when the hostname isn't descriptive enough (for example, in containerized or ephemeral environments). Defaults to the hostname.

```bash
AIKIDO_INSTANCE_NAME="custom-name"
```

## Limit API Discovery Samples <a href="#limit-api-discovery-samples" id="limit-api-discovery-samples"></a>

Limits the number of request samples collected per route. Each sample contains:

* Request body structure and type (for non-GraphQL requests)
* Query parameters
* Authentication methods used

These samples help analyze and document your API's structure, including the expected request format and authentication requirements.

```bash
AIKIDO_MAX_API_DISCOVERY_SAMPLES="10"
```

## Maximum request body size (NodeJS) <a href="#maximum-request-body-size-nodejs" id="maximum-request-body-size-nodejs"></a>

Sets the maximum size (in MB) for incoming HTTP request bodies. This limit helps protect your application from memory exhaustion due to extremely large requests.

```bash
AIKIDO_MAX_BODY_SIZE_MB="10" # MB
```
