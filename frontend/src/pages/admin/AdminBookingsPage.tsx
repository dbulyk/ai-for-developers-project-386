import { useState } from 'react';
import {
  Title, Button, Table, Group, Modal, Text, Loader, Center, Alert,
} from '@mantine/core';
import { notifications } from '@mantine/notifications';
import { useAdminBookings, useCancelBooking, type Booking } from '../../api/hooks';
import { formatDateTime } from '../../lib/time';

export function AdminBookingsPage() {
  const { data: bookings, isLoading, isError } = useAdminBookings();
  const cancelMutation = useCancelBooking();
  const [cancelTarget, setCancelTarget] = useState<Booking | null>(null);

  async function handleCancel() {
    if (!cancelTarget) return;
    try {
      await cancelMutation.mutateAsync(cancelTarget.id);
      setCancelTarget(null);
      notifications.show({ title: 'Cancelled', message: 'Booking has been cancelled.', color: 'green' });
    } catch {
      notifications.show({ title: 'Error', message: 'Failed to cancel booking.', color: 'red' });
    }
  }

  if (isLoading) return <Center h={200}><Loader /></Center>;
  if (isError) return <Alert color="red" title="Error">Failed to load bookings.</Alert>;

  return (
    <>
      <Title order={2} mb="md" data-testid="bookings-title">Bookings</Title>

      {bookings && bookings.length === 0 && (
        <Text c="dimmed">No upcoming bookings.</Text>
      )}

      {bookings && bookings.length > 0 && (
        <Table striped highlightOnHover withTableBorder data-testid="bookings-table">
          <Table.Thead>
            <Table.Tr>
              <Table.Th>Guest</Table.Th>
              <Table.Th>Event Type ID</Table.Th>
              <Table.Th>Start</Table.Th>
              <Table.Th>End</Table.Th>
              <Table.Th>Action</Table.Th>
            </Table.Tr>
          </Table.Thead>
          <Table.Tbody>
            {bookings.map((b) => (
              <Table.Tr key={b.id} data-testid="booking-row">
                <Table.Td data-testid="booking-row-guest">{b.guestName}</Table.Td>
                <Table.Td>{b.eventTypeId}</Table.Td>
                <Table.Td>{formatDateTime(b.startTime)}</Table.Td>
                <Table.Td>{formatDateTime(b.endTime)}</Table.Td>
                <Table.Td>
                  <Button size="xs" color="red" variant="light" onClick={() => setCancelTarget(b)} data-testid="cancel-booking-button">
                    Cancel
                  </Button>
                </Table.Td>
              </Table.Tr>
            ))}
          </Table.Tbody>
        </Table>
      )}

      <Modal opened={cancelTarget !== null} onClose={() => setCancelTarget(null)} title="Cancel Booking">
        <div data-testid="cancel-booking-modal">
          <Text mb="md">
            Cancel booking for <strong>{cancelTarget?.guestName}</strong> at{' '}
            {cancelTarget ? formatDateTime(cancelTarget.startTime) : ''}?
          </Text>
          <Group justify="flex-end">
            <Button variant="default" onClick={() => setCancelTarget(null)}>Keep it</Button>
            <Button color="red" loading={cancelMutation.isPending} onClick={handleCancel} data-testid="confirm-cancel-booking">Cancel booking</Button>
          </Group>
        </div>
      </Modal>
    </>
  );
}
