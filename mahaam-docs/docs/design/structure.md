# App Structure

### Overview

This page discusses codebase template and folding options: feature-based and layer-based approaches.

### Importance

Having a template is important because it defines **what goes where**, making both writing and reading code easier. It saves developers from having to decide where to place each part of the code, and helps newcomers know where to find things, so its important for:

- Maintainability
- Readability
- Avoid overengineering
- Reduce spaghetti code
- Reduce technical debt

### Options

These are common archeticture options:

- **Monolith**: Modules are split by layers, and deployed as one unit.
- **Microservices**: Each module is implemented as a microservice, deployed separately, and communicates between using APIs or message queues.
- **Multi-project solution**: Modules are implemented as separate projects under one solution (common in Java and C#).
- **Modular monolith**: Modules are defined as folders with clear boundaries inside a single project. It is easier to maintain and deploy, and this is Mahaam’s choice.
- **Vertical slice**: micro-level splitting, each functionality has its own slice, e.g., `/Plan/CreatePlan.cs`, which contains both logic and database access for that functionality.

**1. By Feature (recommended)**

- Simple architecture, Group by feature, Modular monolith.

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

**2. By Layer**

```bash
app-be
└── Src/
    ├── Controller/
    │   ├── PlanController.cs      	# Plan APIs
    │   ├── TaskController.cs
    │   └── UserController.cs
    ├── Model/
    │   ├── Plan.cs               	# Plan models
    │   ├── Task.cs
    │   └── User.cs
    ├── Service/
    │   ├── PlanService.cs        	# Plan business logic
    │   ├── TaskService.cs
    │   └── UserService.cs
    ├── Repo/
    │   ├── PlanRepo.cs           	# Plan DB operations
    │   ├── PlanMembersRepo.cs     	# PlanMembers DB operations
    │   ├── TaskRepo.cs
    │   ├── UserRepo.cs
    │   ├── DeviceRepo.cs
    │   └── SuggestedEmailsRepo.cs
    ├── Infra/
    │   ...
    └── Program.cs
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

**Folding by feature in Go** is not straight forward while keeping interfaces in same file such what been did in `C#, Java, TS, and Python`, as it will cause **circular dependences issue**, thats why mahaam chose to fold by layer for Go project and its good change to show the app in layer structure as well.

Another point in `mahaam-api-go` is placing every infra file in its own folder, and that is because in go each folder is a package, and its not allowed two files under same folder to have different packages. Mahaam needs to group and call infra utilities by their modules like: `cache.AppName`, `email.SendMeOtp`, `config.DBUrl`. Keeping all files flat under infra will not give us this, instead it will be `infra.AppName`, `infra.SendMeOtp`, `infra.DBUrl`, which mixes all functions together, thats why each file placed in its own folder.
