// Canonical HTML for the parity tests: the TypeScript twin of internal/canon.
// testdata/canon/cases.json pins both. Change them together.
import { parse, parseFragment } from "parse5";
import type { DefaultTreeAdapterMap } from "parse5";

type ChildNode = DefaultTreeAdapterMap["childNode"];
type Element = DefaultTreeAdapterMap["element"];
type TextNode = DefaultTreeAdapterMap["textNode"];
type Template = DefaultTreeAdapterMap["template"];

const VOID = new Set(["area", "base", "br", "col", "embed", "hr", "img", "input", "link", "meta", "source", "track", "wbr"]);
const WS = /[ \t\n\r\f]+/;

const body: Element = (() => {
  const doc = parse("<!doctype html><html><head></head><body></body></html>");
  const html = doc.childNodes.find((n) => n.nodeName === "html") as Element;
  return html.childNodes.find((n) => n.nodeName === "body") as Element;
})();

const words = (s: string) => s.split(WS).filter(Boolean).join(" ");
const escText = (s: string) => s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
const escAttr = (s: string) => escText(s).replace(/"/g, "&quot;");

export function canon(src: string): string {
  return parseFragment(body, src, {}).childNodes.map(write).join("");
}

function write(n: ChildNode): string {
  if (n.nodeName === "#text") {
    const t = words((n as TextNode).value);
    return t ? escText(t) : "";
  }
  if (!("tagName" in n)) return "";
  const el = n as Element;
  const attrs = [...el.attrs]
    .sort((a, b) => (a.name < b.name ? -1 : a.name > b.name ? 1 : 0))
    .map(({ name, value }) => {
      const v = name === "class" ? words(value) : value;
      return v ? ` ${name}="${escAttr(v)}"` : ` ${name}`;
    })
    .join("");
  const kids = el.tagName === "template" ? (el as Template).content.childNodes : el.childNodes;
  const open = `<${el.tagName}${attrs}>`;
  return VOID.has(el.tagName) ? open : `${open}${kids.map(write).join("")}</${el.tagName}>`;
}
