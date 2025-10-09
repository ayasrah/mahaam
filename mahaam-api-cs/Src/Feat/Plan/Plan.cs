using Mahaam.Feat.Users;

namespace Mahaam.Feat.Plans;

public class Plan
{
	public Guid Id { get; set; }
	public string? Title { get; set; }
	public string? Type { get; set; }
	public int SortOrder { get; set; }
	public DateTime? Starts { get; set; }
	public DateTime? Ends { get; set; }
	public string? DonePercent { get; set; }
	public DateTime? CreatedAt { get; set; }
	public DateTime? UpdatedAt { get; set; }
	public List<User>? Members { get; set; }
	public bool IsShared { get; set; }
	public User User { get; set; }
}

public class PlanIn
{
	public Guid Id { get; set; }
	public string? Title { get; set; }
	public DateTime? Starts { get; set; }
	public DateTime? Ends { get; set; }
}

public static class PlanType
{
	public const string Main = "Main";
	public const string Archived = "Archived";
	public static readonly List<string> All = [Main, Archived];
}

public static class PlanStatus
{
	public const string Open = "Open";
	public const string Closed = "Closed";
	public static readonly List<string> All = [Open, Closed];
}