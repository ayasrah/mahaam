import { Injectable } from '@nestjs/common';
import { DefaultLogRepo, LogRepo } from './monitor/log.repo';
import * as winston from 'winston';
import config from './config';
import { Req } from './req';

// export default interface Log {
//   info(info: string): Promise<void>;
//   error(trafficIdOrError: string | null, error?: string): Promise<void>;
// }

@Injectable()
export default class Log {
  private static logger: winston.Logger;
  private static logRepo: LogRepo;
  private static initialized = false;

  static init(): void {
    if (Log.initialized) return;

    // Initialize logRepo
    Log.logRepo = new DefaultLogRepo();

    const transports: winston.transport[] = [];

    // Add console transport with appropriate formatting
    if (process.env.NODE_ENV !== 'production') {
      transports.push(
        new winston.transports.Console({
          format: winston.format.combine(
            winston.format.timestamp({
              format: 'YYYY-MM-DD HH:mm:ss,SSS',
            }),
            winston.format.printf(({ timestamp, level, message }) => {
              return `${timestamp} ${level.toUpperCase()} ${message}`;
            }),
          ),
        }),
      );
    } else {
      transports.push(new winston.transports.Console());
    }

    if (config.logFile) {
      transports.push(
        new winston.transports.File({
          filename: config.logFile,
          format: winston.format.combine(
            winston.format.timestamp({
              format: 'YYYY-MM-DD HH:mm:ss,SSS',
            }),
            winston.format.printf(({ timestamp, level, message }) => {
              return `${timestamp} ${level.toUpperCase()} ${message}`;
            }),
          ),
        }),
      );
    }

    Log.logger = winston.createLogger({
      level: 'info',
      format: winston.format.combine(
        winston.format.timestamp({
          format: 'YYYY-MM-DD HH:mm:ss,SSS',
        }),
        winston.format.printf(({ timestamp, level, message }) => {
          return `${timestamp} ${level.toUpperCase()} ${message}`;
        }),
      ),
      transports,
    });

    Log.initialized = true;
  }

  static async info(info: string): Promise<void> {
    const trafficId = Req.trafficId;
    const message = trafficId ? `TrafficId: ${trafficId}, ${info}` : info;

    Log.logger.info(message);
    await Log.logRepo.create(trafficId, 'Info', info);
  }

  static async error(error: string): Promise<void> {
    const trafficId = Req.trafficId;
    const message = trafficId ? `TrafficId: ${trafficId}, ${error}` : error;

    Log.logger.error(message);
    await Log.logRepo.create(trafficId, 'Error', error);
  }
}
