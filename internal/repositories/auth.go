package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"weather_back_api_getway/config"
	models "weather_back_api_getway/internal"
)

var (
	ErrUserExists   = errors.New("User already exists")
	ErrUserNotFound = errors.New("User not found")
)

type Database struct {
	db *sql.DB
}

type Auth interface {
	SaveUser(ctx context.Context, uuid string, username string, password []byte, email string) (string, error)
}

func New(d config.Database) *Database {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		d.Server,
		d.Port,
		d.Username,
		d.Password,
		d.Database,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	slog.Info("Successfully connected to database")

	return &Database{db: db}
}

func Stop(db *Database) {
	db.db.Close()
}

func (db *Database) SaveUser(ctx context.Context, uuid string, username string, password []byte, email string) (string, error) {
	stmt, err := db.db.Prepare("INSERT INTO users(uuid, username, email, password) VALUES (?, ?, ?, ?)")
	if err != nil {
		return "", fmt.Errorf("error preparing statement: %w", err)
	}

	res, err := stmt.ExecContext(ctx, uuid, username, string(password), email)
	if err != nil {
		return "", fmt.Errorf("error executing statement: %w", err)
	}

	fmt.Println(res)
	return "", nil
}

func (db *Database) User(ctx context.Context, login string) (models.User, error) {
	stmt, err := db.db.Prepare("SELECT * FROM users WHERE username=?")
	if err != nil {
		return models.User{}, fmt.Errorf("error preparing statement: %w", err)
	}

	row := stmt.QueryRowContext(ctx, login)
	var user models.User
	err = row.Scan(user.UUID, user.Login, user.Email, user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("error scanning row: %w", err)
	}

	return user, nil
}
