//go:build linux

package cms

import (
	"log"

	"github.com/godbus/dbus/v5"
)

// notify shows a native desktop notification via the freedesktop notification
// daemon over the session DBus. It is best-effort: failures are logged and
// otherwise ignored.
func notify(title, message string) {
	go func() {
		conn, err := dbus.SessionBus()
		if err != nil {
			log.Printf("[notify] session bus: %v", err)
			return
		}
		// NOTE: SessionBus returns a process-wide shared connection that the
		// system tray icon (systray) also uses, so it must NOT be closed here —
		// closing it would tear down the tray icon.

		obj := conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
		call := obj.Call("org.freedesktop.Notifications.Notify", 0,
			"Doolox",                  // app_name
			uint32(0),                 // replaces_id
			"doolox",                  // app_icon (theme icon name)
			title,                     // summary
			message,                   // body
			[]string{},                // actions
			map[string]dbus.Variant{}, // hints
			int32(-1),                 // expire_timeout (-1 = server default)
		)
		if call.Err != nil {
			log.Printf("[notify] notification: %v", call.Err)
		}
	}()
}
