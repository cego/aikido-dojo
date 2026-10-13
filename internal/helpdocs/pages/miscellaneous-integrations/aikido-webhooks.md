# Aikido Webhooks

## Webhooks Authentication <a href="#webhooks-authentication" id="webhooks-authentication"></a>

In order to verify that incoming webhooks from Aikido are actually originating from Aikido and that it’s payload has not been altered, Aikido uses hashes from the payload signed with a secret we share with you.

{% hint style="success" %}
Aikido will always send the webhook requests from the following IP address: 52.18.113.172.
{% endhint %}

### Webhook Schema <a href="#webhook-schema" id="webhook-schema"></a>

The webhook schema can be found in [our apidocs](https://apidocs.aikido.dev/reference/webhooks).

### Generating a webhook secret <a href="#generating-a-webhook-secret" id="generating-a-webhook-secret"></a>

1. Go to the [webhooks integration page](https://app.aikido.dev/settings/integrations/api/webhooks)
2. Create a secret by clicking on ‘Add secret’, right above the table on the right-hand side.\
   ​

   [image: Empty webhooks dashboard prompting user to register a new webhook.]
3. When the webhook secret is created, you’ll see a modal where you can copy the webhook secret and store it securely. It’s important that this secret is kept safe and not committed in any of your code repos

### Verifying a Webhook <a href="#verifying-a-webhook" id="verifying-a-webhook"></a>

Whenever Aikido will send you a webhook for an event you configured, we’ll create a hash of the payload of the event and sign it with the secret you have just generated. This unique hash will be included with the via the `X-Aikido-Webhook-Signature` request header. Via this way you can verify that a webhook and its payload is genuine.

In order to protect against replay attacks, we include an epoch timestamp in the webhook payload from right before we dispatch the HTTP request. This timestamp is included as the `dispatched_at` property and should not be older than 30 seconds when you verify the payload.

Regardless of the programming language of your choosing, the process to verify an incoming webhook looks like this:

1. Ensure that the payload is a valid JSON string
2. Get the signature from the `X-Aikido-Webhook-Signature` request header
3. Parse the request body back to a JSON string
4. Create a hmac digest from the stringified request body with the `sha256` algorithm and sign it with your webhook secret from Aikido
5. Validate that the signature from the request header matches with the digest you just generated.
6. Validate that the `dispatched_at` epoch timestamp from the property is not older than 30 seconds

Below we’ll share some pseudo Javascript code showing how you can verify the hash when using the `express` framework. This should be included as a middleware and more validation on the values should be performed to see if they are valid.

```
const { createHmac } = require('node:crypto');

const express = require('express');
const bodyParser = require('body-parser');

const PORT = 4000;

const app = express();

app.use(bodyParser.json());

const isIncomingWebhookValid = (payload, signature) => {
	// get the raw webhook secret from the environment variables
	const aikidoWebhookSecret = process.env.AIKIDO_WEBHOOK_SECRET;

	// lets create the hmac instance
	const hmac = createHmac('sha256', aikidoWebhookSecret);

	// lets convert the payload object to a JSON string
	const rawPayload = JSON.stringify(payload);

	// update the hmac content with the stringified payload
	hmac.update(rawPayload);

	// calculate the digest of the hmac content and return it as a hex value
	const payloadDigest = hmac.digest('hex');

	// if digest does not match the provided signature, the webhook is invalid
	if (payloadDigest !== signature) return false;

	// get the current epoch timestamp
	const currentEpochTimestamp = Math.floor(new Date().getTime() / 1000);
	
	// if the dispatched_at epoch timestamp from the payload is longer than 30 seconds old, the webhook is invalid
	if ((currentEpochTimestamp - (payload['dispatched_at'] ?? 0)) > 30) return false;

	// webhook passed all checks and is valid
	return true;
}

app.post('/webhooks/aikido/issue-created', async (req, res) => {
	const isValid = isIncomingWebhookValid(req.body, req.headers['X-Aikido-Webhook-Signature']);
	if (!isValid) {
		throw new Error(`The request signature is invalid`)
	}
    
	// your business code

    res.status(204);
});

app.listen(PORT, () => {
	console.log(`server listening on port ${port}`);
});
```

***
