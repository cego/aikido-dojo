# Data Exposure Audit (DSPM)

What DSPM is, and how Aikido's Data Exposure Audit takes a code-based approach.

Data security posture management (DSPM) helps you understand where sensitive data lives, who can access it, how it's protected, and where it may be exposed.

Sensitive data doesn't stay in one place. It moves through APIs, databases, logs, queues, SaaS tools, warehouses, backups, and AI workflows, because your application moves it. Traditional DSPM finds where that data landed in storage. For developers, the harder question is often: what code path created the exposure?

## Use cases

* **Map sensitive data in code**: See how PII, secrets, tokens, and customer data move through your app
* **Catch exposure before it ships**: Find unsafe logging, over-broad APIs, and third-party sharing from source
* **Fix the source of the leak**: Get evidence and remediation in code, without connecting your databases
* **Support privacy and compliance reviews**: Tie data handling gaps back to concrete code paths

## Traditional DSPM vs code-based DSPM

Traditional DSPM scans the data estate (buckets, databases, warehouses, SaaS apps, backups) to classify sensitive data and check permissions. Useful for security teams mapping where data already sits.

Code-based DSPM **starts from source instead**. It reasons about how data enters, moves, is stored, and leaves the application, then points at the code change that reduces the risk, **without connecting to your datastores**.

|                   | Code-based DSPM                      | Traditional DSPM                  |
| ----------------- | ------------------------------------ | --------------------------------- |
| **Starts from**   | Application source code              | Storage and connected platforms   |
| **Core question** | What code path created the exposure? | Where does sensitive data live?   |
| **Access needed** | Repositories only                    | Datastores and cloud data systems |

Traditional DSPM finds sensitive data where it lives. Code-based DSPM finds how it gets exposed.

## Data Exposure Audit in Aikido

**Data Exposure Audit** is the data-focused audit type in [AI Code Analysis](/ai-code-analysis/ai-code-analysis.md). It is Aikido's code-based DSPM, alongside Code Security Audit, Cloud Security Audit, and Smart Contract Audit.

A Data Exposure Audit maps sensitive data handling in your repositories and returns findings with evidence and remediation in code. No database access required. To run one, see [Run a Data Exposure Audit](/ai-code-analysis/what-is-dspm/data-exposure-audit.md). For the finding classes in scope, see [What Data Exposure Audit Finds](/ai-code-analysis/what-is-dspm/what-data-exposure-audit-finds.md).

For a longer product perspective, see [What developers need to know about DSPM and data exposure in code](https://www.aikido.dev/blog/what-is-dspm).
