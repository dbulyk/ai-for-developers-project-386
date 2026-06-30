import { expect } from '@playwright/test';
import { test } from '../../fixtures/test';
import { createEventType } from '../../fixtures/eventType';
import { createBooking } from '../../fixtures/booking';
import { getSlots } from '../../fixtures/slots';
import { AdminBookingsPage } from '../../pages/AdminBookingsPage';
import { firstFreeSlot } from '../../helpers/slots';

test('Admin views the list of bookings', async ({ page, request, namespace, cleanup }) => {
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
});
