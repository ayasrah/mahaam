# App Structure

### Overview

This page discusses codebase template and folding options: feature-based and layer-based approaches.

### Importance

Having a template is important because it defines **what goes where**, making both writing and reading code easier. It saves developers from having to decide where to place each part of the code, and helps newcomers know where to find things, so its important for:

- Maintainability and readability
- Avoiding overengineering and spaghetti code

### Architecture Options

Its worthy to mention these common Architectures in this context:

- **Monolith App**: Codebase is split by layers under one project, and deployed as one unit.
- **Modular monolith App**: Codebase is splitted into modules with clear boundaries under a single project, and still deployed as a single unit.
- **Multi-project App**: Codebase is split into many projects under one solution.
- **Microservices App**: Codebase is spit into microservices, deployed separately, and communicates between using APIs or message queues.

### Mahaam App Structure

**1. By Feature**

- Simple architecture, Group by feature.

```bash
app-be
└── Src
    ├── Feat/
    │   ├── Plan/                   # Plan module
    │   │   ├── Plan.cs            	# Plan model
    │   │   ├── PlanController.cs  	# Plan APIs
    │   │   ├── PlanService.cs     	# Plan business logic
    │   │   ├── PlanRepo.cs        	# Plan DB operations
    │   │   └── PlanMembersRepo.cs  # PlanMembers DB operations
    │   ├── Task/                   # Task module
    │   │   ├── Task.cs
    │   │   ├── TaskController.cs
    │   │   ├── TaskService.cs
    │   │   └── TaskRepo.cs
    │   └── User/                    # User module
    │       ├── User.cs
    │       ├── UserController.cs
    │       ├── UserService.cs
    │       ├── UserRepo.cs
    │       ├── DeviceRepo.cs
    │       └── SuggestedEmailsRepo.cs
    ├── Infra/
    │   ...
    └── Program.cs
```

**2. By Layer (In Go implementation)**

```bash
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

### Infra Folder

Infra folder has utility classes that are used by the whole app, like `Cache`, `DB`, `Security`. Monitor functionality is organized in `Infra/Monitor` with controllers, services, and repositories for health checks, logging, and traffic monitoring.

```bash
app-be/
└── Src/
    ├── Feat/
    └── Infra/
        ├── Monitor/
        │   ├── AuditController.cs
        │   ├── Health.cs
        │   ├── HealthController.cs
        │   ├── HealthRepo.cs
        │   ├── HealthService.cs
        │   ├── LogRepo.cs
        │   ├── Models.cs
        │   └── TrafficRepo.cs
        ├── Cache.cs
        ├── Config.cs
        ├── DB.cs
        ├── Email.cs
        ├── Exceptions.cs
        ├── Factory.cs
        ├── Http.cs
        ├── Json.cs
        ├── Log.cs
        ├── Middlewares.cs
        ├── RequestData.cs
        ├── Security.cs
        ├── Starter.cs
        └── Validator.cs
```

### Go case

**Folding by feature in Go** is not straight forward while keeping interfaces in same file such what been did in `C#, Java, TS, and Python`, as it will cause **circular dependences issue** (eg: between user and plan modules), that's why mahaam chose to fold by layer for Go project and its good change to show the app in layer structure as well.

Another point in `mahaam-api-go` is placing every infra file in its own folder, and that is because in go each folder is a package, and its not allowed two files under same folder to have different packages. Mahaam needs to group and call infra utilities by their modules like: `cache.AppName`, `email.SendMeOtp`, `config.DBUrl`. Keeping all files flat under infra will not give us this, instead it will be `infra.AppName`, `infra.SendMeOtp`, `infra.DBUrl`, which mixes all functions together, that's why each file placed in its own folder.
