package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var errUserNotFound = errors.New("user not found")
var errUserNameAlreadyExists = errors.New("user name already exists")

type user struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type authUser struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
}

type userRepository struct {
	db *sql.DB
}

func newUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) create(name, email, passwordHash string) (user, error) {
	result, err := r.db.Exec(
		"INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?)",
		name,
		email,
		passwordHash,
	)
	if err != nil {
		if isUniqueConstraintError(err) {
			return user{}, errUserNameAlreadyExists
		}
		return user{}, fmt.Errorf("insert user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return user{}, fmt.Errorf("read inserted id: %w", err)
	}

	return user{
		ID:    fmt.Sprintf("%d", id),
		Name:  name,
		Email: email,
	}, nil
}

func (r *userRepository) get(id string) (user, error) {
	var foundUser user

	err := r.db.QueryRow(
		"SELECT id, name, email FROM users WHERE id = ? AND deleted_at IS NULL",
		id,
	).Scan(&foundUser.ID, &foundUser.Name, &foundUser.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return user{}, errUserNotFound
	}
	if err != nil {
		return user{}, fmt.Errorf("select user: %w", err)
	}

	return foundUser, nil
}

func (r *userRepository) update(id, name, email, passwordHash string) (user, error) {
	result, err := r.db.Exec(
		"UPDATE users SET name = ?, email = ?, password_hash = ? WHERE id = ? AND deleted_at IS NULL",
		name,
		email,
		passwordHash,
		id,
	)
	if err != nil {
		if isUniqueConstraintError(err) {
			return user{}, errUserNameAlreadyExists
		}
		return user{}, fmt.Errorf("update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return user{}, fmt.Errorf("read updated rows: %w", err)
	}
	if rowsAffected == 0 {
		return user{}, errUserNotFound
	}

	return user{
		ID:    id,
		Name:  name,
		Email: email,
	}, nil
}

func (r *userRepository) findByName(name string) (authUser, error) {
	var foundUser authUser

	err := r.db.QueryRow(
		"SELECT id, name, email, password_hash FROM users WHERE name = ? AND deleted_at IS NULL",
		name,
	).Scan(&foundUser.ID, &foundUser.Name, &foundUser.Email, &foundUser.PasswordHash)
	if errors.Is(err, sql.ErrNoRows) {
		return authUser{}, errUserNotFound
	}
	if err != nil {
		return authUser{}, fmt.Errorf("select user by name: %w", err)
	}

	return foundUser, nil
}

func (r *userRepository) getByName(name string) (user, error) {
	foundUser, err := r.findByName(name)
	if err != nil {
		return user{}, err
	}

	return user{
		ID:    foundUser.ID,
		Name:  foundUser.Name,
		Email: foundUser.Email,
	}, nil
}

func (r *userRepository) list() ([]user, error) {
	rows, err := r.db.Query(
		"SELECT id, name, email FROM users WHERE deleted_at IS NULL ORDER BY id ASC",
	)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	var users []user
	for rows.Next() {
		var foundUser user
		if err := rows.Scan(&foundUser.ID, &foundUser.Name, &foundUser.Email); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}

		users = append(users, foundUser)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}

	return users, nil
}

func (r *userRepository) delete(id string) error {
	result, err := r.db.Exec(
		"UPDATE users SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL",
		id,
	)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read deleted rows: %w", err)
	}
	if rowsAffected == 0 {
		return errUserNotFound
	}

	return nil
}

func isUniqueConstraintError(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "unique")
}
