# Configs

### Overview

Configs are the settings that control the app.

### Stored In

- OS Environment variables
- Vaults
- .env files
- JSON files

### Purpose

- Maintainability: Decouple code from configs.
- Flexibility: Run same code on different environments.
- Scalability: Update configs without changing code.
- Security: Sensitive data should not be hard coded.

### Sample

- Environment (prod, dev, staging)
- Database url,
- log format,
- external service url

### Store

```env
appName=mahaam-api
appEnv=development
dbUrl=db-url-here
httpPort=7023
```

### Model

App should load configs on startup, and either keeps watching the file for changes (recommended) or restart after change.

Mahaam expose these configs from `infra/configs`:

::: code-group

```C#

public class Settings
{
	public ApiSettings Api { get; set; }
	public EmailSettings Email { get; set; }
	public LoggingSettings Logging { get; set; }
	public string DbUrl { get; set; }
	public bool LogReqEnabled { get; set; }
}

public class LoggingSettings
{
	public string File { get; set; }
	public int FileSizeLimit { get; set; }
	public int FileCountLimit { get; set; }
	public string OutputTemplate { get; set; }
}

public class EmailSettings
{
	public string AccountSid { get; set; }
	public string VerificationServiceSid { get; set; }
	public string AuthToken { get; set; }
	public List<string> TestEmails { get; set; }
	public string TestSID { get; set; }
	public string TestOTP { get; set; }
}


public class ApiSettings
{
	public string Name { get; set; }
	public string Version { get; set; }
	public string EnvName { get; set; }
	public int HttpPort { get; set; }
	public string TokenSecretKey { get; set; }
}
```

```Java
@StaticInitSafe
@ConfigMapping(prefix = "mahaam", namingStrategy = ConfigMapping.NamingStrategy.VERBATIM)
public interface Config {
	String apiName();
	String apiVersion();
	String envName();
	String dbUrl();
	String tokenSecretKey();
	String emailAccountSid();
	String emailVerificationServiceSid();
	String emailAuthToken();
	List<String> testEmails();
	String testSID();
	String testOTP();
	Boolean logReqEnabled();
}
```

```Go
type Conf struct {
	ApiName                     string
	ApiVersion                  string
	EnvName                     string
	DBUrl                       string
	LogFile                     string
	LogFileSizeLimit            int
	LogFileCountLimit           int
	LogFileOutputTemplate       string
	LogFileRollingInterval      string
	HTTPPort                    int
	TokenSecretKey              string
	EmailAccountSID             string
	EmailVerificationServiceSID string
	EmailAuthToken              string
	TestEmails                  []string
	TestSID                     string
	TestOTP                     string
	LogReqEnabled               bool
}
```

```TypeScript
const config = {
  get dbUrl(): string {
    return process.env.dbUrl || "";
  },

  get httpPort(): number {
    return parseInt(process.env.httpPort!, 10);
  },
} as const;
```

```Python
class Config(BaseModel):
    dbUrl: str
    httpPort: int
```

:::

### Usage

Mahaam expose configs through single utility as follows
::: code-group

```C#
var cnn = new NpgsqlConnection(settings.DbUrl); // DB conn
builder.WebHost.UseKestrel(opts => {
    opts.Listen(IPAddress.Parse("0.0.0.0"), settings.Api.HttpPort); // Service port
});
```

```Java
if (config.testEmails().contains(email)) {
	verifySid = Config.testSID;
} else {
	verifySid = Email.sendOtp(email);
}
```

```Go
db, err := sqlx.Connect("postgres", configs.DBUrl) // DB conn
	srv := &http.Server{
		Addr: configs.HTTPPort, // Service port
	}
```

```TypeScript
sql = postgres(config.dbUrl); // DB conn
app.listen(config.httpPort); // Service port
```

```Python
DB._engine = create_engine(configs.data.dbUrl) # DB conn
uvicorn.run(app, host="0.0.0.0", port=configs.data.httpPort)
```

:::

### See

- Configuration implementation in: [C#](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-cs/Src/Infra/Settings.cs), [Java](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-java/src/main/java/mahaam/infra/Config.java), [Go](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-go/utils/conf/env.go), [TypeScript](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-ts/src/infra/config.ts), [Python](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-py/infra/configs.py)
