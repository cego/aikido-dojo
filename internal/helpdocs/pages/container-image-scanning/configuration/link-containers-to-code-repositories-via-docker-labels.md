# Link Containers to Code Repositories via Docker Labels

You can link container images to code repositories by adding a Docker label to your images. This is useful when you build images from many different repos and want Aikido to automatically map them.

### Supported labels

Aikido recognizes two Docker labels for linking. You only need one of them.

* `aikido.repository`: Aikido's own label
* `org.opencontainers.image.source`: the standard OCI label

Set the value to the full URL of the code repository the image belongs to.

If both labels are present, `aikido.repository` takes priority.

### How it works

When Aikido scans a container image, it inspects the image's OCI labels. If it finds a supported label, it automatically links that image to the matching code repository.

Once linked, container issues appear inside the relevant repository in your Aikido feed.

{% hint style="info" %}
The label must be set at build time in the Dockerfile. Aikido reads it directly from the image metadata, so it can't be added after the image is pushed.
{% endhint %}

### Unlinking

Once an image is linked via a Docker label, it can't be unlinked through the Aikido UI. To change the link, update the label value in your Dockerfile and rebuild the image.
