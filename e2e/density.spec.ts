import { expect, test, type Page } from "@playwright/test";
import { differences, snapshot } from "./support";

type Box = { height: number; fontSize: string; radius: string };

const box = (page: Page, sel: string): Promise<Box> =>
  page.locator(sel).evaluate((el) => {
    const c = getComputedStyle(el);
    return { height: el.getBoundingClientRect().height, fontSize: c.fontSize, radius: c.borderTopLeftRadius };
  });

const css = (page: Page, sel: string, prop: string): Promise<string> =>
  page.locator(sel).evaluate((el, prop) => getComputedStyle(el).getPropertyValue(prop), prop);

for (const density of ["root", "container"]) {
  test(`compact sizes with the attribute on the ${density}`, async ({ page }) => {
    await page.goto(`/fixture/library/ui.css?density=${density}`);
    expect(await box(page, "#btn")).toEqual({ height: 32, fontSize: "14px", radius: "4px" });
    expect(await box(page, "#btn-sm")).toEqual({ height: 28, fontSize: "13px", radius: "4px" });
    expect(await box(page, "#btn-lg")).toEqual({ height: 44, fontSize: "16px", radius: "4px" });
    expect(await box(page, "#input")).toEqual({ height: 32, fontSize: "14px", radius: "4px" });
    expect(await box(page, "#select")).toEqual({ height: 32, fontSize: "14px", radius: "4px" });
    expect(await box(page, "#area")).toEqual({ height: 80, fontSize: "14px", radius: "4px" });
    for (const sel of ["#label", "#hint", "#error"]) expect(await css(page, sel, "font-size")).toBe("13px");
    expect(await css(page, "#field", "row-gap")).toBe("4px");
    expect(await css(page, "#pill", "border-top-left-radius")).toBe("4px");
    expect(await css(page, "#callout", "border-top-left-radius")).toBe("4px");
    expect(await css(page, "#wrap", "border-top-left-radius")).toBe("4px");
    expect(await css(page, "#card", "border-top-left-radius")).toBe("8px");
    expect(await css(page, "#td", "font-size")).toBe("14px");
    // A dialog stays in its place in the tree, so it inherits compact in the top layer too.
    await page.locator("#dialog").evaluate((d: HTMLDialogElement) => d.showModal());
    expect(await css(page, "#dialog", "border-top-left-radius")).toBe("8px");
    // font-size on :root would shrink rem and every rem-sized control with it.
    expect(await css(page, "html", "font-size")).toBe("16px");
    expect(await css(page, "body", "font-size")).toBe(density === "root" ? "14px" : "16px");
    expect(await css(page, "#library", "font-size")).toBe("14px");
  });
}

test("default sizes", async ({ page }) => {
  await page.goto("/fixture/library/ui.css");
  expect(await box(page, "#btn")).toEqual({ height: 40, fontSize: "16px", radius: "6px" });
  expect(await box(page, "#btn-sm")).toEqual({ height: 32, fontSize: "14px", radius: "6px" });
  expect(await box(page, "#btn-lg")).toEqual({ height: 48, fontSize: "18px", radius: "6px" });
  expect(await css(page, "#card", "border-top-left-radius")).toBe("10px");
  expect(await css(page, "#pill", "border-top-left-radius")).toBe("999px");
});

for (const query of ["", "?density=root"]) {
  test(`fields match buttons (${query || "default"})`, async ({ page }) => {
    await page.goto(`/fixture/library/ui.css${query}`);
    const button = (await box(page, "#btn")).height;
    expect((await box(page, "#input")).height).toBe(button);
    expect((await box(page, "#select")).height).toBe(button);
    // A multi-select lists rows, so it keeps its padding.
    expect(await css(page, "#multi", "padding-top")).toBe("8px");
    expect(await css(page, "#area", "padding-top")).toBe("8px");
  });
}

test("compact has no effect without the attribute", async ({ page }) => {
  await page.goto("/p/blue/light");
  const before = await snapshot(page, "body");
  const removed = await page.evaluate(() => {
    let n = 0;
    const walk = (list: CSSRuleList, remove: (i: number) => void) => {
      for (let i = list.length - 1; i >= 0; i--) {
        const rule = list[i];
        if (rule instanceof CSSStyleRule && rule.selectorText.includes("data-density")) {
          remove(i);
          n++;
        } else if ("cssRules" in rule) {
          const group = rule as CSSGroupingRule;
          walk(group.cssRules, (j) => group.deleteRule(j));
        }
      }
    };
    for (const sheet of Array.from(document.styleSheets)) walk(sheet.cssRules, (i) => sheet.deleteRule(i));
    return n;
  });
  expect(removed, "the page had compact rules to remove").toBeGreaterThan(0);
  expect(differences(before, await snapshot(page, "body"))).toEqual([]);
});
