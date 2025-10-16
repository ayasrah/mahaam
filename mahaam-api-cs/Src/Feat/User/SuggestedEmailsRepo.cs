
using Mahaam.Infra;

namespace Mahaam.Feat.Users;

public interface ISuggestedEmailsRepo
{
	Guid Create(Guid userId, string email);
	int Delete(Guid id);
	List<SuggestedEmail> GetMany(Guid userId);
	public SuggestedEmail GetOne(Guid id);
	int DeleteManyByEmail(string email);

}

class SuggestedEmailsRepo(IDB db) : ISuggestedEmailsRepo
{
	private readonly IDB _db = db;
	public Guid Create(Guid userId, string email)
	{

		var query = @"INSERT INTO suggested_emails (id, user_id, email, created_at) 
			VALUES (@id, @userId, @email, current_timestamp)
			ON CONFLICT (user_id, email) DO NOTHING";
		var id = Guid.NewGuid();
		var updated = _db.Insert(query, new { id, userId, email });
		return updated > 0 ? id : Guid.Empty;
	}

	public int Delete(Guid id)
	{
		var query = "DELETE FROM suggested_emails WHERE id = @id";
		return _db.Delete(query, new { id });
	}

	public int DeleteManyByEmail(string email)
	{
		var query = "DELETE FROM suggested_emails WHERE email = @email";
		return _db.Delete(query, new { email });
	}

	public List<SuggestedEmail> GetMany(Guid userId)
	{
		var query = @"SELECT id, user_id, email, created_at
			FROM suggested_emails WHERE user_id = @userId order by created_at desc";
		return _db.SelectMany<SuggestedEmail>(query, new { userId });
	}

	public SuggestedEmail GetOne(Guid id)
	{
		var query = @"SELECT id, user_id, email, created_at FROM suggested_emails WHERE id = @id";
		return _db.SelectOne<SuggestedEmail>(query, new { id });
	}
}