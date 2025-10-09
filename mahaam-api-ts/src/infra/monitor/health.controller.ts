import { Controller, Get, Res } from '@nestjs/common';
import { Response } from 'express';
import { Health } from './monitor.model';
import Config from '../config';
import { Cache } from '../cache';

export interface HealthController {
  getInfo(res: Response): Promise<void>;
}

@Controller()
export class DefaultHealthController implements HealthController {
  constructor() {}

  @Get('health')
  async getInfo(@Res() res: Response) {
    const result: Health = {
      id: Cache.getHealthId(),
      apiName: Config.apiName,
      apiVersion: Config.apiVersion,
      envName: Config.envName,
      nodeIP: Cache.getNodeIP(),
      nodeName: Cache.getNodeName(),
    };
    res.status(200).json(result);
  }
}
