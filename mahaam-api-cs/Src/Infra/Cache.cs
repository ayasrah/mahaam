using System.Net;
using System.Net.Sockets;
using Mahaam.Infra.Monitoring;


namespace Mahaam.Infra;

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