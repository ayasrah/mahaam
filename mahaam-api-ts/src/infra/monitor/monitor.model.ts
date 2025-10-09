export interface Traffic {
  id: string;
  healthId: string;
  method: string;
  path: string;
  code?: number | null;
  elapsed?: number | null;
  headers?: string | null;
  request?: string | null;
  response?: string | null;
}

export interface TrafficHeaders {
  userId?: string | null;
  deviceId?: string | null;
  appVersion?: string | null;
  appStore?: string | null;
}

export interface Health {
  id: string;
  apiName?: string | null;
  apiVersion?: string | null;
  nodeIP?: string | null;
  nodeName?: string | null;
  envName?: string | null;
}

export interface HealthInfo {
  app: string;
  version: string;
  env: string;
}
