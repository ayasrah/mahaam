# Caching

### Overview

Caching is storing frequently accessed data in **memory**.

### Options

- In-memory variable/dictionary/map (Local to process)
- Distributed cache like Redis and Memcached (Distributed in-memory)
- CDN

### Used for

- Data that rarely changes like: product categories in ecommerce, or countries list.
- In Mahaam, cache stores nodeIP, nodeName and healthID values in memory map.

### Purpose

- **Performance**: Reduces latency
- **Reliability**: Still serve data during brief db downtime

### Model

::: code-group

```C#
class Cache
{
	public static void Init(Health health)
	{
		_health = health;
	}

	private static Health? _health;

	public static string NodeIP => _health?.NodeIP ?? "";
	public static string NodeName => _health?.NodeName ?? "";
	public static string ApiName => _health?.ApiName ?? "";
	public static string ApiVersion => _health?.ApiVersion ?? "";
	public static string EnvName => _health?.EnvName ?? "";
	public static Guid HealthId => _health?.Id ?? Guid.Empty;
}
```

```Java
@ApplicationScoped
public class Cache {

	public static void init(Health health) {
		_health = health;
	}

	private static Health _health;

	public static String getNodeIP() {
		return _health != null ? _health.nodeIP : "";
	}

	public static String getNodeName() {
		return _health != null ? _health.nodeName : "";
	}

	public static String getApiName() {
		return _health != null ? _health.apiName : "";
	}

	public static String getApiVersion() {
		return _health != null ? _health.apiVersion : "";
	}

	public static String getEnvName() {
		return _health != null ? _health.envName : "";
	}

	public static UUID getHealthId() {
		return _health != null ? _health.id : null;
	}
}
```

```Go
var (
	NodeIP     string
	NodeName   string
	ApiName    string
	ApiVersion string
	EnvName    string
	HealthID   uuid.UUID
)

func Init(h *models.Health) {
	NodeIP = h.NodeIP
	NodeName = h.NodeName
	ApiName = h.ApiName
	ApiVersion = h.ApiVersion
	EnvName = h.EnvName
	HealthID = h.ID
}
```

```TypeScript
@Injectable()
export class Cache {
  private static readonly logger = new Logger(Cache.name);
  private static _health: Health | null = null;

  public static init(health: Health): void {
    this._health = health;
  }

  public static getNodeIP(): string {
    return this._health?.nodeIP ?? '';
  }

  public static getNodeName(): string {
    return this._health?.nodeName ?? '';
  }

  public static getApiName(): string {
    return this._health?.apiName ?? '';
  }

  public static getApiVersion(): string {
    return this._health?.apiVersion ?? '';
  }

  public static getEnvName(): string {
    return this._health?.envName ?? 'development';
  }

  public static getHealthId(): string {
    return this._health?.id ?? '';
  }
}
```

```Python
node_ip = ""
node_name = ""
api_name = ""
api_version = ""
env_name = ""
health_id = UUID(int=0)

def init(health: Health):
    """Initialize cache with health object"""
    global node_ip, node_name, api_name, api_version, env_name, health_id
    node_ip = health.node_ip
    node_name = health.node_name
    api_name = health.api_name
    api_version = health.api_version
    env_name = health.env_name
    health_id = health.id

```

:::
