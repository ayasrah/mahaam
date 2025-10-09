/**
 * Validation utility functions
 */

import { InputError } from './errors';

export function required(value: string | null | undefined, name: string): void {
  if (!value || value.trim() === '') {
    throw new InputError(`${name} is required`);
  }
}

export function oneAtLeastRequired(values: (any | null | undefined)[], message: string): void {
  const hasValidValue = values.some((item) => item !== null && item !== undefined && (typeof item !== 'string' || item.trim() !== ''));

  if (!hasValidValue) {
    throw new InputError(message);
  }
}

export function requiredGuid(value: string | null | undefined, name: string): void {
  if (!value || value.trim() === '' || value === '00000000-0000-0000-0000-000000000000') {
    throw new InputError(`${name} is required`);
  }
}

export function requiredBoolean(value: boolean | null | undefined, name: string): void {
  if (value === null || value === undefined) {
    throw new InputError(`${name} is required`);
  }
}

export function requiredNumber(value: number | null | undefined, name: string): void {
  if (value === null || value === undefined) {
    throw new InputError(`${name} is required`);
  }
}

export function requiredObject(value: any | null | undefined, name: string): void {
  if (value === null || value === undefined) {
    throw new InputError(`${name} is required`);
  }
}

export function isIn(item: string, list: string[]): void {
  if (!list.includes(item)) {
    throw new InputError(`${item} is not in [${list.join(',')}]`);
  }
}

export function failIf(condition: boolean, message: string): void {
  if (condition) {
    throw new InputError(message);
  }
}

export function validateEmail(email: string): void {
  required(email, 'email');

  // Email validation regex pattern
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

  if (!emailRegex.test(email)) {
    throw new InputError('Invalid email');
  }
}
