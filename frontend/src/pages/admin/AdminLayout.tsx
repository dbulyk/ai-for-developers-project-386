import { AppShell, NavLink, Text, Group } from '@mantine/core';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';

export function AdminLayout() {
  const navigate = useNavigate();
  const location = useLocation();

  return (
    <AppShell
      header={{ height: 56 }}
      navbar={{ width: 220, breakpoint: 'sm' }}
      padding="md"
    >
      <AppShell.Header>
        <Group h="100%" px="md">
          <Text fw={700} size="lg">Calendar Booking — Admin</Text>
        </Group>
      </AppShell.Header>

      <AppShell.Navbar p="sm">
        <NavLink
          label="Event Types"
          active={location.pathname.startsWith('/admin/event-types')}
          onClick={() => navigate('/admin/event-types')}
          mb="xs"
        />
        <NavLink
          label="Bookings"
          active={location.pathname.startsWith('/admin/bookings')}
          onClick={() => navigate('/admin/bookings')}
        />
      </AppShell.Navbar>

      <AppShell.Main>
        <Outlet />
      </AppShell.Main>
    </AppShell>
  );
}
