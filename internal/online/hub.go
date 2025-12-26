package online

type Hub struct {
	Clients    map[*Client]bool // mapa de clientes
	Broadcast  chan []byte      // canal donde llegan mensajes para todos
	Register   chan *Client     // canal para registrar clientes
	Unregister chan *Client     // canal para desregistrar clientes
}

func NewHub() *Hub {
	hub := &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}

	go hub.Run()

	return hub
}

/// functions of hub

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true // registramos el cliente
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client) // eliminamos el cliente
				close(client.Send)        // y cerramos el canal
			}

		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message: // enviamos el mensaje al canal de cada cliente
				default:

					// Si el canal esta lleno cerramos el canal lleno
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}
