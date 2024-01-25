package main

import "golang.org/x/crypto/bcrypt"

func createUser(username, firstName, lastName, age, gender, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	statement, err := db.Prepare("INSERT INTO users (username, firstName, lastName, age, gender, email, password) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(username, firstName, lastName, age, gender, email, string(hashedPassword))
	return err
}

func checkUser(username, password string) bool {
	row := db.QueryRow("SELECT password FROM users WHERE username = ?", username)

	var hashedPassword string
	err := row.Scan(&hashedPassword)
	if err != nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
