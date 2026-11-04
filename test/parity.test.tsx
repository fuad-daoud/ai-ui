import { readdirSync, readFileSync } from "node:fs";
import { createElement, type FunctionComponent } from "react";
import { renderToStaticMarkup } from "react-dom/server";
import { describe, expect, it } from "vitest";
import * as ui from "../react/index";
import { canon } from "./canon";

type Case = {
  go?: unknown;
  react?: { component: string; props: Record<string, unknown> };
  error?: string;
};

const dir = new URL("../testdata/parity/", import.meta.url);
const files = readdirSync(dir).filter((f) => f.endsWith(".json")).sort();

describe("parity (React)", () => {
  it("has fixtures", () => expect(files.length).toBeGreaterThan(0));
  for (const f of files) {
    const c = JSON.parse(readFileSync(new URL(f, dir), "utf8")) as Case;
    const react = c.react;
    if (!react) continue;
    it(f, () => {
      const C = (ui as unknown as Record<string, unknown>)[react.component];
      expect(typeof C, `ai-ui/react exports ${react.component}`).toBe("function");
      const render = () =>
        renderToStaticMarkup(createElement(C as FunctionComponent<Record<string, unknown>>, react.props));
      if (c.error) {
        expect(render).toThrow(c.error);
        return;
      }
      const golden = readFileSync(new URL(f.replace(/\.json$/, ".html"), dir), "utf8");
      expect(canon(render())).toBe(canon(golden));
    });
  }
});
