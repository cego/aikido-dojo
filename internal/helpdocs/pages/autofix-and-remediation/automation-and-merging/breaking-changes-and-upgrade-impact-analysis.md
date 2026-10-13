# Breaking changes & upgrade impact analysis

Assess breaking changes in dependency upgrades and their impact on your codebase.

## Breaking changes

Aikido makes an assessment on whether a version upgrade contains breaking changes by looking at the library changelogs.

Clicking on an open-source dependency issue in the Aikido feed pops up the issue details, showing the minimal version upgrade required to fix the issue. In the following Spring Security example, CVE-2023-34034 can be fixed by upgrading from version 6.1.0 to version 6.1.2. Aikido has determined that there are no breaking changes in this update, shown by the :white\_check\_mark:.

Clicking the :white\_check\_mark: reveals the details and a link to the changelogs. In this case there are no breaking changes.

If there are breaking changes like for the following Tomcat upgrade from version 8.5.11 to version 9.0.99, Aikido shows a :warning: icon.

Clicking the :warning: reveals a description of the breaking changes and a link to the changelog.

## Upgrade impact analysis

Even if there are breaking changes in the dependency update, this does not mean your code is affected by these changes. The breaking changes could be in a function that is not used in your codebase. Upgrade impact analysis goes one step further by scanning the codebase for usages of the library to determine whether the breaking changes affect your codebase.

Aikido AutoFix performs upgrade impact analysis when a pull request is created. The description of the pull request contains the analysis results, falling into 1 of 3 categories:

* No breaking changes in the library upgrade :white\_check\_mark:
* Breaking changes but they do not affect your codebase :white\_check\_mark:
* Breaking changes that affect your codebase, manual mitigation is required :warning:

In the following example mongoose is updated from version 5.13.21 to 6.13.6. The pull request shows there are 2 breaking changes affecting the codebase and outlines which files are affected and what the impact is:
