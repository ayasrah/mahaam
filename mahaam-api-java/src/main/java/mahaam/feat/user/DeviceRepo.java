package mahaam.feat.user;

import java.util.List;
import java.util.UUID;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import mahaam.feat.user.UserModel.Device;
import mahaam.infra.DB;
import mahaam.infra.Log;
import mahaam.infra.Mapper;

public interface DeviceRepo {
	UUID create(Device device);

	int delete(UUID id);

	int deleteByFingerprint(String fingerprint);

	Device getOne(UUID id);

	List<Device> getMany(UUID userId);

	int updateUserId(UUID deviceId, UUID userId);
}

@ApplicationScoped
class DefaultDeviceRepo implements DeviceRepo {
	@Inject
	DB db;

	@Override
	public UUID create(Device device) {
		try {
			String query = """
					INSERT INTO devices (id, user_id, platform, fingerprint, info, created_at)
					VALUES (:id, :userId, :platform, :fingerprint, :info, current_timestamp)""";
			var deviceId = UUID.randomUUID();
			var params = Mapper.of(
					"id",
					deviceId,
					"userId",
					device.userId,
					"platform",
					device.platform,
					"fingerprint",
					device.fingerprint,
					"info",
					device.info);
			db.insert(query, params);
			return deviceId;
		} catch (Exception e) {
			Log.error("Failed to create device: " + e.getMessage());
			throw new RuntimeException("Failed to create device", e);
		}
	}

	@Override
	public int delete(UUID id) {
		String query = "DELETE FROM devices WHERE id = :id";
		return db.delete(query, Mapper.of("id", id));
	}

	@Override
	public int deleteByFingerprint(String fingerprint) {
		String query = "DELETE FROM devices WHERE fingerprint = :fingerprint";
		return db.delete(query, Mapper.of("fingerprint", fingerprint));
	}

	@Override
	public Device getOne(UUID id) {
		String query = """
				SELECT id d_id, user_id d_userId, platform d_platform, fingerprint d_fingerprint, info d_info, created_at d_createdAt
				FROM devices d WHERE d.id = :id order by d.created_at desc""";
		return db.selectOne(query, Device.class, Mapper.of("id", id));
	}

	@Override
	public List<Device> getMany(UUID userId) {
		String query = """
				SELECT id d_id, user_id d_userId, platform d_platform, fingerprint d_fingerprint, info d_info, created_at d_createdAt
				FROM devices d WHERE d.user_id = :userId order by d.created_at desc""";
		return db.selectList(query, Device.class, Mapper.of("userId", userId));
	}

	@Override
	public int updateUserId(UUID deviceId, UUID userId) {
		String query = "UPDATE devices SET user_id = :userId, updated_at = current_timestamp WHERE id = :deviceId";
		return db.update(query, Mapper.of("deviceId", deviceId, "userId", userId));
	}
}
