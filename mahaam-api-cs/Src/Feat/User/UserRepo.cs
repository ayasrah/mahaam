using Mahaam.Infra;

namespace Mahaam.Feat.Users;

public interface IUserRepo
{
	Task<Guid> Create();
	Task UpdateName(Guid id, string name);
	Task UpdateEmail(Guid id, string email);
	Task<User?> GetOne(string email);
	Task<User> GetOne(Guid id);
	Task<int> Delete(Guid id);
}

class UserRepo(IDB db) : IUserRepo
{
	private readonly IDB _db = db;
	public async Task<Guid> Create()
	{
		const string query = "INSERT INTO users (id, created_at) VALUES (@id, current_timestamp)";
		var id = Guid.NewGuid();
		await _db.Insert(query, new { id });
		return id;
	}

	public async Task UpdateName(Guid id, string name)
	{
		var query = "UPDATE users SET name = @name, updated_at = current_timestamp WHERE id = @id";
		await _db.Update(query, new { id, name });
	}

	public async Task UpdateEmail(Guid id, string email)
	{
		var query = "UPDATE users SET email = @email, updated_at = current_timestamp WHERE id = @id";
		await _db.Update(query, new { id, email });
	}

	public async Task<User?> GetOne(string email)
	{
		var query = "SELECT id, name, email FROM users WHERE email = @email";
		return await _db.SelectOne<User?>(query, new { email });
	}

	public async Task<User> GetOne(Guid id)
	{
		var query = "SELECT id, name, email FROM users WHERE id = @id";
		return await _db.SelectOne<User>(query, new { id });
	}

	public async Task<int> Delete(Guid id)
	{
		var query = "DELETE FROM users WHERE id = @id";
		return await _db.Delete(query, new { id });
	}
}