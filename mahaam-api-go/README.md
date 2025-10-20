# mahaam-api-go

Go implementation of Mahaam backend services.

#### Setup

- Install [Go Runtime](https://mahaam.dev/setup/creation#installing-the-runtime-sdk)
- Install Postgres DB locally or on cloud.
- Create [Mahaam Database Schema](https://github.com/ayasrah/mahaam/blob/main/mahaam-data/mahaam_ddl.sql)
- Rename config.example.json to config.json
- Update dbUrl to map to the new created DB.

#### Configs

Configure the application using `config.json`:

- `tokenSecretKey`
  Generate and fill [API secret key](https://mahaam.dev/infra/security#generating-jwt-secret-key-signing-key)
- `OTP configs`
  In order to get OTP functionality works, either create a Twilio account with SendGrid service or fill emails you want to simulate in `testEmails`. Fill any value in `testSID`, eg: `2ad1a5c27c`, and any number in `testSID`, eg: `549023`

#### Structure

```py
mahaam-api-go/
├── app/                      # Layered modules
│   ├── handler/              # HTTP handlers
│   │   ├── plan.go
│   │   ├── task.go
│   │   ├── user.go
│   │   └── ...
│   ├── models/               # Data models
│   │   ├── plan.go
│   │   ├── task.go
│   │   ├── user.go
│   │   └── ...
│   ├── repo/                 # Repo layer
│   │   ├── plan.go
│   │   ├── task.go
│   │   ├── user.go
│   │   └── ...
│   └── service/              # Business logic layer
│   │   ├── plan.go
│   │   ├── task.go
│   │   ├── user.go
│       └── ...
├── utils/                   # App utils
│   ├── conf/                # Configs
│   ├── email/               # Email service
│   ├── log/                 # Logging service
│   ├── middleware/          # HTTP middlewares
│   └── token/            	 # Token utils
├── config.example.json      # Sample config
├── go.mod
├── go.sum
├── main.go
└── README.md
```

#### Build

```bash
go mod download
go build
```

#### Run

```bash
go run .
```

#### Test

**curl**

```bash
curl --location --request GET 'http://localhost:7023/mahaam-api/health'
```

Sample Response

```json
{
  "id": "a26e78ae-db0b-467d-b1b6-fdce3deca4a5",
  "apiName": "mahaam-api-go",
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
# Build production binary
go build -ldflags="-w -s" -o bin/mahaam-api main.go

# Run production binary
./bin/mahaam-api

# Cross-platform builds
GOOS=linux GOARCH=amd64 go build -o bin/mahaam-api-linux main.go
GOOS=windows GOARCH=amd64 go build -o bin/mahaam-api.exe main.go
GOOS=darwin GOARCH=amd64 go build -o bin/mahaam-api-darwin main.go
```
