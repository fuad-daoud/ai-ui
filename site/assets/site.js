// Copy buttons and Drop it again. Both start hidden, so a page without
// scripts doesn't offer buttons that can't work.
(function () {
  "use strict";

  document.querySelectorAll("[data-copy]").forEach(function (button) {
    var label = button.textContent;
    var timer;
    function say(text) {
      button.textContent = text;
      clearTimeout(timer);
      timer = setTimeout(function () { button.textContent = label; }, 2500);
    }
    button.hidden = false;
    button.addEventListener("click", function () {
      var text = document.getElementById("prompt").textContent.trim();
      var fail = function () { say("Couldn't copy, select the text below"); };
      if (!navigator.clipboard) { fail(); return; }
      navigator.clipboard.writeText(text).then(function () { say("Copied"); }, fail);
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
