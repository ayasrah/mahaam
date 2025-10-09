import { Injectable, Inject } from '@nestjs/common';
import { PlansRepo } from './plans.repo';
import { PlanMembersRepo } from './plan-members.repo';
import { UsersRepo } from '../users/users.repo';
import { SuggestedEmailRepo } from '../users/suggested-email.repo';
import { Plan, PlanIn } from './plans.model';
import { Req } from '../../infra/req';
import { LogicError, InputError, UnauthorizedError } from '../../infra/errors';
import Log from '../../infra/log';
import { DB } from 'src/infra/db';

export interface PlansService {
  getOne(id: string): Promise<Plan | null>;
  getMany(type?: string): Promise<Plan[]>;
  create(plan: PlanIn): Promise<string>;
  update(plan: PlanIn): Promise<void>;
  delete(id: string): Promise<void>;
  share(id: string, email: string): Promise<void>;
  unshare(id: string, email: string): Promise<void>;
  leave(id: string): Promise<void>;
  updateType(id: string, type: string): Promise<void>;
  reOrder(type: string, oldOrder: number, newOrder: number): Promise<void>;
  validateUserOwnsThePlan(planId: string): Promise<void>;
}

@Injectable()
export class DefaultPlansService implements PlansService {
  constructor(
    @Inject('PlansRepo') private readonly plansRepo: PlansRepo,
    @Inject('PlanMembersRepo') private readonly planMembersRepo: PlanMembersRepo,
    @Inject('UsersRepo') private readonly usersRepo: UsersRepo,
    @Inject('SuggestedEmailRepo') private readonly suggestedEmailRepo: SuggestedEmailRepo,
  ) {}

  async getOne(id: string): Promise<Plan | null> {
    return await DB.withTrx(async (trx) => {
      const plan = await this.plansRepo.getOne(trx, id);
      if (plan?.isShared) {
        plan.sharedWith = await this.planMembersRepo.getUsers(trx, id);
      }
      return plan;
    });
  }

  async getMany(type?: string): Promise<Plan[]> {
    const planType = type || 'Main';

    // Plans of the user (shared or not)
    return await DB.withTrx(async (trx) => {
      const plans = await this.plansRepo.getMany(trx, Req.userId, planType);
      if (Req.isLoggedIn) {
        // Plans of others shared with the user
        const planMembers = await this.planMembersRepo.getOtherPlans(trx, Req.userId);
        plans.push(...planMembers);
      }

      return plans;
    });
  }

  async create(plan: PlanIn): Promise<string> {
    const id = await DB.withTrx(async (trx) => {
      const count = await this.plansRepo.getCount(trx, Req.userId, 'Main');
      if (count >= 100) throw new LogicError('max_is_100', 'Max is 100');
      return await this.plansRepo.create(trx, plan);
    });
    return id;
  }

  async update(plan: PlanIn): Promise<void> {
    await this.validateUserOwnsThePlan(plan.id);
    await DB.withTrx(async (trx) => {
      await this.plansRepo.update(trx, plan);
    });
  }

  async delete(id: string): Promise<void> {
    await this.validateUserOwnsThePlan(id);

    await DB.withTrx(async (trx) => {
      await this.plansRepo.removeFromOrder(trx, Req.userId, id);
      await this.plansRepo.delete(trx, id); // This will delete all related records as it is a cascade delete
    });
  }

  async share(id: string, email: string): Promise<void> {
    await this.validateUserLoggedIn();
    await this.validateUserOwnsThePlan(id);

    await DB.withTrx(async (trx) => {
      const user = await this.usersRepo.getOne(trx, email);
      if (!user) {
        throw new InputError(`email:${email} was not found`);
      }

      if (user.id === Req.userId) {
        throw new LogicError('not_allowed_to_share_with_creator', 'Not allowed to share with creator');
      }

      const limit = 20;
      const plan = await this.plansRepo.getOne(trx, id);

      if (plan?.isShared) {
        const sharedWithCount = await this.planMembersRepo.getUsersCount(trx, id);
        if (sharedWithCount >= limit) {
          throw new LogicError('max_is_20', 'Max is 20');
        }
      } else {
        const plansCount = await this.planMembersRepo.getPlansCount(trx, Req.userId);
        if (plansCount >= limit) {
          throw new LogicError('max_is_20', 'Max is 20');
        }
      }

      await this.planMembersRepo.create(trx, id, user.id);
      await this.suggestedEmailRepo.create(trx, Req.userId, email);

      const currentUser = await this.usersRepo.getOneById(trx, Req.userId);
      if (currentUser?.email) {
        await this.suggestedEmailRepo.create(trx, user.id, currentUser.email);
      }
    });
  }

  async unshare(id: string, email: string): Promise<void> {
    await this.validateUserLoggedIn();
    await this.validateUserOwnsThePlan(id);

    await DB.withTrx(async (trx) => {
      const user = await this.usersRepo.getOne(trx, email);
      if (!user) {
        throw new InputError(`email:${email} was not found`);
      }
      await this.planMembersRepo.delete(trx, id, user.id);
    });
  }

  async leave(id: string): Promise<void> {
    await this.validateUserLoggedIn();
    await DB.withTrx(async (trx) => {
      const deletedRecords = await this.planMembersRepo.delete(trx, id, Req.userId);
      if (deletedRecords === 1) Log.info(`user ${Req.userId} left plan ${id}`);
      else throw new InputError(`userId=${Req.userId} cannot leave planId=${id}`);
    });
  }

  async updateType(id: string, type: string): Promise<void> {
    await this.validateUserOwnsThePlan(id);

    await DB.withTrx(async (trx) => {
      const count = await this.plansRepo.getCount(trx, Req.userId, type);
      if (count >= 100) {
        throw new LogicError('max_is_100', 'Max is 100');
      }

      await this.plansRepo.removeFromOrder(trx, Req.userId, id);
      await this.plansRepo.updateType(trx, Req.userId, id, type);
    });
  }

  async reOrder(type: string, oldOrder: number, newOrder: number): Promise<void> {
    await DB.withTrx(async (trx) => {
      const count = await this.plansRepo.getCount(trx, Req.userId, type);
      if (oldOrder > count || newOrder > count) {
        throw new LogicError('oldOrder and newOrder should be less than ' + count);
      }
      await this.plansRepo.updateOrder(trx, Req.userId, type, oldOrder, newOrder);
    });
  }

  async validateUserOwnsThePlan(planId: string): Promise<void> {
    const plan = await DB.withTrx(async (trx) => {
      return await this.plansRepo.getOne(trx, planId);
    });

    if (!plan) throw new InputError('planId not found');
    if (plan.user.id !== Req.userId) throw new UnauthorizedError('User does not own this plan');
  }

  private async validateUserLoggedIn(): Promise<void> {
    const user = await DB.withTrx(async (trx) => {
      return await this.usersRepo.getOneById(trx, Req.userId);
    });
    if (!user?.email) throw new LogicError('you_are_not_logged_in', 'You are not logged In');
  }
}
