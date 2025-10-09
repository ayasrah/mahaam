# mahaam-api-py

Python implementation of Mahaam backend services.

#### Setup

- Install [Python Runtime](https://mahaam.dev/setup/creation#installing-the-runtime-sdk)
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

```bash
mahaam-api-py/
├── feat/                   # Modules
│   ├── plan/               # Plan module
│   ├── task/               # Task module
│   └── user/               # User module
├── infra/                  # Infra utils
├── config.example.json     # Configs
├── main.py                 # Application entry point
├── pyproject.toml          # Python dependencies
└── uv.lock                 # Dependency lock file
```

#### Create virtual environment

```bash
uv venv
```

#### Activate virtual environment

```bash
# Linux/macOS
source .venv/bin/activate

# Fish shell
source .venv/bin/activate.fish

# Windows
.venv\Scripts\activate
```

#### Install dependencies

```bash
uv pip install -r pyproject.toml
```

#### Run

```bash
uv run uvicorn main:app --reload
# or
python3 main.py
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
  "apiName": "mahaam-api-py",
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
# Build production app
uv pip install -r pyproject.toml --no-dev

# Run production app
uv run uvicorn main:app --host 0.0.0.0 --port 7023

# Or run with gunicorn
gunicorn main:app -w 4 -k uvicorn.workers.UvicornWorker --bind 0.0.0.0:7023
```
