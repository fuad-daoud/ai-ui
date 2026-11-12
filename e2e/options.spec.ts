import { expect, test, type Page } from "@playwright/test";

// The v0.3 options, measured on the library fixture at both densities.

const rect = (page: Page, sel: string) =>
  page.locator(sel).evaluate((el) => {
    const r = el.getBoundingClientRect();
    return { left: r.left, right: r.right, top: r.top, bottom: r.bottom, width: r.width, height: r.height };
  });

const css = (page: Page, sel: string, prop: string, pseudo: string | null = null): Promise<string> =>
  page.locator(sel).evaluate((el, [prop, pseudo]) => getComputedStyle(el, pseudo).getPropertyValue(prop), [prop, pseudo] as const);

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

const showModal = (page: Page, sel: string) => page.locator(sel).evaluate((d: HTMLDialogElement) => d.showModal());

for (const [name, query] of [["default", ""], ["compact", "?density=root"]] as const) {
  test.describe(name, () => {
    test.beforeEach(async ({ page }) => page.goto(`/fixture/library/ui.css${query}`));

    test("small and large fields are as tall as the button of their size", async ({ page }) => {
      const sm = (await rect(page, "#btn-sm")).height;
      const lg = (await rect(page, "#btn-lg")).height;
      expect(sm).toBe(name === "compact" ? 28 : 32);
      expect(lg).toBe(name === "compact" ? 44 : 48);
      expect((await rect(page, "#input-sm")).height).toBe(sm);
      expect((await rect(page, "#select-sm")).height).toBe(sm);
      expect((await rect(page, "#input-lg")).height).toBe(lg);
      expect((await rect(page, "#select-lg")).height).toBe(lg);
      expect((await rect(page, "#area-sm")).height).toBe(sm * 2);
      // A select keeps the room its chevron needs.
      expect(await css(page, "#select-sm", "padding-right")).toBe("36px");
    });

    test("a square button is as wide as it is tall", async ({ page }) => {
      for (const [sel, ref] of [["#btn-square", "#btn"], ["#btn-square-sm", "#btn-sm"]]) {
        const r = await rect(page, sel);
        expect(r.width).toBe(r.height);
        expect(r.height).toBe((await rect(page, ref)).height);
      }
    });

    test("a full-width button fills its box and wraps its label", async ({ page }) => {
      const box = await rect(page, "#full-box");
      const b = await rect(page, "#btn-full");
      expect(b.width).toBe(box.width);
      expect(b.height).toBeGreaterThan((await rect(page, "#btn")).height);
    });

    test("an outline button keeps the variant's colour as text on the sheet", async ({ page }) => {
      expect(await css(page, "#btn-outline", "color")).toBe(await token(page, "#library", "--bad"));
      expect(await css(page, "#btn-outline", "background-color")).toBe(await token(page, "#library", "--sheet"));
      expect(await css(page, "#btn-outline", "border-top-color")).toBe(await token(page, "#library", "--rule"));
      expect(await css(page, "#btn-outline-primary", "color")).toBe(await token(page, "#library", "--primary"));
    });

    test("disabled changes colour, not opacity", async ({ page }) => {
      const band = await token(page, "#library", "--band");
      const ink3 = await token(page, "#library", "--ink-3");
      for (const sel of ["#btn-off", "#btn-outline-off", "#link-off"]) {
        expect(await css(page, sel, "opacity"), sel).toBe("1");
        expect(await css(page, sel, "background-color"), sel).toBe(band);
        expect(await css(page, sel, "color"), sel).toBe(ink3);
      }
      expect(await css(page, "#ghost-off", "background-color")).toBe("rgba(0, 0, 0, 0)");
      expect(await css(page, "#ghost-off", "color")).toBe(ink3);
      expect(await css(page, "#btn-off", "cursor")).toBe("not-allowed");
      expect(await css(page, "#btn-off", "pointer-events")).toBe("auto");
      expect(await css(page, "#link-off", "pointer-events")).toBe("none");
    });

    test("the input icon sits inside the input, before its text", async ({ page }) => {
      const input = await rect(page, "#input-icon");
      const icon = await rect(page, "#icon");
      expect(icon.width).toBe(16);
      expect(icon.left).toBeGreaterThan(input.left);
      expect(icon.top).toBeGreaterThan(input.top);
      expect(icon.bottom).toBeLessThan(input.bottom);
      const start = parseFloat(await css(page, "#input-icon", "padding-left"));
      expect(input.left + start).toBeGreaterThan(icon.right);
    });

    test("width knobs set a field's width and a row child's basis", async ({ page }) => {
      expect((await rect(page, "#input-width")).width).toBe(128);
      expect((await rect(page, "#row-narrow")).width).toBe(80);
      expect((await rect(page, "#row-wide")).width).toBeGreaterThan(200);
    });

    test("cards and dialogs are flat only in compact", async ({ page }) => {
      const flat = name === "compact";
      const card = await css(page, "#card", "box-shadow");
      expect(card === "none").toBe(flat);
      await showModal(page, "#dialog");
      const dialog = await css(page, "#dialog", "box-shadow");
      expect(dialog === "none").toBe(flat);
      expect(await css(page, "#dialog", "border-top-width")).toBe("1px");
    });

    test("the eyebrow sits above the title, and a dialog lays out without its ×", async ({ page }) => {
      await showModal(page, "#dialog-eyebrow");
      const eyebrow = await rect(page, "#eyebrow");
      const title = await rect(page, "#dialog-eyebrow-title");
      expect(eyebrow.bottom).toBeLessThanOrEqual(title.top);
      // The eyebrow's text starts where the title's does.
      const start = parseFloat(await css(page, "#eyebrow", "padding-left"));
      expect(Math.abs(eyebrow.left + start - title.left)).toBeLessThan(1);
      expect(await css(page, "#eyebrow", "text-transform")).toBe("uppercase");
      expect(title.width).toBeGreaterThan(0);
    });

    test("a dotted pill draws its dot in its text colour; an outline pill has no fill", async ({ page }) => {
      expect(await css(page, "#pill-dot", "width", "::before")).toBe("6px");
      expect(await css(page, "#pill-dot", "background-color", "::before")).toBe(await css(page, "#pill-dot", "color"));
      expect(await css(page, "#pill-outline", "background-color")).toBe("rgba(0, 0, 0, 0)");
      expect(await css(page, "#pill-outline", "border-top-color")).toBe(await token(page, "#library", "--rule-field"));
      expect(await css(page, "#pill-outline-warn", "border-top-color")).toBe(await token(page, "#library", "--warning"));
    });

    test("a start-marker disclosure draws its chevron first and turns it when open", async ({ page }) => {
      expect(await css(page, "#disc-start-summary", "content", "::after")).toBe("none");
      expect(await css(page, "#disc-start-summary", "width", "::before")).toBe("8px");
      const closed = await css(page, "#disc-start-summary", "transform", "::before");
      await page.locator("#disc-start-summary").click();
      await expect.poll(() => css(page, "#disc-start-summary", "transform", "::before")).not.toBe(closed);
      expect(await css(page, "#disc-norule", "border-bottom-color")).toBe("rgba(0, 0, 0, 0)");
    });
  });
}

