package database

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/regiwitanto/tele-anonymous-chat/internal/models"
)

// DB is the database instance
type DB struct {
	conn *sql.DB
}

// NewDB creates a new database connection
func NewDB(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	db := &DB{conn: conn}
	if err := db.initialize(); err != nil {
		return nil, err
	}

	return db, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// initialize sets up the database tables
func (db *DB) initialize() error {
	query := `
    CREATE TABLE IF NOT EXISTS users (
        user_id INTEGER PRIMARY KEY,
        is_active INTEGER DEFAULT 0,
        current_chat INTEGER,
        last_activity TEXT,
        country TEXT,
        language TEXT,
        gender TEXT
    );
    `

	_, err := db.conn.Exec(query)
	return err
}

// GetUserState retrieves a user's state from the database
func (db *DB) GetUserState(userID int64) (*models.UserState, error) {
	query := `SELECT is_active, current_chat, last_activity, country, language, gender 
              FROM users WHERE user_id = ?`

	row := db.conn.QueryRow(query, userID)

	var isActive int
	var currentChat sql.NullInt64
	var lastActivityStr sql.NullString
	var country, language, gender sql.NullString

	err := row.Scan(&isActive, &currentChat, &lastActivityStr, &country, &language, &gender)
	if err != nil {
		// If no record is found, create a new user state
		if err == sql.ErrNoRows {
			return models.NewUserState(userID), nil
		}
		return nil, err
	}

	lastActivity := time.Now()
	if lastActivityStr.Valid {
		parsedTime, err := time.Parse(time.RFC3339, lastActivityStr.String)
		if err == nil {
			lastActivity = parsedTime
		}
	}

	userState := models.UserState{
		UserID:       userID,
		IsActive:     isActive == 1,
		LastActivity: lastActivity,
		Settings: models.UserSettings{
			Country:  "",
			Language: "",
			Gender:   "",
		},
	}

	if currentChat.Valid {
		userState.CurrentChat = currentChat.Int64
	}

	if country.Valid {
		userState.Settings.Country = country.String
	}

	if language.Valid {
		userState.Settings.Language = language.String
	}

	if gender.Valid {
		userState.Settings.Gender = gender.String
	}

	return &userState, nil
}

// SaveUserState stores a user's state in the database
func (db *DB) SaveUserState(state *models.UserState) error {
	query := `
    INSERT OR REPLACE INTO users 
    (user_id, is_active, current_chat, last_activity, country, language, gender)
    VALUES (?, ?, ?, ?, ?, ?, ?)
    `

	isActive := 0
	if state.IsActive {
		isActive = 1
	}

	lastActivity := state.LastActivity.Format(time.RFC3339)

	_, err := db.conn.Exec(
		query,
		state.UserID,
		isActive,
		state.CurrentChat,
		lastActivity,
		state.Settings.Country,
		state.Settings.Language,
		state.Settings.Gender,
	)

	return err
}

// GetActiveUsers returns the count of active users
func (db *DB) GetActiveUsers() (int, error) {
	query := `SELECT COUNT(*) FROM users WHERE is_active = 1`

	var count int
	err := db.conn.QueryRow(query).Scan(&count)

	return count, err
}

// FindPotentialMatches returns potential matches for a user based on preferences
func (db *DB) FindPotentialMatches(userID int64) ([]int64, error) {
	query := `
    SELECT user_id FROM users 
    WHERE is_active = 1 
      AND current_chat = 0 
      AND user_id != ?
    `

	rows, err := db.conn.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []int64
	for rows.Next() {
		var potentialMatch int64
		if err := rows.Scan(&potentialMatch); err != nil {
			return nil, err
		}
		matches = append(matches, potentialMatch)
	}

	return matches, nil
}
