
using System.Transactions;

namespace Mahaam.Infra.Monitoring;

public interface ILogRepo
{
	void Create(string type, string message, Guid? trafficId);
}

public class LogRepo(IDB db, ICache cache) : ILogRepo
{
	public void Create(string type, string message, Guid? trafficId)
	{
		var err = new
		{
			trafficId,
			type,
			message,
			node_ip = cache.NodeIP(),
		};
		var query = @"INSERT INTO monitor.logs (traffic_id, type, message, node_ip, created_at) 
			VALUES (@trafficId, @type, @message, @node_ip, current_timestamp)";

		Task.Run(() =>
			{
				try
				{
					using var scope = new TransactionScope(TransactionScopeOption.Suppress);
					db.Insert(query, err);
				}
				catch (Exception ex)
				{
					Serilog.Log.Error($"Unable to create log record:({type}: {message}). Cause: {ex.ToString()}");
				}
			}
		);

	}
}