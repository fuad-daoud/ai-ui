// Copy buttons and Drop it again. Both start hidden, so a page without
// scripts doesn't offer buttons that can't work.
(function () {
  "use strict";

  var status = document.getElementById("copy-status");

  document.querySelectorAll("[data-copy]").forEach(function (button) {
    var label = button.textContent;
    var timer;
    function say(text) {
      button.textContent = text;
      status.textContent = text;
      clearTimeout(timer);
      timer = setTimeout(function () { button.textContent = label; }, 2500);
    }
    button.hidden = false;
    button.addEventListener("click", function () {
      var text = document.getElementById("prompt").textContent.trim();
      var fail = function () { say("Couldn't copy, see below"); };
      if (!navigator.clipboard) { fail(); return; }
      navigator.clipboard.writeText(text).then(function () { say("Copied"); }, fail);
    });
  });

  // A theme picker demo themes only its own box (the demo, or the pile): a
  // palette and a mode apply on any element. The site itself stays dark.
  document.querySelectorAll("[data-demo-theme]").forEach(function (row) {
    var demo = row.closest(".demo, [data-demo-scope]");
    var mark = function () {
      row.querySelectorAll(".ui-swatch").forEach(function (b) {
        var on = b.getAttribute("data-demo-set-theme") === (demo.getAttribute("data-theme") || "dark") &&
          b.getAttribute("data-demo-set-palette") === (demo.getAttribute("data-palette") || "blue");
        b.setAttribute("aria-pressed", on ? "true" : "false");
      });
    };
    row.hidden = false;
    mark();
    row.addEventListener("click", function (e) {
      var b = e.target.closest(".ui-swatch");
      if (!b) return;
      demo.setAttribute("data-theme", b.getAttribute("data-demo-set-theme"));
      demo.setAttribute("data-palette", b.getAttribute("data-demo-set-palette"));
      demo.style.color = "var(--ink)";
      mark();
    });
  });

  var pile = document.getElementById("pile");
  var drop = document.getElementById("drop");
  if (pile && drop) {
    drop.hidden = false;
    // Two identical animations, so swapping them restarts the fall without
    // re-creating the components or clearing anything typed into them.
    drop.addEventListener("click", function () {
      pile.classList.toggle("run-a");
      pile.classList.toggle("run-b");
    });
  }
})();
