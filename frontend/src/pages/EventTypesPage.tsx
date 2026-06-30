import { Container, Title, Text, Card, SimpleGrid, Badge, Group, Button, Loader, Center, Alert } from '@mantine/core';
import { Link } from 'react-router-dom';
import { usePublicEventTypes } from '../api/hooks';

export function EventTypesPage() {
  const { data: eventTypes, isLoading, isError } = usePublicEventTypes();

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
        <Alert color="red" title="Error">
          Failed to load event types. Please try again later.
        </Alert>
      </Container>
    );
  }

  return (
    <Container py="xl">
      <Title mb="xs" data-testid="event-types-title">Book a Meeting</Title>
      <Text c="dimmed" mb="xl">
        Choose an event type to see available time slots.
      </Text>

      {eventTypes && eventTypes.length === 0 && (
        <Text c="dimmed">No event types available yet.</Text>
      )}

      <SimpleGrid cols={{ base: 1, sm: 2, md: 3 }}>
        {eventTypes?.map((et) => (
          <Card key={et.id} shadow="sm" padding="lg" radius="md" withBorder data-testid="event-type-card">
            <Group justify="space-between" mb="xs">
              <Text fw={600} size="lg" data-testid="event-type-name">{et.name}</Text>
              <Badge variant="light" data-testid="event-type-duration">{et.durationMinutes} min</Badge>
            </Group>
            {et.description && (
              <Text size="sm" c="dimmed" mb="md">
                {et.description}
              </Text>
            )}
            <Button component={Link} to={`/event-types/${et.id}`} fullWidth mt="auto" data-testid="event-type-book-button">
              Book
            </Button>
          </Card>
        ))}
      </SimpleGrid>
    </Container>
  );
}
