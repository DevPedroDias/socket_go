package main

import (
	"fmt"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
	"github.com/rs/cors"
)

type MessageData struct {
	UserID string `json:"userId"`
	Msg    string `json:"msg"`
}
type RoomHandler struct {
	RoomId string `json:"roomId"`
}

func main() {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("Client disconnected:", s.ID())
		s.LeaveAll()
	})

	server.OnEvent("/", "mensagens", func(s socketio.Conn, data MessageData) {
		fmt.Println("recieved", data.Msg)
		fmt.Println("rooms", server.Rooms("/"))
		server.BroadcastToRoom("/", "pedroEthiago", "responses", data)

	})

	server.OnEvent("/", "join_room", func(s socketio.Conn, input RoomHandler) {
		fmt.Println("joined on room", input)
		s.Join(input.RoomId)
	})

	server.OnEvent("/", "leave_room", func(s socketio.Conn, input RoomHandler) {
		s.Leave(input.RoomId)
	})
	go server.Serve()
	defer server.Close()

	// Configurar CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5500"}, // Substitua pelo URL do seu cliente
		AllowCredentials: true,                              // Permitir credenciais CORS
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	})

	// Configurar roteamento HTTP para o Socket.IO
	http.Handle("/socket.io/", c.Handler(server))

	// Servir outros arquivos estáticos
	http.Handle("/", http.FileServer(http.Dir("./asset")))

	log.Println("Servindo em localhost:4545...")
	log.Fatal(http.ListenAndServe(":4545", nil))
}
