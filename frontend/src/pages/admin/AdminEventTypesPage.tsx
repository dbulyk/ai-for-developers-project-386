import { useState } from 'react';
import {
  Title, Button, Table, Group, Modal, TextInput,
  NumberInput, Textarea, Stack, Text, Loader, Center, Alert,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { notifications } from '@mantine/notifications';
import {
  useAdminEventTypes, useCreateEventType, useUpdateEventType, useDeleteEventType,
  type EventType, type CreateEventTypeRequest,
} from '../../api/hooks';

interface EventTypeFormValues {
  name: string;
  description: string;
  durationMinutes: number;
}

const defaultValues: EventTypeFormValues = { name: '', description: '', durationMinutes: 30 };

export function AdminEventTypesPage() {
  const { data: eventTypes, isLoading, isError } = useAdminEventTypes();
  const createMutation = useCreateEventType();
  const updateMutation = useUpdateEventType();
  const deleteMutation = useDeleteEventType();

  const [createOpen, setCreateOpen] = useState(false);
  const [editTarget, setEditTarget] = useState<EventType | null>(null);
  const [deleteTarget, setDeleteTarget] = useState<EventType | null>(null);

  const form = useForm<EventTypeFormValues>({
    initialValues: defaultValues,
    validate: {
      name: (v) => (v.trim().length < 1 ? 'Name is required' : null),
      durationMinutes: (v) => (v < 1 ? 'Must be at least 1 minute' : null),
    },
  });

  function openCreate() {
    form.setValues(defaultValues);
    setCreateOpen(true);
  }

  function openEdit(et: EventType) {
    form.setValues({ name: et.name, description: et.description, durationMinutes: et.durationMinutes });
    setEditTarget(et);
  }

  async function handleCreate(values: EventTypeFormValues) {
    const body: CreateEventTypeRequest = {
      name: values.name.trim(),
      description: values.description,
      durationMinutes: values.durationMinutes,
    };
    try {
      await createMutation.mutateAsync(body);
      setCreateOpen(false);
      notifications.show({ title: 'Created', message: 'Event type created.', color: 'green' });
    } catch {
      notifications.show({ title: 'Error', message: 'Failed to create event type.', color: 'red' });
    }
  }

  async function handleUpdate(values: EventTypeFormValues) {
    if (!editTarget) return;
    try {
      await updateMutation.mutateAsync({
        id: editTarget.id,
        body: { name: values.name.trim(), description: values.description, durationMinutes: values.durationMinutes },
      });
      setEditTarget(null);
      notifications.show({ title: 'Updated', message: 'Event type updated.', color: 'green' });
    } catch {
      notifications.show({ title: 'Error', message: 'Failed to update event type.', color: 'red' });
    }
  }

  async function handleDelete() {
    if (!deleteTarget) return;
    try {
      await deleteMutation.mutateAsync(deleteTarget.id);
      setDeleteTarget(null);
      notifications.show({ title: 'Deleted', message: 'Event type deleted.', color: 'green' });
    } catch {
      notifications.show({ title: 'Error', message: 'Failed to delete event type.', color: 'red' });
    }
  }

  if (isLoading) return <Center h={200}><Loader /></Center>;
  if (isError) return <Alert color="red" title="Error">Failed to load event types.</Alert>;

  const isPending = createMutation.isPending || updateMutation.isPending;

  return (
    <>
      <Group justify="space-between" mb="md">
        <Title order={2}>Event Types</Title>
        <Button onClick={openCreate}>+ New event type</Button>
      </Group>

      {eventTypes && eventTypes.length === 0 && (
        <Text c="dimmed">No event types yet.</Text>
      )}

      {eventTypes && eventTypes.length > 0 && (
        <Table striped highlightOnHover withTableBorder>
          <Table.Thead>
            <Table.Tr>
              <Table.Th>Name</Table.Th>
              <Table.Th>Description</Table.Th>
              <Table.Th>Duration</Table.Th>
              <Table.Th>Actions</Table.Th>
            </Table.Tr>
          </Table.Thead>
          <Table.Tbody>
            {eventTypes.map((et) => (
              <Table.Tr key={et.id}>
                <Table.Td>{et.name}</Table.Td>
                <Table.Td>{et.description}</Table.Td>
                <Table.Td>{et.durationMinutes} min</Table.Td>
                <Table.Td>
                  <Group gap="xs">
                    <Button size="xs" variant="light" onClick={() => openEdit(et)}>Edit</Button>
                    <Button size="xs" color="red" variant="light" onClick={() => setDeleteTarget(et)}>Delete</Button>
                  </Group>
                </Table.Td>
              </Table.Tr>
            ))}
          </Table.Tbody>
        </Table>
      )}

      {/* Create modal */}
      <Modal opened={createOpen} onClose={() => setCreateOpen(false)} title="New Event Type">
        <form onSubmit={form.onSubmit(handleCreate)}>
          <Stack>
            <TextInput label="Name" placeholder="e.g. 30-min call" required {...form.getInputProps('name')} />
            <Textarea label="Description" placeholder="Optional description" {...form.getInputProps('description')} />
            <NumberInput label="Duration (minutes)" min={1} required {...form.getInputProps('durationMinutes')} />
            <Group justify="flex-end">
              <Button variant="default" onClick={() => setCreateOpen(false)}>Cancel</Button>
              <Button type="submit" loading={isPending}>Create</Button>
            </Group>
          </Stack>
        </form>
      </Modal>

      {/* Edit modal */}
      <Modal opened={editTarget !== null} onClose={() => setEditTarget(null)} title="Edit Event Type">
        <form onSubmit={form.onSubmit(handleUpdate)}>
          <Stack>
            <TextInput label="Name" required {...form.getInputProps('name')} />
            <Textarea label="Description" {...form.getInputProps('description')} />
            <NumberInput label="Duration (minutes)" min={1} required {...form.getInputProps('durationMinutes')} />
            <Group justify="flex-end">
              <Button variant="default" onClick={() => setEditTarget(null)}>Cancel</Button>
              <Button type="submit" loading={isPending}>Save</Button>
            </Group>
          </Stack>
        </form>
      </Modal>

      {/* Delete confirm modal */}
      <Modal opened={deleteTarget !== null} onClose={() => setDeleteTarget(null)} title="Delete Event Type">
        <Text mb="md">
          Delete <strong>{deleteTarget?.name}</strong>? This cannot be undone.
        </Text>
        <Group justify="flex-end">
          <Button variant="default" onClick={() => setDeleteTarget(null)}>Cancel</Button>
          <Button color="red" loading={deleteMutation.isPending} onClick={handleDelete}>Delete</Button>
        </Group>
      </Modal>
    </>
  );
}
