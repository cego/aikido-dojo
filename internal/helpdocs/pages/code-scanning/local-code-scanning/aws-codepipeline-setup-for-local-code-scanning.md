# AWS CodePipeline Setup for Local Code Scanning

The Aikido Security Local Scanner is a tool that enables you to perform Aikido Security scans within your environment, ensuring your code never leaves your premises. The scans take place locally, and the results are then uploaded to the Aikido Security platform. This setup allows easy integration of the Local Scanner into AWS CodePipeline for reporting purposes.

## How to set up Local Scanning <a href="#how-to-set-up-local-scanning" id="how-to-set-up-local-scanning"></a>

**Prerequisite**: make sure to have created an account that allows for Local Scanning. [More information on creating a Local Scanning Account](https://help.aikido.dev/en/articles/9070345-how-to-create-an-account-for-local-scanning-on-aikido).

In this guide, we’ll create a CI pipeline in AWS using AWS CodePipeline and AWS CodeBuild. The pipeline will run a security scan using the Aikido Security local scanner inside a container image.

## 1. Create the CodeBuild Project

Start by creating a CodeBuild project that will run the scan.

Navigate to CodeBuild

1. Open the AWS Console
2. Go to CodeBuild
3. Click Create build project

#### Create an `aikido-security-scan` project.

* Source: Select `No Source`. The source will be provided by the CodePipeline.
* Environment:

  * Image: Custom Image
  * Compute: EC2
  * Environment type: Linux Container
  * Image registry: Other Registry
  * External registry URL: public.ecr.aws/z1o3v2w5/aikidosecurity/local-scanner:latest

  
* Buildspec:

  ```
  version: 0.2

  phases:
    build:
      commands:
        - aikido-local-scanner scan . --apikey $AIKIDO_API_KEY --repositoryname MyRepo --branchname myDefaultBranch
  ```

#### Adding your API key to the CodeBuild Project

* Go to the [Local Scanner setup page](https://app.aikido.dev/settings/integrations/localscan)
* Generate an authentication token and copy. Note that you will only be able to view this token once.
* Add this token to AWS Secrets Manager with key AIKIDO\_API\_KEY and give the it a name like aikido/local-scanner-api-key. Note the full secret ARN; you'll need it shortly.
* Your CodeBuild project runs with an IAM service role. This role needs permission to read the secret. Add the following policy, replacing the ARN with your secret's ARN

```
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "secretsmanager:GetSecretValue"
      ],
      "Resource": "arn:aws:secretsmanager:<region>:<account-id>:secret:aikido/local-scanner-api-key-*"
    }
  ]
}
```

* Now go back to the CodeBuild Project and click Edit > Environment > Additional Configuration. Add an environment variable:
  * Name: AIKIDO\_API\_KEY
  * Value: aikido/local-scanner-api-key
  * Type: Secrets Manager

The value format is: \<secret-name>:\<json-key>. CodeBuild will automatically resolve this at build time and inject the value into the environment. See [AWS Docs](https://docs.aws.amazon.com/codebuild/latest/APIReference/API_EnvironmentVariable.html) for more info.

## 2. Create the CodePipeline

Set up a pipeline for the repository you want to scan.\
Add a source stage that triggers on changes to your default branch. In the Build stage. Choose AWS CodeBuild as the build provider and add the `aikido-security-scan` build project. This connects the pipeline to the build job we created earlier.

## 3. Check your scanning results

After your first scan is done, you can go to the Aikido Feed to check out your results. A repository will have been created, containing all results from the scanning.
