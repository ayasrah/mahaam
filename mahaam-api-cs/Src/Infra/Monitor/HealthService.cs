
namespace Mahaam.Infra.Monitoring;

public interface IHealthService
{
	void ServerStarted(Health health);
	void StartSendingPulses(CancellationToken cancellationToken = default);
	void ServerStopped();
}

class HealthService : IHealthService
{

	public void ServerStarted(Health health)
	{
		App.HealthRepo.Create(health);
	}

	public void StartSendingPulses(CancellationToken cancellationToken = default)
	{
		Task.Run(() =>
		{
			while (!cancellationToken.IsCancellationRequested)
			{
				try
				{
					App.HealthRepo.UpdatePulse(Cache.HealthId);
					Thread.Sleep(1000 * 60); // 1 minute
				}
				catch (Exception e)
				{
					Log.Error(e.ToString());
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
				App.HealthRepo.UpdateStopped(Cache.HealthId);
				var stopMsg = $"✓ {Config.ApiName}-v{Config.ApiVersion}/{Cache.NodeIP}-{Cache.NodeName} stopped with healthID={Cache.HealthId}";
				Log.Info(stopMsg);
			}
			catch (Exception e)
			{
				Log.Error(e.ToString());
			}
		});
		thread.Start();
	}

}
