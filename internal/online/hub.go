package online

type Hub struct {
	Clients    map[*Client]bool // mapa de clientes
	Broadcast  chan []byte      // canal donde llegan mensajes para todos
	counter    int              // canal global para el conteno de los clicks
	Register   chan *Client     // canal para registrar clientes
	Unregister chan *Client     // canal para desregistrar clientes
}

func NewHub() *Hub {
	hub := &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		counter:    0,
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

		case <-h.Broadcast:
			h.counter++ // incrementamos el contador
			for client := range h.Clients {
				select {
				case client.Send <- []byte{byte(h.counter)}: // enviamos el numero de click al cliente
				default:

					// Si el canal esta lleno cerramos el canal lleno
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}
