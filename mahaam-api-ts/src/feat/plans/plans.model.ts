import { User } from '../users/users.model';

export interface Plan {
  id: string;
  title?: string | null;
  type?: string | null;
  sortOrder: number;
  starts?: Date | null;
  ends?: Date | null;
  donePercent?: string | null;
  createdAt?: Date | null;
  updatedAt?: Date | null;
  sharedWith?: User[] | null;
  isShared: boolean;
  user: User;
}

export interface PlanIn {
  id: string;
  title?: string | null;
  starts?: Date | null;
  ends?: Date | null;
}

export class PlanType {
  public static readonly Main = 'Main';
  public static readonly Archived = 'Archived';

  public static readonly All: string[] = [PlanType.Main, PlanType.Archived];
}

export function mapRowToPlan(row: any): Plan {
  return {
    id: row.id,
    title: row.title,
    starts: row.starts,
    ends: row.ends,
    type: row.type,
    sortOrder: row.sort_order,
    donePercent: row.done_percent,
    isShared: row.is_shared,
    user: {
      id: row.user_id,
      email: row.user_email,
      name: row.user_name,
    },
  };
}
