# AutoFix for Infrastructure as Code

AutoFix for Infrastructure as Code supports the following IaC tools:

* Terraform
* CloudFormation
* Serverless framework
* Docker
* Kubernetes
* Helm

AutoFix creates the patch for the IaC files, it **does not automatically apply** the changes to your cloud. The execution of the changes needs to be executed by the user.

On the [AutoFix > IaC](https://app.aikido.dev/issues/fix/iac) page in the Aikido UI, you will find the detected IaC issues. AutoFix can generate a preview of the change or create a Pull Request.

Clicking AutoFix will generate a preview of the change required to remediate the issue.

From here you can

* Create a Pull Request
* Apply the patch directly in your IDE
* Refine the solution with Aikido AI
