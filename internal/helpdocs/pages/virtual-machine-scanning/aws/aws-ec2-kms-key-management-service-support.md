# AWS EC2 KMS (Key Management Service) Support

When using AWS Customer Managed Keys (CMK) for encryption with your EC2 instances, you'll need to grant Aikido's scanning role access to these keys. This guide walks you through the necessary steps.

## Prerequisites <a href="#prerequisites" id="prerequisites"></a>

* You have EC2 instances encrypted with CMK.
* You have access to the key used for EBS encryption in AWS KMS.
* You have [enabled EC2 scanning in Aikido](/virtual-machine-scanning/aws/aws-ec2-virtual-machine-scanning-setup.md#getting-started). You will need the role ARN.

## Steps to Grant Access <a href="#steps-to-grant-access" id="steps-to-grant-access"></a>

1. Navigate to **AWS KMS** → **Customer managed keys**
2. Locate and click on the CMK you use for EC2 encryption
3. In the key details page, scroll to the **Key users** section
4. Click **Add** to include a new key user
5. Search for and select the Aikido VM scanning role
6. Save your changes

[image: AWS KMS key details showing configuration, key policies, and assigned key users.]

### JSON Policy <a href="#json-policy" id="json-policy"></a>

Alternatively, if you manage your KMS keys using infrastructure as code, add the following statements to the JSON policy of the key used for EBS encryption, replacing `arn:aws:iam::112233445566:role/aikido-security-ec2` with your Aikido VM scanning role ARN:

```json
{
  "Sid": "Allow use of the key",
  "Effect": "Allow",
  "Principal": {
    "AWS": "arn:aws:iam::112233445566:role/aikido-security-ec2" // Replace this
  },
  "Action": [
    "kms:Encrypt",
    "kms:Decrypt",
    "kms:ReEncrypt*",
    "kms:GenerateDataKey*",
    "kms:DescribeKey"
  ],
  "Resource": "*"
},
{
  "Sid": "Allow attachment of persistent resources",
  "Effect": "Allow",
  "Principal": {
    "AWS": "arn:aws:iam::112233445566:role/aikido-security-ec2" // Replace this
  },
  "Action": ["kms:CreateGrant", "kms:ListGrants", "kms:RevokeGrant"],
  "Resource": "*",
  "Condition": {
    "Bool": {
      "kms:GrantIsForAWSResource": "true"
    }
  }
}
```

## Verification <a href="#verification" id="verification"></a>

Once you've added the Aikido role as a key user, Aikido should be able to scan your CMK-encrypted EC2 instances.

## Important Notes <a href="#important-notes" id="important-notes"></a>

* Without these permissions, Aikido cannot scan EC2 instances encrypted with your CMK.
* This permission is specifically scoped to allow only the necessary access for scanning purposes.
* If you remove these permissions later, scanning capabilities will be affected.

## Troubleshooting <a href="#troubleshooting" id="troubleshooting"></a>

If you encounter scanning issues with CMK-encrypted instances, verify that:

* The correct CMK is selected
* The Aikido role is properly added as a key user
* The key policy hasn't been modified to remove the permission

## Managing which VMs are scanned <a href="#managing-which-vms-are-scanned" id="managing-which-vms-are-scanned"></a>

Aikido supports [an inclusion and exclusion model for VM scanning](/virtual-machine-scanning/aws/managing-which-vms-are-scanned.md).

## Scanning Frequency

Each day, 10% of all VMs enabled for scanning are rescanned, so on average, every VM is scanned once every 10 days.
