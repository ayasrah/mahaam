
namespace Mahaam.Infra.Monitoring;

public interface IHealthService
{
	void ServerStarted(Health health);
	void StartSendingPulses(CancellationToken cancellationToken = default);
	void ServerStopped();
}

public class HealthService(IHealthRepo healthRepo, ILog log) : IHealthService
{
	private readonly IHealthRepo _healthRepo = healthRepo;
	private readonly ILog _log = log;
	public void ServerStarted(Health health)
	{
		_healthRepo.Create(health);
	}

	public void StartSendingPulses(CancellationToken cancellationToken = default)
	{
		Task.Run(() =>
		{
			while (!cancellationToken.IsCancellationRequested)
			{
				try
				{
					_healthRepo.UpdatePulse(Cache.HealthId);
					Thread.Sleep(1000 * 60); // 1 minute
				}
				catch (Exception e)
				{
					_log.Error(e.ToString());
				}
			}
		}, cancellationToken);
	}

	public void ServerStopped()
	{
		var thread = new Thread(() =>
		{
			try
			{
				_healthRepo.UpdateStopped(Cache.HealthId);
				var stopMsg = $"✓ {Config.ApiName}-v{Config.ApiVersion}/{Cache.NodeIP}-{Cache.NodeName} stopped with healthID={Cache.HealthId}";
				_log.Info(stopMsg);
			}
			catch (Exception e)
			{
				_log.Error(e.ToString());
			}
		});
		thread.Start();
	}

}
