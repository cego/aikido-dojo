# Drata Integration

The Drata integration automatically pushes Evidence to Drata for SOC2 and ISO27001:2022. Activate Drata from the Aikido Integrations page.

### API Key <a href="#api-key" id="api-key"></a>

You'll need to create a Drata API Key. You can generate such key in Drata via [Username > Settings > API Keys > Create API Key](https://app.drata.com/account-settings/api-keys/add).

#### Details <a href="#details" id="details"></a>

Make sure the Expiration is set to 'Never Expires'.

[image: API key setup screen with fields for name, expiration, and allowed IP addresses.]

#### Scopes <a href="#scopes" id="scopes"></a>

The access for the scopes can be set to 'Custom' with at least following scopes.

**Controls:**

* **List controls:** *Read*
* **Create control:** *Write*
* **Map external evidence:** *Read, Write*
* **Delete mapped external evidence:** *Write*

**Workspaces:**

* **List workspaces:** *Read*

**Frameworks:**

* **List frameworks:** *Read*
* **List framework requirements:** *Read*

[image: User permissions matrix for controls, workspaces, and frameworks management.]

#### Save <a href="#save" id="save"></a>

Next, click 'Save' and copy your generated API Key.

Back in Aikido, paste the API Key and click 'Next'. After that, choose your Drata workspace and click 'Save'.

[image: Enter your Drata API Key to integrate data and proceed to the next step.]

### Done <a href="#done" id="done"></a>

Aikido will now daily create a PDF report and sync this as 'external evidence' to Drata. We'll create a control with code 'AIKIDO' and link the relevant SOC2 and ISO27001 requirements. You can search for this control [here](https://app.drata.com/compliance/controls/inscope?q=Aikido).

Under 'Control evidence', the Aikido PDF will be attached every month

#### &#x20;<a href="#set-up-drata-integration" id="set-up-drata-integration"></a>
