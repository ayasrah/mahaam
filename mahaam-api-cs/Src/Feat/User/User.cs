
namespace Mahaam.Feat.Users;

public class User
{
	public Guid Id { set; get; }
	public string? Email { get; set; }
	public string? Name { get; set; }
}


public class Device
{
	public Guid Id { set; get; }
	public Guid UserId { set; get; }
	public string? Platform { set; get; }
	public string Fingerprint { set; get; }
	public string? Info { set; get; }
	public DateTime? CreatedAt { get; set; }
}

public class SuggestedEmail
{
	public Guid Id { set; get; }
	public Guid UserId { set; get; }
	public string? Email { set; get; }
	public DateTime? CreatedAt { get; set; }
}

public struct VerifiedUser
{
	public Guid UserId { set; get; }
	public Guid DeviceId { set; get; }
	public string Jwt { set; get; }
	public string? UserFullName { set; get; }
	public string? Email { get; set; }
}

public struct CreatedUser
{
	public Guid Id { set; get; }
	public Guid DeviceId { set; get; }
	public string Jwt { set; get; }
}