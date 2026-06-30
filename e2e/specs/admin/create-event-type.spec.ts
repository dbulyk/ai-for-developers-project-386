import { expect } from '@playwright/test';
import { test } from '../../fixtures/test';
import { listEventTypes } from '../../fixtures/eventType';
import { AdminEventTypesPage } from '../../pages/AdminEventTypesPage';
import { HomePage } from '../../pages/HomePage';
import { expectNotification } from '../../helpers/notifications';

test('Admin creates a new event type', async ({ page, request, namespace, cleanup }) => {
  const adminPage = new AdminEventTypesPage(page);
  await adminPage.goto();
  await adminPage.expectLoaded();

  await adminPage.clickNewEventType();
  const name = namespace.name('Consultation');
  await adminPage.fillForm(name, 'Detailed consultation', 45);
  await adminPage.submitCreate();

  await expectNotification(page, 'Event type created.');
  await adminPage.expectRow(name);

  // Verify it appears on the public home page
  const home = new HomePage(page);
  await home.goto();
  await home.expectEventType(name, '45 min');

  // Schedule cleanup by finding the created event type via API
  const eventTypes = await listEventTypes(request);
  const created = eventTypes.find((et) => et.name === name);
  if (created) cleanup.eventTypes.push(created.id);
});
