package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	todo "rest_API"
)

type AuthPostgres struct {
	db *pgx.Conn
}

func NewAuthPostgres(db *pgx.Conn) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user todo.User) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s (name, username, password) values ('%s', '%s', '%s') RETURNING id", usersTable, user.Name, user.Username, user.Password)
	row := r.db.QueryRow(context.Background(), query)
	row.Scan(&id)

	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (todo.User, error) {
	var user todo.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username='%s' AND password='%s'", usersTable, username, password)
	row := r.db.QueryRow(context.Background(), query)
	row.Scan(&user.Id)

	return user, nil
}
