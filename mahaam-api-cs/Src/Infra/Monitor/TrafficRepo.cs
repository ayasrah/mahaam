
using System.Transactions;

namespace Mahaam.Infra.Monitoring;

public interface ITrafficRepo
{
	void Create(Traffic traffic);
}

public class TrafficRepo(IDB db, ILog log) : ITrafficRepo
{
	public void Create(Traffic traffic)
	{
		Task.Run(() =>
		{
			try
			{
				var query = @"INSERT INTO monitor.traffic (id, health_id, method, path, code, elapsed, headers, request, response, created_at)
					VALUES(@Id, @HealthId, @Method, @Path, @Code, @Elapsed, @Headers, @Request, @Response, current_timestamp)";
				using var scope = new TransactionScope(TransactionScopeOption.Suppress);
				db.Insert(query, traffic);
			}
			catch (Exception e)
			{
				log.Error("error creating traffic record: " + e.ToString());
			}
		});
	}
}