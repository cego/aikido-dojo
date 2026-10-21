# Add Custom Code Rules

Custom Code Rules allow you to enforce team-specific coding standards and patterns unique to your organization. Create rules that detect specific code patterns, enforce naming conventions, or flag architectural violations that aren't covered by default checks.

### What are Custom Code Rules?

Custom Code Rules are AI-powered checks that you define to match your team's specific requirements. Unlike default checks that apply general best practices, custom rules enforce standards unique to your codebase, architecture, or business logic.

{% hint style="info" %}
If you want to implement global or repo-specific guidelines, check out our [Code Context](/code-quality/add-extra-code-context.md) functionality.
{% endhint %}

### Creating a Custom Code Rule

{% stepper %}
{% step %}
**Define your rule**

1. Navigate to [**Code Quality** > **Checks**](https://app.aikido.dev/code-quality/checks) tab
2. Click **Add Custom Code Rule**
3. Write a clear description of what the rule should detect:

Example:

```
Allow only Alpine base images in Dockerfiles as base images
```

{% endstep %}

{% step %}
**Select target languages**

Choose which programming languages this rule applies to
{% endstep %}

{% step %}
**Generate and refine examples**

* Click **Generate Examples** to let AI create initial code samples
* Review and modify the generated examples
* Provide both compliant and non-compliant examples:

**Compliant example:**

```dockerfile
# ✅ Code that follows the rule
FROM alpine:3.18
RUN apk add --no-cache nodejs
```

**Non-compliant example:**

```dockerfile
# ❌ Code that violates the rule
FROM ubuntu:latest
RUN apt-get update && apt-get install nodejs
```

{% endstep %}

{% step %}
**Validate your rule**

1. Click **Validate Rule** to test your examples
2. The system will verify that:
   * Your compliant examples pass the rule
   * Your non-compliant examples are correctly flagged
   * The rule logic is consistent and clear
3. Adjust examples if validation fails
   {% endstep %}

{% step %}
**Configure rule details**

Once validated, provide additional information about your rule:

1. **Title:** Give your rule a clear, descriptive name

   ```
   Use Alpine base images in Docker containers
   ```
2. **TL;DR:** Write a brief summary of the issue

   ```
   Non-Alpine base images increase container size and attack surface
   ```
3. **How to fix:** Provide actionable guidance for developers

   <pre data-overflow="wrap"><code>Replace your base image with an Alpine Linux variant. For example, change 'FROM node:18' to 'FROM node:18-alpine'. You may need to adjust package installation commands from apt-get to apk.
   </code></pre>
4. Click **Save Rule** to activate it

Your custom rule will now appear in the Checks tab and begin scanning new pull requests in enabled repositories.
{% endstep %}
{% endstepper %}

### Writing effective Custom Rules

#### Be specific and clear

❌ **Too vague:**

```
Use proper error handling
```

✅ **Specific and actionable:**

```
All API endpoints must wrap database calls in try-catch blocks and return 
standardized error responses with status codes
```

#### Focus on patterns, not style

Custom rules work best for detecting logical patterns rather than formatting:

✅ **Good custom rule candidates:**

* API authentication requirements
* Database transaction patterns
* Security header implementations
* Business logic validations

❌ **Better handled by linters:**

* Indentation and spacing
* Bracket placement
* Variable naming style

### Common Custom Rule examples

#### Security rules

```
All SQL queries must use parameterized statements. Direct string 
concatenation in SQL queries is not allowed.
```

#### Architecture rules

```
Controllers should not directly access the database. All database 
operations must go through a service or repository layer.
```

#### API standards

```
All REST API endpoints must include rate limiting headers 
(X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset).
```

#### Testing requirements

```
Every exported function must have at least one corresponding test 
in the __tests__ directory with the same file name pattern.
```

#### Documentation standards

```
All public API methods must include JSDoc comments with @param, 
@returns, and @throws annotations.
```
