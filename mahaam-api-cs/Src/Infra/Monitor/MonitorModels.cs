namespace Mahaam.Infra.Monitoring;

public class Traffic
{
	public Guid Id { set; get; }
	public Guid HealthId { set; get; }
	public required string Method { set; get; }
	public required string Path { set; get; }
	public int? Code { set; get; }
	public long? Elapsed { set; get; }
	public string? Headers { set; get; }
	public string? Request { set; get; }
	public string? Response { set; get; }
}

public class TrafficHeaders
{
	public Guid? UserId { set; get; }
	public Guid? DeviceId { set; get; }
	public string? AppVersion { set; get; }
	public string? AppStore { set; get; }
}

public class Health
{
	public Guid Id { get; set; }
	public string? ApiName { get; set; }
	public string? ApiVersion { get; set; }
	public string? NodeIP { get; set; }
	public string? NodeName { get; set; }
	public string? EnvName { get; set; }
}
