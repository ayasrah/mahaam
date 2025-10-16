
using Mahaam.Infra;

namespace Mahaam.Feat.Users;

public interface IDeviceRepo
{
	Guid Create(Device device);
	int Delete(Guid id);
	int DeleteByFingerprint(string fingerprint);
	Device GetOne(Guid id);
	List<Device> GetMany(Guid userId);
	int UpdateUserId(Guid id, Guid userId);
}

public class DeviceRepo(IDB db) : IDeviceRepo
{
	private readonly IDB _db = db;
	public Guid Create(Device device)
	{
		const string query = @"INSERT INTO devices (id, user_id, platform, fingerprint, info, created_at) 
			VALUES (@Id, @UserId, @Platform, @Fingerprint, @Info, current_timestamp)";
		device.Id = Guid.NewGuid();
		_db.Insert(query, device);
		return device.Id;
	}

	public int DeleteByFingerprint(string fingerprint)
	{
		var query = "DELETE FROM devices WHERE fingerprint = @fingerprint";
		return _db.Delete(query, new { fingerprint });
	}

	public int Delete(Guid id)
	{
		var query = "DELETE FROM devices WHERE id = @id";
		return _db.Delete(query, new { id });
	}

	public Device GetOne(Guid id)
	{
		var query = @"SELECT id, user_id, platform, fingerprint, info, created_at
			FROM devices WHERE id = @id order by created_at desc";
		return _db.SelectOne<Device>(query, new { id });
	}

	public List<Device> GetMany(Guid userId)
	{
		var query = @"SELECT id, user_id, platform, fingerprint, info, created_at
			FROM devices WHERE user_id = @userId order by created_at desc";
		return _db.SelectMany<Device>(query, new { userId });
	}

	public int UpdateUserId(Guid id, Guid userId)
	{
		var query = "UPDATE devices SET user_id = @userId, updated_at = current_timestamp WHERE id = @id";
		return _db.Update(query, new { id, userId });
	}
}