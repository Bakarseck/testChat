package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]string)

var upgrader = websocket.Upgrader{}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	username, err := GetUsernameFromRequest(r)
	if err != nil {
		log.Printf("Échec de récupération de l'utilisateur: %v", err)
		return
	}

	clients[ws] = username
	broadcastUserList()

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}

	delete(clients, ws)
	broadcastUserList()
}

func broadcastUserList() {
	userList := GetOnlineUsers()
	userListJSON, err := json.Marshal(userList)
	if err != nil {
		log.Printf("Erreur lors de la sérialisation de la liste des utilisateurs: %v", err)
		return
	}

	msg := Message{
		Sender:   "Serveur",
		Receiver: "All",
		Content:  string(userListJSON),
	}

	for client := range clients {
		err := client.WriteJSON(msg)
		if err != nil {
			log.Printf("Erreur lors de l'envoi de la liste des utilisateurs: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}
