# Introduction to Reachability Analysis

Learn how Aikido helps you identify which vulnerabilities are exploitable.

Aikido’s reachability analysis is focused on identifying whether a vulnerability is actually exploitable in the context of your application. Instead of flagging every vulnerable dependency or risky pattern as equally critical, Aikido builds a graph of modules, functions, and data flows within your codebase. It then asks a key question: ***Is there an execution path from real application behavior to the vulnerable code?*** Only when such a path exists is the finding considered truly security-relevant.

### Conceptual Overview

At a high level, Aikido constructs an abstract program graph in which:

* **Nodes** represent functions, methods, modules, files, and sometimes higher-level components.
* **Edges** represent relations such as function calls, imports, dependency inclusion, and data-flow between variables and parameters.

Vulnerabilities are anchored at specific nodes in this graph (for example, a known vulnerable function in a third-party library, or a sink such as a SQL execution point). Aikido then computes whether there exists a path from **entry points** (HTTP handlers, CLI commands, background jobs, etc.) to these vulnerable nodes. If no such path exists under the current code and configuration, the vulnerability is treated as non-reachable for that project revision.

In this way, reachability transforms a flat list of alerts into a structured subset of issues that are provably connected to real program behavior.

### Types of Reachability

Aikido uses several complementary notions of reachability. They share the same foundation, but operate at different levels of abstraction:

| Type                     | Scope                       | Main Question                                            | Typical Signals                                    |
| ------------------------ | --------------------------- | -------------------------------------------------------- | -------------------------------------------------- |
| **Dependency-level**     | SCA / package dependencies  | “Does the project ever call into the vulnerable symbol?” | Imports, symbol references, dependency manifests   |
| **Function-level**       | SAST / taint analysis       | “Can untrusted input flow into a dangerous sink?”        | Sources, sanitizers, sinks, inter-procedural flows |
| **Contextual (runtime)** | Build & runtime environment | “Does this code execute in the actual runtime context?”  | Prod vs dev deps, dead code, build-only tooling    |

These categories are not mutually exclusive. A typical Aikido finding may be filtered first by dependency-level reachability, then refined by function-level reachability, and finally reinterpreted in the context of how the project is built and deployed.

### How Reachability Works

#### 1. Graph construction

For each project, Aikido constructs a lightweight, language- and framework-aware graph:

* **Structural relations**: imports, requires, module references, inheritance, and other static relationships.
* **Call relations**: function calls, method invocations, event handlers.
* **Data-flow relations**: propagation of potentially tainted data through variables, parameters, and return values.

This graph is intentionally approximate but designed to preserve the properties needed to answer “is there a path?” questions efficiently.

#### 2. Anchoring vulnerabilities

The scanners then map each raw finding to one or more *anchor points* in the graph:

* For **SCA**, the anchor is usually a vulnerable symbol (a specific function, class, method, or module) associated with a CVE or advisory.
* For **SAST**, the anchor is a code location of interest, often a security sink such as `exec`, SQL queries, file system access, or deserialization.
* For **infrastructure scans**, anchor points may correspond to images, packages, or configuration elements that participate in the runtime execution environment.

#### 3. Reachability queries

Given a set of entry points and anchor nodes, Aikido asks reachability questions of the form:

> “Is there at least one valid execution path from any entry point to this anchor, under the current build and configuration?”

Algorithmically, this reduces to graph reachability (forward or backward traversal) constrained by language semantics and configuration knowledge. For data-flow analysis, the reachability relation is enriched with t**aint information**: a path is only considered security-relevant if untrusted data is able to traverse it without being fully neutralized by sanitization.

#### 4. Integration with triage and prioritization

Reachability is applied *before* higher-level triage:

* Findings with **no reachable path** are removed or heavily downgraded.
* Findings with a **clear, explainable path** are retained and passed to later stages (e.g., risk scoring, business impact analysis).
* Developers can inspect the concrete path (call stack and data-flow) to understand how the issue becomes exploitable.

This pipeline ensures that most “cry-wolf” findings are filtered out at the structural level rather than simply hidden behind severity tuning.

### Effects on Findings and Workflow

Reachability analysis has several practical consequences for how teams experience Aikido:

1. **Noise reduction**

   Many alerts produced by traditional tools stem from dependencies that are present but never invoked, or from code paths that are effectively dead. By eliminating issues without a demonstrable execution path, Aikido substantially reduces false positives and redundant tickets.
2. **More meaningful severity**

   Severity is no longer just a property of the CVE or rule; it is a property of the CVE *plus* the project’s usage. A medium-severity library issue that is deeply embedded on a hot path may be more urgent than a high-severity issue in a dev-only tool that is never executed in production.
3. **Developer-aligned explanations**

   Showing the concrete path (“request handler → service → library function → vulnerable symbol”) reframes security findings as recognizable program behavior. This lowers the cognitive overhead for developers and makes remediation efforts more targeted and efficient.
4. **Stable signal over time**

   Because reachability is recomputed as code and dependencies change, the issue backlog tracks actual architectural evolution. As dead dependencies are removed or code is refactored, previously reachable findings may transition to non-reachable, and vice versa, providing a dynamic, architecture-aware view of risk.

### Limitations and Design Choices

Aikido’s reachability engine uses a ‘conservative **under-approximation**,’ so it only marks code as unreachable when it can be determined with high confidence. Dynamic language features like reflection, dynamic imports, metaprogramming, or complex build steps can make parts of the call graph unclear or invisible. In these ambiguous cases, Aikido errs on the side of caution: instead of suppressing a potentially relevant issue, it will either keep it as-is or lower its severity.

This trade-off strikes a balance between two goals:

* Reducing noise by filtering out issues that are definitely not reachable at runtime.
* Preventing a false sense of security when reachability can’t be determined with certainty.

### Summary

Reachability analysis in Aikido maps how different parts of your application are connected by tracing function calls, data flows, and interactions between components across source code, dependencies, and infrastructure. A vulnerability is only treated as a real risk if there’s a clear, executable path from an application entry point to the vulnerable behavior. By embedding this reasoning into every scanner and integrating it tightly with triage, Aikido transforms a large, noisy stream of raw alerts into a smaller, high-confidence set of findings that more accurately reflect how attackers could exploit the system in practice.
