import { Injectable, Inject } from '@nestjs/common';
import { HealthRepo } from './health.repo';
import { Cache } from '../cache';
import Config from '../config';
import Log from '../log';

export interface HealthService {
  serverStarted(): Promise<void>;
  serverStopped(): Promise<void>;
}

@Injectable()
export class DefaultHealthService implements HealthService {
  private healthID: string | null = null;
  private pulseInterval: NodeJS.Timeout | null = null;

  constructor(@Inject('HealthRepo') private readonly healthRepo: HealthRepo) {}

  async serverStarted(): Promise<void> {
    this.healthID = await this.healthRepo.create();

    // Initialize cache with health data
    const healthData = await this.healthRepo.getById(this.healthID);
    if (healthData) {
      Cache.init(healthData);
    }

    // Wait 2 seconds before starting pulses
    setTimeout(() => {
      this.startSendingPulses();
    }, 2000);
  }

  private startSendingPulses(): void {
    if (!this.healthID) {
      Log.error('Cannot start sending pulses: healthID is null');
      return;
    }

    // Send pulse every minute (60,000 ms)
    this.pulseInterval = setInterval(async () => {
      try {
        if (this.healthID) {
          await this.healthRepo.updatePulse(this.healthID);
        }
      } catch (error) {
        Log.error(error instanceof Error ? error.toString() : String(error));
      }
    }, 60 * 1000); // 1 minute
  }

  async serverStopped(): Promise<void> {
    // Stop sending pulses
    if (this.pulseInterval) {
      clearInterval(this.pulseInterval);
      this.pulseInterval = null;
    }

    try {
      if (this.healthID) {
        await this.healthRepo.updateStopped(this.healthID);
        const stopMsg = `✓ ${Config.apiName}-v${Config.apiVersion}/${Cache.getNodeIP()}-${Cache.getNodeName()} stopped with healthID=${this.healthID}`;

        Log.info(stopMsg);
      }
    } catch (error) {
      Log.error(error instanceof Error ? error.toString() : String(error));
    }
  }
}
