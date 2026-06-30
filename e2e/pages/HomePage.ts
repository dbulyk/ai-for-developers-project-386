import type { Page, Locator } from '@playwright/test';
import { expect } from '@playwright/test';
import { Selectors } from '../helpers/selectors';

export class HomePage {
  readonly page: Page;

  constructor(page: Page) {
    this.page = page;
  }

  async goto(): Promise<void> {
    await this.page.goto('/');
  }

  async expectLoaded(): Promise<void> {
    await expect(this.page.getByTestId(Selectors.eventTypesTitle)).toBeVisible();
    await expect(this.page.getByRole('heading', { name: 'Book a Meeting' })).toBeVisible();
  }

  eventTypeCard(name: string): Locator {
    return this.page
      .getByTestId(Selectors.eventTypeCard)
      .filter({ has: this.page.getByTestId(Selectors.eventTypeName).filter({ hasText: name }) });
  }

  async expectEventType(name: string, duration: string): Promise<void> {
    const card = this.eventTypeCard(name);
    await expect(card).toBeVisible();
    await expect(card.getByTestId(Selectors.eventTypeDuration)).toHaveText(duration);
    await expect(card.getByTestId(Selectors.eventTypeBookButton)).toBeVisible();
  }

  async clickBook(name: string): Promise<void> {
    await this.eventTypeCard(name).getByTestId(Selectors.eventTypeBookButton).click();
    await this.page.waitForURL(/\/event-types\/.+/);
  }
}
