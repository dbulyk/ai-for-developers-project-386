import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { components } from './schema.d.ts';
import { client } from './client';

export type EventType = components['schemas']['EventType'];
export type Booking = components['schemas']['Booking'];
export type AvailableDay = components['schemas']['AvailableDay'];
export type CreateBookingRequest = components['schemas']['CreateBookingRequest'];
export type CreateEventTypeRequest = components['schemas']['CreateEventTypeRequest'];
export type UpdateEventTypeRequest = components['schemas']['UpdateEventTypeRequest'];

// --- Public ---

export function usePublicEventTypes() {
  return useQuery({
    queryKey: ['public', 'event-types'],
    queryFn: async () => {
      const { data } = await client.GET('/public/event-types');
      return data ?? [];
    },
  });
}

export function useSlots(eventTypeId: string) {
  return useQuery({
    queryKey: ['public', 'slots', eventTypeId],
    queryFn: async () => {
      const { data, response } = await client.GET('/public/event-types/{id}/slots', {
        params: { path: { id: eventTypeId } },
      });
      if (response.status === 404) throw Object.assign(new Error('Event type not found'), { status: 404 });
      return data ?? [];
    },
  });
}

export function useCreateBooking() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (body: CreateBookingRequest) => {
      const { data, response } = await client.POST('/public/bookings', { body });
      if (response.status === 409) throw Object.assign(new Error('Slot already taken'), { status: 409 });
      if (response.status === 404) throw Object.assign(new Error('Event type not found'), { status: 404 });
      if (!data) throw new Error('Booking failed');
      return data;
    },
    onSuccess: (_data, variables) => {
      void qc.invalidateQueries({ queryKey: ['public', 'slots', variables.eventTypeId] });
    },
  });
}

// --- Admin ---

export function useAdminEventTypes() {
  return useQuery({
    queryKey: ['admin', 'event-types'],
    queryFn: async () => {
      const { data } = await client.GET('/admin/event-types');
      return data ?? [];
    },
  });
}

export function useCreateEventType() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (body: CreateEventTypeRequest) => {
      const { data } = await client.POST('/admin/event-types', { body });
      if (!data) throw new Error('Failed to create event type');
      return data;
    },
    onSuccess: () => void qc.invalidateQueries({ queryKey: ['admin', 'event-types'] }),
  });
}

export function useUpdateEventType() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async ({ id, body }: { id: string; body: UpdateEventTypeRequest }) => {
      const { data, response } = await client.PUT('/admin/event-types/{id}', {
        params: { path: { id } },
        body,
      });
      if (response.status === 404) throw Object.assign(new Error('Event type not found'), { status: 404 });
      if (!data) throw new Error('Failed to update event type');
      return data;
    },
    onSuccess: () => void qc.invalidateQueries({ queryKey: ['admin', 'event-types'] }),
  });
}

export function useDeleteEventType() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (id: string) => {
      const { response } = await client.DELETE('/admin/event-types/{id}', {
        params: { path: { id } },
      });
      if (response.status === 404) throw Object.assign(new Error('Event type not found'), { status: 404 });
    },
    onSuccess: () => void qc.invalidateQueries({ queryKey: ['admin', 'event-types'] }),
  });
}

export function useAdminBookings() {
  return useQuery({
    queryKey: ['admin', 'bookings'],
    queryFn: async () => {
      const { data } = await client.GET('/admin/bookings');
      return data ?? [];
    },
  });
}

export function useCancelBooking() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (id: string) => {
      const { response } = await client.DELETE('/admin/bookings/{id}', {
        params: { path: { id } },
      });
      if (response.status === 404) throw Object.assign(new Error('Booking not found'), { status: 404 });
    },
    onSuccess: () => void qc.invalidateQueries({ queryKey: ['admin', 'bookings'] }),
  });
}
