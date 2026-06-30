import { expect } from '@playwright/test';
import { test } from '../../fixtures/test';
import { createEventType } from '../../fixtures/eventType';
import { createBooking } from '../../fixtures/booking';
import { getSlots } from '../../fixtures/slots';
import { HomePage } from '../../pages/HomePage';
import { BookingPage } from '../../pages/BookingPage';
import { firstFreeSlot, formatSlotTime } from '../../helpers/slots';
import { expectNotification } from '../../helpers/notifications';

test('Guest gets 409 when booking a slot that became taken', async ({ page, request, namespace, cleanup }) => {
  const eventType = await createEventType(request, namespace, {
    name: namespace.name('30-min call'),
    durationMinutes: 30,
  });
  cleanup.eventTypes.push(eventType.id);

  const days = await getSlots(request, eventType.id);
  const slot = firstFreeSlot(days);
  expect(slot).toBeDefined();

  const home = new HomePage(page);
  await home.goto();
  await home.clickBook(eventType.name);

  const bookingPage = new BookingPage(page);
  await bookingPage.expectLoaded(eventType.name);

  const day = days.find((d) => d.slots.some((s) => s.startTime === slot!.startTime))!;
  await bookingPage.selectDay(day.date);
  await bookingPage.openBookingModal(formatSlotTime(slot!.startTime));
  await bookingPage.fillGuestName('Иван');

  // Another "user" books the same slot via API while the modal is open
  const booking = await createBooking(request, {
    eventTypeId: eventType.id,
    guestName: 'Другой',
    startTime: slot!.startTime,
  });
  cleanup.bookings.push(booking.id);

  await bookingPage.confirmBooking();

  await expectNotification(page, 'Slot already taken');
  await bookingPage.expectModalHidden();

  // Reselect the day to pick up the refetched slot statuses
  // (the selectedDay object in React state is stale after refetch)
  await bookingPage.selectDay(day.date);

  // Slot list should now show the slot as taken
  await bookingPage.expectSlotStatus(formatSlotTime(slot!.startTime), 'taken');
});
