import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import { Routes, Route } from 'react-router-dom';
import { BookingPage } from '../pages/BookingPage';
import { renderWithProviders } from './utils';
import * as hooks from '../api/hooks';

vi.mock('../api/hooks', async (importOriginal) => {
  const actual = await importOriginal<typeof hooks>();
  return { ...actual };
});

const mockEventTypes: hooks.EventType[] = [
  { id: 'et-1', name: '30-min call', description: 'Quick call', durationMinutes: 30 },
];

const mockAvailableDays: hooks.AvailableDay[] = [
  { date: '2026-06-30', slots: ['2026-06-30T06:00:00Z', '2026-06-30T06:30:00Z'] },
];

function renderBookingPage() {
  return renderWithProviders(
    <Routes>
      <Route path="/event-types/:id" element={<BookingPage />} />
    </Routes>,
    { initialEntries: ['/event-types/et-1'] },
  );
}

describe('BookingFlow', () => {
  it('shows available slots and allows booking', async () => {
    const user = userEvent.setup();

    vi.spyOn(hooks, 'usePublicEventTypes').mockReturnValue({
      data: mockEventTypes,
      isLoading: false,
      isError: false,
    } as ReturnType<typeof hooks.usePublicEventTypes>);

    vi.spyOn(hooks, 'useSlots').mockReturnValue({
      data: mockAvailableDays,
      isLoading: false,
      isError: false,
      refetch: vi.fn(),
    } as unknown as ReturnType<typeof hooks.useSlots>);

    const mutateAsync = vi.fn().mockResolvedValue({
      id: 'b-1',
      eventTypeId: 'et-1',
      guestName: 'Alice',
      startTime: '2026-06-30T06:00:00Z',
      endTime: '2026-06-30T06:30:00Z',
      createdAt: '2026-06-26T10:00:00Z',
    });

    vi.spyOn(hooks, 'useCreateBooking').mockReturnValue({
      mutateAsync,
      isPending: false,
    } as unknown as ReturnType<typeof hooks.useCreateBooking>);

    renderBookingPage();

    expect(screen.getByText('2026-06-30')).toBeInTheDocument();
    await user.click(screen.getByText('2026-06-30'));

    expect(screen.getAllByText(/^\d{2}:\d{2}$/).length).toBeGreaterThan(0);
    await user.click(screen.getAllByText(/^\d{2}:\d{2}$/)[0]);

    await waitFor(() => {
      expect(screen.getByLabelText(/your name/i)).toBeInTheDocument();
    });

    await user.type(screen.getByLabelText(/your name/i), 'Alice');
    await user.click(screen.getByRole('button', { name: /confirm booking/i }));

    await waitFor(() => {
      expect(mutateAsync).toHaveBeenCalledWith(
        expect.objectContaining({ guestName: 'Alice', eventTypeId: 'et-1' }),
      );
    });
  });

  it('shows conflict notification on 409', async () => {
    const user = userEvent.setup();

    vi.spyOn(hooks, 'usePublicEventTypes').mockReturnValue({
      data: mockEventTypes,
      isLoading: false,
      isError: false,
    } as ReturnType<typeof hooks.usePublicEventTypes>);

    vi.spyOn(hooks, 'useSlots').mockReturnValue({
      data: mockAvailableDays,
      isLoading: false,
      isError: false,
      refetch: vi.fn(),
    } as unknown as ReturnType<typeof hooks.useSlots>);

    const conflictError = Object.assign(new Error('Slot already taken'), { status: 409 });

    vi.spyOn(hooks, 'useCreateBooking').mockReturnValue({
      mutateAsync: vi.fn().mockRejectedValue(conflictError),
      isPending: false,
    } as unknown as ReturnType<typeof hooks.useCreateBooking>);

    renderBookingPage();

    await user.click(screen.getByText('2026-06-30'));
    await user.click(screen.getAllByText(/^\d{2}:\d{2}$/)[0]);

    await waitFor(() => {
      expect(screen.getByLabelText(/your name/i)).toBeInTheDocument();
    });

    await user.type(screen.getByLabelText(/your name/i), 'Bob');
    await user.click(screen.getByRole('button', { name: /confirm booking/i }));

    await waitFor(() => {
      expect(screen.getByText(/slot already taken/i)).toBeInTheDocument();
    });
  });
});
