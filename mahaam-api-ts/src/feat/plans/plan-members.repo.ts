import { Injectable } from '@nestjs/common';
import { Plan, mapRowToPlan } from './plans.model';
import { User } from '../users/users.model';
import { DB, Trx } from '../../infra/db';

export interface PlanMembersRepo {
  create(trx: Trx, planId: string, userId: string): Promise<void>;
  delete(trx: Trx, planId: string, userId: string): Promise<number>;
  getUsers(trx: Trx, planId: string): Promise<User[]>;
  getUsersCount(trx: Trx, planId: string): Promise<number>;
  getOtherPlans(trx: Trx, userId: string): Promise<Plan[]>;
  getPlansCount(trx: Trx, userId: string): Promise<number>;
}

@Injectable()
export class DefaultPlanMembersRepo implements PlanMembersRepo {
  async create(trx: Trx, planId: string, userId: string): Promise<void> {
    await trx`
      INSERT INTO plan_members (plan_id, user_id, created_at) 
      VALUES (${planId}, ${userId}, current_timestamp)`;
  }

  async delete(trx: Trx, planId: string, userId: string): Promise<number> {
    const result = await trx`DELETE FROM plan_members WHERE plan_id = ${planId} AND user_id = ${userId}`;
    return result.count || 0;
  }

  async getUsers(trx: Trx, planId: string): Promise<User[]> {
    const result = await trx`
      SELECT u.id, u.name, u.email
	  FROM plan_members cm
      JOIN users u ON cm.user_id = u.id
      WHERE cm.plan_id = ${planId}
      ORDER BY u.name
    `;
    return result.map((row) => DB.as<User>(row));
  }

  async getUsersCount(trx: Trx, planId: string): Promise<number> {
    const result = await trx`
      SELECT COUNT(*) as count
      FROM plan_members
      WHERE plan_id = ${planId}
    `;
    return Number(result[0].count) || 0;
  }

  async getOtherPlans(trx: Trx, userId: string): Promise<Plan[]> {
    const result = await trx`
      SELECT c.id, c.title, c.starts, c.ends, c.type, c.done_percent, c.sort_order, true as is_shared,
        u.id as user_id, u.email as user_email, u.name as user_name
      FROM plan_members cm
      JOIN plans c ON cm.plan_id = c.id
      JOIN users u ON c.user_id = u.id
      WHERE cm.user_id = ${userId}
      ORDER BY c.sort_order ASC
    `;
    return result.map((row) => mapRowToPlan(row));
  }

  async getPlansCount(trx: Trx, userId: string): Promise<number> {
    const result = await trx`
      SELECT COUNT(DISTINCT cm.plan_id) as count
      FROM plan_members cm
      JOIN plans c ON cm.plan_id = c.id
      WHERE c.user_id = ${userId}
    `;
    return Number(result[0].count) || 0;
  }
}
