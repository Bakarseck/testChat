package main

import (
	"fmt"
	"net/http"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	lastName := r.FormValue("lastName")
	firstName := r.FormValue("firstName")
	age := r.FormValue("age")
	gender := r.FormValue("gender")
	email := r.FormValue("email")
	password := r.FormValue("password")

	err := createUser(username, firstName, lastName, age, gender, email, password)
	if err != nil {
		http.Error(w, "Erreur lors de la création de l'utilisateur" + err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Inscription réussie"))
}

func Signin(w http.ResponseWriter, r *http.Request) {
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
		broadcastUserList()
	} else {
		http.Error(w, "Identifiants invalides", http.StatusUnauthorized)
	}
}

func VerifySession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "no cookie", http.StatusBadRequest)
		return
	}

	for token := range Sessions {
		if token == cookie.Value {
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	w.WriteHeader(http.StatusBadRequest)
}
