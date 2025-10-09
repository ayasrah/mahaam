import { Injectable } from '@nestjs/common';
import { Plan, PlanIn, PlanType, mapRowToPlan } from './plans.model';
import { DB, Trx } from '../../infra/db';
import { Req } from '../../infra/req';
import { randomUUID } from 'crypto';
import Log from 'src/infra/log';

const GROUP_STATUS = {
  OPEN: 'Open',
  CLOSED: 'Closed',
} as const;

const INITIAL_DONE_PERCENT = '0/0';

export interface PlansRepo {
  getOne(trx: Trx, id: string): Promise<Plan | null>;
  getMany(trx: Trx, userId: string, type: string): Promise<Plan[]>;
  create(trx: Trx, plan: PlanIn): Promise<string>;
  update(trx: Trx, plan: PlanIn): Promise<void>;
  delete(trx: Trx, id: string): Promise<void>;
  updateDonePercent(trx: Trx, id: string): Promise<void>;
  removeFromOrder(trx: Trx, userId: string, id: string): Promise<void>;
  updateOrder(trx: Trx, userId: string, type: string, oldOrder: number, newOrder: number): Promise<void>;
  updateType(trx: Trx, userId: string, id: string, type: string): Promise<void>;
  getCount(trx: Trx, userId: string, type: string): Promise<number>;
  updateUserId(trx: Trx, oldUserId: string, newUserId: string): Promise<number>;
}

@Injectable()
export class DefaultPlansRepo implements PlansRepo {
  async create(trx: Trx, plan: PlanIn): Promise<string> {
    const id = randomUUID();

    await trx`
      INSERT INTO plans (id, user_id, title, starts, ends, type, status, done_percent, sort_order, created_at)
      VALUES (${id}, ${Req.userId}, ${plan.title || null}, ${plan.starts || null}, ${plan.ends || null}, 
	  ${PlanType.Main}, ${GROUP_STATUS.OPEN}, ${INITIAL_DONE_PERCENT},
	  (SELECT COUNT(1) FROM plans WHERE user_id = ${Req.userId} AND type = ${PlanType.Main}), 
	  current_timestamp)
    `;

    return id;
  }

  async update(trx: Trx, plan: PlanIn): Promise<void> {
    await trx`
      UPDATE plans 
      SET title = ${plan.title || null}, starts = ${plan.starts || null}, ends = ${plan.ends || null}, updated_at = current_timestamp 
      WHERE id = ${plan.id}
    `;
  }

  async getOne(trx: Trx, id: string): Promise<Plan | null> {
    const result = await trx`
		SELECT c.id, c.title, c.starts, c.ends, c.type, c.done_percent, c.sort_order, c.user_id,
			EXISTS(SELECT 1 FROM plan_members cm WHERE cm.plan_id = c.id) AS is_shared,
			u.id user_id, u.email user_email, u.name user_name
		FROM plans c
		LEFT JOIN users u ON c.user_id = u.id
		WHERE c.id = ${id}`;

    if (result.length === 0) return null;

    return mapRowToPlan(result[0]);
  }

  async getMany(trx: Trx, userId: string, type: string): Promise<Plan[]> {
    try {
      const result = await trx`
		SELECT c.id, c.title, c.starts, c.ends, c.type, c.done_percent, c.sort_order, c.user_id,
			EXISTS(SELECT 1 FROM plan_members cm WHERE cm.plan_id = c.id) AS is_shared,
			u.id user_id, u.email user_email, u.name user_name
		FROM plans c
		LEFT JOIN users u ON c.user_id = u.id
		WHERE c.user_id = ${userId} AND c.type = ${type}
		ORDER BY c.sort_order DESC`;

      return result.map((row) => mapRowToPlan(row));
    } catch (e) {
      console.error(e);
      return [];
    }
  }

  async delete(trx: Trx, id: string): Promise<void> {
    const result = await trx`DELETE FROM plans WHERE id = ${id}`;
    if (result.count > 0) {
      Log.info(`Plan ${id} deleted`);
    }
  }

  async updateDonePercent(trx: Trx, id: string): Promise<void> {
    const tasks = await trx`SELECT * FROM tasks WHERE plan_id = ${id}`;
    const taskList = DB.as<any[]>(tasks);

    const done = taskList.filter((task) => task.done).length;
    const total = taskList.length;
    const donePercent = `${done}/${total}`;

    await trx`UPDATE plans SET done_percent = ${donePercent} WHERE id = ${id}`;
  }

  async removeFromOrder(trx: Trx, userId: string, id: string): Promise<void> {
    await trx`
      UPDATE plans SET sort_order = sort_order - 1
      WHERE user_id = ${userId} AND 
        type = (SELECT type FROM plans WHERE id = ${id}) AND
        sort_order > (SELECT sort_order FROM plans WHERE id = ${id})
    `;
  }

  async updateOrder(trx: Trx, userId: string, type: string, oldOrder: number, newOrder: number): Promise<void> {
    await trx`
      UPDATE plans SET sort_order = 
        CASE 
          WHEN sort_order = ${oldOrder} THEN ${newOrder}
          WHEN sort_order > ${oldOrder} AND sort_order <= ${newOrder} THEN sort_order - 1
          WHEN sort_order >= ${newOrder} AND sort_order < ${oldOrder} THEN sort_order + 1
          ELSE sort_order
        END
      WHERE 
        user_id = ${userId} AND 
        type = ${type}
    `;
  }

  async updateType(trx: Trx, userId: string, id: string, type: string): Promise<void> {
    await trx`
      UPDATE plans 
      SET type = ${type}, 
      sort_order = (SELECT COUNT(1) FROM plans WHERE user_id = ${userId} AND type = ${type}),
      updated_at = current_timestamp 
      WHERE id = ${id}
    `;
  }

  async getCount(trx: Trx, userId: string, type: string): Promise<number> {
    const result = await trx`SELECT COUNT(*) as count FROM plans WHERE user_id = ${userId} AND type = ${type}`;
    return DB.as<{ count: number }[]>(result)[0].count;
  }

  async updateUserId(trx: Trx, oldUserId: string, newUserId: string): Promise<number> {
    const result = await trx`
      UPDATE plans 
      SET user_id = ${newUserId},
      sort_order = (sort_order + (SELECT COUNT(1) FROM plans WHERE user_id = ${newUserId})),
      updated_at = current_timestamp 
      WHERE user_id = ${oldUserId}
    `;

    return result.count;
  }
}
