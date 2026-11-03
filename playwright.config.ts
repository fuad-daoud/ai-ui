import { existsSync } from "node:fs";
import { defineConfig } from "@playwright/test";

// contabo has /usr/bin/chromium. Set PW_CHROMIUM=bundled to use Playwright's
// own build instead (after `npx playwright install chromium`).
const pw = process.env.PW_CHROMIUM;
const executablePath = pw === "bundled" ? undefined : pw ?? (existsSync("/usr/bin/chromium") ? "/usr/bin/chromium" : undefined);

export default defineConfig({
  testDir: "e2e",
  snapshotPathTemplate: "{testDir}/__screenshots__/{arg}{ext}",
  workers: 1,
  use: {
    baseURL: "http://127.0.0.1:4317",
    browserName: "chromium",
    viewport: { width: 1100, height: 900 },
    launchOptions: { executablePath },
  },
  expect: { toHaveScreenshot: { maxDiffPixelRatio: 0.001 } },
  webServer: {
    command: "go run ./gallery -addr 127.0.0.1:4317",
    url: "http://127.0.0.1:4317/p/blue/light",
    reuseExistingServer: false,
    timeout: 120_000,
  },
});
