import type { Page, Locator } from '@playwright/test';
import { expect } from '@playwright/test';
import { Selectors, exactText } from '../helpers/selectors';

export class AdminBookingsPage {
  readonly page: Page;

  constructor(page: Page) {
    this.page = page;
  }

  async goto(): Promise<void> {
    await this.page.goto('/admin/bookings');
  }

  async expectLoaded(): Promise<void> {
    await expect(this.page.getByTestId(Selectors.bookingsTitle)).toBeVisible();
    await expect(this.page.getByRole('heading', { name: 'Bookings' })).toBeVisible();
  }

  rowByGuest(guestName: string): Locator {
    return this.page
      .getByTestId(Selectors.bookingRow)
      .filter({ has: this.page.getByTestId(Selectors.bookingRowGuest).filter({ hasText: exactText(guestName) }) });
  }

  async expectBooking(guestName: string): Promise<void> {
    await expect(this.rowByGuest(guestName)).toBeVisible();
  }

  async expectBookingHidden(guestName: string): Promise<void> {
    await expect(this.rowByGuest(guestName)).toBeHidden();
  }

  async clickCancel(guestName: string): Promise<void> {
    await this.rowByGuest(guestName).getByTestId(Selectors.cancelBookingButton).click();
    await expect(this.page.getByTestId(Selectors.cancelBookingModal)).toBeVisible();
  }

  async confirmCancel(): Promise<void> {
    await this.page.getByTestId(Selectors.confirmCancelBooking).click();
    await expect(this.page.getByTestId(Selectors.cancelBookingModal)).toBeHidden();
  }
}
