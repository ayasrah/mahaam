
using Mahaam.Infra;

namespace Mahaam.Feat.Users;

public interface ISuggestedEmailsRepo
{
	Task<Guid> Create(Guid userId, string email);
	Task<int> Delete(Guid id);
	Task<List<SuggestedEmail>> GetMany(Guid userId);
	Task<SuggestedEmail> GetOne(Guid id);
	Task<int> DeleteManyByEmail(string email);

}

class SuggestedEmailsRepo(IDB db) : ISuggestedEmailsRepo
{
	private readonly IDB _db = db;
	public async Task<Guid> Create(Guid userId, string email)
	{

		var query = @"INSERT INTO suggested_emails (id, user_id, email, created_at) 
			VALUES (@id, @userId, @email, current_timestamp)
			ON CONFLICT (user_id, email) DO NOTHING";
		var id = Guid.NewGuid();
		var updated = await _db.Insert(query, new { id, userId, email });
		return updated > 0 ? id : Guid.Empty;
	}

	public async Task<int> Delete(Guid id)
	{
		var query = "DELETE FROM suggested_emails WHERE id = @id";
		return await _db.Delete(query, new { id });
	}

	public async Task<int> DeleteManyByEmail(string email)
	{
		var query = "DELETE FROM suggested_emails WHERE email = @email";
		return await _db.Delete(query, new { email });
	}

	public async Task<List<SuggestedEmail>> GetMany(Guid userId)
	{
		var query = @"SELECT id, user_id, email, created_at
			FROM suggested_emails WHERE user_id = @userId order by created_at desc";
		return await _db.SelectMany<SuggestedEmail>(query, new { userId });
	}

	public async Task<SuggestedEmail> GetOne(Guid id)
	{
		var query = @"SELECT id, user_id, email, created_at FROM suggested_emails WHERE id = @id";
		return await _db.SelectOne<SuggestedEmail>(query, new { id });
	}
}