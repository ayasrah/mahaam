export interface User {
  id: string;
  email?: string | null;
  status?: string | null;
  name?: string | null;
}

export interface Device {
  id: string;
  userId: string;
  platform?: string | null;
  fingerprint: string;
  info?: string | null;
  createdAt?: Date | null;
}

export interface SuggestedEmail {
  id: string;
  userId: string;
  email?: string | null;
  createdAt?: Date | null;
}

export interface VerifiedUser {
  userId: string;
  deviceId: string;
  jwt: string;
  userFullName?: string | null;
  email?: string | null;
}

export interface CreatedUser {
  id: string;
  deviceId: string;
  jwt: string;
}
