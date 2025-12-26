package clickOnline

import (
	"log"
	"net/http"

	"clickOnline/internal/online"
)

func (a *ApiServer) serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := a.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &online.Client{
		Hub:  a.hub,
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
