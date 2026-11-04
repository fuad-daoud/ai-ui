import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";
import { canon } from "./canon";

const cases = JSON.parse(
  readFileSync(new URL("../testdata/canon/cases.json", import.meta.url), "utf8"),
) as { name: string; in: string; out: string }[];

describe("canon (shared vectors)", () => {
  it("has vectors", () => expect(cases.length).toBeGreaterThan(0));
  for (const c of cases) {
    it(c.name, () => expect(canon(c.in)).toBe(c.out));
  }
});
