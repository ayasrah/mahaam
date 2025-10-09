using Mahaam.Infra;

namespace Mahaam.Feat.Users;

public interface IUserRepo
{
	Guid Create();
	void UpdateName(Guid id, string name);
	void UpdateEmail(Guid id, string email);
	User? GetOne(string email);
	User GetOne(Guid id);
	int Delete(Guid id);
}

class UserRepo : IUserRepo
{
	public Guid Create()
	{
		const string query = "INSERT INTO users (id, created_at) VALUES (@id, current_timestamp)";
		var id = Guid.NewGuid();
		DB.Insert(query, new { id });
		return id;
	}

	public void UpdateName(Guid id, string name)
	{
		var query = "UPDATE users SET name = @name, updated_at = current_timestamp WHERE id = @id";
		DB.Update(query, new { id, name });
	}

	public void UpdateEmail(Guid id, string email)
	{
		var query = "UPDATE users SET email = @email, updated_at = current_timestamp WHERE id = @id";
		DB.Update(query, new { id, email });
	}

	public User? GetOne(string email)
	{
		var query = "SELECT id, name, email FROM users WHERE email = @email";
		return DB.SelectOne<User?>(query, new { email });
	}

	public User GetOne(Guid id)
	{
		var query = "SELECT id, name, email FROM users WHERE id = @id";
		return DB.SelectOne<User>(query, new { id });
	}

	public int Delete(Guid id)
	{
		var query = "DELETE FROM users WHERE id = @id";
		return DB.Delete(query, new { id });
	}
}