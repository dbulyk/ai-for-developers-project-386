import type { APIRequestContext } from '@playwright/test';
import { expect } from '@playwright/test';
import { apiClient } from './api';

export interface Booking {
  id: string;
  eventTypeId: string;
  guestName: string;
  startTime: string;
  endTime: string;
  createdAt: string;
}

export interface CreateBookingRequest {
  eventTypeId: string;
  guestName: string;
  startTime: string;
}

export async function createBooking(
  request: APIRequestContext,
  body: CreateBookingRequest
): Promise<Booking> {
  const response = await apiClient(request).post('/public/bookings', body);
  expect(response.ok()).toBeTruthy();
  return (await response.json()) as Booking;
}

export async function deleteBooking(request: APIRequestContext, id: string): Promise<void> {
  const response = await apiClient(request).delete(`/admin/bookings/${id}`);
  expect([204, 404]).toContain(response.status());
}

export async function listBookings(request: APIRequestContext): Promise<Booking[]> {
  const response = await apiClient(request).get('/admin/bookings');
  expect(response.ok()).toBeTruthy();
  return (await response.json()) as Booking[];
}

export async function findBookingId(
  request: APIRequestContext,
  filters: { eventTypeId?: string; guestName?: string; startTime?: string }
): Promise<string | undefined> {
  const bookings = await listBookings(request);
  const match = bookings.find(
    (b) =>
      (!filters.eventTypeId || b.eventTypeId === filters.eventTypeId) &&
      (!filters.guestName || b.guestName === filters.guestName) &&
      (!filters.startTime || b.startTime === filters.startTime)
  );
  return match?.id;
}
