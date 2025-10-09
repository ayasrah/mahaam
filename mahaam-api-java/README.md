# mahaam-api-java

Java implementation of Mahaam backend services.

#### Setup

- Install [JDK](https://mahaam.dev/setup/creation#installing-the-runtime-sdk)
- Install [Quarkus cli](https://quarkus.io/get-started/)
- Install Postgres DB locally or on cloud.
- Create [Mahaam Database Schema](https://github.com/ayasrah/mahaam/blob/main/mahaam-data/mahaam_ddl.sql)
- Rename `.env.example` to `.env`.
- Update dbUrl to map to the new created DB.

#### Configs

Configure the application using `src/main/resources/application.properties` and `.env`:

- `tokenSecretKey`
  Generate and fill [API secret key](https://mahaam.dev/infra/security#generating-jwt-secret-key-signing-key)
- `OTP configs`
  In order to get OTP functionality works, either create a Twilio account with SendGrid service or fill emails you want to simulate in `testEmails`. Fill any value in `testSID`, eg: `2ad1a5c27c`, and any number in `testSID`, eg: `549023`

#### Structure

```bash
mahaam-api-java/
├── src/
│   └── main/
│       ├── java/mahaam/
│       │   ├── feat/               # Modules
│       │   │   ├── plan/           # Plan module
│       │   │   ├── task/           # Task module
│       │   │   └── user/           # User module
│       │   └── infra/              # Infra utils
│       └── resources/
│           └── application.properties # Configs
├── .env.example					# Configs
├── build.gradle                    # Gradle build configs
├── gradle.properties               # Gradle properties
├── gradlew                         # Gradle wrapper Unix
├── gradlew.bat                     # Gradle wrapper Windows
└── settings.gradle                 # Gradle settings
```

#### Build

```bash
./gradlew build
```

#### Run

```bash
./gradlew quarkusDev # Dev Mode with live reload
./gradlew quarkusRun # Standard mode
quarkus dev	     # using quarkus cli
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
  "apiName": "mahaam-api-java",
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
# Build production JAR
./gradlew build -Dquarkus.package.jar.type=uber-jar

# Run production JAR
java -jar build/quarkus-app/quarkus-run.jar

# Or run uber JAR
java -jar build/*-runner.jar
```
