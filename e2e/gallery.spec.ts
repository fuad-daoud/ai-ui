import AxeBuilder from "@axe-core/playwright";
import { expect, test } from "@playwright/test";

const paths = ["blue", "mono"].flatMap((p) =>
  ["light", "dark"].flatMap((m) => [`/p/${p}/${m}`, `/p/${p}/${m}/dialog`]),
);

for (const path of paths) {
  test(`gallery ${path}`, async ({ page }) => {
    await page.goto(path);
    if (path.endsWith("/dialog")) {
      // The reset's `* { margin: 0 }` must not undo the UA's centring of a modal.
      const box = await page.locator("#confirm").boundingBox();
      const view = page.viewportSize();
      expect(box && view, "dialog and viewport measured").toBeTruthy();
      expect(Math.abs(box!.x + box!.width / 2 - view!.width / 2)).toBeLessThan(2);
      expect(Math.abs(box!.y + box!.height / 2 - view!.height / 2)).toBeLessThan(2);
    }
    await expect(page).toHaveScreenshot(`${path.slice(3).replaceAll("/", "-")}.png`, {
      fullPage: true,
      animations: "disabled",
    });
    const { violations } = await new AxeBuilder({ page }).analyze();
    expect(violations.map((v) => `${v.id}: ${v.nodes.map((n) => n.target.join(" ")).join(", ")}`)).toEqual([]);
  });
}
