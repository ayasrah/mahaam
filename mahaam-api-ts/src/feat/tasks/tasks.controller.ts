import {
  Controller,
  Get,
  Post,
  Delete,
  Patch,
  Inject,
  Param,
  Body,
  Res,
  ParseBoolPipe,
  ParseIntPipe,
} from '@nestjs/common';
import { Response } from 'express';
import { TasksService } from './tasks.service';
import * as rule from '../../infra/rule';

export interface TasksController {
  create(planId: string, title: string, res: Response): Promise<void>;
  delete(planId: string, id: string, res: Response): Promise<void>;
  updateDone(planId: string, id: string, done: boolean, res: Response): Promise<void>;
  updateTitle(id: string, title: string, res: Response): Promise<void>;
  reOrder(planId: string, oldOrder: number, newOrder: number, res: Response): Promise<void>;
  getMany(planId: string, res: Response): Promise<void>;
}

@Controller('plans/:planId/tasks')
export class DefaultTasksController implements TasksController {
  constructor(@Inject('TasksService') private readonly tasksService: TasksService) {}

  @Post()
  async create(@Param('planId') planId: string, @Body('title') title: string, @Res() res: Response) {
    rule.required(planId, 'planId');
    rule.required(title, 'title');
    const id = await this.tasksService.create(planId, title);
    res.status(201).json(id);
  }

  @Delete(':id')
  async delete(@Param('planId') planId: string, @Param('id') id: string, @Res() res: Response) {
    rule.required(planId, 'planId');
    rule.required(id, 'id');
    await this.tasksService.delete(planId, id);
    res.sendStatus(204);
  }

  @Patch(':id/done')
  async updateDone(
    @Param('planId') planId: string,
    @Param('id') id: string,
    @Body('done', ParseBoolPipe) done: boolean,
    @Res() res: Response,
  ) {
    rule.required(planId, 'planId');
    rule.required(id, 'id');
    rule.requiredBoolean(done, 'done');
    await this.tasksService.updateDone(planId, id, done);
    res.sendStatus(200);
  }

  @Patch(':id/title')
  async updateTitle(@Param('id') id: string, @Body('title') title: string, @Res() res: Response) {
    rule.required(id, 'id');
    rule.required(title, 'title');
    await this.tasksService.updateTitle(id, title);
    res.sendStatus(200);
  }

  @Patch('reorder')
  async reOrder(
    @Param('planId') planId: string,
    @Body('oldOrder', ParseIntPipe) oldOrder: number,
    @Body('newOrder', ParseIntPipe) newOrder: number,
    @Res() res: Response,
  ) {
    rule.required(planId, 'planId');
    rule.requiredNumber(oldOrder, 'oldOrder');
    rule.requiredNumber(newOrder, 'newOrder');
    await this.tasksService.reOrder(planId, oldOrder, newOrder);
    res.sendStatus(200);
  }

  @Get()
  async getMany(@Param('planId') planId: string, @Res() res: Response) {
    rule.required(planId, 'planId');
    const result = await this.tasksService.getList(planId);
    res.status(200).json(result);
  }
}
