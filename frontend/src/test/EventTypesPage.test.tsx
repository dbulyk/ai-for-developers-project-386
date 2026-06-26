import { screen } from '@testing-library/react';
import { vi } from 'vitest';
import { EventTypesPage } from '../pages/EventTypesPage';
import { renderWithProviders } from './utils';
import * as hooks from '../api/hooks';

vi.mock('../api/hooks', async (importOriginal) => {
  const actual = await importOriginal<typeof hooks>();
  return { ...actual };
});

describe('EventTypesPage', () => {
  it('renders list of event types from the API', () => {
    vi.spyOn(hooks, 'usePublicEventTypes').mockReturnValue({
      data: [
        { id: 'et-1', name: '30-min call', description: 'Quick call', durationMinutes: 30 },
        { id: 'et-2', name: '1-hour session', description: 'Deep dive', durationMinutes: 60 },
      ],
      isLoading: false,
      isError: false,
    } as ReturnType<typeof hooks.usePublicEventTypes>);

    renderWithProviders(<EventTypesPage />);

    expect(screen.getByText('30-min call')).toBeInTheDocument();
    expect(screen.getByText('1-hour session')).toBeInTheDocument();
    expect(screen.getAllByRole('link', { name: /book/i })).toHaveLength(2);
  });

  it('shows loader while loading', () => {
    vi.spyOn(hooks, 'usePublicEventTypes').mockReturnValue({
      data: undefined,
      isLoading: true,
      isError: false,
    } as ReturnType<typeof hooks.usePublicEventTypes>);

    renderWithProviders(<EventTypesPage />);
    expect(document.querySelector('.mantine-Loader-root')).toBeInTheDocument();
  });
});
