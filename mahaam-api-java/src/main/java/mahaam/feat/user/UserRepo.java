package mahaam.feat.user;

import java.util.UUID;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import mahaam.feat.user.UserModel.User;
import mahaam.infra.DB;
import mahaam.infra.Mapper;

public interface UserRepo {
	UUID create();

	void updateName(UUID id, String name);

	void updateEmail(UUID id, String email);

	User getOne(String email);

	User getOne(UUID id);

	int delete(UUID id);
}

@ApplicationScoped
class DefaultUserRepo implements UserRepo {

	@Inject
	DB db;

	@Override
	public UUID create() {
		String query = "INSERT INTO users (id, created_at) VALUES (:id, current_timestamp)";
		UUID id = UUID.randomUUID();
		db.insert(query, Mapper.of("id", id));
		return id;
	}

	@Override
	public void updateName(UUID id, String name) {
		String query = "UPDATE users SET name = :name, updated_at = current_timestamp WHERE id = :id";
		db.update(query, Mapper.of("id", id, "name", name));
	}

	@Override
	public void updateEmail(UUID id, String email) {
		String query = "UPDATE users SET email = :email, updated_at = current_timestamp WHERE id = :id";
		db.update(query, Mapper.of("id", id, "email", email));
	}

	@Override
	public User getOne(String email) {
		String query = "SELECT id u_id, name u_name, email u_email FROM users u WHERE u.email = :email";
		return db.selectOne(query, User.class, Mapper.of("email", email));
	}

	@Override
	public User getOne(UUID id) {
		String query = "SELECT id u_id, name u_name, email u_email FROM users u WHERE u.id = :id";
		return db.selectOne(query, User.class, Mapper.of("id", id));
	}

	@Override
	public int delete(UUID id) {
		String query = "DELETE FROM users WHERE id = :id";
		return db.delete(query, Mapper.of("id", id));
	}
}
