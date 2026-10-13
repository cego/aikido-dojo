# Connect Google Cloud Project

Securing your cloud infrastructure is crucial to protecting your data. You can leverage Aikido's security checks to detect and address any misconfigurations in your infrastructure.

To view the list of security checks performed by Aikido on your cloud environment, go to the 'checks' tab on the [cloud overview page](https://app.aikido.dev/clouds). Filter to GCP to see specific checks performed on your connected GCP project(s).

To get started, head to the [cloud overview page](https://app.aikido.dev/clouds) on Aikido and click 'Connect cloud.' Follow the step-by-step setup wizard to connect your GCP project with Aikido.

First, you'll need to provide the project ID to help identify the correct project. Then, you'll be prompted to enable API access to specific GCP services. This is a critical step that enables Aikido to inspect the security of your cloud resources through API requests.

After enabling API access, the setup wizard will guide you through creating a service account with limited, read-only access to specific services in your GCP project. This account will be associated with the necessary roles and permissions, all of which are read-only. This ensures that Aikido can perform its security checks without the risk of unintended modifications to your resources.

{% hint style="info" %}
Consider [Workload Identity Federation](https://help.aikido.dev/cloud-scanning/connect-your-cloud/gcp/google-cloud-workload-identity-federation-setup) instead of the service account-based setup. You can choose it during the onboarding process.
{% endhint %}

### Custom Role

If you do not want to assign the suggested roles during the setup, you can [create](https://docs.cloud.google.com/iam/docs/creating-custom-roles#creating) the following custom role and assign it the service account:

```
title: Aikido Custom Role
description: Minimal permissions for Aikido Security
stage: GA
includedPermissions:
  - alloydb.clusters.list
  - alloydb.instances.getIamPolicy
  - alloydb.instances.list
  - apikeys.keys.list
  - appengine.applications.get
  - appengine.applications.getIamPolicy
  - appengine.services.list
  - appengine.versions.list
  - artifactregistry.dockerimages.get
  - artifactregistry.dockerimages.list
  - artifactregistry.files.get
  - artifactregistry.files.list
  - artifactregistry.locations.get
  - artifactregistry.locations.list
  - artifactregistry.packages.get
  - artifactregistry.packages.list
  - artifactregistry.repositories.get
  - artifactregistry.repositories.list
  - artifactregistry.tags.get
  - artifactregistry.tags.list
  - artifactregistry.versions.get
  - artifactregistry.versions.list
  - backupdr.backupPlanAssociations.list
  - batch.jobs.list
  - bigquery.datasets.get
  - bigquery.datasets.getIamPolicy
  - bigquery.tables.list
  - bigtable.instances.list
  - cloudasset.assets.listResource
  - cloudfunctions.functions.list
  - cloudkms.cryptoKeyVersions.list
  - cloudkms.cryptoKeys.getIamPolicy
  - cloudkms.cryptoKeys.list
  - cloudkms.cryptoKeyVersions.getIamPolicy
  - cloudkms.keyRings.getIamPolicy
  - cloudkms.keyRings.list
  - cloudfunctions.functions.getIamPolicy
  - cloudsql.backupRuns.list
  - cloudsql.instances.list
  - cloudsql.users.list
  - composer.environments.list
  - compute.addresses.list
  - compute.autoscalers.list
  - compute.backendBuckets.list
  - compute.backendServices.list
  - compute.disks.getIamPolicy
  - compute.disks.list
  - compute.firewalls.list
  - compute.forwardingRules.list
  - compute.globalAddresses.list
  - compute.globalForwardingRules.list
  - compute.healthChecks.list
  - compute.images.getIamPolicy
  - compute.images.list
  - compute.instanceGroupManagers.list
  - compute.instanceGroups.get
  - compute.instanceGroups.list
  - compute.instanceTemplates.list
  - compute.instances.getIamPolicy
  - compute.instances.list
  - compute.machineImages.list
  - compute.networks.list
  - compute.nodeGroups.list
  - compute.nodeTemplates.list
  - compute.projects.get
  - compute.regions.list
  - compute.regionSslPolicies.list
  - compute.resourcePolicies.getIamPolicy
  - compute.resourcePolicies.list
  - compute.routers.list
  - compute.routes.list
  - compute.securityPolicies.list
  - compute.snapshots.list
  - compute.sslPolicies.list
  - compute.subnetworks.getIamPolicy
  - compute.subnetworks.list
  - compute.targetHttpProxies.list
  - compute.targetHttpsProxies.list
  - compute.targetPools.list
  - compute.targetSslProxies.list
  - compute.targetVpnGateways.list
  - compute.urlMaps.list
  - compute.vpnGateways.list
  - compute.vpnTunnels.list
  - container.clusters.get
  - container.clusters.getIamPolicy
  - container.clusters.list
  - dataplex.lakes.list
  - dataplex.tasks.list
  - dataplex.zones.list
  - dataproc.clusters.getIamPolicy
  - dataproc.clusters.list
  - datastore.databases.list
  - deploymentmanager.deployments.list
  - dns.managedZones.list
  - dns.policies.list
  - dns.resourceRecordSets.list
  - iam.roles.get
  - iam.roles.list
  - iam.serviceAccountKeys.get
  - iam.serviceAccountKeys.list
  - iam.serviceAccounts.getIamPolicy
  - iam.serviceAccounts.list
  - logging.buckets.list
  - logging.logMetrics.list
  - logging.sinks.list
  - metastore.services.list
  - monitoring.alertPolicies.list
  - monitoring.groups.list
  - monitoring.notificationChannels.list
  - orgpolicy.policies.list
  - orgpolicy.policy.get
  - pubsub.snapshots.list
  - pubsub.subscriptions.list
  - pubsub.topics.getIamPolicy
  - pubsub.topics.list
  - redis.clusters.list
  - redis.instances.list
  - resourcemanager.folders.get
  - resourcemanager.folders.list
  - resourcemanager.projects.get
  - resourcemanager.projects.getIamPolicy
  - resourcemanager.projects.list
  - run.jobs.getIamPolicy
  - run.jobs.list
  - run.locations.list
  - run.services.getIamPolicy
  - run.services.list
  - secretmanager.secrets.getIamPolicy
  - secretmanager.secrets.list
  - serviceusage.services.list
  - spanner.instances.list
  - storage.buckets.getIamPolicy
  - storage.buckets.list
  - tpu.nodes.list
  - vpcaccess.connectors.list

```

Once the service account is created, you'll need to generate an access key and upload it to Aikido. This key will be used by Aikido to make the necessary API requests to scan your resources.

Finally, you can name your connected project in Aikido and specify the environment it operates in. This information helps Aikido prioritize findings based on the severity and impact to your business.

Within 1-2 minutes after connecting your account, Aikido will report misconfigurations that could pose a threat.

### Advanced Rules

Besides the checks mentioned above, Aikido offers a suite of complementary checks/rules that you can enable. We call these advanced cloud rules and you can find them [here](https://app.aikido.dev/clouds/checks?cloudCheckType=advanced). After enabling any of these rules, you can expect to see the results (as new issues in [the feed](https://app.aikido.dev/queue)) within a few seconds.

Just like the standard checks, these are evaluated with each scan of your cloud environments. Moreover, they are mapped to the compliance reports. By default, the advanced rules will appear as *disabled* in the compliance reports, unless you activate them.
