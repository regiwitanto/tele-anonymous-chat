package database

import (
	"time"
)

// ActiveChat represents a chat between two users with their last activity times
type ActiveChat struct {
	User1ID       int64
	User2ID       int64
	LastActivity1 time.Time
	LastActivity2 time.Time
}

// GetActiveChats retrieves all active chats from the database
func (db *DB) GetActiveChats() ([]ActiveChat, error) {
	query := `
	SELECT u1.user_id, u2.user_id, u1.last_activity, u2.last_activity 
	FROM users u1, users u2
	WHERE u1.current_chat = u2.user_id 
	  AND u2.current_chat = u1.user_id
	  AND u1.user_id < u2.user_id`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []ActiveChat
	for rows.Next() {
		var user1ID, user2ID int64
		var lastActivity1, lastActivity2 string

		if err := rows.Scan(&user1ID, &user2ID, &lastActivity1, &lastActivity2); err != nil {
			return nil, err
		}

		// Parse timestamps
		t1, err := time.Parse(time.RFC3339, lastActivity1)
		if err != nil {
			return nil, err
		}

		t2, err := time.Parse(time.RFC3339, lastActivity2)
		if err != nil {
			return nil, err
		}

		chat := ActiveChat{
			User1ID:       user1ID,
			User2ID:       user2ID,
			LastActivity1: t1,
			LastActivity2: t2,
		}

		chats = append(chats, chat)
	}

	return chats, nil
}
