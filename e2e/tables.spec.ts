import AxeBuilder from "@axe-core/playwright";
import { expect, test, type Page } from "@playwright/test";

// The v0.3.1 table options, measured on the library fixture.

const css = (page: Page, sel: string, prop: string, pseudo: string | null = null): Promise<string> =>
  page.locator(sel).evaluate((el, [prop, pseudo]) => getComputedStyle(el, pseudo).getPropertyValue(prop), [prop, pseudo] as const);

const token = (page: Page, sel: string, name: string): Promise<string> =>
  page.locator(sel).evaluate((el, name) => {
    const probe = document.createElement("span");
    probe.style.color = `var(${name})`;
    el.appendChild(probe);
    const c = getComputedStyle(probe).color;
    probe.remove();
    return c;
  }, name);

const box = (page: Page, sel: string) =>
  page.locator(sel).evaluate((el) => {
    const r = el.getBoundingClientRect();
    return { top: r.top, left: r.left, width: r.width, height: r.height };
  });

test.describe("at desktop width", () => {
  test.beforeEach(async ({ page }) => page.goto("/fixture/library/ui.css"));

  test("an unframed table has no frame, and its outer columns sit flush", async ({ page }) => {
    expect(await css(page, "#plain-wrap", "border-top-width")).toBe("0px");
    expect(await css(page, "#plain-wrap", "border-top-left-radius")).toBe("0px");
    expect(await css(page, "#plain-wrap", "background-color")).toBe("rgba(0, 0, 0, 0)");
    expect(await css(page, "#plain-first", "padding-left")).toBe("0px");
    expect(await css(page, "#plain-last", "padding-right")).toBe("0px");
    expect(await css(page, "#plain-th", "padding-left")).toBe("0px");
  });

  test("a quiet head has no fill and small grey text, and its knobs take effect", async ({ page }) => {
    expect(await css(page, "#plain-th", "background-color")).toBe("rgba(0, 0, 0, 0)");
    expect(await css(page, "#plain-th", "color")).toBe(await token(page, "#library", "--ink-3"));
    expect(await css(page, "#plain-th", "font-size")).toBe("12px");
    expect(await css(page, "#plain-th", "text-transform")).toBe("none");
    expect(await css(page, "#caps-th", "text-transform")).toBe("uppercase");
    expect(await css(page, "#caps-th", "border-bottom-color")).toBe(await token(page, "#library", "--rule"));
    // A default head keeps its fill.
    expect(await css(page, "#wrap thead th >> nth=0", "background-color")).toBe(await token(page, "#library", "--band"));
  });

  test("a sticky head stays at the top of a bounded wrap as it scrolls", async ({ page }) => {
    const wrap = await box(page, "#sticky-wrap");
    expect(wrap.height).toBeLessThanOrEqual(128);
    await page.locator("#sticky-wrap").evaluate((el) => el.scrollTo(0, 120));
    expect(await page.locator("#sticky-wrap").evaluate((el) => el.scrollTop)).toBeGreaterThan(0);
    const th = await box(page, "#sticky-th");
    expect(Math.abs(th.top - (wrap.top + 1))).toBeLessThan(1);
    expect(await css(page, "#sticky-th", "background-color")).toBe(await token(page, "#library", "--sheet"));
  });

  test("a stacked table is still a table at desktop width", async ({ page }) => {
    expect(await css(page, "#stack", "display")).toBe("table");
    expect(await css(page, "#stack-row", "display")).toBe("table-row");
    expect(await css(page, "#stack-role", "content", "::before")).toBe("none");
  });

  test("a hovered row is tinted where the pointer can hover", async ({ page }) => {
    await page.locator("#td").hover();
    expect(await css(page, "#td", "background-color")).not.toBe("rgba(0, 0, 0, 0)");
  });
});

test.describe("on a phone", () => {
  test.use({ viewport: { width: 390, height: 844 } });

  test("rows stack into labelled lines", async ({ page }) => {
    await page.goto("/fixture/library/ui.css");
    expect(await css(page, "#stack", "display")).toBe("block");
    expect(await css(page, "#stack-row", "display")).toBe("block");
    // A labelled cell shows its label; the name cell has none and spans the row.
    expect(await css(page, "#stack-role", "content", "::before")).toBe('"Role"');
    expect(await page.locator("#stack-role").evaluate((el) => getComputedStyle(el, "::before").getPropertyValue("width"))).toBe("128px");
    expect(await css(page, "#stack-name", "content", "::before")).toBe("none");
    const row = await box(page, "#stack-row");
    expect((await box(page, "#stack-name")).width).toBe(row.width);
    expect(await css(page, "#stack-rowhead", "display")).toBe("block");
    // The header row is hidden from sight but stays in the markup.
    expect((await box(page, "#stack-head")).width).toBeLessThanOrEqual(1);
    expect(await page.locator("#stack-head th").count()).toBe(3);
    // A number sits beside its label.
    expect(await css(page, "#stack-num", "text-align")).toBe("start");
    // Nothing is wider than the screen.
    expect(await css(page, "#stack-wrap", "overflow-x")).toBe("visible");
    expect(row.width).toBeLessThanOrEqual(390);
  });

  test("axe passes on the gallery with its tables stacked", async ({ page }) => {
    await page.goto("/p/blue/light/compact");
    const { violations } = await new AxeBuilder({ page }).analyze();
    expect(violations.map((v) => `${v.id}: ${v.nodes.map((n) => n.target.join(" ")).join(", ")}`)).toEqual([]);
  });
});

test.describe("on a touch screen", () => {
  test.use({ hasTouch: true, isMobile: true, viewport: { width: 1100, height: 900 } });

  test("a tapped row keeps no tint", async ({ page }) => {
    await page.goto("/fixture/library/ui.css");
    expect(await page.evaluate(() => matchMedia("(hover: hover)").matches), "the page sees no hover").toBe(false);
    await page.locator("#td").tap();
    expect(await css(page, "#td", "background-color")).toBe("rgba(0, 0, 0, 0)");
  });
});
