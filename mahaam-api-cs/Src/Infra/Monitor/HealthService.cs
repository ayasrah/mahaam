
namespace Mahaam.Infra.Monitoring;

public interface IHealthService
{
	void ServerStarted(Health health);
	void StartSendingPulses(CancellationToken cancellationToken = default);
	void ServerStopped();
}

public class HealthService(IHealthRepo healthRepo, ILog log, Settings settings, ICache cache) : IHealthService
{

	public void ServerStarted(Health health)
	{
		healthRepo.Create(health).GetAwaiter().GetResult();
	}

	public void StartSendingPulses(CancellationToken cancellationToken = default)
	{
		Task.Run(() =>
		{
			while (!cancellationToken.IsCancellationRequested)
			{
				try
				{
					healthRepo.UpdatePulse(cache.HealthId());
					Thread.Sleep(1000 * 60); // 1 minute
				}
				catch (Exception e)
				{
					log.Error(e.ToString());
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
				healthRepo.UpdateStopped(cache.HealthId());
				var stopMsg = $"✓ {settings.Api.Name}-v{settings.Api.Version}/{cache.NodeIP()}-{cache.NodeName()} stopped with healthID={cache.HealthId()}";
				log.Info(stopMsg);
			}
			catch (Exception e)
			{
				log.Error(e.ToString());
			}
		});
		thread.Start();
	}

}
