import { useState } from 'react';
import {
  Container, Title, Text, Button, Group, Stack, Modal,
  TextInput, Loader, Center, Alert, Badge, SimpleGrid, Card, Anchor,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { notifications } from '@mantine/notifications';
import { useParams, Link } from 'react-router-dom';
import { usePublicEventTypes, useSlots, useCreateBooking, type AvailableDay } from '../api/hooks';
import { formatSlotTime, formatDateTime } from '../lib/time';

interface BookingFormValues {
  guestName: string;
}

export function BookingPage() {
  const { id } = useParams<{ id: string }>();
  const [selectedDay, setSelectedDay] = useState<AvailableDay | null>(null);
  const [selectedSlot, setSelectedSlot] = useState<string | null>(null);
  const [modalOpen, setModalOpen] = useState(false);

  const { data: eventTypes } = usePublicEventTypes();
  const eventType = eventTypes?.find((e) => e.id === id);

  const { data: availableDays, isLoading, isError, refetch } = useSlots(id ?? '');

  const createBooking = useCreateBooking();

  const form = useForm<BookingFormValues>({
    initialValues: { guestName: '' },
    validate: {
      guestName: (v) => (v.trim().length < 1 ? 'Name is required' : null),
    },
  });

  function openBookingModal(slot: string) {
    setSelectedSlot(slot);
    form.reset();
    setModalOpen(true);
  }

  async function handleSubmit(values: BookingFormValues) {
    if (!selectedSlot || !id) return;
    try {
      await createBooking.mutateAsync({
        eventTypeId: id,
        guestName: values.guestName.trim(),
        startTime: selectedSlot,
      });
      setModalOpen(false);
      setSelectedSlot(null);
      notifications.show({
        title: 'Booking confirmed!',
        message: `Your booking for ${formatDateTime(selectedSlot)} has been confirmed.`,
        color: 'green',
      });
    } catch (err: unknown) {
      const status = (err as { status?: number }).status;
      if (status === 409) {
        notifications.show({
          title: 'Slot already taken',
          message: 'This slot was just booked by someone else. Please choose another time.',
          color: 'orange',
        });
        setModalOpen(false);
        void refetch();
      } else if (status === 404) {
        notifications.show({
          title: 'Not found',
          message: 'Event type not found.',
          color: 'red',
        });
      } else {
        notifications.show({
          title: 'Error',
          message: 'Something went wrong. Please try again.',
          color: 'red',
        });
      }
    }
  }

  if (isLoading) {
    return (
      <Center h={300}>
        <Loader />
      </Center>
    );
  }

  if (isError) {
    return (
      <Container py="xl">
        <Alert color="red" title="Not found">
          Event type not found. <Anchor component={Link} to="/">Go back</Anchor>
        </Alert>
      </Container>
    );
  }

  return (
    <Container py="xl">
      <Anchor component={Link} to="/" mb="md" display="block">
        ← Back to event types
      </Anchor>

      <Title mb="xs">{eventType?.name ?? 'Book a slot'}</Title>
      {eventType && (
        <Group mb="xl">
          <Badge variant="light">{eventType.durationMinutes} min</Badge>
          {eventType.description && <Text c="dimmed">{eventType.description}</Text>}
        </Group>
      )}

      {availableDays && availableDays.length === 0 && (
        <Text c="dimmed">No available slots in the next 14 days.</Text>
      )}

      {availableDays && availableDays.length > 0 && (
        <>
          <Text fw={500} mb="sm">Select a date:</Text>
          <SimpleGrid cols={{ base: 2, sm: 4, md: 7 }} mb="xl">
            {availableDays.map((day) => (
              <Button
                key={day.date}
                variant={selectedDay?.date === day.date ? 'filled' : 'outline'}
                onClick={() => setSelectedDay(day)}
              >
                {day.date}
              </Button>
            ))}
          </SimpleGrid>

          {selectedDay && (
            <>
              <Text fw={500} mb="sm">Available times on {selectedDay.date}:</Text>
              <SimpleGrid cols={{ base: 3, sm: 4, md: 6 }}>
                {selectedDay.slots.map((slot) => (
                  <Card
                    key={slot.startTime}
                    shadow="xs"
                    padding="sm"
                    radius="md"
                    withBorder
                    style={{
                      cursor: slot.status === 'taken' ? 'default' : 'pointer',
                      opacity: slot.status === 'taken' ? 0.4 : 1,
                    }}
                    onClick={() => slot.status === 'free' && openBookingModal(slot.startTime)}
                  >
                    <Text ta="center" fw={500}>
                      {formatSlotTime(slot.startTime)}
                    </Text>
                  </Card>
                ))}
              </SimpleGrid>
            </>
          )}
        </>
      )}

      <Modal
        opened={modalOpen}
        onClose={() => setModalOpen(false)}
        title={selectedSlot ? `Book at ${formatDateTime(selectedSlot)}` : 'Book a slot'}
      >
        <Stack>
          <form onSubmit={form.onSubmit(handleSubmit)}>
            <TextInput
              label="Your name"
              placeholder="Jane Doe"
              required
              {...form.getInputProps('guestName')}
            />
            <Group justify="flex-end" mt="md">
              <Button variant="default" onClick={() => setModalOpen(false)}>
                Cancel
              </Button>
              <Button type="submit" loading={createBooking.isPending}>
                Confirm booking
              </Button>
            </Group>
          </form>
        </Stack>
      </Modal>
    </Container>
  );
}
