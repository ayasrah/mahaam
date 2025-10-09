
using System.Transactions;

namespace Mahaam.Infra.Monitoring;

public interface ITrafficRepo
{
	void Create(Traffic traffic);
}

public class TrafficRepo : ITrafficRepo
{

	public void Create(Traffic traffic)
	{
		Task.Run(() =>
		{
			try
			{
				var query = @"INSERT INTO x_traffic (id, health_id, method, path, code, elapsed, headers, request, response, created_at)
					VALUES(@Id, @HealthId, @Method, @Path, @Code, @Elapsed, @Headers, @Request, @Response, current_timestamp)";
				using var scope = new TransactionScope(TransactionScopeOption.Suppress);
				DB.Insert(query, traffic);
			}
			catch (Exception e)
			{
				Log.Error("error creating traffic record: " + e.ToString());
			}
		});
	}
}