package mahaam.feat.task;

import java.time.Instant;
import java.util.UUID;

public class Task {
	public UUID id;
	public UUID planId;
	public String title;
	public boolean done;
	public int sortOrder;
	public Instant createdAt;
	public Instant updatedAt;
}
