/* ai-ui ui.js: load it once, deferred:
     <script src="{{href "ui.js"}}" defer></script>
   React: import "ai-ui/ui.js" once, in a client component.

   It adds what the platform leaves to script for the library's overlays:
   - a menu (.ui-menu, a popover) focuses its checked or first item on open;
     Down/Up/Home/End move through its items; Tab closes it; activating a
     menuitem closes it; where CSS anchor positioning is missing it places
     the menu under its opener;
   - Left/Right/Home/End move between swatches (.ui-swatches);
   - a dialog or sheet with closedby="any" closes on a backdrop click where
     the browser doesn't do that itself.
   It listens on document, so markup added later needs nothing. The menu and
   sheet still open, close and dismiss without it. */
(function () {
  "use strict";
  if (typeof document === "undefined" || (window.aiui && window.aiui.ui)) return;

  var ITEM = ".ui-menu-item:not(:disabled):not([aria-disabled=\"true\"]), .ui-swatch:not(:disabled)";
  var anchored = typeof CSS !== "undefined" && CSS.supports && CSS.supports("position-area", "block-end");
  var closedBy = typeof HTMLDialogElement !== "undefined" && "closedBy" in HTMLDialogElement.prototype;
  var lastOpener = null;

  function visible(el) {
    return el.offsetParent !== null || el.getClientRects().length > 0;
  }

  function items(menu) {
    return Array.prototype.filter.call(menu.querySelectorAll(ITEM), visible);
  }

  function move(list, from, key) {
    var i = list.indexOf(from);
    var n = list.length;
    if (!n) return;
    var next = key === "Home" ? 0 : key === "End" ? n - 1 : key === "ArrowDown" || key === "ArrowRight" ? (i + 1) % n : (i - 1 + n) % n;
    list[next].focus();
  }

  // place puts a menu under its opener, its end edge on the opener's, when
  // the browser can't anchor it with CSS.
  function place(menu) {
    var opener = lastOpener && lastOpener.getAttribute("popovertarget") === menu.id ? lastOpener : document.querySelector('[popovertarget="' + menu.id + '"]');
    var sheet = menu.getAttribute("data-sheet") === "phone" && window.matchMedia("(max-width: 799.98px)").matches;
    if (!opener || sheet) return;
    var r = opener.getBoundingClientRect();
    var m = menu.getBoundingClientRect();
    var gap = 4;
    var top = r.bottom + gap;
    if (top + m.height > window.innerHeight && r.top - gap - m.height >= 0) top = r.top - gap - m.height;
    var start = menu.getAttribute("data-align") === "start";
    var left = start ? r.left : r.right - m.width;
    left = Math.max(gap, Math.min(left, window.innerWidth - m.width - gap));
    menu.style.position = "fixed";
    menu.style.inset = "auto";
    menu.style.margin = "0";
    menu.style.top = top + "px";
    menu.style.left = left + "px";
  }

  document.addEventListener("click", function (e) {
    var t = e.target;
    if (!t || !t.closest) return;
    var opener = t.closest("[popovertarget]");
    if (opener) lastOpener = opener;
    var item = t.closest(".ui-menu-item[role=\"menuitem\"]");
    var menu = item && item.closest(".ui-menu");
    if (menu && !item.hasAttribute("popovertarget") && menu.matches(":popover-open")) {
      // After the click's own work (a link, a form) has started.
      setTimeout(function () {
        if (menu.matches(":popover-open")) menu.hidePopover();
      }, 0);
    }
  });

  document.addEventListener(
    "toggle",
    function (e) {
      var menu = e.target;
      if (!menu.classList || !menu.classList.contains("ui-menu") || e.newState !== "open") return;
      if (!anchored) place(menu);
      var list = items(menu);
      var checked = list.filter(function (el) {
        return el.getAttribute("aria-checked") === "true" || el.getAttribute("aria-pressed") === "true";
      })[0];
      var first = checked || list[0];
      if (first) first.focus();
    },
    true
  );

  document.addEventListener("keydown", function (e) {
    var t = e.target;
    if (!t || !t.closest || e.altKey || e.ctrlKey || e.metaKey) return;
    var menu = t.closest(".ui-menu");
    var row = t.closest(".ui-swatches");
    var k = e.key;
    if (row && t.matches(".ui-swatch") && (k === "ArrowLeft" || k === "ArrowRight" || (!menu && (k === "Home" || k === "End")))) {
      var swatches = Array.prototype.filter.call(row.querySelectorAll(".ui-swatch:not(:disabled)"), visible);
      if (getComputedStyle(row).direction === "rtl" && (k === "ArrowLeft" || k === "ArrowRight")) k = k === "ArrowLeft" ? "ArrowRight" : "ArrowLeft";
      move(swatches, t, k);
      e.preventDefault();
      return;
    }
    if (!menu) return;
    if (k === "ArrowDown" || k === "ArrowUp" || k === "Home" || k === "End") {
      move(items(menu), t, k);
      e.preventDefault();
    } else if (k === "Tab" && menu.matches(":popover-open")) {
      // The browser returns focus to the opener; Tab then moves on from it.
      menu.hidePopover();
    }
  });

  // closedby="any": a click on the backdrop (the dialog element itself,
  // outside its box) closes it.
  if (!closedBy) {
    document.addEventListener("click", function (e) {
      var d = e.target;
      if (!d || d.tagName !== "DIALOG" || !d.open || d.getAttribute("closedby") !== "any") return;
      var r = d.getBoundingClientRect();
      if (e.clientX < r.left || e.clientX > r.right || e.clientY < r.top || e.clientY > r.bottom) d.close();
    });
  }

  window.aiui = window.aiui || {};
  window.aiui.ui = { version: 1 };
})();
