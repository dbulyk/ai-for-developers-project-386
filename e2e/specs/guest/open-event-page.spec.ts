import { expect } from '@playwright/test';
import { test } from '../../fixtures/test';
import { createEventType } from '../../fixtures/eventType';
import { HomePage } from '../../pages/HomePage';
import { BookingPage } from '../../pages/BookingPage';

test('Guest opens event page and sees available slots', async ({ page, request, namespace, cleanup }) => {
  const eventType = await createEventType(request, namespace, {
    name: namespace.name('30-min call'),
    durationMinutes: 30,
  });
  cleanup.eventTypes.push(eventType.id);

  const home = new HomePage(page);
  await home.goto();
  await home.clickBook(eventType.name);

  const bookingPage = new BookingPage(page);
  await bookingPage.expectLoaded(eventType.name);
  await bookingPage.expectDuration('30 min');

  const firstDayButton = page.getByTestId('day-button').first();
  await expect(firstDayButton).toBeVisible();
});
