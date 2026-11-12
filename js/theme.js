/* ai-ui theme.js: load it blocking in <head>, before the stylesheet:
     <script src="{{href "theme.js"}}"></script>
   It applies the visitor's stored theme and palette before first paint and
   persists new choices. Until the visitor chooses, the server's own
   data-theme / data-palette on <html> stay untouched.

   API:    window.aiui.theme.get()   -> { theme: "light"|"dark"|"system", palette: string }
           window.aiui.theme.set({ theme, palette })   null clears a choice
           window.aiui.theme.sync(root?)   re-marks the choice buttons under root
   Markup: <button data-ui-set-theme="dark">   <button data-ui-set-palette="mono">
           an empty value clears that choice. Each such element is marked
           aria-pressed (aria-checked for a radio role) when it matches.
   Event:  "aiui:themechange" on document, detail = get(), after every set()
           and when the OS flips light/dark while the mode follows it. */
(function () {
  "use strict";
  var KEY = "ai-ui:theme";
  var THEMES = ["light", "dark", "system"];
  var PALETTE = /^[a-z0-9-]{1,32}$/;
  var root = document.documentElement;
  var server = { theme: root.getAttribute("data-theme"), palette: root.getAttribute("data-palette") };
  var choice = clean(read());

  function read() {
    try {
      var v = JSON.parse(window.localStorage.getItem(KEY) || "{}");
      return v && typeof v === "object" ? v : {};
    } catch (e) {
      return {};
    }
  }

  function write() {
    try {
      window.localStorage.setItem(KEY, JSON.stringify(choice));
    } catch (e) {
      /* storage is blocked: the choice still applies to this page */
    }
  }

  function validTheme(v) {
    return THEMES.indexOf(v) >= 0;
  }

  function validPalette(v) {
    return typeof v === "string" && PALETTE.test(v);
  }

  function clean(v) {
    var c = {};
    if (validTheme(v.theme)) c.theme = v.theme;
    if (validPalette(v.palette)) c.palette = v.palette;
    return c;
  }

  function setAttr(name, value) {
    if (value) root.setAttribute(name, value);
    else root.removeAttribute(name);
  }

  function apply() {
    var theme = choice.theme ? (choice.theme === "system" ? null : choice.theme) : server.theme;
    setAttr("data-theme", theme);
    setAttr("data-palette", choice.palette || server.palette);
  }

  function get() {
    return {
      theme: choice.theme || server.theme || "system",
      palette: choice.palette || server.palette || "blue"
    };
  }

  // matches says whether a choice element names the current theme and palette.
  // An empty value means system, or the page's default palette.
  function matches(el, now) {
    if (el.hasAttribute("data-ui-set-theme") && (el.getAttribute("data-ui-set-theme") || "system") !== now.theme) return false;
    if (el.hasAttribute("data-ui-set-palette") && (el.getAttribute("data-ui-set-palette") || server.palette || "blue") !== now.palette) return false;
    return true;
  }

  var CHECKED = /^(menuitemradio|menuitemcheckbox|radio)$/;

  function sync(scope) {
    var now = get();
    var els = (scope || document).querySelectorAll("[data-ui-set-theme],[data-ui-set-palette]");
    for (var i = 0; i < els.length; i++) {
      var el = els[i];
      el.setAttribute(CHECKED.test(el.getAttribute("role") || "") ? "aria-checked" : "aria-pressed", matches(el, now) ? "true" : "false");
    }
  }

  function announce() {
    sync();
    var detail = get();
    var e;
    try {
      e = new CustomEvent("aiui:themechange", { detail: detail });
    } catch (err) {
      e = document.createEvent("CustomEvent");
      e.initCustomEvent("aiui:themechange", false, false, detail);
    }
    document.dispatchEvent(e);
  }

  function set(patch) {
    if (!patch || typeof patch !== "object") return get();
    if (patch.theme === null) delete choice.theme;
    else if (validTheme(patch.theme)) choice.theme = patch.theme;
    if (patch.palette === null) delete choice.palette;
    else if (validPalette(patch.palette)) choice.palette = patch.palette;
    apply();
    write();
    announce();
    return get();
  }

  apply();

  if (document.readyState === "loading") document.addEventListener("DOMContentLoaded", function () { sync(); });
  else sync();

  // A page that follows the OS changes look when the OS does; say so.
  if (window.matchMedia) {
    var dark = window.matchMedia("(prefers-color-scheme: dark)");
    var onFlip = function () {
      if (get().theme === "system") announce();
    };
    if (dark.addEventListener) dark.addEventListener("change", onFlip);
    else if (dark.addListener) dark.addListener(onFlip);
  }

  document.addEventListener("click", function (e) {
    var t = e.target;
    var el = t && t.closest ? t.closest("[data-ui-set-theme],[data-ui-set-palette]") : null;
    if (!el) return;
    var patch = {};
    if (el.hasAttribute("data-ui-set-theme")) patch.theme = el.getAttribute("data-ui-set-theme") || null;
    if (el.hasAttribute("data-ui-set-palette")) patch.palette = el.getAttribute("data-ui-set-palette") || null;
    set(patch);
  });

  window.aiui = window.aiui || {};
  window.aiui.theme = { get: get, set: set, sync: sync };
})();
