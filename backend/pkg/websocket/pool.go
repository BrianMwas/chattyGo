package websocket

import "fmt"

type Pool struct {
	// This channel will send out "New user registered"
	Register chan *Client
	// Will un register a user and notify the pool when a user disconnects
	Unregister chan *Client
	// Map of clients to a boolean value. The boolean value can be used to dictate active and inactive but not disconnected
	Clients map[*Client]bool
	// Loops through all clients in the pool and send the message through the socket connection
	Broadcast chan Message
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("Size of connection pool: ", len(pool.Clients))

			for client, _ := range pool.Clients {
				fmt.Println(client)
				client.Conn.WriteJSON(Message{Type: 1, Body: "New User just joined"})
			}
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of connection pool", len(pool.Clients))
			for client, _ := range pool.Clients {
				client.Conn.WriteJSON(Message{Type: 1, Body: "User disconnected"})
			}
			break

		case message := <-pool.Broadcast:
			fmt.Println("Sending message to all clients in Pool")

			for client, _ := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}

		}
	}
}

// Use of & is to access its memory address
// while * resolves what follows it goes and gets it from the memory address
func NewPool() *Pool {

	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}
