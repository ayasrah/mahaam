# mahaam-api-cs

C# implementation of Mahaam backend services.

#### Setup

- Install [.NET Runtime](https://mahaam.dev/setup/creation#installing-the-runtime-sdk)
- Install Postgres DB locally or on cloud.
- Create [Mahaam Database Schema](https://github.com/ayasrah/mahaam/blob/main/mahaam-data/mahaam_ddl.sql)
- Rename config.example.json to config.json
- Update dbUrl to map to the new created DB.

#### Configs

- `tokenSecretKey`
  Generate and fill [API secret key](https://mahaam.dev/infra/security#generating-jwt-secret-key-signing-key)
- `OTP configs`
  In order to get OTP functionality works, either create a Twilio account with SendGrid service or fill emails you want to simulate in `testEmails`. Fill any value in `testSID`, eg: `2ad1a5c27c`, and any number in `testSID`, eg: `549023`

#### Structure

```bash
mahaam-api-cs/
├── Src/
│   ├── Feat/               # Modules
│   │   ├── Plan/           # Plan module
│   │   ├── Task/           # Task module
│   │   └── User/           # User module
│   ├── Infra/              # Infra utils
│   └── Program.cs          # App entry point
├── config.example.json     # Configs
├── mahaam.csproj
└── mahaam.sln
```

#### Build

```bash
dotnet restore
dotnet build
```

#### Run

```bash
dotnet run
```

#### Test

**curl**

```bash
curl --location --request GET 'http://localhost:7023/mahaam-api/health'
```

Sample Response

```json
{
  "id": "d56c5e3c-27fa-4e8a-9a04-edaf5e545fec",
  "apiName": "mahaam-api-cs",
  "apiVersion": "1.0",
  "nodeIP": "192.168.100.22",
  "nodeName": "ayasrah-pc",
  "envName": "local"
}
```

**Postman**

Test using [Postman](https://mahaam.dev/test/test).

**Swagger**

Check [swagger docs](https://mahaam.dev/infra/swagger).

```
http://localhost:7023/mahaam-api/swagger/index.html
```

#### Production

```bash
# Build and run
dotnet build --configuration Release
dotnet run --configuration Release

# Or run published app
dotnet publish --configuration Release
cd bin/Release/net8.0/publish
dotnet mahaam.dll
```
