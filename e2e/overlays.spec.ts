import AxeBuilder from "@axe-core/playwright";
import { expect, test, type Page } from "@playwright/test";

// The v0.4 menu, sheet and theme picker, measured on the library fixture.

const rect = (page: Page, sel: string) =>
  page.locator(sel).evaluate((el) => {
    const r = el.getBoundingClientRect();
    return { left: r.left, right: r.right, top: r.top, bottom: r.bottom, width: r.width, height: r.height };
  });

const css = (page: Page, sel: string, prop: string, pseudo: string | null = null): Promise<string> =>
  page.locator(sel).evaluate((el, [prop, pseudo]) => getComputedStyle(el, pseudo).getPropertyValue(prop), [prop, pseudo] as const);

// settle waits for a sheet's slide-in to finish before it is measured.
const settle = (page: Page, sel: string) =>
  page.locator(sel).evaluate((el) => Promise.all(el.getAnimations().map((a) => a.finished)).then(() => undefined));

// open centres an opener in the viewport, so its menu has room below, then clicks it.
const open = async (page: Page, sel: string) => {
  await page.locator(sel).evaluate((el) => el.scrollIntoView({ block: "center" }));
  await page.click(sel);
};

const focused = (page: Page) => page.evaluate(() => document.activeElement?.id ?? "");
const isOpen = (page: Page, sel: string) => page.locator(sel).evaluate((el) => el.matches(":popover-open, [open]"));

test.describe("menu", () => {
  test.beforeEach(async ({ page }) => page.goto("/fixture/library/ui.css"));

  test("opens under its opener, end edges lined up, and focuses its first item", async ({ page }) => {
    await open(page, "#edge-btn");
    const edge = await rect(page, "#edge-btn");
    const edgeMenu = await rect(page, "#edge-menu");
    expect(Math.abs(edgeMenu.top - (edge.bottom + 4))).toBeLessThan(1);
    expect(Math.abs(edgeMenu.right - edge.right)).toBeLessThan(1);
    expect(edgeMenu.width).toBe(256);
    await page.keyboard.press("Escape");
    // Near the start edge there's no room to end-align it: it flips to start-aligned.
    await open(page, "#menu-btn");
    expect(await isOpen(page, "#menu")).toBe(true);
    const btn = await rect(page, "#menu-btn");
    const menu = await rect(page, "#menu");
    expect(Math.abs(menu.top - (btn.bottom + 4))).toBeLessThan(1);
    expect(Math.abs(menu.left - btn.left)).toBeLessThan(1);
    await expect.poll(() => focused(page)).toBe("mi-account");
    // The identity row is one line with an ellipsis.
    expect(await css(page, "#menu-label", "text-overflow")).toBe("ellipsis");
    expect((await rect(page, "#menu-label")).height).toBeLessThan(40);
    expect(await page.locator("#menu-label").evaluate((el) => el.scrollWidth > el.clientWidth)).toBe(true);
  });

  test("arrows, Home and End move through the items and wrap", async ({ page }) => {
    await open(page, "#menu-btn");
    await expect.poll(() => focused(page)).toBe("mi-account");
    const steps: [string, string][] = [
      ["ArrowDown", "mi-rename"],
      ["End", "mi-delete"],
      ["ArrowDown", "mi-account"],
      ["ArrowUp", "mi-delete"],
      ["Home", "mi-account"],
    ];
    for (const [key, want] of steps) {
      await page.keyboard.press(key);
      expect(await focused(page), key).toBe(want);
    }
  });

  test("a menuitem closes it; Escape closes it and returns focus to the opener", async ({ page }) => {
    await open(page, "#menu-btn");
    await page.click("#mi-rename");
    await expect.poll(() => isOpen(page, "#menu")).toBe(false);
    await open(page, "#menu-btn");
    await expect.poll(() => focused(page)).toBe("mi-account");
    await page.keyboard.press("Escape");
    expect(await isOpen(page, "#menu")).toBe(false);
    expect(await focused(page)).toBe("menu-btn");
  });

  test("Tab closes it and moves on from the opener", async ({ page }) => {
    await open(page, "#menu-btn");
    await expect.poll(() => focused(page)).toBe("mi-account");
    await page.keyboard.press("Tab");
    expect(await isOpen(page, "#menu")).toBe(false);
    expect(await focused(page)).toBe("theme-btn");
  });

  test("a radio item is checked by theme.js, focused on open, and keeps the menu open", async ({ page }) => {
    // Nothing is chosen yet: the mode follows the system and the palette is blue.
    expect(await page.locator("#tm-system").getAttribute("aria-checked")).toBe("true");
    expect(await page.locator("#tm-blue").getAttribute("aria-checked")).toBe("false");
    expect(await css(page, "#tm-system", "visibility", "::after")).toBe("visible");
    expect(await css(page, "#tm-blue", "visibility", "::after")).toBe("hidden");
    await open(page, "#theme-btn");
    await expect.poll(() => focused(page)).toBe("tm-system");
    await page.click("#tm-mono-dark");
    expect(await page.evaluate(() => [document.documentElement.dataset.theme, document.documentElement.dataset.palette])).toEqual(["dark", "mono"]);
    expect(await page.locator("#tm-mono-dark").getAttribute("aria-checked")).toBe("true");
    expect(await isOpen(page, "#theme-menu")).toBe(true);
  });

  test("data-align=start lines up the start edges; the sheet head is hidden at desktop width", async ({ page }) => {
    await open(page, "#theme-btn");
    const btn = await rect(page, "#theme-btn");
    const menu = await rect(page, "#theme-menu");
    expect(Math.abs(menu.left - btn.left)).toBeLessThan(1);
    expect(await css(page, "#theme-menu-head", "display")).toBe("none");
  });

  test("flips above its opener when there's no room below", async ({ page }) => {
    await page.locator("#low-btn").evaluate((el) => el.scrollIntoView({ block: "end" }));
    await page.click("#low-btn");
    const btn = await rect(page, "#low-btn");
    const menu = await rect(page, "#low-menu");
    expect(Math.abs(menu.bottom - (btn.top - 4))).toBeLessThan(1);
  });

  test("still opens and closes without ui.js", async ({ page }) => {
    await page.route("**/_ui/ui.js*", (r) => r.fulfill({ status: 200, contentType: "text/javascript", body: "" }));
    await page.goto("/fixture/library/ui.css");
    await open(page, "#edge-btn");
    expect(await isOpen(page, "#edge-menu")).toBe(true);
    const btn = await rect(page, "#edge-btn");
    expect(Math.abs((await rect(page, "#edge-menu")).right - btn.right)).toBeLessThan(1);
    await page.mouse.click(5, 5);
    expect(await isOpen(page, "#edge-menu")).toBe(false);
  });

  test("axe passes with each menu open", async ({ page }) => {
    for (const [btn, menu] of [["#menu-btn", "#menu"], ["#theme-btn", "#theme-menu"]]) {
      await open(page, btn);
      const { violations } = await new AxeBuilder({ page }).include(menu).analyze();
      expect(violations.map((v) => `${v.id}: ${v.nodes.map((n) => n.target.join(" ")).join(", ")}`), menu).toEqual([]);
      await page.keyboard.press("Escape");
    }
  });
});

