# Service Creation

### Overview

This page shows how to create and run a **backend API service**.

### Steps

In general these are the steps to set up an API service:

- Install the runtime.
- Create or clone a project starter.
- Install the dependencies.

### Installing the runtime/SDK

The first step of developing backend services is to install the runtime/sdk of the language:

::: code-group

```bash [C#]
# The runtime/sdk is .NET, you can download from:
# https://dotnet.microsoft.com/en-us/download or using cli:
# Windows
winget install Microsoft.DotNet.SDK.8
# macOS
brew install dotnet
# Linux
apt install -y dotnet-sdk-8.0
```

```bash [Java]
# The runtime/sdk JRE or JDK, you can download from:
# https://adoptium.net/en-GB/temurin/releases or using cli:
# Windows
winget install Eclipse.Temurin.21.JDK
# macOS
brew install openjdk@21
# Linux
apt install openjdk-21-jdk
```

```bash [Go]
# The runtime/sdk is Go runtime, you can download from:
# https://go.dev/doc/install or using cli:
# Windows
winget install GoLang.Go
# macOS
brew install go
# Linux
apt install golang-go
```

```bash [TypeScript]
# The runtime/sdk is Node, you can download from:
# https://nodejs.org/en/download or using cli:
# Windows
winget install OpenJS.NodeJS
npm install -g pnpm
# macOS
brew install node pnpm
# Linux
curl -o- https://fnm.vercel.app/install | bash
fnm install 22
```

```bash [Python]
# The runtime/sdk is Python, you can download from:
# https://www.python.org/downloads/ or using cli:
# Windows
winget install Python.Python.3.12
pip install uv
# macOS
brew install python@3.12 uv
# Linux
apt install python3.12 python3-pip
pip install uv
```

:::

### Creating a project

The first option is to create a backend service project using the runtime/framework cli:

::: code-group

```bash [C#]
# Create directory
mkdir mahaam-api-cs
cd mahaam-api-cs
# Create .NET Service
dotnet new webapi -n mahaam
# Create solution file
dotnet new sln -n mahaam
# Add project to solution
dotnet sln add mahaam/mahaam.csproj
```

```bash [Java]
# Create new directory
mkdir mahaam-api-java
cd mahaam-api-java
# Initialize Quarkus project
sdk install quarkus
quarkus create app mahaam-api-java --gradle
```

```bash [Go]
# Create new directory
mkdir mahaam-api-go
cd mahaam-api-go
# Initialize Go module
go mod init mahaam-api-go
# Create main.go
touch main.go
# Install Gin framework
go get github.com/gin-gonic/gin
```

```bash [TypeScript]
# Create new directory
mkdir mahaam-api-ts
cd mahaam-api-ts
# Initialize NestJS project
npm i -g @nestjs/cli
nest new mahaam-api-ts --package-manager=pnpm
```

```bash [Python]
# Create new directory
mkdir mahaam-api-py
cd mahaam-api-py
# Initialize Python project with uv
uv init mahaam-api-py
cd mahaam-api-py
source venv/bin/activate
# Add FastAPI dependencies
uv add "fastapi[standard]"
```

:::

### Cloning a project template

Cloning a project starter template could be very helpful, and there are many available with predefined configurations. Even within a company, its a good practice to prepare a starter project that matches business needs and could be cloned by developers.

::: code-group

```bash [C#]
# Clone a .NET Web API template
# git clone <template-repo-url>
```

```bash [Java]
# Clone Quarkus starter template
git clone https://github.com/quarkusio/quarkus-quickstarts.git
cp -r quarkus-quickstarts/getting-started mahaam-api-java
cd mahaam-api-java
# Or use Quarkus CLI to create project
# quarkus create app mahaam-api-java --template=rest
```

```bash [Go]
# Clone Go web server template
git clone https://github.com/gin-gonic/examples.git
cp -r examples/gin-gonic/gin-basic mahaam-api-go
cd mahaam-api-go
# Or create manually with go mod init
# mkdir mahaam-api-go && cd mahaam-api-go
# go mod init mahaam-api-go
```

```bash [TypeScript]
git clone https://github.com/nestjs/typescript-starter.git mahaam-api-ts
cd mahaam-api-ts
pnpm install
pnpm run start
```

```bash [Python]
# Clone FastAPI starter template
git clone https://github.com/tiangolo/full-stack-fastapi-postgresql.git
cp -r full-stack-fastapi-postgresql/backend mahaam-api-py
cd mahaam-api-py
# Or create manually with FastAPI
# mkdir mahaam-api-py && cd mahaam-api-py
# uv init
```

:::

### Installing dependencies

In order to install the service dependencies, you can run:

::: code-group

```bash [C#]
# Restore NuGet packages
dotnet restore
# Add additional packages if needed
dotnet add package Dapper
dotnet add package Microsoft.Data.SqlClient
```

```bash [Java]
# Add following to build.gradle
# implementation 'io.quarkus:quarkus-rest'
# Then build and download dependencies
./gradlew build
# Or just download dependencies without building
./gradlew dependencies

```

```bash [Go]
# Download and tidy dependencies
go mod tidy
# Add additional dependencies if needed
go get github.com/gin-gonic/gin
```

```bash [TypeScript]
# Install dependencies
pnpm install
# Add additional dependencies if needed
pnpm add @nestjs/typeorm typeorm pg
pnpm add @nestjs/jwt @nestjs/passport passport passport-jwt
pnpm add class-validator class-transformer
```

```bash [Python]
# Install dependencies
uv sync
# Add additional dependencies if needed
uv add "fastapi[all]"
```

:::

### Building and running

You can build and run the service by running:

::: code-group

```bash [C#]
# Restore packages
dotnet restore
# Build the project
dotnet build
# Run the application
dotnet run
```

```bash [Java]
# Build the project
./gradlew build
# Run in development mode
./gradlew quarkusDev
# Or run the built JAR
java -jar build/quarkus-app/quarkus-run.jar
```

```bash [Go]
# Download dependencies
go mod tidy
# Build the project
go build -o mahaam-api-go
# Run the application
./mahaam-api-go
# Or run directly
go run main.go
```

```bash [TypeScript]
# Install dependencies
pnpm install
# Run in development mode
pnpm start:dev
# Build for production
pnpm build
# Run production build
pnpm start:prod
```

```bash [Python]
# Install dependencies (if not already done)
uv sync
# Run the application
uv run fastapi dev main.py
# Or with uvicorn directly
uv run uvicorn main:app --reload --port 8000
```

:::

### Access

Once running, the API will be available at the configured HTTP port and base path.

:::: code-group

```md [C#]
API: http://localhost:7023/mahaam-api
Swagger: http://localhost:7023/mahaam-api/swagger/index.html
```

```md [Java]
API: http://localhost:7023/mahaam-api
Swagger: http://localhost:7023/q/swagger-ui/
```

```md [Go]
API: http://localhost:7023/mahaam-api
Swagger: http://localhost:7023/swagger/index.html
```

```md [TypeScript]
API: http://localhost:7023/mahaam-api
Swagger: http://localhost:7023/mahaam-docs
```

```md [Python]
API: http://localhost:7023/mahaam-api
Docs: http://localhost:7023/mahaam-api/docs
```

::::
