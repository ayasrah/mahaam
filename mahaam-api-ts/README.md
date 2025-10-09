# mahaam-api-ts

TypeScript implementation of Mahaam backend services.

#### Setup

- Install [Node.js Runtime](https://mahaam.dev/setup/creation#installing-the-runtime-sdk)
- Install Postgres DB locally or on cloud.
- Create [Mahaam Database Schema](https://github.com/ayasrah/mahaam/blob/main/mahaam-data/mahaam_ddl.sql)
- Rename `.env.example` to `.env`
- Update dbUrl to map to the new created DB.

#### Configs

Configure the application using `.env` and `src/config/`:

- `tokenSecretKey`
  Generate and fill [API secret key](https://mahaam.dev/infra/security#generating-jwt-secret-key-signing-key)
- `OTP configs`
  In order to get OTP functionality works, either create a Twilio account with SendGrid service or fill emails you want to simulate in `testEmails`. Fill any value in `testSID`, eg: `2ad1a5c27c`, and any number in `testSID`, eg: `549023`

#### Structure

```bash
mahaam-api-ts/
├── src/
│   ├── feat/                   # Modules
│   │   ├── plans/              # Plan module
│   │   ├── tasks/              # Task module
│   │   └── users/              # User module
│   ├── infra/                  # Infra utils
│   └── main.ts                 # App entry point
├── .env.example                # Configs
├── package.json                # App manifest
└── tsconfig.json               # TypeScript configs
```

#### Build

```bash
pnpm install
pnpm run build
```

#### Run

```bash
pnpm run start:dev
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
  "apiName": "mahaam-api-ts",
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
# Run with PM2
pm2 start dist/main.js --name mahaam-api-ts
```
