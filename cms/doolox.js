// Doolox live-preview helper.
//
// This file is loaded only while Doolox is serving the site for preview
// (http://localhost:5000); it is never included in the generated static site.
// It registers two shortcuts:
//   - Ctrl+G (Cmd+G on macOS) — generates the static site via /generate.
//   - Ctrl+P (Cmd+P on macOS) — generates the site, then commits and pushes
//     it with git via /publish.
//
// The server shows a native system notification with the result.
(function () {
    "use strict";

    function request(url, action) {
        fetch(url)
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
                console.error("[Doolox] " + action + " failed: " + err);
            });
    }

    function generate() {
        request("/generate", "site generation");
    }

    function publish() {
        request("/publish", "publish");
    }

    document.addEventListener("keydown", function (event) {
        if (!(event.ctrlKey || event.metaKey) || event.altKey) {
            return;
        }

        var key = event.key.toLowerCase();
        if (key === "g") {
            event.preventDefault();
            generate();
        } else if (key === "p") {
            event.preventDefault();
            publish();
        }
    });
})();


