# Configs

#### Overview

Configs are the settings that control the app.

#### Stored In

- OS Environment variables
- Vaults
- .env files
- JSON files

#### Purpose

- Maintainability: Decouple code from configs.
- Flexibility: Run same code on different environments.
- Scalability: Update configs without changing code.
- Security: Sensitive data should not be hard coded.

#### Sample

- Environment (prod, dev, staging)
- Database url,
- log format,
- external service url

#### Store

```env
appName=mahaam-api
appEnv=development
dbUrl=db-url-here
httpPort=7023
```

#### Model

App should load configs on startup, and either keeps watching the file for changes (recommended) or restart after change.

Mahaam expose these configs from `infra/configs`:

::: code-group

```C#

public class Config
{
    public static string DbUrl => GetValue("dbUrl");
    public static string HttpPort => GetValue("httpPort");

    private static string GetValue(string key)
    {
        return _configuration[key] ?? throw new ArgumentException($"Config key '{key}' not found.");
    }
}
```

```Java
public class Config {
    public static final String dbUrl = ConfigProvider.getConfig().getValue("dbUrl", String.class);
    public static final String httpPort = ConfigProvider.getConfig().getValue("httpPort", String.class);
}
```

```Go
type config struct {
    DBUrl          string
    HTTPPort       int
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

#### Usage

Mahaam expose configs through single utility as follows
::: code-group

```C#
var cnn = new NpgsqlConnection(Config.DbUrl); // DB conn
builder.WebHost.UseKestrel(opts => {
    opts.Listen(IPAddress.Parse("0.0.0.0"), Config.port); // Service port
});
```

```Java
if (Config.testEmails.contains(email)) {
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
