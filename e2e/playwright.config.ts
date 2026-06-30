import { defineConfig, devices } from '@playwright/test';

const baseURL = process.env.BASE_URL ?? 'http://localhost:5173';
const apiURL = process.env.API_URL ?? 'http://localhost:8080';
const reuseExistingServer = !process.env.CI;

export default defineConfig({
  testDir: './specs',
  fullyParallel: false,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: 1,
  reporter: process.env.CI ? [['html'], ['github']] : 'list',
  use: {
    baseURL,
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
  },

  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],

  webServer: [
    {
      command: `cd ../backend && CORS_ALLOWED_ORIGINS=${baseURL} go run ./cmd/server`,
      url: `${apiURL}/public/event-types`,
      timeout: 120 * 1000,
      reuseExistingServer,
      stdout: 'pipe',
      stderr: 'pipe',
    },
    {
      command: `cd ../frontend && VITE_API_BASE_URL=${apiURL} npm run dev`,
      url: baseURL,
      timeout: 120 * 1000,
      reuseExistingServer,
      stdout: 'pipe',
      stderr: 'pipe',
    },
  ],
});
