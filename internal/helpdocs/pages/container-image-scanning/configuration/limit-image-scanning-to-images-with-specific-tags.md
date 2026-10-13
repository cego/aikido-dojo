# Limit Image Scanning to Images with Specific Tags

By default, Aikido will always scan **the last image that was uploaded** in your docker repository.

In some cases it might be interesting to make sure only an image tagged with a specific tag (such as 'production-') is scanned.

Aikido allows you to set a tag filter per image. Go to [container settings ](https://app.aikido.dev/settings/container-image-registry)and click the triple dots action menu to enable this

[image: Docker repo management interface with options to add, scan, or edit tag filters.]

It's possible to scan only images that start with specific words such as 'production-12345678'. To do this, enter the tag production-\* using \* as a wildcard.

***
