# Support for Specific Branch Scanning

### Introduction <a href="#introduction" id="introduction"></a>

Utilize Aikido's manual scanning feature to directly compare a branch with your current project state/main branch, ideal for situations where you have not integrated your CI with Aikido yet.

This functionality allows you compare a specific branch / pull request / tag with the current state in Aikido, highlighting the changes with the main branch. Scanning takes typically **1-2mins** so no need to leave the UI.

**Note.** In case you have old legacy branches that needs daily scanning (eg branch V3, branch V4), check out our [Multi-Branch Scanning Feature](https://help.aikido.dev/en/articles/8979512-how-to-set-up-multi-branch-scanning).

### How to scan a specific branch <a href="#how-to-scan-a-specific-branch" id="how-to-scan-a-specific-branch"></a>

**Step 1**: Navigate to a specific repository detail page within Aikido

**Step 2**: Click on the Scan Branch button.\
​

[image: Repository dashboard showing critical security issues and scan options for the "terragoat" project.]

**Step 3:** In the prompted field, enter the **name of the branch or tag** you wish to scan. Make sure to type the exact name of the branch / tag to avoid any errors. You can select which types of scans to execute.

[image: Initiate a security scan for a specified code branch with scan type options.]

**Step 4.** In the bottom right corner, you will be able to follow the progress of the scanning. Once the scanning is done, click 'View Diffs'\
​

[image: Scan completed notification with "View Diffs" button displayed for a demo application.]

**Step 5.** Check which new issues are introduced and resolved on the comparison page.\
​

[image: Branch security scan summary showing newly introduced critical and high severity vulnerabilities.]

By following these steps, you can effectively conduct a manual scan of your branch/tag in Aikido, which will allow you to review changes, identify new or resolved issues, and make informed decisions about integrating the branch into your main codebase.

***
