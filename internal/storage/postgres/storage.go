package postgres

import (
	"database/sql"
	"errors"

	"github.com/kiricle/file-uploader/internal/lib/hash"
	"github.com/kiricle/file-uploader/internal/models"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(DB_URL string) *Storage {
	db, err := connectDB(DB_URL)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	return &Storage{
		db: db,
	}
}

func connectDB(DB_URL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", DB_URL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (s *Storage) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	query := "SELECT * FROM users WHERE email = $1"
	result := s.db.QueryRow(query, email)

	if err := result.Scan(&user.ID, &user.Email, &user.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, nil
		}
		return models.User{}, err
	}

	return user, nil
}

func (s *Storage) SaveUser(signUpDto models.SignUpDTO) (int64, error) {
	passwordHash, err := hash.HashPassword(signUpDto.Password)
	if err != nil {
		return 0, err
	}

	query := "INSERT INTO users (email, password_hash) VALUES ($1, $2)"
	_, sqlErr := s.db.Exec(query, signUpDto.Email, passwordHash)
	if sqlErr != nil {
		return 0, sqlErr
	}

	return 1, nil
}
