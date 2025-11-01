using System.Net;
using System.Net.Sockets;
using Mahaam.Infra.Monitoring;


namespace Mahaam.Infra;

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