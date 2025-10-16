using Dapper;
using Npgsql;

namespace Mahaam.Infra;

public interface IDB
{
	T SelectOne<T>(string query, object? param = null);
	T SelectOne<T1, T2, T>(string query, Func<T1, T2, T> map, object? param = null, string splitOn = "id");
	List<T> SelectMany<T>(string query, object? param = null);
	List<T> SelectMany<T1, T2, T>(string query, Func<T1, T2, T> map, object? param = null, string splitOn = "id");
	int Insert(string query, object? param = null);
	int Update(string query, object? param = null);
	int Delete(string query, object? param = null);
}

class DB : IDB
{
	public T SelectOne<T>(string query, object? param = null)
	{
		using var cnn = new NpgsqlConnection(Config.DbUrl);
		return cnn.QuerySingleOrDefault<T>(query, param);
	}

	public T SelectOne<T1, T2, T>(string query, Func<T1, T2, T> map, object? param = null, string splitOn = "id")
	{
		using var cnn = new NpgsqlConnection(Config.DbUrl);
		var result = cnn.Query(query, map, param, splitOn: splitOn);
		return result.FirstOrDefault();
	}

	public List<T> SelectMany<T>(string query, object? param = null)
	{
		using var cnn = new NpgsqlConnection(Config.DbUrl);
		var result = cnn.Query<T>(query, param);
		return [.. result];
	}

	public List<T> SelectMany<T1, T2, T>(string query, Func<T1, T2, T> map, object? param = null, string splitOn = "id")
	{
		using var cnn = new NpgsqlConnection(Config.DbUrl);
		var result = cnn.Query(query, map, param, splitOn: splitOn);
		return [.. result];
	}

	public int Insert(string query, object? param = null) => Execute(query, param);
	public int Update(string query, object? param = null) => Execute(query, param);
	public int Delete(string query, object? param = null) => Execute(query, param);

	private int Execute(string query, object? param = null)
	{
		using var cnn = new NpgsqlConnection(Config.DbUrl);
		return cnn.Execute(query, param);
	}
}
