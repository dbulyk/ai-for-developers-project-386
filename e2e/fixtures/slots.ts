import type { APIRequestContext } from '@playwright/test';
import { expect } from '@playwright/test';
import { apiClient } from './api';
import type { AvailableDay } from '../helpers/slots';

export async function getSlots(
  request: APIRequestContext,
  eventTypeId: string
): Promise<AvailableDay[]> {
  const response = await apiClient(request).get(`/public/event-types/${eventTypeId}/slots`);
  expect(response.ok()).toBeTruthy();
  return (await response.json()) as AvailableDay[];
}