test.describe("on a touch screen", () => {
  test.use({ hasTouch: true, isMobile: true, viewport: { width: 1100, height: 900 } });

  test("compact controls grow to finger sizes; default ones don't change", async ({ page }) => {
    await page.goto("/fixture/library/ui.css?density=root");
    expect(await page.evaluate(() => matchMedia("(pointer: coarse)").matches), "the page sees a coarse pointer").toBe(true);
    expect((await rect(page, "#btn")).height).toBe(40);
    expect((await rect(page, "#btn-sm")).height).toBe(34);
    expect((await rect(page, "#btn-lg")).height).toBe(44);
    expect((await rect(page, "#input")).height).toBe(40);
    expect((await rect(page, "#input-sm")).height).toBe(34);
    expect((await rect(page, "#select-lg")).height).toBe(44);
    await page.goto("/fixture/library/ui.css");
    expect((await rect(page, "#btn")).height).toBe(40);
    expect((await rect(page, "#btn-sm")).height).toBe(32);
    expect((await rect(page, "#input-sm")).height).toBe(32);
  });
});

test.describe("on a phone", () => {
  test.use({ viewport: { width: 360, height: 740 } });

  test("a long title wraps inside its dialog and keeps the × in view", async ({ page }) => {
    await page.goto("/fixture/library/ui.css");
    await showModal(page, "#dialog-long");
    const dialog = await rect(page, "#dialog-long");
    const title = await rect(page, "#dialog-long-title");
    const x = await rect(page, "#dialog-long-x");
    expect(dialog.right).toBeLessThanOrEqual(360);
    // Without the fix the title keeps the email's width (about 500px) and the dialog scrolls sideways.
    expect(title.left).toBeGreaterThanOrEqual(dialog.left);
    expect(title.right).toBeLessThanOrEqual(x.left);
    expect(x.right).toBeLessThanOrEqual(dialog.right);
    expect(await page.locator("#dialog-long").evaluate((el) => el.scrollWidth <= el.clientWidth)).toBe(true);
  });
});
