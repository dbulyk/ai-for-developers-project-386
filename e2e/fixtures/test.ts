import { test as base } from '@playwright/test';
import type { APIRequestContext } from '@playwright/test';
import { createNamespace, type Namespace } from './namespace';
import { apiClient } from './api';

export interface CleanupTracker {
  eventTypes: string[];
  bookings: string[];
}

class CleanupTrackerImpl implements CleanupTracker {
  eventTypes: string[] = [];
  bookings: string[] = [];
  private request: APIRequestContext;

  constructor(request: APIRequestContext) {
    this.request = request;
  }

  async run(): Promise<void> {
    for (const id of this.bookings) {
      try {
        await apiClient(this.request).delete(`/admin/bookings/${id}`);
      } catch {
        // ignore cleanup errors
      }
    }
    for (const id of this.eventTypes) {
      try {
        await apiClient(this.request).delete(`/admin/event-types/${id}`);
      } catch {
        // ignore cleanup errors
      }
    }
  }
}

export interface TestFixtures {
  namespace: Namespace;
  cleanup: CleanupTracker;
}

export const test = base.extend<TestFixtures>({
  namespace: async ({}, use, testInfo) => {
    const slug = testInfo.title.replace(/[^a-zA-Z0-9_-]+/g, '-').replace(/^-|-$/g, '').slice(0, 40);
    await use(createNamespace(slug || 'test'));
  },
  cleanup: async ({ request }, use) => {
    const tracker = new CleanupTrackerImpl(request);
    await use(tracker);
    await tracker.run();
  },
});
