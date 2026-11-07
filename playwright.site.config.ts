import { existsSync } from "node:fs";
import { defineConfig } from "@playwright/test";

// The site's own config, so the gallery's server and baselines stay separate.
// It serves under /ai-ui/ to match GitHub Pages and prove the relative links.
const pw = process.env.PW_CHROMIUM;
const executablePath = pw === "bundled" ? undefined : pw ?? (existsSync("/usr/bin/chromium") ? "/usr/bin/chromium" : undefined);

export default defineConfig({
  testDir: "site/e2e",
  workers: 1,
  use: {
    baseURL: "http://127.0.0.1:4318/ai-ui/",
    browserName: "chromium",
    viewport: { width: 1280, height: 900 },
    launchOptions: { executablePath },
  },
  webServer: {
    command: "go run . -out out -serve 127.0.0.1:4318",
    cwd: "site",
    url: "http://127.0.0.1:4318/ai-ui/",
    reuseExistingServer: false,
    timeout: 120_000,
  },
});
