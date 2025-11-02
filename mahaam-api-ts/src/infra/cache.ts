import { Injectable, Logger } from '@nestjs/common';
import { Health } from './monitor/monitor.model';

@Injectable()
export class Cache {
  private static readonly logger = new Logger(Cache.name);
  private static _health: Health | null = null;

  public static init(health: Health): void {
    this._health = health;
    this.logger.log(`Cache initialized with health ID: ${health.id}`);
  }

  public static getNodeIP(): string {
    return this._health?.nodeIP ?? '';
  }

  public static getNodeName(): string {
    return this._health?.nodeName ?? '';
  }

  public static getHealthId(): string {
    return this._health?.id ?? '';
  }
}
