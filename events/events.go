package events

import (
	"github.com/kataras/go-events"
)

// General events
const (
	RPC_CLIENT_DISCONNECT = `clientDisconnect`
	DOWNLOAD_UPDATE       = `downloadUpdate`
	UPDATE                = `update`
	CLEANUP_COMPLETED     = `cleanupCompleted`
)

var _events = events.New()

func Get() events.EventEmmiter {
	return _events
}
