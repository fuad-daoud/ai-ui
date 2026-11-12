import AxeBuilder from "@axe-core/playwright";
import { expect, test, type Page } from "@playwright/test";

const home = "./";
const releases = "releases/";
const controls = "button, input, select, summary, a";
const components = "components/";
const guides = ["guides/", "guides/getting-started/", "guides/theming/", "guides/density/", "guides/existing-apps/"];
const docs = [components, ...guides];

// Wait until nothing is animating, so measurements see the settled pile.
const settled = (page: Page) =>
  page.waitForFunction(() => document.getAnimations().every((a) => a.playState !== "running"));

test("every asset loads, and the links work under /ai-ui/", async ({ page }) => {
  const bad: string[] = [];
  page.on("response", (r) => { if (r.status() >= 400) bad.push(`${r.status()} ${r.url()}`); });
  page.on("requestfailed", (r) => bad.push(`failed ${r.url()}`));
  await page.goto(home);
  await expect(page.locator("h1")).toContainText("the whole UI kit");
  const font = await page.evaluate(async () => { await document.fonts.ready; return document.fonts.check('800 16px Outfit'); });
  expect(font).toBe(true);
  await page.getByRole("link", { name: "Release notes" }).click();
  await expect(page).toHaveURL(/\/ai-ui\/releases\/$/);
  await expect(page.locator("main h1")).toHaveText("Changelog");
  await page.getByRole("link", { name: "ai-ui", exact: true }).click();
  await expect(page.locator("#pile")).toBeVisible();
  expect(bad).toEqual([]);
});

for (const [name, path] of [["home", home], ["release notes", releases]]) {
  test(`axe passes on ${name}`, async ({ page }) => {
    await page.goto(path);
    await settled(page);
    const { violations } = await new AxeBuilder({ page }).analyze();
    expect(violations.map((v) => `${v.id}: ${v.nodes.map((n) => n.target).join(", ")}`)).toEqual([]);
  });

  test(`${name} doesn't scroll sideways at 360px`, async ({ page }) => {
    await page.setViewportSize({ width: 360, height: 800 });
    await page.goto(path);
    await settled(page);
    const { scroll, client } = await page.evaluate(() => ({ scroll: document.documentElement.scrollWidth, client: document.documentElement.clientWidth }));
    expect(scroll).toBeLessThanOrEqual(client);
  });
}

for (const path of [home, releases, ...docs]) {
  test(`no duplicate IDs on ${path}`, async ({ page }) => {
    await page.goto(path);
    const ids = await page.evaluate(() => [...document.querySelectorAll("[id]")].map((e) => e.id));
    expect(ids.filter((id, i) => ids.indexOf(id) !== i)).toEqual([]);
  });
}

for (const path of docs) {
  test(`axe passes on ${path}`, async ({ page }) => {
    await page.goto(path);
    const { violations } = await new AxeBuilder({ page }).analyze();
    expect(violations.map((v) => `${v.id}: ${v.nodes.map((n) => n.target).join(", ")}`)).toEqual([]);
  });

  for (const width of [320, 360]) {
    test(`${path} doesn't scroll sideways at ${width}px`, async ({ page }) => {
      await page.setViewportSize({ width, height: 800 });
      await page.goto(path);
      const { scroll, client } = await page.evaluate(() => ({ scroll: document.documentElement.scrollWidth, client: document.documentElement.clientWidth }));
      expect(scroll).toBeLessThanOrEqual(client);
    });
  }
}

// body clips sideways overflow, so a page-width check alone can't see a table
// pushing past the screen: measure the frame, and scroll it from the keyboard.
test("a table wider than a phone scrolls inside its frame", async ({ page }) => {
  await page.setViewportSize({ width: 320, height: 800 });
  await page.goto("guides/density/");
  const frame = page.locator(".prose .ui-table-wrap").first();
  const r = await frame.evaluate((wrap) => {
    wrap.querySelector<HTMLElement>("td")!.style.minWidth = "60rem";
    return { scrolls: wrap.scrollWidth > wrap.clientWidth, right: wrap.getBoundingClientRect().right, view: document.documentElement.clientWidth };
  });
  expect(r.scrolls).toBe(true);
  expect(r.right).toBeLessThanOrEqual(r.view);
  await frame.focus();
  await page.keyboard.press("ArrowRight");
  await expect.poll(() => frame.evaluate((wrap) => wrap.scrollLeft)).toBeGreaterThan(0);
});

