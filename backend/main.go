package main

import (
	"fmt"
	"net/http"

	"github.com/brianMwas/gochatapp/pkg/websocket"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Home page reached")
}

func serveWS(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Print("Websocket endpoint reached")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%v\n", err)
	}

	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
}

func setupRoutes() {
	pool := websocket.NewPool()

	go pool.Start()

	http.HandleFunc("/", home)
	http.HandleFunc("/ws", func(rw http.ResponseWriter, r *http.Request) {
		serveWS(pool, rw, r)
	})
}

func main() {
	setupRoutes()
	http.ListenAndServe(":8082", nil)
}
