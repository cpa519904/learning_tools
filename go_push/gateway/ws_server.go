package gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"strings"
)

var wsUpgrader = websocket.Upgrader{
	// 允许所有CORS跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func getConn(res http.ResponseWriter, req *http.Request) {
	var (
		err    error
		wsConn *websocket.Conn
	)
	if wsConn, err = wsUpgrader.Upgrade(res, req, nil); err != nil {
		return
	}
	gin.Context{}.ClientIP()
	ws := NewWsConnection(wsConn)
	ws.SetIp(ClientIP(req))
	ws.WsHandle()
}

func ClientIP(c *http.Request) string {
	clientIP := c.Header.Get("X-Forwarded-For")
	clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
	if clientIP == "" {
		clientIP = strings.TrimSpace(c.Header.Get(("X-Real-Ip")))
	}
	if clientIP != "" {
		return clientIP
	}
	addr := c.Header.Get("X-Appengine-Remote-Addr")
	if addr != "" {
		return addr
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(c.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}
