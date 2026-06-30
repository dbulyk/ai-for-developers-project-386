import type { Page, Locator } from '@playwright/test';
import { expect } from '@playwright/test';
import { Selectors } from '../helpers/selectors';
import { formatSlotTime } from '../helpers/slots';

export class BookingPage {
  readonly page: Page;

  constructor(page: Page) {
    this.page = page;
  }

  async goto(eventTypeId: string): Promise<void> {
    await this.page.goto(`/event-types/${eventTypeId}`);
  }

  async expectLoaded(eventTypeName: string): Promise<void> {
    await expect(this.page.getByTestId(Selectors.bookingPageTitle)).toHaveText(eventTypeName);
  }

  async expectDuration(duration: string): Promise<void> {
    await expect(this.page.getByTestId(Selectors.bookingPageDuration)).toHaveText(duration);
  }

  dayButton(date: string): Locator {
    return this.page.getByTestId(Selectors.dayButton).filter({ hasText: date });
  }

  async selectDay(date: string): Promise<void> {
    await this.dayButton(date).click();
  }

  slotCardByTime(time: string): Locator {
    return this.page
      .getByTestId(Selectors.slotCard)
      .filter({ hasText: time });
  }

  async clickSlot(time: string): Promise<void> {
    await this.slotCardByTime(time).click();
  }

  async clickSlotByUtc(utcIso: string): Promise<void> {
    await this.clickSlot(formatSlotTime(utcIso));
  }

  async expectSlotStatus(time: string, status: 'free' | 'taken'): Promise<void> {
    const card = this.slotCardByTime(time);
    await expect(card).toHaveAttribute(Selectors.slotStatusAttr, status);
  }

  async expectSlotDisabled(time: string): Promise<void> {
    const card = this.slotCardByTime(time);
    await expect(card).toHaveCSS('opacity', '0.4');
    await expect(card).toHaveCSS('cursor', /default|auto/);
  }

  async openBookingModal(time: string): Promise<void> {
    await this.clickSlot(time);
    await expect(this.modal()).toBeVisible();
  }

  modal(): Locator {
    return this.page.getByTestId(Selectors.bookingModal);
  }

  async expectModalVisible(): Promise<void> {
    await expect(this.modal()).toBeVisible();
  }

  async expectModalHidden(): Promise<void> {
    await expect(this.modal()).toBeHidden();
  }

  async fillGuestName(name: string): Promise<void> {
    await this.page.getByTestId(Selectors.guestNameInput).fill(name);
  }

  async confirmBooking(): Promise<void> {
    await this.page.getByTestId(Selectors.confirmBookingButton).click();
  }

  async cancelModal(): Promise<void> {
    await this.page.getByTestId(Selectors.cancelBookingModalButton).click();
  }
}
