import { Injectable, Logger } from '@nestjs/common';
import { Health } from './monitor/monitor.model';

@Injectable()
export class Cache {
  private static _health: Health | null = null;

  public static init(health: Health): void {
    this._health = health;
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
