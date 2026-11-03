import { execFileSync } from "node:child_process";
import { fileURLToPath } from "node:url";
import { describe, expect, it } from "vitest";

const root = fileURLToPath(new URL("..", import.meta.url));
const out = execFileSync("npx", ["@tailwindcss/cli", "-i", "test/tw-fixture/in.css"], {
  cwd: root,
  encoding: "utf8",
});

describe("react/theme.css with Tailwind v4", () => {
  it("puts the ui layer between base and components", () => {
    const ours = out.indexOf("@layer theme, base, ui, components, utilities;");
    expect(ours).toBeGreaterThanOrEqual(0);
    const theirs = out.indexOf("@layer theme, base, components, utilities;");
    if (theirs >= 0) expect(ours).toBeLessThan(theirs);
  });

  it("maps colour tokens to utilities", () => {
    expect(out).toMatch(/\.bg-sheet\s*\{\s*background-color:\s*var\(--sheet\)/);
    expect(out).toMatch(/\.text-ink-2\s*\{\s*color:\s*var\(--ink-2\)/);
    expect(out).toMatch(/\.border-rule\s*\{\s*border-color:\s*var\(--rule\)/);
  });

  it("keeps the library inside its own layers", () => {
    expect(out).toMatch(/@layer ui\.components\s*\{[\s\S]*\.ui-btn/);
  });
});
