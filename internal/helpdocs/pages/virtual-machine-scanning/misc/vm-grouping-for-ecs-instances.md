# VM Grouping for EC2 instances

To optimize scanning efficiency, Aikido groups certain EC2 instances and scans only one instance from each group. Grouping works as follows:

* **Auto Scaling Groups (ASG)**: All EC2 instances in the same AWS Auto Scaling Group are shown as a single VM group in Aikido. The VM group's name matches the ASG name.
* **Karpenter Node Pools**: All EC2 instances that belong to the same Karpenter node pool within an EKS cluster are grouped together. The VM group's name matches the instance name pattern used by Karpenter.
* **No grouping**: EC2 instances that are not part of an ASG or Karpenter node pool are treated as standalone VMs and scanned individually.
