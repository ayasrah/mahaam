package mahaam.feat.user;

import java.util.List;
import java.util.UUID;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import mahaam.feat.user.UserModel.SuggestedEmail;
import mahaam.infra.DB;
import mahaam.infra.Mapper;

public interface SuggestedEmailsRepo {
	UUID create(UUID userId, String email);

	int delete(UUID id);

	List<SuggestedEmail> getMany(UUID userId);

	SuggestedEmail getOne(UUID id);

	int deleteManyByEmail(String email);
}

@ApplicationScoped
class DefaultSuggestedEmailsRepo implements SuggestedEmailsRepo {

	@Inject
	DB db;

	@Override
	public UUID create(UUID userId, String email) {
		String query = """
				INSERT INTO suggested_emails (id, user_id, email, created_at)
				VALUES (:id, :userId, :email, current_timestamp)
				ON CONFLICT (user_id, email) DO NOTHING""";
		UUID id = UUID.randomUUID();
		int updated = db.insert(query, Mapper.of("id", id, "userId", userId, "email", email));
		return updated > 0 ? id : null;

	}

	@Override
	public int delete(UUID id) {
		String query = "DELETE FROM suggested_emails WHERE id = :id";
		return db.delete(query, Mapper.of("id", id));
	}

	@Override
	public List<SuggestedEmail> getMany(UUID userId) {
		String query = """
				SELECT id s_id, user_id s_userId, email s_email, created_at s_createdAt
				FROM suggested_emails s
				WHERE s.user_id = :userId ORDER BY s.created_at DESC""";
		return db.selectList(query, SuggestedEmail.class, Mapper.of("userId", userId));
	}

	@Override
	public SuggestedEmail getOne(UUID id) {
		String query = """
				SELECT id s_id, user_id s_userId, email s_email, created_at s_createdAt
				FROM suggested_emails s
				WHERE s.id = :id""";
		return db.selectOne(query, SuggestedEmail.class, Mapper.of("id", id));
	}

	@Override
	public int deleteManyByEmail(String email) {
		String query = "DELETE FROM suggested_emails WHERE email = :email";
		return db.delete(query, Mapper.of("email", email));
	}
}
