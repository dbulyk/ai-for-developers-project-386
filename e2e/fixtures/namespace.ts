import { randomBytes } from 'crypto';

export interface Namespace {
  prefix: string;
  name: (base: string) => string;
}

export function createNamespace(label: string): Namespace {
  const prefix = `e2e-${label}-${Date.now()}-${randomBytes(4).toString('hex')}`;
  return {
    prefix,
    name: (base: string) => `${prefix}::${base}`,
  };
}