test("the nav reaches Components and Guides on the site", async ({ page }) => {
  await page.goto(home);
  const nav = page.getByRole("navigation", { name: "Main" });
  await nav.getByRole("link", { name: "Components" }).click();
  await expect(page).toHaveURL(/\/ai-ui\/components\/$/);
  await expect(page.locator("main h1")).toHaveText("Components");
  await nav.getByRole("link", { name: "Guides" }).click();
  await expect(page).toHaveURL(/\/ai-ui\/guides\/$/);
  await page.getByRole("main").getByRole("link", { name: "Theming" }).click();
  await expect(page).toHaveURL(/\/ai-ui\/guides\/theming\/$/);
  await expect(page.locator("main h1")).toHaveText("Theming");
  await expect(page.getByRole("navigation", { name: "Guides" }).getByRole("link", { name: "Theming" })).toHaveAttribute("aria-current", "page");
});

for (const scripts of [true, false]) {
  test(`the dialog demo opens and closes${scripts ? "" : " without scripts"}`, async ({ browser, baseURL }) => {
    const context = await browser.newContext({ javaScriptEnabled: scripts, baseURL });
    const page = await context.newPage();
    await page.goto(components);
    const opener = page.getByRole("button", { name: "Delete invoice" });
    await opener.click();
    const dialog = page.getByRole("dialog", { name: "Delete invoice?" });
    await expect(dialog).toBeVisible();
    await dialog.getByRole("button", { name: "Cancel" }).click();
    await expect(dialog).toBeHidden();
    await expect(opener).toBeFocused();
    await context.close();
  });
}

test("Tab reaches the side list and every demo control on the Components page", async ({ page }) => {
  await page.goto(components);
  // A radio group is one tab stop, at its checked radio; disabled controls and a closed dialog's or menu's aren't stops.
  const total = await page.evaluate(() => {
    const els = [...document.querySelectorAll<HTMLElement>(".side a, .demo :is(a[href], button, input, select, textarea, summary)")]
      .filter((el) => !(el as HTMLInputElement).disabled && !el.closest("dialog, [popover]")
        && !(el instanceof HTMLInputElement && el.type === "radio" && !el.checked));
    els.forEach((el, i) => { el.dataset.want = String(i); });
    return els.length;
  });
  const seen = new Set<string>();
  for (let i = 0; i < total * 2 + 20; i++) {
    await page.keyboard.press("Tab");
    const want = await page.evaluate(() => (document.activeElement as HTMLElement | null)?.dataset.want ?? "");
    if (want) seen.add(want);
  }
  expect(seen.size).toBe(total);
});

test("the fall ends within five seconds", async ({ page }) => {
  await page.goto(home);
  const end = await page.evaluate(() =>
    Math.max(...document.getAnimations().map((a) => Number(a.effect?.getComputedTiming().endTime ?? 0))),
  );
  expect(end).toBeGreaterThan(0);
  expect(end).toBeLessThan(5000);
});

test("Drop it again replays the fall and keeps what was typed", async ({ page }) => {
  await page.goto(home);
  await settled(page);
  await page.locator("#pile-email").fill("grace@example.com");
  await page.locator("#drop").click();
  const running = await page.evaluate(() => document.getAnimations().filter((a) => a.playState === "running").length);
  expect(running).toBeGreaterThan(0);
  await settled(page);
  await expect(page.locator("#pile-email")).toHaveValue("grace@example.com");
});

// The heap is placed by hand, so this is what keeps it usable.
for (const width of [320, 360, 768, 1151, 1152, 1280, 1920]) {
  test(`at ${width}px no control in the pile is covered or cut off`, async ({ page }) => {
    await page.setViewportSize({ width, height: 900 });
    await page.goto(home);
    await settled(page);
    const bad = await page.evaluate((selector) => {
      const pile = document.getElementById("pile")!;
      const problems: string[] = [];
      for (const el of pile.querySelectorAll<HTMLElement>(selector)) {
        el.scrollIntoView({ block: "center" });
        const r = el.getBoundingClientRect();
        const p = pile.getBoundingClientRect();
        const name = (el.textContent || el.id).trim();
        if (r.left < p.left - 1 || r.right > p.right + 1) problems.push(`${name}: past the pile's edge`);
        const hit = document.elementFromPoint(r.left + r.width / 2, r.top + r.height / 2);
        if (hit !== el && !el.contains(hit)) problems.push(`${name}: covered by ${hit?.textContent?.trim().slice(0, 30)}`);
      }
      return problems;
    }, controls);
    expect(bad).toEqual([]);
  });
}

