# Managing Which VMs Are Scanned

Aikido gives you precise control over which virtual machines (VMs) are included in security scans through AWS tags. This allows you to customize your scanning scope based on your specific security needs.

### Available Tags

#### AIKIDO\_INCLUDE\_VM

When this tag is applied to any VM in a region with the value `true`, Aikido switches to an "opt-in" scanning model, where **only** VMs with this tag will be scanned in that region.

#### AIKIDO\_EXCLUDE\_VM

When this tag is applied to a VM with the value `true`, Aikido will skip scanning this specific VM, regardless of other scan configurations.

{% hint style="info" %}
**Important Notes**

* The AIKIDO\_INCLUDE\_VM tag changes scanning behavior for the entire region when used
* AIKIDO\_EXCLUDE\_VM takes precedence over AIKIDO\_INCLUDE\_VM if both are applied
* If no AIKIDO\_INCLUDE\_VM tags exist in a region, Aikido scans all VMs by default (except those with AIKIDO\_EXCLUDE\_VM)

By leveraging these tags, you can create a tailored scanning strategy that aligns with your security priorities and resource management needs.
{% endhint %}

### Use Cases

#### 1. Scanning Only Production Environments

If you want to focus your security scans on production workloads while ignoring development environments:

<pre><code><strong># Tag all production VMs
</strong>aws ec2 create-tags --resources i-1234567890abcdef0 i-0987654321fedcba0 --tags Key=AIKIDO_INCLUDE_VM,Value=true
</code></pre>

#### 2. Excluding Sensitive Systems

For VMs that contain sensitive data or require specialized handling:

```
# Exclude a sensitive database server
aws ec2 create-tags --resources i-db12345678901234 --tags Key=AIKIDO_EXCLUDE_VM,Value=true
```

#### 3. Region-Specific Scanning Policies

Since the `AIKIDO_INCLUDE_VM` tag is region-scoped, you can implement different scanning policies per region:

* In us-east-1: Tag only critical workloads for scanning
* In us-west-2: Allow scanning of all VMs by not using the `AIKIDO_INCLUDE_VM` tag

### Scanning Frequency <a href="#scanning-frequency" id="scanning-frequency"></a>

Each day, 10% of all VMs enabled for scanning are rescanned, so on average, every VM is scanned once every 10 days.
