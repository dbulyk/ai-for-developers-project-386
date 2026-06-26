import createClient from 'openapi-fetch';
import type { paths } from './schema.d.ts';

export const client = createClient<paths>({
  baseUrl: import.meta.env.VITE_API_BASE_URL,
});
