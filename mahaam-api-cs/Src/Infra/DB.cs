using Dapper;
using Npgsql;

namespace Mahaam.Infra;

public interface IDB
{
	Task<T> SelectOne<T>(string query, object? param = null);
	Task<T> SelectOne<T1, T2, T>(string query, Func<T1, T2, T> map, object? param = null, string splitOn = "id");
	Task<List<T>> SelectMany<T>(string query, object? param = null);
	Task<List<T>> SelectMany<T1, T2, T>(string query, Func<T1, T2, T> map, object? param = null, string splitOn = "id");
	Task<int> Insert(string query, object? param = null);
	Task<int> Update(string query, object? param = null);
	Task<int> Delete(string query, object? param = null);
}

class DB(Settings settings) : IDB
{
	private readonly Settings _settings = settings;
	public async Task<T> SelectOne<T>(string query, object? param = null)
	{
		await using var cnn = new NpgsqlConnection(_settings.DbUrl);
		return await cnn.QuerySingleOrDefaultAsync<T>(query, param);
	}

	public async Task<T> SelectOne<T1, T2, T>(string query, Func<T1, T2, T> map, object? param = null, string splitOn = "id")
	{
		await using var cnn = new NpgsqlConnection(_settings.DbUrl);
		var result = await cnn.QueryAsync(query, map, param, splitOn: splitOn);
		return result.FirstOrDefault();
	}

	public async Task<List<T>> SelectMany<T>(string query, object? param = null)
	{
		await using var cnn = new NpgsqlConnection(_settings.DbUrl);
		var result = await cnn.QueryAsync<T>(query, param);
		return [.. result];
	}

	public async Task<List<T>> SelectMany<T1, T2, T>(string query, Func<T1, T2, T> map, object? param = null, string splitOn = "id")
	{
		await using var cnn = new NpgsqlConnection(_settings.DbUrl);
		var result = await cnn.QueryAsync(query, map, param, splitOn: splitOn);
		return [.. result];
	}

	public async Task<int> Insert(string query, object? param = null) => await ExecuteAsync(query, param);
	public async Task<int> Update(string query, object? param = null) => await ExecuteAsync(query, param);
	public async Task<int> Delete(string query, object? param = null) => await ExecuteAsync(query, param);

	private async Task<int> ExecuteAsync(string query, object? param = null)
	{
		await using var cnn = new NpgsqlConnection(_settings.DbUrl);
		return await cnn.ExecuteAsync(query, param);
	}
}


