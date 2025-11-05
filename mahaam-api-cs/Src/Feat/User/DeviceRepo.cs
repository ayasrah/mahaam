
using Mahaam.Infra;

namespace Mahaam.Feat.Users;

public interface IDeviceRepo
{
	Task<Guid> Create(Device device);
	Task<int> Delete(Guid id);
	Task<int> DeleteByFingerprint(string fingerprint);
	Task<Device> GetOne(Guid id);
	Task<List<Device>> GetMany(Guid userId);
	Task<int> UpdateUserId(Guid id, Guid userId);
}

public class DeviceRepo(IDB db) : IDeviceRepo
{
	public async Task<Guid> Create(Device device)
	{
		const string query = @"INSERT INTO devices (id, user_id, platform, fingerprint, info, created_at) 
			VALUES (@Id, @UserId, @Platform, @Fingerprint, @Info, current_timestamp)";
		device.Id = Guid.NewGuid();
		await db.Insert(query, device);
		return device.Id;
	}

	public async Task<int> DeleteByFingerprint(string fingerprint)
	{
		var query = "DELETE FROM devices WHERE fingerprint = @fingerprint";
		return await db.Delete(query, new { fingerprint });
	}

	public async Task<int> Delete(Guid id)
	{
		var query = "DELETE FROM devices WHERE id = @id";
		return await db.Delete(query, new { id });
	}

	public async Task<Device> GetOne(Guid id)
	{
		var query = @"SELECT id, user_id, platform, fingerprint, info, created_at
			FROM devices WHERE id = @id order by created_at desc";
		return await db.SelectOne<Device>(query, new { id });
	}

	public async Task<List<Device>> GetMany(Guid userId)
	{
		var query = @"SELECT id, user_id, platform, fingerprint, info, created_at
			FROM devices WHERE user_id = @userId order by created_at desc";
		return await db.SelectMany<Device>(query, new { userId });
	}

	public async Task<int> UpdateUserId(Guid id, Guid userId)
	{
		var query = "UPDATE devices SET user_id = @userId, updated_at = current_timestamp WHERE id = @id";
		return await db.Update(query, new { id, userId });
	}
}