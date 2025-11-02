
using System.Transactions;

namespace Mahaam.Infra.Monitoring;

public interface ITrafficRepo
{
	void Create(Traffic traffic);
}

public class TrafficRepo(IDB db, ILog log) : ITrafficRepo
{
	private readonly IDB _db = db;
	private readonly ILog _log = log;
	public void Create(Traffic traffic)
	{
		Task.Run(() =>
		{
			try
			{
				var query = @"INSERT INTO monitor.traffic (id, health_id, method, path, code, elapsed, headers, request, response, created_at)
					VALUES(@Id, @HealthId, @Method, @Path, @Code, @Elapsed, @Headers, @Request, @Response, current_timestamp)";
				using var scope = new TransactionScope(TransactionScopeOption.Suppress);
				_db.Insert(query, traffic);
			}
			catch (Exception e)
			{
				_log.Error("error creating traffic record: " + e.ToString());
			}
		});
	}
}