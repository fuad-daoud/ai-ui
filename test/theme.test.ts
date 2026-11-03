import { readFileSync } from "node:fs";
import { JSDOM } from "jsdom";
import { describe, expect, it } from "vitest";

const src = readFileSync(new URL("../js/theme.js", import.meta.url), "utf8");

type Choice = { theme?: string | null; palette?: string | null };
type ThemeAPI = { get(): { theme: string; palette: string }; set(patch: Choice): unknown };

function boot(opts: { htmlAttrs?: string; stored?: string; body?: string; blockStorage?: boolean } = {}) {
  const dom = new JSDOM(
    `<!doctype html><html${opts.htmlAttrs ?? ""}><head></head><body>${opts.body ?? ""}</body></html>`,
    { url: "https://ai-ui.test/", runScripts: "outside-only" },
  );
  const w = dom.window;
  if (opts.stored !== undefined) w.localStorage.setItem("ai-ui:theme", opts.stored);
  if (opts.blockStorage) {
    w.Storage.prototype.getItem = () => {
      throw new w.DOMException("blocked", "SecurityError");
    };
    w.Storage.prototype.setItem = () => {
      throw new w.DOMException("blocked", "SecurityError");
    };
  }
  w.eval(src);
  const api = (w as unknown as { aiui: { theme: ThemeAPI } }).aiui.theme;
  const stored = () => (opts.blockStorage ? null : w.localStorage.getItem("ai-ui:theme"));
  return { dom, html: w.document.documentElement, api, stored };
}

describe("theme.js", () => {
  it("leaves the server's attributes alone until the visitor chooses", () => {
    const { html, api } = boot({ htmlAttrs: ' data-theme="dark" data-palette="whiteboard"' });
    expect(html.getAttribute("data-theme")).toBe("dark");
    expect(html.getAttribute("data-palette")).toBe("whiteboard");
    expect(api.get()).toEqual({ theme: "dark", palette: "whiteboard" });
  });

  it("applies a stored choice as soon as it runs", () => {
    const { html } = boot({ stored: '{"theme":"dark","palette":"mono"}' });
    expect(html.getAttribute("data-theme")).toBe("dark");
    expect(html.getAttribute("data-palette")).toBe("mono");
  });

  it("ignores corrupt storage", () => {
    const { html, api } = boot({ stored: "{oops", htmlAttrs: ' data-palette="mono"' });
    expect(html.hasAttribute("data-theme")).toBe(false);
    expect(api.get()).toEqual({ theme: "system", palette: "mono" });
  });

  it("ignores stored values outside the allowed sets", () => {
    const { html } = boot({ stored: '{"theme":"purple","palette":"mono\\"] {x"}' });
    expect(html.hasAttribute("data-theme")).toBe(false);
    expect(html.hasAttribute("data-palette")).toBe(false);
  });

  it("still applies choices when storage throws", () => {
    const { html, api } = boot({ blockStorage: true });
    expect(() => api.set({ theme: "dark" })).not.toThrow();
    expect(html.getAttribute("data-theme")).toBe("dark");
  });

  it("lets system clear a server-forced theme, and null restore it", () => {
    const { html, api, stored } = boot({ htmlAttrs: ' data-theme="dark"' });
    api.set({ theme: "system" });
    expect(html.hasAttribute("data-theme")).toBe(false);
    expect(JSON.parse(stored() ?? "null")).toEqual({ theme: "system" });
    api.set({ theme: null });
    expect(html.getAttribute("data-theme")).toBe("dark");
    expect(JSON.parse(stored() ?? "null")).toEqual({});
  });

  it("applies data-ui-set-* clicks, including clicks on a child of the button", () => {
    const { dom, html, stored } = boot({
      body: '<button data-ui-set-palette="mono"><span id="inner">Mono</span></button><button id="reset" data-ui-set-palette="">Default</button>',
    });
    dom.window.document.getElementById("inner")!.click();
    expect(html.getAttribute("data-palette")).toBe("mono");
    expect(JSON.parse(stored() ?? "null")).toEqual({ palette: "mono" });
    dom.window.document.getElementById("reset")!.click();
    expect(html.hasAttribute("data-palette")).toBe(false);
  });
});
