# Request Context

### Overview

Request context is a thread-safe, request-scoped storage that maintains data for a single HTTP request.

### Sample

Mahaam request context includes:

- **trafficId**: Request tracking ID
- **userId**: Current user ID
- **deviceId**: Client device ID
- **appStore/appVersion**: App metadata
- **isLoggedIn**: Auth status

Context clears when request ends. No data leaks between requests.

### Purpose

Solves key problems:

- **Tracing**: `TrafficId` enables end-to-end request tracking.
- **Logging**: Auto-include context in all logs.
- **Monitoring**: Request-scoped monitoring data.
- **Cleaner code**: Context variables (`UserId, DeviceId`,...) are available throughout the request pipeline without parameter passing. Eliminates the need to pass context parameters through every function call.

Without request context, you would need to manually pass these values through every service, repository, and utility function, leading to polluted method signatures.

### Implementation

Languages uses async storage to implement the request context:

- **C#**: AsyncLocal for async/await
- **Java**: Vert.x Context for reactive
- **TypeScript**: AsyncLocalStorage for Node.js
- **Python**: contextvars for async context
- **Go**: Not available, explicit context passed via methods

**Setup Pattern:**

1. Middleware extracts from headers
2. Store in request context
3. Access throughout request
4. Auto cleanup after response

::: code-group

```C#
static class ReqContext<T>
{
	private static readonly ConcurrentDictionary<string, AsyncLocal<T>> state = new();

	public static void Set(string name, T data)
		=> state.GetOrAdd(name, _ => new AsyncLocal<T>()).Value = data;

	public static T? Get(string name) =>
		state.TryGetValue(name, out AsyncLocal<T>? data) ? data.Value : default;
}

public static class Req
{
	public static Guid TrafficId
	{
		get => ReqContext<Guid>.Get("trafficId")!;
		set => ReqContext<Guid>.Set("trafficId", value);
	}
	public static Guid UserId
	{
		get => ReqContext<Guid>.Get("userId");
		set => ReqContext<Guid>.Set("userId", value);
	}
	// Other context variables...
}
```

```Java
class ReqContext {
	public static <T> void set(String key, T value) {
		var ctx = Vertx.currentContext();
		if (ctx != null)
			ctx.putLocal(key, value);
	}

	@SuppressWarnings("unchecked")
	public static <T> T get(String key) {
		var ctx = Vertx.currentContext();
		return ctx != null ? (T) ctx.getLocal(key) : null;
	}

	public static void clear(String key) {
		var ctx = Vertx.currentContext();
		if (ctx != null)
			ctx.removeLocal(key);
	}
}

public class Req {
	public static UUID getTrafficId() {
		return ReqContext.get("trafficId");
	}
	public static void setTrafficId(UUID value) {
		ReqContext.set("trafficId", value);
	}
	public static UUID getUserId() {
		return ReqContext.get("userId");
	}
	// Other context variables...
}
```

```Go
// Context passed explicitly through handlers
```

```TypeScript
class ReqCtx {
  private static readonly storage = new AsyncLocalStorage<Map<string, any>>();

  public static set<T>(name: string, data: T): void {
    const context = this.storage.getStore();
    if (context) context.set(name, data);
  }

  public static get<T>(name: string): T | undefined {
    const context = this.storage.getStore();
    if (context) return context.get(name);
    return undefined;
  }

  public static run<T>(callback: () => T): T {
    const contextMap = new Map<string, any>();
    return this.storage.run(contextMap, callback);
  }
}

export class Req {
  public static run<T>(callback: () => T): T {
    return ReqCtx.run(callback);
  }

  public static get trafficId(): string {
    return ReqCtx.get<string>("trafficId") || "";
  }

  public static set trafficId(value: string) {
    ReqCtx.set("trafficId", value);
  }
  // Other context variables...
}
```

```Python
_request_context: contextvars.ContextVar[dict[str, Any]] = contextvars.ContextVar("_request_context", default={})

class ReqCtx:
    @staticmethod
    def run(callback):
        token = _request_context.set({})
        try:
            return callback()
        finally:
            _request_context.reset(token)

    @staticmethod
    def set(name: str, value: Any) -> None:
        ctx = _request_context.get()
        ctx[name] = value

    @staticmethod
    def get(name: str) -> Optional[Any]:
        return _request_context.get().get(name)

class Req:
    @staticmethod
    def run(callback):
        return ReqCtx.run(callback)

    @property
    def traffic_id(self) -> str:
        return ReqCtx.get("trafficId") or ""

    @traffic_id.setter
    def traffic_id(self, value: str) -> None:
        ReqCtx.set("trafficId", value)
    # Other context variables...
```

:::

### Usage Examples

**Middleware Setup:**

```C#
// C# middleware
public async Task InvokeAsync(HttpContext context, RequestDelegate next)
{
    Req.TrafficId = Guid.NewGuid();
    Req.UserId = ExtractUserIdFromToken(context);
    await next(context);
}
```

**Business Logic Access:**

```C#
// Anywhere in request pipeline
public void LogUserAction(string action)
{
    logger.LogInfo($"User {Req.UserId} performed {action} (TrafficId: {Req.TrafficId})");
}
```

Request context params are stored in traffic table, headers column.

### See

- Request context implementation in: [C#](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-cs/Src/Infra/Req.cs), [Java](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-java/src/main/java/mahaam/infra/Req.java), [TypeScript](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-ts/src/infra/req.ts), [Python](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-py/infra/req.py)
