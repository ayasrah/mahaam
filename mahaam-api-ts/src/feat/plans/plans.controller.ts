import {
  Controller,
  Get,
  Post,
  Put,
  Delete,
  Patch,
  Inject,
  Query,
  Param,
  Body,
  Res,
  ParseIntPipe,
} from '@nestjs/common';
import { Response } from 'express';
import { PlansService } from './plans.service';
import { PlanType, PlanIn } from './plans.model';
import * as rule from '../../infra/rule';

export interface PlansController {
  create(plan: PlanIn, res: Response): Promise<void>;
  update(plan: PlanIn, res: Response): Promise<void>;
  delete(id: string, res: Response): Promise<void>;
  share(id: string, email: string, res: Response): Promise<void>;
  unshare(id: string, email: string, res: Response): Promise<void>;
  leave(id: string, res: Response): Promise<void>;
  updateType(id: string, type: string, res: Response): Promise<void>;
  reOrder(type: string, oldOrder: number, newOrder: number, res: Response): Promise<void>;
  getOne(planId: string, res: Response): Promise<void>;
  getMany(res: Response, type?: string): Promise<void>;
}

@Controller('plans')
export class DefaultPlansController implements PlansController {
  constructor(@Inject('PlansService') private readonly plansService: PlansService) {}

  @Post()
  async create(@Body() plan: PlanIn, @Res() res: Response) {
    rule.oneAtLeastRequired([plan.title, plan.starts, plan.ends], 'title or starts or ends is required');
    const id = await this.plansService.create(plan);
    res.status(201).json(id);
  }

  @Put()
  async update(@Body() plan: PlanIn, @Res() res: Response) {
    rule.required(plan.id, 'Id');
    rule.oneAtLeastRequired([plan.title, plan.starts, plan.ends], 'title or starts or ends is required');
    await this.plansService.update(plan);
    res.sendStatus(200);
  }

  @Delete(':id')
  async delete(@Param('id') id: string, @Res() res: Response) {
    rule.required(id, 'id');
    await this.plansService.delete(id);
    res.sendStatus(204);
  }

  @Patch(':id/share')
  async share(@Param('id') id: string, @Body('email') email: string, @Res() res: Response) {
    rule.required(id, 'id');
    rule.required(email, 'email');
    await this.plansService.share(id, email);
    res.sendStatus(200);
  }

  @Patch(':id/unshare')
  async unshare(@Param('id') id: string, @Body('email') email: string, @Res() res: Response) {
    rule.required(id, 'id');
    rule.required(email, 'email');
    await this.plansService.unshare(id, email);
    res.sendStatus(200);
  }

  @Patch(':id/leave')
  async leave(@Param('id') id: string, @Res() res: Response) {
    rule.required(id, 'id');
    await this.plansService.leave(id);
    res.sendStatus(200);
  }

  @Patch(':id/type')
  async updateType(@Param('id') id: string, @Body('type') type: string, @Res() res: Response) {
    rule.required(id, 'id');
    rule.required(type, 'type');
    rule.isIn(type, PlanType.All);
    await this.plansService.updateType(id, type);
    res.sendStatus(200);
  }

  @Patch('reorder')
  async reOrder(
    @Body('type') type: string,
    @Body('oldOrder', ParseIntPipe) oldOrder: number,
    @Body('newOrder', ParseIntPipe) newOrder: number,
    @Res() res: Response,
  ) {
    rule.required(type, 'type');
    rule.isIn(type, PlanType.All);
    rule.requiredNumber(oldOrder, 'oldOrder');
    rule.requiredNumber(newOrder, 'newOrder');
    await this.plansService.reOrder(type, oldOrder, newOrder);
    res.sendStatus(200);
  }

  @Get(':planId')
  async getOne(@Param('planId') planId: string, @Res() res: Response) {
    rule.required(planId, 'planId');
    const plan = await this.plansService.getOne(planId);
    res.status(200).json(plan);
  }

  @Get()
  async getMany(@Res() res: Response, @Query('type') type?: string) {
    if (!type) type = PlanType.Main;
    else rule.isIn(type, PlanType.All);

    const plans = await this.plansService.getMany(type);
    res.status(200).json(plans);
  }
}