test.describe("menu on a phone", () => {
  test.use({ viewport: { width: 390, height: 844 } });

  test("data-sheet=phone makes a bottom sheet with its head and a scrim; × closes it", async ({ page }) => {
    await page.goto("/fixture/library/ui.css");
    await open(page, "#theme-btn");
    const menu = await rect(page, "#theme-menu");
    expect(menu.left).toBe(0);
    expect(menu.width).toBe(390);
    expect(Math.abs(menu.bottom - 844)).toBeLessThan(1);
    expect(await css(page, "#theme-menu-head", "display")).toBe("flex");
    expect(await css(page, "#theme-menu", "background-color", "::backdrop")).not.toBe("rgba(0, 0, 0, 0)");
    expect((await rect(page, "#tm-blue")).height).toBe(44);
    await page.click("#theme-menu-x");
    expect(await isOpen(page, "#theme-menu")).toBe(false);
  });

  test("a menu without data-sheet stays anchored, inside the viewport", async ({ page }) => {
    await page.goto("/fixture/library/ui.css");
    await open(page, "#menu-btn");
    const menu = await rect(page, "#menu");
    expect(menu.left).toBeGreaterThanOrEqual(0);
    expect(menu.right).toBeLessThanOrEqual(390);
    expect(await css(page, "#menu", "background-color", "::backdrop")).toBe("rgba(0, 0, 0, 0)");
  });
});

