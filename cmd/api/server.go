package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"clickOnline/internal/online"
)

type ApiServer struct {
	router   *mux.Router
	upgrader *websocket.Upgrader
	hub      *online.Hub
	port     string
}

func newApiServer(port string, router *mux.Router, upgrader *websocket.Upgrader, hub *online.Hub) *ApiServer {
	return &ApiServer{
		router:   router,
		upgrader: upgrader,
		hub:      hub,
		port:     port,
	}
}

func newUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,

		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func main() {
	router := mux.NewRouter()
	upgrader := newUpgrader()
	hub := online.NewHub()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := newApiServer(port, router, upgrader, hub)
	server.route()

	log.Fatal(http.ListenAndServe(":"+port, server.router))
}
