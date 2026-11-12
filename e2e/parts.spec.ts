import { expect, test, type Page } from "@playwright/test";

// The v0.5 page parts, measured on the library fixture.

const rect = (page: Page, sel: string) =>
  page.locator(sel).evaluate((el) => {
    const r = el.getBoundingClientRect();
    return { left: r.left, right: r.right, top: r.top, bottom: r.bottom, width: r.width, height: r.height };
  });

const css = (page: Page, sel: string, prop: string): Promise<string> =>
  page.locator(sel).evaluate((el, prop) => getComputedStyle(el).getPropertyValue(prop), prop);

// token resolves a colour token the way the browser paints it, inside sel.
const token = (page: Page, sel: string, name: string): Promise<string> =>
  page.locator(sel).evaluate((el, name) => {
    const probe = document.createElement("span");
    probe.style.color = `var(${name})`;
    el.appendChild(probe);
    const c = getComputedStyle(probe).color;
    probe.remove();
    return c;
  }, name);

test.describe("at desktop width", () => {
  test.beforeEach(async ({ page }) => page.goto("/fixture/library/ui.css"));

  test("a page head puts its aside at the end, on the title's baseline row", async ({ page }) => {
    const head = await rect(page, "#ph");
    const main = await rect(page, "#ph-main");
    const aside = await rect(page, "#ph-aside");
    expect(Math.abs(aside.right - head.right)).toBeLessThan(1);
    expect(Math.abs(aside.bottom - main.bottom)).toBeLessThan(1);
    expect(await css(page, "#ph", "border-bottom-width")).toBe("1px");
    const size = parseFloat(await css(page, "#ph-title", "font-size"));
    expect(size).toBeGreaterThanOrEqual(28);
    expect(size).toBeLessThanOrEqual(38);
  });

  test("an eyebrow is small, uppercase and grey", async ({ page }) => {
    expect(await css(page, "#ph-eyebrow", "font-size")).toBe("12px");
    expect(await css(page, "#ph-eyebrow", "text-transform")).toBe("uppercase");
    expect(await css(page, "#ph-eyebrow", "color")).toBe(await token(page, "#library", "--ink-3"));
  });

  test("a section head pushes its aside to the end", async ({ page }) => {
    expect(Math.abs((await rect(page, "#sh-aside")).right - (await rect(page, "#sh")).right)).toBeLessThan(1);
    expect(await css(page, "#sh", "border-bottom-width")).toBe("1px");
  });

  test("stat tiles share the row; values are tabular and take their tile's colour", async ({ page }) => {
    const a = await rect(page, "#stat-a");
    const b = await rect(page, "#stat-b");
    expect(Math.abs(a.width - b.width)).toBeLessThan(1);
    expect(Math.abs(a.top - b.top)).toBeLessThan(1);
    expect(await css(page, "#stat-a-value", "font-variant-numeric")).toBe("tabular-nums");
    expect(await css(page, "#stat-a-note", "color")).toBe(await token(page, "#library", "--bad"));
    expect(await css(page, "#stat-c-value", "color")).toBe(await token(page, "#library", "--good"));
    expect(await css(page, "#stats-plain", "border-top-width")).toBe("0px");
    expect(await css(page, "#stat-plain-2", "border-left-width")).toBe("1px");
  });

  test("the current chip and a pressed chip are marked; others aren't", async ({ page }) => {
    const primary = await token(page, "#library", "--primary");
    expect(await css(page, "#chip-current", "border-top-color")).toBe(primary);
    expect(await css(page, "#chip-current", "font-weight")).toBe("600");
    expect(await css(page, "#chip-pressed", "border-top-color")).toBe(primary);
    expect(await css(page, "#chip-off", "border-top-color")).toBe(await token(page, "#library", "--rule"));
    expect((await rect(page, "#chip-off")).height).toBe(28);
    expect(await css(page, "#badge", "border-top-left-radius")).toBe("999px");
  });

  test("a segmented control fills the pressed segment and greys the unavailable one", async ({ page }) => {
    expect(await css(page, "#seg-on", "background-color")).toBe(await token(page, "#library", "--primary"));
    expect(await css(page, "#seg-on", "color")).toBe(await token(page, "#library", "--on-primary"));
    expect(await css(page, "#seg-off", "background-color")).toBe("rgba(0, 0, 0, 0)");
    expect(await css(page, "#seg-disabled", "background-color")).toBe(await token(page, "#library", "--band"));
  });

  test("an avatar is a circle; inside a square button it fills the button", async ({ page }) => {
    expect(await css(page, "#avatar", "border-top-left-radius")).toBe("50%");
    const btn = await rect(page, "#avatar-btn");
    expect(btn.width).toBe(btn.height);
    const inner = await rect(page, "#avatar-in");
    expect(inner.width).toBeGreaterThan(btn.width - 4);
  });

  test("a meter is a progressbar; one without a value slides", async ({ page }) => {
    await expect(page.getByRole("progressbar", { name: "Import" })).toHaveAttribute("value", "40");
    expect((await rect(page, "#meter")).height).toBe(8);
    expect(await css(page, "#meter-ind", "animation-name")).toBe("ui-meter-slide");
  });

  test("an empty state is framed and centred; data-frame=none is flush", async ({ page }) => {
    expect(await css(page, "#empty", "border-top-style")).toBe("dashed");
    expect(await css(page, "#empty", "text-align")).toBe("center");
    expect(await css(page, "#empty-plain", "border-top-width")).toBe("0px");
    expect(await css(page, "#empty-plain", "padding-top")).toBe("0px");
    expect(await css(page, "#empty-plain", "text-align")).toBe("start");
  });

  test("the current row keeps its tint under the pointer", async ({ page }) => {
    const mark = await token(page, "#library", "--mark-bg");
    expect(await css(page, "#row-current", "background-color")).toBe(mark);
    await page.locator("#row-current").hover();
    expect(await css(page, "#row-current", "background-color")).toBe(mark);
    expect(await css(page, "#row-plain", "background-color")).toBe("rgba(0, 0, 0, 0)");
  });
});

test("under reduced motion an indeterminate meter stands still", async ({ page }) => {
  await page.emulateMedia({ reducedMotion: "reduce" });
  await page.goto("/fixture/library/ui.css");
  expect(await css(page, "#meter-ind", "animation-name")).toBe("none");
});

test.describe("on a phone", () => {
  test.use({ viewport: { width: 390, height: 844 } });

  test("a page head wraps its aside under the title, and the title shrinks", async ({ page }) => {
    await page.goto("/fixture/library/ui.css");
    const main = await rect(page, "#ph-main");
    const aside = await rect(page, "#ph-aside");
    expect(aside.top).toBeGreaterThanOrEqual(main.bottom);
    expect(await css(page, "#ph-title", "font-size")).toBe("28px");
  });
});

test.describe("on a touch screen", () => {
  test.use({ hasTouch: true, isMobile: true, viewport: { width: 1100, height: 900 } });

  test("chips and segments grow to finger size", async ({ page }) => {
    await page.goto("/fixture/library/ui.css");
    expect((await rect(page, "#chip-off")).height).toBe(40);
    expect((await rect(page, "#seg-off")).height).toBeGreaterThanOrEqual(40);
  });
});
