import { Controller, Post, Body, Res } from '@nestjs/common';
import { Response } from 'express';
import Log from '../log';

export interface AuditController {
  error(error: string, res: Response): Promise<void>;
  info(info: string, res: Response): Promise<void>;
}

@Controller('audit')
export class DefaultAuditController implements AuditController {
  @Post('error')
  async error(@Body('error') error: string, @Res() res: Response): Promise<void> {
    Log.error('mahaam-mb: ' + error);
    res.sendStatus(201);
  }

  @Post('info')
  async info(@Body('info') info: string, @Res() res: Response): Promise<void> {
    Log.info('mahaam-mb: ' + info);
    res.sendStatus(201);
  }
}
