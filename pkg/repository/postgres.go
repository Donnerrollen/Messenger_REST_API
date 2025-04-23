package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
)

const (
	usersTable     = "users"
	usersChatTable = "users_chat"
	chatsTable     = "chats"
	messagesTable  = "messages"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPosgresDB(cfg Config) (*pgx.Conn, error) {
	db, err := pgx.Connect(context.Background(), fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode))

	if err != nil {
		return nil, err
	}

	err = db.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return db, nil
}
