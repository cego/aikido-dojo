# Webhooks

In addition to Slack and Microsoft Teams notifications it's also possible to setup Webhooks for Zen events. To create or modify a Webhook, open the [Webhooks page](https://app.aikido.dev/settings/integrations/api/webhooks) .

Aikido will attempt to re-submit the notification 5 times if the request is not successful. A request is deemed unsuccessful if HTTP status other than 200 is returned or if a timeout occurs. [This page](https://apidocs.aikido.dev/reference/verirfying-webhooks) describes how webhook events can be verified.

### Webhook Events

The following events are currently sent:

<details>

<summary>Attack detected (Sample)</summary>

{% code overflow="wrap" %}

```json
{
   "event_type": "zen.attack",
   "created_at": 1760015974,
   "payload": {
      "url": "http://app.aikido.dev/runtime/services/X/events?groupId=X&eventId=X",
      "eventId": 123,
      "attack": {
         "blocked": false,
         "kind": "sql_injection"
      },
      "request": {
         "method": "POST",
         "route":"/add",
         "ip_address": "::ffff:127.0.0.1",
         "user_agent": "Mozilla/5.0 ..."
      },
      "service": {
         "id": 1,
         "name": "Test",
         "environment": "production"
      },
      "instance":{
         "id": 1,
         "ip_address": "172.17.0.123",
         "hostname": "hostname.local"
      },
      "user": {
         "id": "123abc",
         "blocked_at": 0
      },
      "message": "SQL injection detected in Test"
   },
   "dispatched_at": 1760016033,
   "workspace_id": 123
}
```

{% endcode %}

</details>

<details>

<summary>Attack wave detected (Sample)</summary>

```json
{
   "event_type": "zen.attack_wave",
   "created_at": 1760015974,
   "payload": {
      "url": "http://app.aikido.dev/runtime/services/X/events?groupId=X&eventId=X",
      "eventId": 123,
      "attack": [],
      "request": {
         "ip_address": "::ffff:127.0.0.1",
         "user_agent": "Mozilla/5.0 ..."
      },
      "service": {
         "id": 1,
         "name": "Test",
         "environment": "production"
      },
      "instance":{
         "id": 1,
         "ip_address": "172.17.0.123",
         "hostname": "hostname.local"
      },
      "user": {
         "id": "123abc",
         "blocked_at": 0
      },
      "message": "Attack wave detected in Test"
   },
   "dispatched_at": 1760016033,
   "workspace_id": 123
}
```

</details>

<details>

<summary>New outbound connection (Sample)</summary>

```json
{
   "event_type": "zen.outbound.discovered",
   "created_at": 1760016813,
   "payload": {
      "hostname": "aikido.dev",
      "port": 443,
      "service": {
         "id": 123,
         "name": "Test",
         "environment": "production"
      },
      "message": "New outbound connection detected in Test"
   },
   "timestamp": 1757599746000
}
```

</details>
