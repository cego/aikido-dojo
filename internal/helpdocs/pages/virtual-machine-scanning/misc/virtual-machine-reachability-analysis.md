# Virtual Machine Reachability Analysis

Virtual Machine Reachability Analysis shows how a VM can be reached inside your cloud network. It builds an interactive diagram of the network path to an instance (for example: Internet → Load Balancer → VM), including the ports involved, so you can quickly spot unintended exposure.

## What it helps you do

* **Confirm public exposure:** See if a VM is reachable from the internet, and through which components.
* **Understand the attack path:** Visualize the exact route traffic takes (load balancers, security groups, network hops) instead of just seeing a "public" or "private" label.
* **Reduce risk faster:** Identify surprising entry points and tighten access where needed.

## Where to find it

You can open the reachability diagram in two ways:

1. In [Virtual Machines](https://app.aikido.dev/virtual-machines), hover over your VM and select **"View Virtual Machine Reachability"**.
2. From the [feed](https://app.aikido.dev/queue), open any VM issue and click the reachability diagram link directly from the issue detail view.

### How reachability affects your issues

Aikido uses reachability data to automatically adjust the severity of VM issues, so you spend less time triaging findings that can't actually be exploited.

* **Severity downgrade:** If the affected package isn't reachable from the internet, Aikido lowers the severity. You'll see *"No network reachability"* in the scoring breakdown.
* **Severity upgrade:** If other risk factors apply (for example, a proof-of-concept exploit is publicly available), Aikido raises the severity accordingly.
* **Auto-ignore:** Aikido automatically ignores issues affecting specific packages if the VM is not reachable on the package-specific port(s). The issue shows an ignore reason of *"No network reachability"* in the feed.\
  \
  Here are some examples of packages and the ports Aikido checks:

  | Package(s)                       | Port(s)       |
  | -------------------------------- | ------------- |
  | `openssh`, `openssh-server`      | 22            |
  | `telnet`, `telnetd`, `inetutils` | 23            |
  | `redis`                          | 6379          |
  | `mongodb`, `mongod`              | 27017         |
  | `n8n`                            | 80, 443, 5678 |

  This list continues to grow as Aikido adds support for more packages.
