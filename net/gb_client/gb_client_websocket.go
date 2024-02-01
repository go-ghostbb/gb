package gbclient

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketClient wraps the underlying websocket client connection
// and provides convenient functions.
type WebSocketClient struct {
	*websocket.Dialer
}

// NewWebSocket creates and returns a new WebSocketClient object.
func NewWebSocket() *WebSocketClient {
	return &WebSocketClient{
		&websocket.Dialer{
			Proxy:            http.ProxyFromEnvironment,
			HandshakeTimeout: 45 * time.Second,
		},
	}
}
