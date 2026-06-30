import { expect } from '@playwright/test';
import { test } from '../../fixtures/test';
import { createEventType } from '../../fixtures/eventType';
import { AdminEventTypesPage } from '../../pages/AdminEventTypesPage';
import { expectNotification } from '../../helpers/notifications';

test('Admin edits an event type', async ({ page, request, namespace, cleanup }) => {
  const originalName = namespace.name('Consultation');
  const eventType = await createEventType(request, namespace, {
    name: originalName,
    description: 'Detailed consultation',
    durationMinutes: 45,
  });
  cleanup.eventTypes.push(eventType.id);

  const adminPage = new AdminEventTypesPage(page);
  await adminPage.goto();
  await adminPage.expectLoaded();
  await adminPage.expectRow(originalName);

  await adminPage.clickEdit(originalName);
  const updatedName = namespace.name('Consultation Updated');
  await adminPage.fillForm(updatedName, 'Detailed consultation', 45);
  await adminPage.submitSave();

  await expectNotification(page, 'Event type updated.');
  await adminPage.expectRow(updatedName);
  await adminPage.expectRowHidden(originalName);
});
