import { expect, test, type Page } from "@playwright/test";
import { differences, snapshot, type Snapshot } from "./support";

// These two are set on html and body only, so they can't match on a page ui-scoped.css leaves alone.
const PAGE_LEVEL = ["text-size-adjust", "-webkit-font-smoothing"];

async function library(page: Page, sheet: string, query = ""): Promise<Snapshot> {
  await page.goto(`/fixture/library/${sheet}${query}`);
  await page.locator("#dialog").evaluate((d: HTMLDialogElement) => d.showModal());
  // Freeze what moves (an indeterminate meter), so the two snapshots see the same frame.
  await page.evaluate(() => document.getAnimations().forEach((a) => { a.pause(); a.currentTime = 0; }));
  return snapshot(page, "#library");
}

for (const query of ["", "?density=container"]) {
  test(`library markup matches under both sheets (${query || "default"})`, async ({ page }) => {
    const full = await library(page, "ui.css", query);
    const scoped = await library(page, "ui-scoped.css", query);
    expect(full.length).toBeGreaterThan(100);
    expect(differences(full, scoped, PAGE_LEVEL)).toEqual([]);
  });
}

test("a focused control draws the same ring under both sheets", async ({ page }) => {
  const ring = async (sheet: string) => {
    await page.goto(`/fixture/library/${sheet}`);
    await page.keyboard.press("Tab");
    await page.keyboard.press("Tab");
    return page.evaluate(() => {
      const el = document.activeElement!;
      const c = getComputedStyle(el);
      return { id: el.id, visible: el.matches(":focus-visible"), style: c.outlineStyle, width: c.outlineWidth, offset: c.outlineOffset, color: c.outlineColor };
    });
  };
  const full = await ring("ui.css");
  expect(full).toMatchObject({ id: "sec", visible: true, style: "solid", width: "2px" });
  expect(await ring("ui-scoped.css")).toEqual(full);
});

for (const colorScheme of ["light", "dark"] as const) {
  test.describe(`a document beside the library, ${colorScheme}`, () => {
    test.use({ colorScheme });

    const load = async (page: Page, sheet: string) => {
      await page.goto(`/fixture/document/${sheet}`);
      return { style: await snapshot(page, "#document"), shot: await page.screenshot({ fullPage: true }) };
    };

    test("ui-scoped.css leaves it untouched", async ({ page }) => {
      const none = await load(page, "none");
      const scoped = await load(page, "ui-scoped.css");
      expect(differences(none.style, scoped.style)).toEqual([]);
      expect(Buffer.compare(none.shot, scoped.shot)).toBe(0);
    });

    // Proves the fixture can see a leak: if ui.css stopped changing it, the test above would prove nothing.
    test("ui.css restyles it", async ({ page }) => {
      const none = await load(page, "none");
      const full = await load(page, "ui.css");
      expect(differences(none.style, full.style).length).toBeGreaterThan(50);
      expect(Buffer.compare(none.shot, full.shot)).not.toBe(0);
    });
  });
}
