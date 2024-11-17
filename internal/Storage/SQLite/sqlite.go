package SQLite

import (
	Storage2 "SupChat/internal/Storage"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS tickets(
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
		    openDate STRING NOT NULL,
		    closed BOOLEAN NOT NULL,
		    closeDate STRING
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}
	if _, err := stmt.Exec(); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}
	stmt, err = db.Prepare(`
		CREATE TABLE IF NOT EXISTS messages(
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
		    message STRING NOT NULL,
		    'from' STRING NOT NULL,
		    'date' STRING NOT NULL,
		    ticketsId INTEGER REFERENCES tickets(id)
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}
	if _, err := stmt.Exec(); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) OpenTicket() (int, error) {
	stmt, err := s.db.Prepare("INSERT INTO tickets(openDate, closed) VALUES (?, ?)")
	if err != nil {
		return -1, fmt.Errorf("failed to prepare statement: %w", err)
	}
	res, err := stmt.Exec(time.Now().Format("15:01:05 02.01.2006"), false)
	if err != nil {
		return -1, fmt.Errorf("failed to execute statement: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("failed to get last insert ID: %w", err)
	}
	return int(id), nil
}

func (s *Storage) CloseTicket(id int) error {
	stmt, err := s.db.Prepare("UPDATE tickets SET closeDate=? WHERE id=?")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	_, err = stmt.Exec(time.Now().Format("15:01:05 02.01.2006"), id)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}
	stmt, err = s.db.Prepare("UPDATE tickets SET closed=? WHERE id=?")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	_, err = stmt.Exec(true, id)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}
	return nil
}

func (s *Storage) CreateMessage(message, from string, ticketsId int) error {
	stmt, err := s.db.Prepare("INSERT INTO messages(message, `from`, date, ticketsId) VALUES(?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	_, err = stmt.Exec(message, from, time.Now().Format("15:01:05 02.01.2006"), ticketsId)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}
	return nil
}

func (s *Storage) GetTickets() ([]int, error) {
	var ticketsId []int
	rows, err := s.db.Query("SELECT id FROM tickets where closed != true")
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		ticketsId = append(ticketsId, id)
	}
	return ticketsId, nil
}

func (s *Storage) GetMessages(ticketId int) ([]Storage2.Message, error) {
	var messages []Storage2.Message
	rows, err := s.db.Query("SELECT id, message, `from`, date FROM messages WHERE ticketsId=?", ticketId)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	for rows.Next() {
		var message Storage2.Message
		err = rows.Scan(&message.Id, &message.Message, &message.From, &message.Date)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (s *Storage) CheckOpenness(ticketId int) (bool, error) {
	row := s.db.QueryRow("SELECT closed FROM tickets where id=?", ticketId)
	var closed bool
	if err := row.Scan(&closed); err != nil {
		return false, fmt.Errorf("failed to scan row: %w", err)
	}
	return closed, nil
}
