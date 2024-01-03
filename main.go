package main

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./chat.db")
	if err != nil {
		log.Fatal(err)
	}

	createUserTable := `CREATE TABLE IF NOT EXISTS users (
        username TEXT PRIMARY KEY,
        password TEXT
    );`
	statement, _ := db.Prepare(createUserTable)
	statement.Exec()

	createMessageTable := `CREATE TABLE IF NOT EXISTS messages (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        sender TEXT,
        receiver TEXT,
        content TEXT,
        timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	statement, _ = db.Prepare(createMessageTable)
	statement.Exec()
}

var Sessions = make(map[string]string)
var sessionMutex sync.Mutex

func GenerateSessionID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

var clients = make(map[*websocket.Conn]bool)

var onlineUsers = make(map[string]bool)
var upgrader = websocket.Upgrader{}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	clients[ws] = true

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
}

func GetUsernameFromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return "", err
	}

	sessionToken := cookie.Value
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	username, ok := Sessions[sessionToken]
	if !ok {
		return "", fmt.Errorf("session invalide ou expirée")
	}

	return username, nil
}

func main() {

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		err := createUser(username, password)
		if err != nil {
			http.Error(w, "Erreur lors de la création de l'utilisateur", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Inscription réussie"))
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		if checkUser(username, password) {

			sessionID, err := GenerateSessionID()
			if err != nil {
				http.Error(w, "Erreur lors de la création de la session", http.StatusInternalServerError)
				return
			}

			cookie := http.Cookie{
				Name:   "session_token",
				Value:  sessionID,
				MaxAge: 30 * 3600,
			}
			http.SetCookie(w, &cookie)

			Sessions[sessionID] = username

			w.Write([]byte("Connexion réussie"))
		} else {
			http.Error(w, "Identifiants invalides", http.StatusUnauthorized)
		}
	})

	http.HandleFunc("/ws", handleConnections)

	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
