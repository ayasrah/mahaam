import { Injectable, Inject } from '@nestjs/common';
import { TasksRepo } from './tasks.repo';
import { PlansRepo } from '../plans/plans.repo';
import { Task } from './tasks.model';
import { LogicError } from '../../infra/errors';
import { DB, Trx } from 'src/infra/db';

export interface TasksService {
  create(planId: string, title: string): Promise<string>;
  getList(planId: string): Promise<Task[]>;
  delete(planId: string, id: string): Promise<void>;
  updateDone(planId: string, id: string, done: boolean): Promise<void>;
  updateTitle(id: string, title: string): Promise<void>;
  reOrder(planId: string, oldOrder: number, newOrder: number): Promise<void>;
}

@Injectable()
export class DefaultTasksService implements TasksService {
  constructor(
    @Inject('TasksRepo') private readonly tasksRepo: TasksRepo,
    @Inject('PlansRepo') private readonly plansRepo: PlansRepo,
  ) {}

  async create(planId: string, title: string): Promise<string> {
    return await DB.withTrx(async (trx) => {
      const count = await this.tasksRepo.getCount(trx, planId);
      if (count >= 100) throw new LogicError('max_is_100', 'Max is 100');

      const id = await this.tasksRepo.create(trx, planId, title);
      await this.plansRepo.updateDonePercent(trx, planId);

      return id;
    });
  }

  async getList(planId: string): Promise<Task[]> {
    return DB.withTrx((trx) => this.tasksRepo.getAll(trx, planId));
  }

  async delete(planId: string, id: string): Promise<void> {
    await DB.withTrx(async (trx) => {
      await this.tasksRepo.updateOrderBeforeDelete(trx, planId, id);
      await this.tasksRepo.deleteOne(trx, id);
      await this.plansRepo.updateDonePercent(trx, planId);
    });
  }

  async updateDone(planId: string, id: string, done: boolean): Promise<void> {
    await DB.withTrx(async (trx) => {
      await this.tasksRepo.updateDone(trx, id, done);
      await this.plansRepo.updateDonePercent(trx, planId);

      var count = await this.tasksRepo.getCount(trx, planId);
      var task = await this.tasksRepo.getOne(trx, id);
      console.log(`task=${JSON.stringify(task)}, count=${count}`);
      await this.reOrder(planId, task.sortOrder, done ? 0 : count - 1, trx);
    });
  }

  async updateTitle(id: string, title: string): Promise<void> {
    await DB.withTrx((trx) => this.tasksRepo.updateTitle(trx, id, title));
  }

  async reOrder(planId: string, oldOrder: number, newOrder: number, trx?: Trx): Promise<void> {
    if (oldOrder === newOrder) return;

    const execute = async (transaction: Trx) => {
      const count = await this.tasksRepo.getCount(transaction, planId);
      if (oldOrder > count || newOrder > count) {
        throw new Error(`oldOrder and newOrder should be less than ${count}`);
      }
      await this.tasksRepo.updateOrder(transaction, planId, oldOrder, newOrder);
    };

    if (trx) await execute(trx);
    else await DB.withTrx(execute);
  }
}
