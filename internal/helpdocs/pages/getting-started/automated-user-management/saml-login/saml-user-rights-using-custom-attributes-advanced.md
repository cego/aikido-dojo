# SAML User Rights: Custom Attributes (Advanced)

> These are the advanced way of setting up user rights. We recommend using [SAML Access Profiles](/getting-started/automated-user-management/saml-login/saml-user-rights-access-profiles-recommended.md)
>
> [https://help.aikido.dev/doc/saml-user-rights-access-profiles-recommended/docVaVb0VPy1](/getting-started/automated-user-management/saml-login/saml-user-rights-access-profiles-recommended.md)

This guide provides detailed instructions on how to configure and manage user rights within Aikido using SAML custom attributes. By leveraging attributes such as `aikido_role`, `aikido_data_edit_rights`, `aikido_can_ignore`, `aikido_can_snooze`, `aikido_can_change_severity`, `aikido_can_manage_teams`, and `aikido_teams`, you can control user permissions and roles from within your identity provider. This approach ensures that users have the same access in Aikido as set up in your identity provider.

* **aikido\_access\_profile:** [**More info**](/getting-started/automated-user-management/saml-login/saml-user-rights-access-profiles-recommended.md)\
  When setting up SAML Access Profiles, this is the claim to use.

  ```xml
  <saml:Attribute Name="aikido_access_profile">
      <saml:AttributeValue xsi:type="xs:anyType">My Access Profile</saml:AttributeValue>
  </saml:Attribute>
  ```
* **aikido\_username:** You can define the name of the user in Aikido

  ```xml
  <saml:Attribute Name="aikido_username">
      <saml:AttributeValue xsi:type="xs:anyType">John Doe</saml:AttributeValue>
  </saml:Attribute>
  ```
* **aikido\_role:** `admin`, `default`, `team_only`

  ```xml
  <saml:Attribute Name="aikido_role">
      <saml:AttributeValue xsi:type="xs:anyType">default</saml:AttributeValue>
  </saml:Attribute>
  ```
* **aikido\_data\_edit\_rights:** `standard`, `read_only`

  ```xml
  <saml:Attribute Name="aikido_data_edit_rights">
      <saml:AttributeValue xsi:type="xs:anyType">standard</saml:AttributeValue>
  </saml:Attribute>
  ```
* **aikido\_can\_ignore:** `true`, `false`

  ```xml
  <saml:Attribute Name="aikido_can_ignore">
      <saml:AttributeValue xsi:type="xs:anyType">true</saml:AttributeValue>
  </saml:Attribute>
  ```
* **aikido\_can\_snooze:** `true`, `false`

  ```xml
  <saml:Attribute Name="aikido_can_snooze">
      <saml:AttributeValue xsi:type="xs:anyType">true</saml:AttributeValue>
  </saml:Attribute>
  ```
* **aikido\_can\_change\_severity:** `true`, `false`

  ```xml
  <saml:Attribute Name="aikido_can_change_severity">
      <saml:AttributeValue xsi:type="xs:anyType">true</saml:AttributeValue>
  </saml:Attribute>
  ```
* **aikido\_can\_manage\_teams:** `true`, `false`

  ```xml
  <saml:Attribute Name="aikido_can_manage_teams">
      <saml:AttributeValue xsi:type="xs:anyType">true</saml:AttributeValue>
  </saml:Attribute>
  ```
* **aikido\_can\_export\_data:** `true`, `false`

  ```xml
  <saml:Attribute Name="aikido_can_export_data">
      <saml:AttributeValue xsi:type="xs:anyType">true</saml:AttributeValue>
  </saml:Attribute>
  ```
* **aikido\_can\_manage\_clouds:** `true`, `false`

  ```xml
  <saml:Attribute Name="aikido_can_manage_clouds">
      <saml:AttributeValue xsi:type="xs:anyType">true</saml:AttributeValue>
  </saml:Attribute>
  ```
* **aikido\_can\_manage\_containers:** `true`, `false`

  ```xml
  <saml:Attribute Name="aikido_can_manage_containers">
      <saml:AttributeValue xsi:type="xs:anyType">true</saml:AttributeValue>
  </saml:Attribute>
  ```
* **aikido\_can\_manage\_domains:** `true`, `false`

  ```xml
  <saml:Attribute Name="aikido_can_manage_domains">
      <saml:AttributeValue xsi:type="xs:anyType">true</saml:AttributeValue>
  </saml:Attribute>
  ```
* **aikido\_can\_manage\_pentests:** `true`, `false`

  ```xml
  <saml:Attribute Name="aikido_can_manage_pentests">
      <saml:AttributeValue xsi:type="xs:anyType">true</saml:AttributeValue>
  </saml:Attribute>
  ```
* **aikido\_can\_manage\_code\_quality:** `true`, `false`

  <pre class="language-xml"><code class="lang-xml">&#x3C;saml:Attribute Name="aikido_can_manage_code_quality">
  <strong>    &#x3C;saml:AttributeValue xsi:type="xs:anyType">true&#x3C;/saml:AttributeValue>
  </strong>&#x3C;/saml:Attribute>
  </code></pre>
* **aikido\_can\_manage\_repos:** `true`, `false`

  <pre class="language-xml"><code class="lang-xml">&#x3C;saml:Attribute Name="aikido_can_manage_repos">
  <strong>    &#x3C;saml:AttributeValue xsi:type="xs:anyType">true&#x3C;/saml:AttributeValue>
  </strong>&#x3C;/saml:Attribute>
  </code></pre>
* **aikido\_teams:** You can define the different teams where the user is a part of here. If the team(s) do not exist in Aikido, it will be created. The user will auto-join these given teams. The user will be removed from all other teams if this is set up.

  ```xml
  <saml:Attribute Name="aikido_teams">
      <saml:AttributeValue xsi:type="xs:anyType">team1</saml:AttributeValue>
      <saml:AttributeValue xsi:type="xs:anyType">team2</saml:AttributeValue>
  </saml:Attribute>
  ```
* **aikido\_workspace\_ids:** You can define the different Aikido workspaces where the user is a part of here. The user will auto-join these given workspaces. The user will be removed from all other workspaces if this field is set up.

  ```xml
  <saml:Attribute Name="aikido_workspace_ids">
      <saml:AttributeValue xsi:type="xs:anyType">1233</saml:AttributeValue>
      <saml:AttributeValue xsi:type="xs:anyType">2511</saml:AttributeValue>
  </saml:Attribute>
  ```
* **github\_samlidentity\_nameid:** (github specific) Team member sync based on SAML users currently only works if users also log into GitHub via SAML. By sending the `github_samlidentity_nameid` SAML claim containing the user’s email address, triggering team-sync will also sync the SAML team members.

  ```xml
  <saml:Attribute Name="github_samlidentity_nameid">
      <saml:AttributeValue xsi:type="xs:anyType">example@email.com</saml:AttributeValue>
  </saml:Attribute>
  ```
* **github\_user\_slug:** (github specific) Team member sync based on SAML users currently only works if users also log into GitHub via SAML. By sending the `github_user_slug` SAML claim containing the user’s github slug, triggering team-sync will also sync the SAML team members.

  ```xml
  <saml:Attribute Name="github_user_slug">
      <saml:AttributeValue xsi:type="xs:anyType"johndoe</saml:AttributeValue>
  </saml:Attribute>
  ```
