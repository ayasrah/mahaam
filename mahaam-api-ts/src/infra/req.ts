import { AsyncLocalStorage } from 'async_hooks';

class ReqCtx {
  private static readonly storage = new AsyncLocalStorage<Map<string, any>>();

  /**
   * Middlewares below can get this value. The above middlewares cannot
   */
  public static set<T>(name: string, data: T): void {
    const context = this.storage.getStore();
    if (context) context.set(name, data);
  }

  public static get<T>(name: string): T | undefined {
    const context = this.storage.getStore();
    if (context) return context.get(name);
    return undefined;
  }

  /**
   * Initialize async context - should be called at the beginning of request processing
   * This establishes the async context that persists through the entire request chain
   * It clears the context after the callback is executed
   */
  public static run<T>(callback: () => T): T {
    const contextMap = new Map<string, any>();
    return this.storage.run(contextMap, callback);
  }
}

export class Req {
  public static run<T>(callback: () => T): T {
    return ReqCtx.run(callback);
  }

  public static get trafficId(): string {
    return ReqCtx.get<string>('trafficId') || '';
  }

  public static set trafficId(value: string) {
    ReqCtx.set('trafficId', value);
  }

  public static get userId(): string {
    return ReqCtx.get<string>('userId') || '';
  }

  public static set userId(value: string) {
    ReqCtx.set('userId', value);
  }

  public static get deviceId(): string {
    return ReqCtx.get<string>('deviceId') || '';
  }

  public static set deviceId(value: string) {
    ReqCtx.set('deviceId', value);
  }

  public static get appStore(): string {
    return ReqCtx.get<string>('appStore') || '';
  }

  public static set appStore(value: string) {
    ReqCtx.set('appStore', value);
  }

  public static get appVersion(): string {
    return ReqCtx.get<string>('appVersion') || '';
  }

  public static set appVersion(value: string) {
    ReqCtx.set('appVersion', value);
  }

  public static get isLoggedIn(): boolean {
    return ReqCtx.get<boolean>('isLoggedIn') || false;
  }

  public static set isLoggedIn(value: boolean) {
    ReqCtx.set('isLoggedIn', value);
  }

  public static get startTime(): number {
    return ReqCtx.get<number>('startTime') || 0;
  }

  public static set startTime(value: number) {
    ReqCtx.set('startTime', value);
  }

  public static clear(): void {
    ReqCtx.set('trafficId', '');
    ReqCtx.set('userId', '');
    ReqCtx.set('deviceId', '');
    ReqCtx.set('appStore', '');
    ReqCtx.set('appVersion', '');
    ReqCtx.set('isLoggedIn', false);
    ReqCtx.set('startTime', 0);
  }
}
