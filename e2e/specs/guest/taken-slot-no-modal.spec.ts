import { expect } from '@playwright/test';
import { test } from '../../fixtures/test';
import { createEventType } from '../../fixtures/eventType';
import { createBooking } from '../../fixtures/booking';
import { getSlots } from '../../fixtures/slots';
import { HomePage } from '../../pages/HomePage';
import { BookingPage } from '../../pages/BookingPage';
import { firstTakenSlot, formatSlotTime } from '../../helpers/slots';

test('Clicking a taken slot does not open booking modal', async ({ page, request, namespace, cleanup }) => {
  const eventType = await createEventType(request, namespace, {
    name: namespace.name('30-min call'),
    durationMinutes: 30,
  });
  cleanup.eventTypes.push(eventType.id);

  const daysBefore = await getSlots(request, eventType.id);
  const freeSlot = daysBefore.flatMap((d) => d.slots).find((s) => s.status === 'free');
  expect(freeSlot).toBeDefined();

  const booking = await createBooking(request, {
    eventTypeId: eventType.id,
    guestName: 'Другой',
    startTime: freeSlot!.startTime,
  });
  cleanup.bookings.push(booking.id);

  const daysAfter = await getSlots(request, eventType.id);
  const takenSlot = firstTakenSlot(daysAfter);
  expect(takenSlot).toBeDefined();

  const home = new HomePage(page);
  await home.goto();
  await home.clickBook(eventType.name);

  const bookingPage = new BookingPage(page);
  await bookingPage.expectLoaded(eventType.name);

  const day = daysAfter.find((d) => d.slots.some((s) => s.startTime === takenSlot!.startTime))!;
  await bookingPage.selectDay(day.date);

  await bookingPage.clickSlot(formatSlotTime(takenSlot!.startTime));
  await bookingPage.expectModalHidden();

  // Guest name input should not be present outside the closed modal
  await expect(page.getByTestId('guest-name-input')).toBeHidden();
});
