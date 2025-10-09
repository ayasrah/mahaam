import { Injectable } from '@nestjs/common';
import { Task } from './tasks.model';
import { DB, Trx } from '../../infra/db';
import { randomUUID } from 'crypto';
import { NotFoundError } from 'src/infra/errors';

export interface TasksRepo {
  getAll(trx: Trx, planId: string): Promise<Task[]>;
  getOne(trx: Trx, id: string): Promise<Task>;
  create(trx: Trx, planId: string, title: string): Promise<string>;
  deleteOne(trx: Trx, id: string): Promise<void>;
  deleteAll(trx: Trx, planId: string): Promise<void>;
  updateDone(trx: Trx, id: string, done: boolean): Promise<void>;
  updateTitle(trx: Trx, id: string, title: string): Promise<void>;
  updateOrder(trx: Trx, planId: string, oldOrder: number, newOrder: number): Promise<void>;
  updateOrderBeforeDelete(trx: Trx, planId: string, id: string): Promise<void>;
  getCount(trx: Trx, planId: string): Promise<number>;
}

@Injectable()
export class DefaultTasksRepo implements TasksRepo {
  async getAll(trx: Trx, planId: string): Promise<Task[]> {
    const result =
      await trx`SELECT id, plan_id, title, done, sort_order, created_at, updated_at FROM tasks WHERE plan_id = ${planId} ORDER BY sort_order ASC`;
    return DB.as<Task[]>(result);
  }

  async getOne(trx: Trx, id: string): Promise<Task> {
    const result =
      await trx`SELECT id, plan_id, title, done, sort_order, created_at, updated_at FROM tasks WHERE id = ${id}`;
    if (result.length === 0) {
      throw new NotFoundError(`task id=${id} not found`);
    }
    return DB.as<Task>(result[0]);
  }

  async create(trx: Trx, planId: string, title: string): Promise<string> {
    const id = randomUUID();
    await trx`INSERT INTO tasks (id, plan_id, title, done, sort_order, created_at) VALUES (${id}, ${planId}, ${title}, ${false},
	  (SELECT COUNT(1) FROM tasks WHERE plan_id = ${planId}), current_timestamp)`;
    return id;
  }

  async deleteOne(trx: Trx, id: string): Promise<void> {
    const result = await trx`DELETE FROM tasks WHERE id = ${id}`;
    if ((result.count || 0) === 0) {
      throw new NotFoundError(`task id=${id} not found`);
    }
  }

  async deleteAll(trx: Trx, planId: string): Promise<void> {
    await trx`DELETE FROM tasks WHERE plan_id = ${planId}`;
  }

  async updateDone(trx: Trx, id: string, done: boolean): Promise<void> {
    const result = await trx`UPDATE tasks SET done = ${done}, updated_at = current_timestamp WHERE id = ${id}`;
    if ((result.count || 0) === 0) {
      throw new NotFoundError(`task id=${id} not found`);
    }
  }

  async updateTitle(trx: Trx, id: string, title: string): Promise<void> {
    const result = await trx`UPDATE tasks SET title = ${title}, updated_at = current_timestamp WHERE id = ${id}`;
    if ((result.count || 0) === 0) {
      throw new NotFoundError(`task id=${id} not found`);
    }
  }

  async updateOrderBeforeDelete(trx: Trx, planId: string, id: string): Promise<void> {
    await trx`UPDATE tasks SET sort_order = sort_order - 1 WHERE plan_id = ${planId} AND sort_order > (SELECT sort_order FROM tasks WHERE id = ${id})`;
  }

  async updateOrder(trx: Trx, planId: string, oldOrder: number, newOrder: number): Promise<void> {
    console.log(`planId=${planId}, oldOrder=${oldOrder}, newOrder=${newOrder}`);
    await trx`UPDATE tasks SET sort_order = 
	CASE 
		WHEN sort_order = ${oldOrder} THEN ${newOrder} 
		WHEN sort_order > ${oldOrder} AND sort_order <= ${newOrder} THEN sort_order - 1
		WHEN sort_order >= ${newOrder} AND sort_order < ${oldOrder} THEN sort_order + 1 
	ELSE sort_order END 
	WHERE plan_id = ${planId}`;
  }

  async getCount(trx: Trx, planId: string): Promise<number> {
    const result = await trx`SELECT COUNT(1) as count FROM tasks WHERE plan_id = ${planId}`;
    return DB.as<{ count: number }[]>(result)[0].count;
  }
}
