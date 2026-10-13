# Scala: Dependency Management and Scanning for build.sbt

## The Challenge <a href="#the-challenge" id="the-challenge"></a>

Aikido can find known vulnerabilities (CVE) in your Scala dependencies as well as malware and dangerous licenses being used by those dependencies.

How does Aikido find those dependencies and their transitive subdependencies?

In case of Scala we scan the `build.sbt` file for dependencies. It should be noted that `build.sbt` files might not contain exact versions for some of your dependencies. This can cause Aikido to not find the full range of risks in your application.

It's therefore recommended to use `build.sbt.lock` lockfiles that contain an exact version for each dependency as well as the subdependencies.

There are other reasons to use lockfiles besides making security scanning easier for Aikido:

* Using a lockfile protects you against supply chain attacks via malicious packages. This kind of attack is becoming more popular
* Using a lockfiles makes your build more predictable as everyone is using the exact same minor version of packages. Less chance of 'works on my machine'
* Faster build times: no need for dependency resolution anymore

## Example <a href="#example" id="example"></a>

Let's look at a recent case that highlights these challenges.

```
// Original build.sbt without lock files
libraryDependencies ++= Seq(
  "com.typesafe.akka" %% "akka-http" % "10.2.+",
  "org.apache.spark" %% "spark-core" % "3.+",
  "com.datastax.cassandra" % "cassandra-driver-core" % "latest.release"
)
```

Aikido reported:

* 3 Critical CVEs in Akka HTTP
* 2 High-severity vulnerabilities in Spark
* 1 Critical vulnerability in the Cassandra driver

After investigation, all turned out to be false positives. The scanner was checking against different versions than what was actually being used in production:

```
// What the scanner thought was running
akka-http 10.2.0   // Vulnerable
spark-core 3.0.0   // Vulnerable
cassandra-driver 4.0.0  // Vulnerable

// What was actually running in production
akka-http 10.2.10  // Secure
spark-core 3.3.2   // Secure
cassandra-driver 4.15.0  // Secure
```

## Solution: Add a lockfile <a href="#solution-add-a-lockfile" id="solution-add-a-lockfile"></a>

Use the [SBT Dependency Lock](https://github.com/stringbean/sbt-dependency-lock) plugin to generate a lockfile for your project.

**Step 1: Add the SBT Dependency Lock Plugin**

```
// In plugins.sbt
addSbtPlugin("software.purpledragon" % "sbt-dependency-lock" % "1.5.1")
```

**Step 2: Generate Lock Files**

```
// Run this command to resolve dependencies and generate lock files
sbt "dependencyLockWrite"
```

The resulting lock file (build.sbt.lock) explicitly defines all dependencies:

```
{
  "com.typesafe.akka:akka-http_2.13": "10.2.10",
  "org.apache.spark:spark-core_2.13": "3.3.2",
  "com.datastax.cassandra:cassandra-driver-core": "4.15.0"
}
```

**Step 3: Enforce Locked Dependencies**

```
// Run this command to resolve dependencies and validate against lockfile
sbt "dependencyLockCheck"
```

### Alternative Workaround: Container Scanning <a href="#alternative-workaround-container-scanning" id="alternative-workaround-container-scanning"></a>

While lock files provide excellent dependency control at the source level, container scanning offers another powerful approach to security validation. Since containers include the compiled artifacts, they represent the actual state of your production environment.

#### Benefits of Container Scanning <a href="#benefits-of-container-scanning" id="benefits-of-container-scanning"></a>

```
# Example Dockerfile showing what gets scanned
FROM openjdk:11-jre-slim

# Your compiled artifacts are here
COPY target/scala-2.13/your-app.jar /app/
COPY target/scala-2.13/lib/* /app/lib/

# These are the actual versions that will run in production
RUN ls -la /app/lib/

# akka-http_2.13-10.2.10.jar
# spark-core_2.13-3.3.2.jar
# cassandra-driver-core-4.15.0.jar
```

#### When to Use Container Scanning <a href="#when-to-use-container-scanning" id="when-to-use-container-scanning"></a>

Container scanning is particularly valuable when:

* You need to validate production-ready artifacts
* Your build process involves multiple stages
* You want to verify the exact versions running in production
* You need to scan both your application and its runtime environment

## Container Scanning vs. Lock Files <a href="#container-scanning-vs-lock-files" id="container-scanning-vs-lock-files"></a>

| **Aspect**     | **Container Scanning** | **Lock Files**       |
| -------------- | ---------------------- | -------------------- |
| When to verify | Post-build             | Pre-build            |
| What's checked | Compiled artifacts     | Source dependencies  |
| Accuracy       | Production-exact       | Development-exact    |
| Integration    | CI/CD pipeline         | Development workflow |

> **Best Practice**: Use both approaches. While container scanning is valuable, it's most effective when used alongside lock files.
