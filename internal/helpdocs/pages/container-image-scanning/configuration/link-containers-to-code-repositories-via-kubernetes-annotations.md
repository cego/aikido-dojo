# Link Containers to Code Repositories via Kubernetes Annotations

Linking containers to code repositories keeps related vulnerabilities in one place. This is useful when you manage many repositories and images.

Aikido can link containers automatically by reading a specific annotation on your Kubernetes workloads.

#### Standard annotation

Add the `aikido.dev/repository` annotation to the workload's top-level `metadata.annotations`. Set the value to your own repository's URL.

For a Deployment, it looks like this:

{% code title="deployment.yaml" %}
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: demo-app
  annotations:
    aikido.dev/repository: https://github.com/AikidoSec/demo-app-1
spec:
  replicas: 1
  selector:
    matchLabels:
      app: demo-app
  template:
    metadata:
      labels:
        app: demo-app
    spec:
      containers:
        - name: demo-app
          image: demo-app:latest
```
{% endcode %}

#### Custom annotation

If your organization already has its own convention (for example, `repo-url`), Aikido can map that key to the standard one for you, so you don't need to change your workload manifests.

{% hint style="info" %}
To use a custom annotation key, reach out to Aikido support to have it mapped for your account.
{% endhint %}

#### How it works

* Aikido reads the annotation off the Kubernetes workload (standard key, or your mapped custom key).
* Aikido uses this annotation to automatically link the container to the matching code repository.&#x20;
* If the container is already linked to a repository, this won't override it.
