import type { APIRequestContext } from '@playwright/test';
import { expect } from '@playwright/test';
import { apiClient } from './api';
import type { Namespace } from './namespace';

export interface EventType {
  id: string;
  name: string;
  description: string;
  durationMinutes: number;
}

export interface CreateEventTypeRequest {
  name: string;
  description: string;
  durationMinutes: number;
}

export async function createEventType(
  request: APIRequestContext,
  ns: Namespace,
  overrides: Partial<CreateEventTypeRequest> = {}
): Promise<EventType> {
  const body: CreateEventTypeRequest = {
    name: ns.name('30-min call'),
    description: 'E2E test event type',
    durationMinutes: 30,
    ...overrides,
  };

  const response = await apiClient(request).post('/admin/event-types', body);
  expect(response.ok()).toBeTruthy();
  return (await response.json()) as EventType;
}

export async function deleteEventType(request: APIRequestContext, id: string): Promise<void> {
  const response = await apiClient(request).delete(`/admin/event-types/${id}`);
  expect([204, 404]).toContain(response.status());
}

export async function listEventTypes(request: APIRequestContext): Promise<EventType[]> {
  const response = await apiClient(request).get('/admin/event-types');
  expect(response.ok()).toBeTruthy();
  return (await response.json()) as EventType[];
}
