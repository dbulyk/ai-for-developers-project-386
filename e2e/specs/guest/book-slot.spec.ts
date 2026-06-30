import { expect } from '@playwright/test';
import { test } from '../../fixtures/test';
import { createEventType } from '../../fixtures/eventType';
import { findBookingId } from '../../fixtures/booking';
import { getSlots } from '../../fixtures/slots';
import { HomePage } from '../../pages/HomePage';
import { BookingPage } from '../../pages/BookingPage';
import { firstFreeSlot, formatSlotTime } from '../../helpers/slots';
import { expectNotification } from '../../helpers/notifications';

test('Guest successfully books a free slot', async ({ page, request, namespace, cleanup }) => {
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
  await bookingPage.confirmBooking();

  await expectNotification(page, 'Booking confirmed!');
  await bookingPage.expectModalHidden();

  // Find and schedule cleanup for the UI-created booking
  const bookingId = await findBookingId(request, {
    eventTypeId: eventType.id,
    guestName: 'Иван',
    startTime: slot!.startTime,
  });
  if (bookingId) cleanup.bookings.push(bookingId);

  // Verify the slot became taken after refetch
  await page.reload();
  await bookingPage.selectDay(day.date);
  await bookingPage.expectSlotStatus(formatSlotTime(slot!.startTime), 'taken');
});
