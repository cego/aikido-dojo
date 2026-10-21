# Kaniko Image Scanning with Local Scanner

When performing image scans using the Aikido Local Scanner, the Docker daemon is used to find the image to perform the scan on. Since [Kaniko](https://github.com/GoogleContainerTools/kaniko) does not depend on a Docker daemon, an alternative approach is needed. That approach is to create a tarball and perform the scan on that tarball.

You can create a tarball from your image using the `--tar-path` flag in your kaniko build command. Alternatively, you can download a tar version of the image from the repository.

Once the tarball is available, use this command to perform the scan.

```shellscript
aikido-local-scanner image-scan /my-image.tar --image-name my-image-name --apikey AIK-xxx
```
