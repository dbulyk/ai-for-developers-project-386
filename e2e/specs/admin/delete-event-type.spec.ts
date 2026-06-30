import { expect } from '@playwright/test';
import { test } from '../../fixtures/test';
import { createEventType } from '../../fixtures/eventType';
import { AdminEventTypesPage } from '../../pages/AdminEventTypesPage';
import { HomePage } from '../../pages/HomePage';
import { expectNotification } from '../../helpers/notifications';

test('Admin deletes an event type', async ({ page, request, namespace, cleanup }) => {
  const name = namespace.name('Consultation Updated');
  const eventType = await createEventType(request, namespace, {
    name,
    description: 'Detailed consultation',
    durationMinutes: 45,
  });
  cleanup.eventTypes.push(eventType.id);

  const adminPage = new AdminEventTypesPage(page);
  await adminPage.goto();
  await adminPage.expectLoaded();
  await adminPage.expectRow(name);

  await adminPage.clickDelete(name);
  await adminPage.confirmDelete();

  await expectNotification(page, 'Event type deleted.');
  await adminPage.expectRowHidden(name);

  // Verify it disappeared from the public home page
  const home = new HomePage(page);
  await home.goto();
  await home.expectLoaded();
  await expect(home.eventTypeCard(name)).toBeHidden();
});
