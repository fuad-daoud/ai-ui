import type { Page } from "@playwright/test";

export type Snapshot = { key: string; style: Record<string, string> }[];

// Every computed property of every element under root, and of its ::before and ::after.
export function snapshot(page: Page, root: string): Promise<Snapshot> {
  return page.evaluate((root) => {
    const out: { key: string; style: Record<string, string> }[] = [];
    const els = document.querySelectorAll(`${root}, ${root} *`);
    els.forEach((el, i) => {
      for (const pseudo of [null, "::before", "::after"]) {
        const c = getComputedStyle(el, pseudo);
        const style: Record<string, string> = {};
        for (let n = 0; n < c.length; n++) {
          const prop = c.item(n);
          if (!prop.startsWith("--")) style[prop] = c.getPropertyValue(prop);
        }
        out.push({ key: `${i}:${el.tagName.toLowerCase()}#${el.id}.${el.getAttribute("class") ?? ""}${pseudo ?? ""}`, style });
      }
    });
    return out;
  }, root);
}

export function differences(a: Snapshot, b: Snapshot, ignore: string[] = []): string[] {
  if (a.length !== b.length) return [`different element counts: ${a.length} and ${b.length}`];
  const out: string[] = [];
  a.forEach((row, i) => {
    for (const prop of Object.keys(row.style)) {
      if (ignore.includes(prop)) continue;
      if (row.style[prop] !== b[i].style[prop]) out.push(`${row.key} ${prop}: ${row.style[prop]} -> ${b[i].style[prop]}`);
    }
  });
  return out;
}
