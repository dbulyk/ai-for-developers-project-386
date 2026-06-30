import type { Page, Locator } from '@playwright/test';
import { expect } from '@playwright/test';
import { Selectors, exactText } from '../helpers/selectors';

export class AdminEventTypesPage {
  readonly page: Page;

  constructor(page: Page) {
    this.page = page;
  }

  async goto(): Promise<void> {
    await this.page.goto('/admin/event-types');
  }

  async expectLoaded(): Promise<void> {
    await expect(this.page.getByRole('heading', { name: 'Event Types' })).toBeVisible();
    await expect(this.page.getByTestId(Selectors.newEventTypeButton)).toBeVisible();
  }

  async clickNewEventType(): Promise<void> {
    await this.page.getByTestId(Selectors.newEventTypeButton).click();
    await expect(this.page.getByTestId(Selectors.createEventTypeModal)).toBeVisible();
  }

  async fillForm(name: string, description: string, duration: number): Promise<void> {
    await this.page.getByTestId(Selectors.eventTypeNameInput).fill(name);
    await this.page.getByTestId(Selectors.eventTypeDescriptionInput).fill(description);
    await this.page.getByTestId(Selectors.eventTypeDurationInput).fill(String(duration));
  }

  async submitCreate(): Promise<void> {
    await this.page.getByTestId(Selectors.createEventTypeSubmit).click();
    await expect(this.page.getByTestId(Selectors.createEventTypeModal)).toBeHidden();
  }

  async submitSave(): Promise<void> {
    await this.page.getByTestId(Selectors.saveEventTypeSubmit).click();
    await expect(this.page.getByTestId(Selectors.editEventTypeModal)).toBeHidden();
  }

  rowByName(name: string): Locator {
    return this.page
      .getByTestId(Selectors.eventTypeRow)
      .filter({ has: this.page.getByTestId(Selectors.eventTypeRowName).filter({ hasText: exactText(name) }) });
  }

  async expectRow(name: string): Promise<void> {
    await expect(this.rowByName(name)).toBeVisible();
  }

  async expectRowHidden(name: string): Promise<void> {
    await expect(this.rowByName(name)).toBeHidden();
  }

  async clickEdit(name: string): Promise<void> {
    await this.rowByName(name).getByTestId(Selectors.editEventTypeButton).click();
    await expect(this.page.getByTestId(Selectors.editEventTypeModal)).toBeVisible();
  }

  async clickDelete(name: string): Promise<void> {
    await this.rowByName(name).getByTestId(Selectors.deleteEventTypeButton).click();
    await expect(this.page.getByTestId(Selectors.deleteEventTypeModal)).toBeVisible();
  }

  async confirmDelete(): Promise<void> {
    await this.page.getByTestId(Selectors.confirmDeleteEventType).click();
    await expect(this.page.getByTestId(Selectors.deleteEventTypeModal)).toBeHidden();
  }
}
