import type { Page } from '@playwright/test';

export async function expectNotification(page: Page, title: string): Promise<void> {
  const notification = page.locator('[role="alert"]', { hasText: title }).first();
  await notification.waitFor({ state: 'visible', timeout: 5000 });
}

export async function expectNoNotification(page: Page, title: string): Promise<boolean> {
  const notification = page.locator('[role="alert"]', { hasText: title }).first();
  return (await notification.count()) === 0;
}
