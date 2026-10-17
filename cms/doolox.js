// Doolox live-preview helper.
//
// This file is loaded only while Doolox is serving the site for preview
// (http://localhost:5000); it is never included in the generated static site.
// It registers a Ctrl+G (Cmd+G on macOS) shortcut that triggers site
// generation via the /generate endpoint. The server shows a native system
// notification with the result (the same notification the tray menu's
// "Generate" action uses).
(function () {
    "use strict";

    function generate() {
        fetch("/generate")
            .then(function (resp) {
                if (!resp.ok) {
                    throw new Error("HTTP " + resp.status);
                }
                return resp.text();
            })
            .then(function (text) {
                console.log("[Doolox] " + text.trim());
            })
            .catch(function (err) {
                console.error("[Doolox] site generation failed: " + err);
            });
    }

    document.addEventListener("keydown", function (event) {
        if ((event.ctrlKey || event.metaKey) && !event.altKey && event.key.toLowerCase() === "g") {
            event.preventDefault();
            generate();
        }
    });
})();


