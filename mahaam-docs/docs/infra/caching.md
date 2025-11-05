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
public interface ICache
{
	void Init(Health health);
	string NodeIP();
	string NodeName();
	Guid HealthId();
}

class Cache : ICache
{
	private Health? _health;

	public void Init(Health health)
	{
		_health = health;
	}

	public string NodeIP() => _health?.NodeIP ?? "";
	public string NodeName() => _health?.NodeName ?? "";
	public Guid HealthId() => _health?.Id ?? Guid.Empty;
}
```

```Java
@ApplicationScoped
public class Cache {

	public void init(Health health) {
		_health = health;
	}

	private Health _health;

	public String nodeIP() {
		return _health != null ? _health.nodeIP : "";
	}

	public String nodeName() {
		return _health != null ? _health.nodeName : "";
	}

	public UUID healthId() {
		return _health != null ? _health.id : null;
	}
}
```

```Go
var env *Environment

func Env() *Environment {
	return env
}

func NewEnvironment(h *models.Health) {
	env = &Environment{
		NodeIP:   h.NodeIP,
		NodeName: h.NodeName,
		HealthID: h.ID,
	}
}

type Environment struct {
	NodeIP   string
	NodeName string
	HealthID uuid.UUID
}
```

```TypeScript
@Injectable()
export class Cache {
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

  public static getHealthId(): string {
    return this._health?.id ?? '';
  }
}
```

```Python
class Cache:
    _health: Health | None = None

    @classmethod
    def init(cls, health: Health) -> None:
        cls._health = health

    @classmethod
    def node_ip(cls) -> str:
        return cls._health.node_ip if cls._health else ""

    @classmethod
    def node_name(cls) -> str:
        return cls._health.node_name if cls._health else ""

    @classmethod
    def health_id(cls) -> UUID:
        return cls._health.id if cls._health else UUID(int=0)
```

:::

### See

- Cache implementation in: [C#](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-cs/Src/Infra/Cache.cs), [Java](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-java/src/main/java/mahaam/infra/Cache.java), [Go](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-go/utils/cache/cache.go), [TypeScript](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-ts/src/infra/cache.ts), [Python](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-py/infra/cache.py)
