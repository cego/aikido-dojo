# Aikido Images

Updating to the latest version of a base image can be difficult because it may require changes to your application. When updating to a newer base image is not a viable option, you can stay secure by using **Aikido Images**.

Aikido Images is a registry of zero-CVE base images: drop-in replacements that carry patched versions of libraries with reported **CRITICAL** or **HIGH** severity security issues.

{% hint style="info" %}
Aikido Images cover OS-level CVEs under a defined [patch creation SLA](/autofix-and-remediation/aikido-libraries-overview/patch-creation-sla.md). Application-level package vulnerabilities (Python, npm, Java, Go, etc.) are also fixed, but on a best-effort basis with no guaranteed SLA. For guaranteed application-level coverage, see [Aikido Libraries](/autofix-and-remediation/aikido-libraries-overview.md).
{% endhint %}

For example, [CVE-2025-4373](https://security-tracker.debian.org/tracker/CVE-2025-4373) is fixed by Debian in Trixie and Sid but not in Bookworm. The Aikido Images variant of `debian:bookworm` contains a patched version of `glib2.0` that fixes this vulnerability. Using this image avoids breaking changes while maintaining a good security posture.

## Aikido Images availability <a href="#aikido-images-availability" id="aikido-images-availability"></a>

You can check the available images on <https://images.aikido.dev/> or in the Aikido UI under [AutoFix > Aikido Images](https://app.aikido.dev/issues/aikido-images).

A subset of the Aikido Images is available on [the public ECR gallery](https://gallery.ecr.aws/). These images are identical to the images in the `docker.aikido.io` registry and are updated on the same cadence. The public images can be pull from the public ECR and `docker.aikido.io`.

## See also

{% content-ref url="/pages/NgYUbx12E86BLPyg3v58" %}
[AutoFix for Containers: Aikido Images](/autofix-and-remediation/aikido-images/autofix-for-containers-using-aikido-images.md)
{% endcontent-ref %}
