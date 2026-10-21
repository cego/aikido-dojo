# Link Containers to Code Repositories via Kubernetes Annotations

Linking containers to code repositories keeps related vulnerabilities in one place. This is useful when you manage many repositories and images.

Aikido can link containers automatically by reading a specific label or annotation on your Kubernetes workloads.

#### Standard annotation

Add this label or annotation to your Kubernetes workload, using the key `aikido.dev/repository`. Set the value to your own repository's URL, for example:

```
aikido.dev/repository: 
https://github.com/AikidoSec/demo-app-1
```

{% hint style="info" %}
If both a label and an annotation with this key are set on the same workload, the annotation value is used.
{% endhint %}

#### Custom annotation

If your organization already has its own convention (for example, `repo-url`), Aikido can map that key to the standard one for you, so you don't need to change your workload manifests.

{% hint style="info" %}
To use a custom annotation key, reach out to Aikido support to have it mapped for your account.
{% endhint %}

#### How it works

* Aikido reads the relevant label or annotation off the Kubernetes workload (standard key, or your mapped custom key).
* Aikido uses this tag to automatically link the container to the matching code repository.
* If the container is already linked to a repository, this won't override it.
