# Azure DevOps Server Setup for Local Image Scanning

The Aikido Security Local Scanner is a tool that enables you to perform Aikido Security scans within your environment, ensuring your data never leaves your premises. The scans take place locally, and the results are then uploaded to the Aikido Security platform. This setup allows you to scan any images locally on your own machine or in your CI/CD pipeline before they get sent to the registry.

## How to set up Local Scanning for images <a href="#how-to-set-up-local-scanning-for-images" id="how-to-set-up-local-scanning-for-images"></a>

### 1. Get your authentication token <a href="#id-1-get-your-authentication-token" id="id-1-get-your-authentication-token"></a>

1. Go to the [Local Scanner Image setup page](https://app.aikido.dev/settings/container-image-registry/add/localscan)
2. Generate an authentication token and copy. Note that you will only be able to view this token once.
3. You'll need this token in the Pipeline setup.

### 2. Running the Local Scanner <a href="#id-2-running-the-local-scanner" id="id-2-running-the-local-scanner"></a>

Navigate to Pipelines in the left hand navigation, click *Pipelines* and select *New Pipeline* or *Create Pipeline*. In the Configure step, select *Starter Pipeline*.

```yaml
trigger:
- master

pool:
  vmImage: ubuntu-latest

steps:
- script: docker build -t my-image .
  displayName: 'Build docker image'  
- script: docker run --rm -v /var/run/docker.sock:/var/run/docker.sock aikidosecurity/local-scanner image-scan my-image --apikey $(AIKIDO_API_KEY)
  displayName: 'Run Aikido scan'
```

You can also separate the building and scanning in separate stages or jobs. Use the example that best fits your use case and adapt to fit your needs.

```yaml
trigger:
  - master

pool:
  vmImage: ubuntu-latest

stages:
- stage: Build
  displayName: Build Docker Image
  jobs:
  - job: BuildAndSave
    displayName: Build and Save Docker Image
    steps:
      - script: |
          docker build -t my-image .
          docker save my-image -o $(Pipeline.Workspace)/my-image.tar
        displayName: 'Build and save Docker image to TAR'
      - task: PublishPipelineArtifact@1
        inputs:
          targetPath: '$(Pipeline.Workspace)/my-image.tar'
          artifactName: 'my-image-tar'
        displayName: 'Publish Docker image TAR as artifact'

- stage: Scan
  displayName: Scan Docker Image
  dependsOn: Build
  jobs:
  - job: RunScanner
    displayName: Scan with Aikido Local Scanner
    steps:
      - task: DownloadPipelineArtifact@2
        inputs:
          artifact: 'my-image-tar'
          path: '$(Pipeline.Workspace)'
        displayName: 'Download the saved image TAR'
      - script: |
          docker run --rm -v $(Pipeline.Workspace):$(Pipeline.Workspace) aikidosecurity/local-scanner image-scan $(Pipeline.Workspace)/my-image.tar --image-name my-image --apikey $(AIKIDO_API_KEY)
        displayName: 'Scan Docker image with Aikido'
```

\
Now it is time to add the api key as a secret variable. Click variables in the right hand corner and add a new variable. Fill in AIKIDO\_API\_KEY as the name. Check the *Keep this value secret* checkbox. Now save and run your pipeline.

### 3. Check your scanning results <a href="#id-3-check-your-scanning-results" id="id-3-check-your-scanning-results"></a>

After your first scan is done, you can go to the Aikido Feed to check out your results. An image with the name you specified will have been created in the Aikido UI, containing all results from the scanning.
