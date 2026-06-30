import { expect } from '@playwright/test';
import { test } from '../../fixtures/test';
import { createEventType } from '../../fixtures/eventType';
import { HomePage } from '../../pages/HomePage';

test('Guest views the list of event types', async ({ page, request, namespace, cleanup }) => {
  const eventType = await createEventType(request, namespace, {
    name: namespace.name('30-min call'),
    durationMinutes: 30,
  });
  cleanup.eventTypes.push(eventType.id);

  const home = new HomePage(page);
  await home.goto();
  await home.expectLoaded();
  await home.expectEventType(eventType.name, '30 min');
});
