
namespace Mahaam.Feat.Tasks;

public class Task
{
	public Guid Id { get; set; }
	public Guid PlanId { get; set; }
	public string Title { get; set; }
	public bool Done { get; set; }
	public int SortOrder { get; set; }
	public DateTime CreatedAt { get; set; }
	public DateTime? UpdatedAt { get; set; }
}
