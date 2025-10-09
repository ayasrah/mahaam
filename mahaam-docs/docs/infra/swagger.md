# Swagger

### Overview

Swagger is live docs for REST API services.

### Purpose

- Clear and live docs for APIs.
- Built-in lightweight testing via browser.
- API is self-documented.

### Mahaam Swagger

Swagger exposes the API docs at `/swagger` or `/docs` endpoints, eg, in `mahaam-api-cs` the docs endpoint is: `http://localhost:7023/mahaam-api/swagger`

#### Exposed Modules

Following are the main exposed modules for Mahaam:

<img src="/swagger1.png" alt="Swagger 1"  width="400" style="border-radius:5px;"/>

#### Plan Module APIs

Following are the exposed APIs for Plan module:

<img src="/swagger2.png" alt="Swagger 2"  width="400" style="border-radius:5px;"/>

#### Try POST a plan

You can try an endpoint just from the browser, following is POST a plan:

<img src="/swagger3.png" alt="Swagger 3"  width="400" style="border-radius:5px;"/>

### Configure

::: code-group

```C#
// In Program.cs
builder.Services.AddSwaggerGen(c =>
{
    c.SwaggerDoc("v1", new OpenApiInfo
    {
        Title = "Mahaam API",
        Version = "v1"
    });
});

app.UseSwagger();
app.UseSwaggerUI(c =>
{
    c.SwaggerEndpoint("/swagger/v1/swagger.json", "Mahaam API V1");
});
```

```Java
// Using Quarkus OpenAPI
@OpenAPIDefinition(
    info = @Info(
        title = "Mahaam API",
        version = "1.0.0"
    )
)
@ApplicationScoped
public class SwaggerConfig {
    // Annotations on controllers generate docs
}
```

```Go
// Using swaggo/swag
// @title Mahaam API
// @version 1.0
// @description Mahaam API documentation
// @host localhost:7023
// @BasePath /api

// @Summary Get all users
// @Description Retrieve list of all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} User
// @Router /users [get]
func GetUsers(c *gin.Context) {
    // handler implementation
}
```

```TypeScript
// Using @nestjs/swagger
import { DocumentBuilder, SwaggerModule } from "@nestjs/swagger";

const config = new DocumentBuilder()
  .setTitle("Mahaam API")
  .setDescription("Mahaam API documentation")
  .setVersion("1.0")
  .addBearerAuth()
  .build();

const document = SwaggerModule.createDocument(app, config);
SwaggerModule.setup("swagger", app, document);
```

```Python
# Using FastAPI (built-in Swagger)
from fastapi import FastAPI

app = FastAPI(
    title="Mahaam API",
    description="Mahaam API documentation",
    version="1.0.0",
    docs_url="/swagger",
    redoc_url="/redoc"
)

@app.get("/users", response_model=List[User])
async def get_users():
    """Get all users"""
    return users
```

:::

### Usage

::: code-group

```C#
// Controller with Swagger annotations
[ApiController]
[Route("api/[controller]")]
public class UsersController : ControllerBase
{
    [HttpGet]
    [ProducesResponseType(typeof(List<User>), 200)]
    public async Task<IActionResult> GetUsers()
    {
        // implementation
    }
}
```

```Java
// JAX-RS with OpenAPI annotations
@Path("/users")
public class UserController {
    @GET
    @Operation(summary = "Get all users")
    @APIResponse(responseCode = "200", description = "List of users")
    public List<User> getUsers() {
        // implementation
    }
}
```

```Go
// Gin with swag comments
// @Summary Get user by ID
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Router /users/{id} [get]
func GetUser(c *gin.Context) {
    // implementation
}
```

```TypeScript
// NestJS with decorators
@Controller("users")
@ApiTags("users")
export class UsersController {
  @Get()
  @ApiOperation({ summary: "Get all users" })
  @ApiResponse({ status: 200, type: [User] })
  async getUsers(): Promise<User[]> {
    // implementation
  }
}
```

```Python
# FastAPI with type hints
@app.get("/users/{user_id}", response_model=User)
async def get_user(
    user_id: int = Path(..., description="The ID of the user")
):
    """Get a specific user by ID"""
    # implementation
```

:::
