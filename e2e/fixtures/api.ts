import type { APIRequestContext } from '@playwright/test';

const apiURL = process.env.API_URL ?? 'http://localhost:8080';

export function apiClient(request: APIRequestContext) {
  return {
    get: (path: string) => request.get(`${apiURL}${path}`),
    post: (path: string, data: unknown) =>
      request.post(`${apiURL}${path}`, {
        data,
        headers: { 'Content-Type': 'application/json' },
      }),
    put: (path: string, data: unknown) =>
      request.put(`${apiURL}${path}`, {
        data,
        headers: { 'Content-Type': 'application/json' },
      }),
    delete: (path: string) => request.delete(`${apiURL}${path}`),
  };
}
