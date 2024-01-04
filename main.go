package main

import (
	"database/sql"
	"log"
	"net/http"

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

func GetOnlineUsers() []*OnlineUsers {
	var users []*OnlineUsers

	rows, err := db.Query("SELECT username FROM users")
	if err != nil {
		log.Println("Erreur lors de la récupération des utilisateurs:", err)
		return nil
	}
	defer rows.Close()

	onlineUsers := make(map[string]bool)
	for conn, username := range clients {
		if conn != nil {
			onlineUsers[username] = true
		}
	}

	for rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			log.Println("Erreur lors de la lecture de la ligne:", err)
			continue
		}

		user := &OnlineUsers{
			Username: username,
			Online:   onlineUsers[username],
		}
		users = append(users, user)
	}

	return users
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/signup", Signup)

	http.HandleFunc("/login", Signin)

	http.HandleFunc("/ws", handleConnections)

	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
