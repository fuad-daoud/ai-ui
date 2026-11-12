import { readdirSync, readFileSync } from "node:fs";
import { createElement, type FunctionComponent } from "react";
import { renderToStaticMarkup } from "react-dom/server";
import { describe, expect, it } from "vitest";
import * as ui from "../react/index";
import { canon } from "./canon";

type Node = { component?: string; tag?: string; props?: Record<string, unknown> };
type Case = {
  go?: unknown;
  react?: { component: string; props: Record<string, unknown> };
  error?: string;
};

// build turns a fixture node into a React element. A child, or a prop such as
// icon, is a string, a library component ({ component, props }) or a plain
// element ({ tag, props }).
function build(node: Node, key?: number): ReturnType<typeof createElement> {
  const props: Record<string, unknown> = { ...(node.props ?? {}) };
  if (key !== undefined) props.key = key;
  for (const [k, v] of Object.entries(props)) if (k === "children" || isNode(v)) props[k] = kids(v);
  if (node.tag) return createElement(node.tag, props);
  const C = (ui as unknown as Record<string, unknown>)[node.component ?? ""];
  if (typeof C !== "function") throw new Error(`ai-ui/react doesn't export ${node.component}`);
  return createElement(C as FunctionComponent<Record<string, unknown>>, props);
}

function isNode(v: unknown): boolean {
  return !!v && typeof v === "object" && !Array.isArray(v) && ("component" in v || "tag" in v);
}

function kids(v: unknown): unknown {
  if (Array.isArray(v)) return v.map((c, i) => (c && typeof c === "object" ? build(c as Node, i) : c));
  if (v && typeof v === "object") return build(v as Node);
  return v;
}

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
      const render = () => renderToStaticMarkup(build(react));
      if (c.error) {
        expect(render).toThrow(c.error);
        return;
      }
      const golden = readFileSync(new URL(f.replace(/\.json$/, ".html"), dir), "utf8");
      expect(canon(render())).toBe(canon(golden));
    });
  }
});
