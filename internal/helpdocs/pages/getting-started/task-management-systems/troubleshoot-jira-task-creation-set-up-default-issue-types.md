# Troubleshoot Jira Task Creation: Set Up Default Issue Types

Task creation inside Aikido can fail, **particularly when the issue types that have been setup in Jira include required fields that are not available in Aikido**. The fix is to create a new issue type in Jira without any required fields, and inform Aikido about this new name.

### 1. Create a new Issue Type in JIRA <a href="#id-1-create-a-new-issue-type-in-jira" id="id-1-create-a-new-issue-type-in-jira"></a>

**Step 1:** Log in to JIRA and go to Settings > Issues > **Issue Types**

**Step 2: Add a New Issue Type**

Click *Add Issue Type*. Name it something indicative like **Security Fix**, and ensure it does not mandate any field that previously hindered Aikido task creation. This issue type must not include any required fields that Aikido cannot sync.

### 2. Configuring Aikido to Use the New Issue Type <a href="#id-2-configuring-aikido-to-use-the-new-issue-type" id="id-2-configuring-aikido-to-use-the-new-issue-type"></a>

**Step 1:** Go to [**JIRA Integration Settings**](https://app.aikido.dev/settings/integrations/tasktracker) in Aikido

**Step 2: Input the New Issue Type Name** in the default task type entry.

Enter the name of the issue type you created in JIRA (e.g., Security Fix). This links Aikido's ticket creation process with the new JIRA issue type.

[image: Set a default task type to "Task 2" for automatic task creation.]

***
