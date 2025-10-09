import { Injectable, Inject } from '@nestjs/common';
import { randomUUID } from 'crypto';
import { DB } from './db';
import { Cache } from './cache';
import { HealthService } from './monitor/health.service';
import { HealthRepo } from './monitor/health.repo';
import Config from './config';
import Log from './log';

@Injectable()
export class Starter {
  constructor(
    @Inject('HealthService') private readonly healthService: HealthService,
    @Inject('HealthRepo') private readonly healthRepo: HealthRepo,
  ) {}

  async start(): Promise<void> {
    await this.initDB();

    // Create health object
    const health = {
      id: randomUUID(),
      apiName: Config.apiName,
      apiVersion: Config.apiVersion,
      envName: Config.envName,
      nodeIP: this.getNodeIP(),
      nodeName: this.getNodeName(),
    };

    // Initialize health service and cache
    await this.healthService.serverStarted();
    Cache.init(health);

    // Log startup message
    const startMsg = `✓ ${Cache.getApiName()}-v${Cache.getApiVersion()}/${Cache.getNodeIP()}-${Cache.getNodeName()} started with healthID=${Cache.getHealthId()}`;
    Log.info(startMsg);

    // Health service already handles starting pulses internally
  }

  private async initDB(): Promise<void> {
    try {
      await DB.init();

      // Extract host from database URL for logging
      const dbUrl = Config.dbUrl;
      const pattern = /:\/\/[^@]+@([^:/]+)/;
      const match = dbUrl.match(pattern);
      const host = match ? match[1] : 'unknown';

      Log.info(`✓ Connected to DB on server ${host}`);
    } catch (error) {
      Log.error(`Failed to connect to database: ${error}`);
      throw error;
    }
  }

  private getNodeIP(): string {
    try {
      // For Node.js, we'll use a simple approach to get local IP
      // This is a simplified version - in production you might want to use a more robust solution
      const os = require('os');
      const interfaces = os.networkInterfaces();

      for (const name of Object.keys(interfaces)) {
        for (const interface_ of interfaces[name]) {
          if (interface_.family === 'IPv4' && !interface_.internal) {
            return interface_.address;
          }
        }
      }

      return '127.0.0.1';
    } catch (error) {
      Log.error(`An error occurred while getting the local IP address: ${error}`);
      return '127.0.0.1';
    }
  }

  private getNodeName(): string {
    try {
      const os = require('os');
      return os.hostname();
    } catch (error) {
      Log.error(`Failed to get node name: ${error}`);
      return 'unknown';
    }
  }
}