test.describe("sheet", () => {
  test("docks to the bottom at full width, and a backdrop click closes it", async ({ page }) => {
    await page.goto("/fixture/library/ui.css");
    await page.click("#sheet-btn");
    await settle(page, "#sheet");
    const view = page.viewportSize()!;
    const sheet = await rect(page, "#sheet");
    expect(sheet.left).toBe(0);
    expect(sheet.width).toBe(view.width);
    expect(Math.abs(sheet.bottom - view.height)).toBeLessThan(1);
    expect(sheet.height).toBeLessThanOrEqual(view.height * 0.7);
    // The page behind is inert: a click there lands on the backdrop and closes the sheet.
    await page.mouse.click(view.width / 2, 10);
    await expect.poll(() => isOpen(page, "#sheet")).toBe(false);
  });

  test("closes on Escape and returns focus; the × closes it", async ({ page }) => {
    await page.goto("/fixture/library/ui.css");
    await page.click("#sheet-btn");
    await page.keyboard.press("Escape");
    expect(await isOpen(page, "#sheet")).toBe(false);
    expect(await focused(page)).toBe("sheet-btn");
    await page.click("#sheet-btn");
    await page.click("#sheet-x");
    expect(await isOpen(page, "#sheet")).toBe(false);
  });

  test("data-side=end docks a drawer to the end, full height", async ({ page }) => {
    await page.goto("/fixture/library/ui.css");
    await page.click("#drawer-btn");
    await settle(page, "#drawer");
    const view = page.viewportSize()!;
    const drawer = await rect(page, "#drawer");
    expect(Math.abs(drawer.right - view.width)).toBeLessThan(1);
    expect(drawer.top).toBe(0);
    expect(drawer.height).toBe(view.height);
    expect(drawer.width).toBe(448);
  });

  test("slides in, unless reduced motion is asked for", async ({ page }) => {
    await page.goto("/fixture/library/ui.css");
    await page.click("#sheet-btn");
    expect(await css(page, "#sheet", "transition-duration")).toBe("0.15s");
    await page.emulateMedia({ reducedMotion: "reduce" });
    await page.goto("/fixture/library/ui.css");
    await page.click("#sheet-btn");
    expect(await css(page, "#sheet", "transition-duration")).toBe("0s");
  });

  test("compact sheets are flat; menus keep their shadow", async ({ page }) => {
    await page.goto("/fixture/library/ui.css?density=root");
    await page.click("#sheet-btn");
    expect(await css(page, "#sheet", "box-shadow")).toBe("none");
    await page.keyboard.press("Escape");
    await open(page, "#menu-btn");
    expect(await css(page, "#menu", "box-shadow")).not.toBe("none");
    expect((await rect(page, "#mi-account")).height).toBe(34);
  });

  test("axe passes with a sheet open", async ({ page }) => {
    await page.goto("/fixture/library/ui.css");
    await page.click("#sheet-btn");
    const { violations } = await new AxeBuilder({ page }).include("#sheet").analyze();
    expect(violations.map((v) => v.id)).toEqual([]);
  });
});

test.describe("theme picker", () => {
  test.beforeEach(async ({ page }) => page.goto("/fixture/library/ui.css"));

  test("a face paints its own palette and mode inside another", async ({ page }) => {
    // Blue dark's paper is oklch(18% 0.02 262), inside a mono light region.
    expect(await css(page, "#face-blue-dark", "background-image")).toContain("oklch(0.18 0.02 262)");
    expect(await css(page, "#mono-region", "color-scheme")).toBe("light");
  });

  test("the system face shows light over dark", async ({ page }) => {
    expect(await css(page, "#face-system", "background-image", "::before")).toContain("oklch(0.982 0.01 262)");
    expect(await css(page, "#face-system", "background-image", "::after")).toContain("oklch(0.18 0.02 262)");
  });

  test("a click sets the theme, marks the swatch pressed and says so", async ({ page }) => {
    await page.evaluate(() => {
      (window as unknown as { seen: unknown[] }).seen = [];
      document.addEventListener("aiui:themechange", (e) => (window as unknown as { seen: unknown[] }).seen.push((e as CustomEvent).detail));
    });
    expect(await page.locator("#sw-system").getAttribute("aria-pressed")).toBe("true");
    await page.click("#sw-blue-dark");
    expect(await page.evaluate(() => document.documentElement.dataset.theme)).toBe("dark");
    expect(await page.locator("#sw-blue-dark").getAttribute("aria-pressed")).toBe("true");
    expect(await page.locator("#sw-system").getAttribute("aria-pressed")).toBe("false");
    expect(await css(page, "#sw-blue-dark", "border-top-width")).toBe("2px");
    expect(await page.evaluate(() => (window as unknown as { seen: unknown[] }).seen)).toEqual([{ theme: "dark", palette: "blue" }]);
    // The choice is stored and applied before paint on the next page.
    await page.reload();
    expect(await page.evaluate(() => document.documentElement.dataset.theme)).toBe("dark");
  });

  test("Left, Right, Home and End move between swatches", async ({ page }) => {
    await page.focus("#sw-system");
    const steps: [string, string][] = [
      ["ArrowRight", "sw-blue"],
      ["End", "sw-rose"],
      ["ArrowRight", "sw-system"],
      ["ArrowLeft", "sw-rose"],
      ["Home", "sw-system"],
    ];
    for (const [key, want] of steps) {
      await page.keyboard.press(key);
      expect(await focused(page), key).toBe(want);
    }
  });

  test("a light-only palette stays light when dark is asked for", async ({ page }) => {
    expect(await css(page, "#rose-dark", "color-scheme")).toBe("light");
    expect(await css(page, "#rose-dark", "color")).toBe(await css(page, "#rose-light", "color"));
    expect(await css(page, "#rose-dark", "background-color")).toBe(await css(page, "#rose-light", "background-color"));
  });
});
