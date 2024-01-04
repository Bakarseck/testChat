package main

type Message struct {
	Sender    string
	Receiver  string
	Content   string
	Timestamp string
}

type User struct {
	Username string
	Password string
}

type OnlineUsers struct {
	Username string
	Online   bool
}