// Controls aren't the only thing the heap can bury: a screenshot at 1152px showed
// the Email card covering the table's amounts, with every control still reachable.
for (const width of [320, 360, 768, 1151, 1152, 1280, 1920]) {
  test(`at ${width}px no item is cut off and no text is covered`, async ({ page }) => {
    await page.setViewportSize({ width, height: 900 });
    await page.goto(home);
    await settled(page);
    const bad = await page.evaluate(() => {
      const pile = document.getElementById("pile")!;
      const problems: string[] = [];
      for (const item of pile.querySelectorAll<HTMLElement>(".it")) {
        if (item.offsetLeft < 0 || item.offsetLeft + item.offsetWidth > pile.clientWidth) {
          problems.push(`${item.className}: past the pile's edge`);
        }
      }
      const texts = ".ui-num, th, .ui-callout-title, .ui-callout p, .ui-dialog-title, .ui-dialog-body, .ui-field-label, .ui-field-error, .ui-check-label, .ui-pill";
      for (const el of pile.querySelectorAll<HTMLElement>(texts)) {
        if (el.closest(".ui-vh")) continue;
        el.scrollIntoView({ block: "center" });
        const r = el.getBoundingClientRect();
        const hit = document.elementFromPoint(r.left + r.width / 2, r.top + r.height / 2);
        // Its own item is fine. Anything else means it's covered, or clipped below the floor.
        const own = hit === el || el.contains(hit) || (hit?.contains(el) && hit.closest(".it") === el.closest(".it"));
        if (!own) {
          problems.push(`${el.textContent!.trim()}: covered by ${hit?.textContent?.trim().slice(0, 30)}`);
        }
      }
      return problems;
    });
    expect(bad).toEqual([]);
  });
}

test("Tab reaches every control in the pile", async ({ page }) => {
  await page.goto(home);
  await settled(page);
  const total = await page.locator(`#pile :is(${controls})`).count();
  const seen = new Set<string>();
  for (let i = 0; i < 40; i++) {
    await page.keyboard.press("Tab");
    const id = await page.evaluate(() => {
      const el = document.activeElement as HTMLElement;
      return el.closest("#pile") ? `${el.tagName}:${el.id}:${el.textContent}` : "";
    });
    if (id) seen.add(id);
  }
  expect(seen.size).toBe(total);
});

test("under reduced motion the components are in place and Drop it again is hidden", async ({ page }) => {
  await page.emulateMedia({ reducedMotion: "reduce" });
  await page.goto(home);
  // Not "none running": a slow load could read that before the fall starts.
  const animations = await page.evaluate(() => document.getAnimations().length);
  expect(animations).toBe(0);
  await expect(page.locator("#drop")).toBeHidden();
  await expect(page.locator("#pile-email")).toBeInViewport({ ratio: 0 });
  const name = await page.locator(".i-save").evaluate((el) => getComputedStyle(el).animationName);
  expect(name).toBe("none");
});

test("Copy puts the exact prompt on the clipboard", async ({ page, context }) => {
  await context.grantPermissions(["clipboard-read", "clipboard-write"]);
  await page.goto(home);
  const prompt = (await page.locator("#prompt").textContent())!.trim();
  expect(prompt).toContain("AGENTS.md");
  await page.locator("[data-copy]").first().click();
  await expect(page.locator("[data-copy]").first()).toHaveText("Copied");
  await expect(page.locator("#copy-status")).toHaveText("Copied");
  expect(await page.evaluate(() => navigator.clipboard.readText())).toBe(prompt);
});

test("a refused clipboard says so, even on a phone, and the prompt stays visible", async ({ page }) => {
  await page.setViewportSize({ width: 360, height: 800 });
  await page.addInitScript(() => {
    navigator.clipboard.writeText = () => Promise.reject(new Error("refused"));
  });
  await page.goto(home);
  await page.locator("[data-copy]").first().click();
  const button = page.locator("[data-copy]").first();
  await expect(button).toHaveText(/see below/);
  const box = (await button.boundingBox())!;
  expect(box.x).toBeGreaterThanOrEqual(0);
  expect(box.x + box.width).toBeLessThanOrEqual(360);
  await expect(page.locator("#copy-status")).toHaveText(/see below/);
  await expect(page.locator("#prompt")).toBeVisible();
});

test("without scripts the page reads, and its script-only buttons are hidden", async ({ browser, baseURL }) => {
  const context = await browser.newContext({ javaScriptEnabled: false, baseURL });
  const page = await context.newPage();
  await page.goto(home);
  await expect(page.locator("h1")).toBeVisible();
  await expect(page.locator("#prompt")).toBeVisible();
  await expect(page.locator("#pile-email")).toBeVisible();
  await expect(page.locator("#drop")).toBeHidden();
  await expect(page.locator("[data-copy]").first()).toBeHidden();
  await context.close();
});
