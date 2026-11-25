# Compiled artifact scanning with Local Scanner

Use Aikido Local Scanner in image mode to scan a directory of compiled binaries or a mounted filesystem. The scan inventories packages in that directory, creates a container in your workspace, and can write a CycloneDX SBOM to a local file.

Use this when the software is not a container image and has no registry. Typical targets are a Windows build directory, or a Linux SD card from an embedded board after you mount it.

The repository [`scan`](https://help.aikido.dev/code-scanning/local-code-scanning/cli-options-for-local-scanner) command reads source code and dependency files. For a folder of compiled binaries, run `image-scan` on that folder.

## Scan the directory

{% tabs %}
{% tab title="Linux" %}
Mount the SD card or filesystem, then scan the mount:

```shellscript
aikido-local-scanner image-scan /mnt/sdcard \
  --image-name board-sdcard \
  --output-cyclonedx-json sbom.json \
  --apikey AIK_CI_xxx
```
{% endtab %}
{% tab title="Windows" %}
Scan the directory that contains the binaries:

```shellscript
aikido-local-scanner.exe image-scan "C:\build\output" --image-name product-binaries --output-cyclonedx-json sbom.json --apikey AIK_CI_xxx
```
{% endtab %}
{% endtabs %}

`--image-name` is the container name Aikido creates for this scan. `--output-cyclonedx-json` writes a CycloneDX JSON file locally. The full flag list is in [CLI options for local image scanning](https://help.aikido.dev/container-image-scanning/local-image-scanning/cli-options-for-local-image-scanning).

## Where the results go

The scan creates a container in your workspace under the image name you set. Findings appear in the feed with your other container results.

You can download the SBOM again later:

* On the container page, open **Actions** and choose **Export RAW SBOM**. See [Export RAW SBOM of Your Containers](https://help.aikido.dev/container-image-scanning/configuration/export-raw-sbom-of-your-containers).
* From **Reports > Licenses & SBOM**. See [Exporting an SBOM](https://help.aikido.dev/code-scanning/miscellaneous/generate-sbom-based-on-open-source-packages).

## What the SBOM includes

Aikido inventories package metadata it finds in the directory, the same way it does for a container image. That includes OS package databases and dependency manifests present on the filesystem.

A directory of proprietary binaries with no package metadata produces a thin SBOM. Aikido does not recover library versions from compiled `.so`, `.exe`, or `.dll` files.

[Lockfile-less C/C++ scanning](https://help.aikido.dev/code-scanning/scanning-practices/c-c-lockfile-less-scanning) runs on a repository scan. It hashes C and C++ source files in the tree. It does not identify libraries inside compiled binaries.

## Buildroot and Yocto

For an embedded system built with Buildroot or Yocto, the SBOM generated at build time is more complete than a scan of the finished filesystem. Upload that file with the [container SBOM API](https://help.aikido.dev/code-scanning/miscellaneous/generate-sbom-based-on-open-source-packages#upload-the-sbom-to-aikido).
