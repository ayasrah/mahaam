import { Injectable } from '@nestjs/common';
import { DB } from '../db';
import { Cache } from '../cache';
import config from '../config';
import { randomUUID } from 'crypto';
import { Health } from './monitor.model';

export interface HealthRepo {
  create(): Promise<string>;
  getById(id: string): Promise<Health | null>;
  updatePulse(id: string): Promise<void>;
  updateStopped(id: string): Promise<void>;
}

@Injectable()
export class DefaultHealthRepo implements HealthRepo {
  async create(): Promise<string> {
    const id = randomUUID();
    await DB.sql`INSERT INTO monitor.health (id, api_name, api_version, env_name, node_ip, node_name, started_at) VALUES (${id}, ${config.apiName}, ${config.apiVersion}, ${config.envName}, ${Cache.getNodeIP()}, ${Cache.getNodeName()}, current_timestamp)`;
    return id;
  }

  async getById(id: string): Promise<Health | null> {
    const result = await DB.sql<Health[]>`
      SELECT id, api_name as "apiName", api_version as "apiVersion", env_name as "envName", node_ip as "nodeIP", node_name as "nodeName"
      FROM monitor.health 
      WHERE id = ${id}
    `;
    return result.length > 0 ? result[0] : null;
  }

  async updatePulse(id: string): Promise<void> {
    await DB.sql`UPDATE monitor.health SET pulsed_at = current_timestamp WHERE id = ${id}`;
  }

  async updateStopped(id: string): Promise<void> {
    await DB.sql`UPDATE monitor.health SET stopped_at = current_timestamp WHERE id = ${id}`;
  }
}
