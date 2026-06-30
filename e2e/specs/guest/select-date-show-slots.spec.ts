import { expect } from '@playwright/test';
import { test } from '../../fixtures/test';
import { createEventType } from '../../fixtures/eventType';
import { createBooking } from '../../fixtures/booking';
import { getSlots } from '../../fixtures/slots';
import { HomePage } from '../../pages/HomePage';
import { BookingPage } from '../../pages/BookingPage';
import { firstFreeSlot, formatSlotTime } from '../../helpers/slots';

test('Guest selects a date and sees free and taken slots', async ({ page, request, namespace, cleanup }) => {
  const eventType = await createEventType(request, namespace, {
    name: namespace.name('30-min call'),
    durationMinutes: 30,
  });
  cleanup.eventTypes.push(eventType.id);

  const days = await getSlots(request, eventType.id);
  const freeSlot = firstFreeSlot(days);
  expect(freeSlot).toBeDefined();

  // Occupy the second free slot to have a taken slot in the list
  const secondFreeSlot = days
    .flatMap((d) => d.slots)
    .filter((s) => s.status === 'free')[1];
  if (secondFreeSlot) {
    const booking = await createBooking(request, {
      eventTypeId: eventType.id,
      guestName: 'Other Guest',
      startTime: secondFreeSlot.startTime,
    });
    cleanup.bookings.push(booking.id);
  }

  const home = new HomePage(page);
  await home.goto();
  await home.clickBook(eventType.name);

  const bookingPage = new BookingPage(page);
  await bookingPage.expectLoaded(eventType.name);

  const day = days.find((d) => d.slots.some((s) => s.startTime === freeSlot!.startTime))!;
  await bookingPage.selectDay(day.date);

  // Free slot is visible and clickable
  const freeCard = page.locator('[data-testid="slot-card"][data-status="free"]').first();
  await expect(freeCard).toBeVisible();

  // Taken slot is visible and disabled
  if (secondFreeSlot) {
    await bookingPage.expectSlotStatus(formatSlotTime(secondFreeSlot.startTime), 'taken');
    await bookingPage.expectSlotDisabled(formatSlotTime(secondFreeSlot.startTime));
  }
});
