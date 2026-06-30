import { expect } from '@playwright/test';
import { test } from '../../fixtures/test';
import { createEventType } from '../../fixtures/eventType';
import { createBooking } from '../../fixtures/booking';
import { getSlots } from '../../fixtures/slots';
import { AdminBookingsPage } from '../../pages/AdminBookingsPage';
import { BookingPage } from '../../pages/BookingPage';
import { firstFreeSlot, formatSlotTime } from '../../helpers/slots';
import { expectNotification } from '../../helpers/notifications';

test('Admin cancels a booking', async ({ page, request, namespace, cleanup }) => {
  const eventType = await createEventType(request, namespace, {
    name: namespace.name('30-min call'),
    durationMinutes: 30,
  });
  cleanup.eventTypes.push(eventType.id);

  const days = await getSlots(request, eventType.id);
  const slot = firstFreeSlot(days);
  expect(slot).toBeDefined();

  const booking = await createBooking(request, {
    eventTypeId: eventType.id,
    guestName: 'Иван',
    startTime: slot!.startTime,
  });
  cleanup.bookings.push(booking.id);

  const adminPage = new AdminBookingsPage(page);
  await adminPage.goto();
  await adminPage.expectLoaded();
  await adminPage.expectBooking('Иван');

  await adminPage.clickCancel('Иван');
  await adminPage.confirmCancel();

  await expectNotification(page, 'Booking has been cancelled.');
  await adminPage.expectBookingHidden('Иван');

  // Verify the corresponding slot is available again on the public page
  const bookingPage = new BookingPage(page);
  await bookingPage.goto(eventType.id);
  await bookingPage.expectLoaded(eventType.name);

  const day = days.find((d) => d.slots.some((s) => s.startTime === slot!.startTime))!;
  await bookingPage.selectDay(day.date);
  await bookingPage.expectSlotStatus(formatSlotTime(slot!.startTime), 'free');
});
