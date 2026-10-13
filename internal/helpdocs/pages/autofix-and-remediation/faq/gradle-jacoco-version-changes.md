# Gradle - Jacoco version changes

**TL;DR** Make sure to commit `gradlew` to your repository for AutoFix to use the correct Gradle version.

If no `gradlew` script is present in the repo, AutoFix will us Gradle 8.14.4. Most likely this is not the Gradle version required for the repository. This can result in random version bumps in libraries like Jacoco.

Use the following commands to add `gradlew` to your repository:

```
## Navigate to your repo
cd <repository>

## Generate the gradlew wrapper
gradle wrapper

## Check that the gradlew wrapper is created and make it executable
ls gradlew
chmod +x gradlew

## Commit the wrapper to the repo
git add gradlew gradlew.bat gradle/wrapper/gradle-wrapper.jar gradle/wrapper/gradle-wrapper.properties
git commit -m 'Add the gradle wrapper'
```

See the [gradle wrapper documentation](https://docs.gradle.org/current/userguide/gradle_wrapper.html) for more information.
