package mahaam.infra;

import java.util.List;
import java.util.Map;

import javax.sql.DataSource;

import org.jdbi.v3.core.Jdbi;
import org.jdbi.v3.core.mapper.reflect.FieldMapper;
import org.jdbi.v3.postgres.PostgresPlugin;

import jakarta.inject.Inject;
import jakarta.inject.Singleton;
import mahaam.feat.plan.PlanModel.Plan;
import mahaam.feat.task.Task;
import mahaam.feat.user.UserModel.Device;
import mahaam.feat.user.UserModel.SuggestedEmail;
import mahaam.feat.user.UserModel.User;
import mahaam.infra.monitor.MonitorModel.Health;
import mahaam.infra.monitor.MonitorModel.Traffic;

@Singleton
public class DB {

	// Pair class to hold two generic objects
	public static class Pair<T, U> {
		private final T first;
		private final U second;

		public Pair(T first, U second) {
			this.first = first;
			this.second = second;
		}

		public T getFirst() {
			return first;
		}

		public U getSecond() {
			return second;
		}
	}

	private Jdbi jdbi;

	@Inject
	DataSource dataSource;

	public Jdbi getJdbi() {
		return jdbi;
	}

	public void init() {
		// jdbi = Jdbi.create(Config.dbUrl);
		jdbi = Jdbi.create(dataSource);
		jdbi.installPlugin(new PostgresPlugin());
		// Use BeanMapper for nested bean mapping

		jdbi.registerRowMapper(FieldMapper.factory(Plan.class, "c"));
		jdbi.registerRowMapper(FieldMapper.factory(User.class, "u"));
		jdbi.registerRowMapper(FieldMapper.factory(Traffic.class, "t"));
		jdbi.registerRowMapper(FieldMapper.factory(Health.class, "h"));
		jdbi.registerRowMapper(FieldMapper.factory(Log.class, "l"));
		jdbi.registerRowMapper(FieldMapper.factory(SuggestedEmail.class, "s"));
		jdbi.registerRowMapper(FieldMapper.factory(Device.class, "d"));
		jdbi.registerRowMapper(FieldMapper.factory(Task.class, "t"));
		Log.info("DB initialized successfully");
	}

	public <T> List<T> selectList(String sql, Class<T> clazz, Map<String, Object> params) {

		return jdbi.withHandle(
				handle -> {
					return handle.createQuery(sql).bindMap(params).mapTo(clazz).list();
				});
	}

	public <T> T selectOne(String sql, Class<T> clazz, Map<String, Object> params) {

		return jdbi.withHandle(
				handle -> {
					return handle.createQuery(sql).bindMap(params).mapTo(clazz).findOne().orElse(null);
				});
	}

	public long selectCount(String sql, Map<String, Object> params) {

		return jdbi.withHandle(
				handle -> {
					return handle.createQuery(sql).bindMap(params).mapTo(Long.class).one();
				});
	}

	public int insert(String sql, Map<String, ?> params) {
		return execute(sql, params);
	}

	public int update(String sql, Map<String, ?> params) {
		return execute(sql, params);
	}

	public int delete(String sql, Map<String, ?> params) {
		return execute(sql, params);
	}

	private int execute(String sql, Map<String, ?> params) {
		return jdbi.withHandle(
				handle -> {
					return handle.createUpdate(sql).bindMap(params).execute();
				});
	}
}
