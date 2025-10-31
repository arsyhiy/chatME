package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

var DB *sql.DB // экспортируемая глобальная переменная

type Storage struct {
	DB *sql.DB
}

// NewStorage создаёт базу (или открывает существующую)
func NewStorage(path string) (*Storage, error) {
	var err error
	DB, err = sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	// создаём таблицы, если их нет
	createUsers := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE
	);
	`
	createMessages := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sender TEXT,
		recipient TEXT,  -- NULL для общих сообщений
		text TEXT,
		created_at DATETIME
	);
	`
	if _, err := DB.Exec(createUsers); err != nil {
		return nil, err
	}
	if _, err := DB.Exec(createMessages); err != nil {
		return nil, err
	}

	return &Storage{DB: DB}, nil
}

// Close закрывает базу
func (s *Storage) Close() {
	s.DB.Close()
}

// Добавить пользователя
func (s *Storage) AddUser(name string) error {
	_, err := s.DB.Exec("INSERT OR IGNORE INTO users(name) VALUES(?)", name)
	return err
}

// Сохранить сообщение
func (s *Storage) AddMessage(sender, recipient, text string) error {
	_, err := s.DB.Exec(
		"INSERT INTO messages(sender, recipient, text, created_at) VALUES(?, ?, ?, ?)",
		sender, recipient, text, time.Now(),
	)
	return err
}

// Получить последние N сообщений
func (s *Storage) GetLastMessages(limit int) ([]string, error) {
	rows, err := s.DB.Query(`
		SELECT sender, recipient, text, created_at
		FROM messages
		ORDER BY created_at DESC
		LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []string{}
	for rows.Next() {
		var sender, recipient, text string
		var ts time.Time
		if err := rows.Scan(&sender, &recipient, &text, &ts); err != nil {
			return nil, err
		}
		line := fmt.Sprintf("[%s] %s: %s", ts.Format("15:04"), sender, text)
		if recipient != "" {
			line = fmt.Sprintf("[%s] %s → %s: %s", ts.Format("15:04"), sender, recipient, text)
		}
		result = append([]string{line}, result...) // добавляем в начало
	}
	return result, nil

}
