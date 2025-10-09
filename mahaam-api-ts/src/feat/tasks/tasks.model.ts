export interface Task {
  id: string;
  planId: string;
  title: string;
  done: boolean;
  sortOrder: number;
  createdAt: Date;
  updatedAt?: Date | null;
}
