package main

func CreateMessage(sender, receiver, content string) error {
	statement, err := db.Prepare("INSERT INTO messages (sender, receiver, content) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(sender, receiver, content)
	return err
}

func GetMessages(sender, receiver string) ([]Message, error) {
	rows, err := db.Query("SELECT sender, receiver, content, timestamp FROM messages WHERE (sender = ? AND receiver = ?) OR (sender = ? AND receiver = ?) ORDER BY timestamp", sender, receiver, receiver, sender)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.Sender, &msg.Receiver, &msg.Content, &msg.Timestamp); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
