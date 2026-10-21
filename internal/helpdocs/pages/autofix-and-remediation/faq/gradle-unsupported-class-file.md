# Gradle - unsupported class file

When the Java toolchain version is not defined in the `build.gradle`, AutoFix might use the wrong Java version resulting in the following error:

{% code title="AutoFix Output" %}

```
Unsupported class file major version xx
```

{% endcode %}

This error can be solved by specifying the Java version in the `build.gradle` (or `build.gradle.kts`) file by adding:

{% code title="" %}

```
java {
    toolchain {
        languageVersion = JavaLanguageVersion.of(21)
    }
}
```

{% endcode %}

**Note:** Make sure to replace `21` with the actual Java version your project requires.
