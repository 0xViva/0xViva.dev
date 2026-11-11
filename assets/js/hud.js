(function () {
  "use strict";

  function pad(value) {
    return String(value).padStart(2, "0");
  }

  function renderClock() {
    var now = new Date();
    var text =
      pad(now.getUTCHours()) +
      ":" +
      pad(now.getUTCMinutes()) +
      ":" +
      pad(now.getUTCSeconds()) +
      " UTC";
    var clocks = document.querySelectorAll("[data-hud-clock]");
    for (var i = 0; i < clocks.length; i++) {
      if (clocks[i].textContent !== text) {
        clocks[i].textContent = text;
      }
    }
  }

  renderClock();
  window.setInterval(renderClock, 1000);

  var coarse = window.matchMedia("(pointer: coarse)");
  var reduced = window.matchMedia("(prefers-reduced-motion: reduce)");
  var root = document.documentElement;

  if (!coarse.matches && !reduced.matches) {
    root.classList.add("pointer-live");
  }

  document.addEventListener(
    "pointermove",
    function (event) {
      if (coarse.matches || reduced.matches) {
        return;
      }
      root.style.setProperty("--pointer-x", event.clientX + "px");
      root.style.setProperty("--pointer-y", event.clientY + "px");
    },
    { passive: true }
  );
})();
