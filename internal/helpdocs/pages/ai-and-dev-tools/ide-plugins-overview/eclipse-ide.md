# Eclipse IDE

Aikido automatically scans your projects for hardcoded secrets (API keys, tokens) and insecure code patterns (SQL injections, path traversal, ..) so you can catch issues early and keep your codebase safe.

Scans run automatically whenever you open a file or save changes, making it easy to catch issues early in development.

When security issues are found, they're highlighted directly in your code and listed in the Aikido window.

## Installation and Authentication

{% stepper %}
{% step %}
**Go to the** [**Eclipse marketplace**](https://marketplace.eclipse.org/content/aikido-security) **and install the Aikido Plugin**

In Eclipse, go to `Help` > `Eclipse Marketplace...` and look for 'Aikid&#x6F;*'* in the search bar. Click 'Install'.

After installation, restart Eclipse.
{% endstep %}

{% step %}
**Authenticate with Aikido**

In Aikido, go to the [Eclipse Integration Screen](https://app.aikido.dev/settings/integrations/ide/eclipse) and create your token.

After installation the Aikido Security View will open automatically in Eclipse. If this is not the case, go to `Window` > `Show View` > `Other...` > Look for Aikido and select the 'Aikido Security' view. In the Aikido Security view, click `Connect to Aikido` and enter your token.

{% endstep %}

{% step %}
**Try out an example**

Below you can find an example `sample.java` file that can be used to verify if the extension is working correctly, it should detect one SAST issue on line 29 and one exposed secret on line 23.

```java
import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.ResultSet;
import java.sql.Statement;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import com.google.gson.Gson;
import java.io.IOException;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class UserServlet extends HttpServlet {
    
    @Override
    protected void doGet(HttpServletRequest req, HttpServletResponse res) 
            throws IOException {
        String connStr = "jdbc:sqlserver://myserver.database.windows.net:1433;" +
                        "database=mydb;" +
                        "user=myuser;" +
                        "password=$uperSecret123!@#;" +
                        "encrypt=true;" +
                        "trustServerCertificate=false;" +
                        "loginTimeout=30;";
        
        String username = req.getParameter("username");
        String unsafeQuery = "SELECT * FROM users WHERE username = '" + username + "'";
        
        try (Connection conn = DriverManager.getConnection(connStr);
             Statement stmt = conn.createStatement();
             ResultSet result = stmt.executeQuery(unsafeQuery)) {
            
            // Convert ResultSet to JSON-like structure
            List<Map<String, Object>> resultList = new ArrayList<>();
            int columnCount = result.getMetaData().getColumnCount();
            
            while (result.next()) {
                Map<String, Object> row = new HashMap<>();
                for (int i = 1; i <= columnCount; i++) {
                    row.put(result.getMetaData().getColumnName(i), 
                           result.getObject(i));
                }
                resultList.add(row);
            }
            
            res.setStatus(200);
            res.setContentType("application/json");
            res.getWriter().write(new Gson().toJson(resultList));
            
        } catch (Exception e) {
            res.setStatus(500);
            res.getWriter().write("Error: " + e.getMessage());
        }
    }
}
```

{% endstep %}
{% endstepper %}

### Uninstall

To uninstall the Aikido Eclipse plugin:

Go to `Help` > `Eclipse Marketplace...` > `Installed` > Look for Aikido Security and click `Uninstall`
