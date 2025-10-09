using Dapper;
using Npgsql;

namespace Mahaam.Infra
{
	public class DB
	{
		public static T SelectOne<T>(string query, object? param = null)
		{
			using var cnn = new NpgsqlConnection(Config.DbUrl);
			return cnn.QuerySingleOrDefault<T>(query, param);
		}

		public static T SelectOne<T1, T2, T>(string query, Func<T1, T2, T> map, object? param = null, string splitOn = "id")
		{
			using var cnn = new NpgsqlConnection(Config.DbUrl);
			var result = cnn.Query(query, map, param, splitOn: splitOn);
			return result.FirstOrDefault();
		}

		public static List<T> SelectMany<T>(string query, object? param = null)
		{
			using var cnn = new NpgsqlConnection(Config.DbUrl);
			var result = cnn.Query<T>(query, param);
			return [.. result];
		}

		public static List<T> SelectMany<T1, T2, T>(string query, Func<T1, T2, T> map, object? param = null, string splitOn = "id")
		{
			using var cnn = new NpgsqlConnection(Config.DbUrl);
			var result = cnn.Query(query, map, param, splitOn: splitOn);
			return [.. result];
		}

		public static int Insert(string query, object? param = null) => Execute(query, param);
		public static int Update(string query, object? param = null) => Execute(query, param);
		public static int Delete(string query, object? param = null) => Execute(query, param);

		private static int Execute(string query, object? param = null)
		{
			using var cnn = new NpgsqlConnection(Config.DbUrl);
			return cnn.Execute(query, param);
		}
	}
}