package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type usersStore interface {
	List() (map[int64]User, error)
	Get(id int64) (User, error)
	Add(id int64, user User) error
	Update(id int64, user User) error
	Remove(id int64) error
}

type MemStore struct {
	data map[int64]User
}

func NewMemStore() *MemStore {
	return &MemStore{
		make(map[int64]User),
	}
}

func (m *MemStore) List() (map[int64]User, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "SELECT * FROM users")

	users := make(map[int64]User)
	for rows.Next() {
		var user User

		if err := rows.Scan(&user.Id, &user.Name, &user.Status); err != nil {
			fmt.Fprintf(os.Stderr, "Scan failed: %v\n", err)
			return users, err
		}
		users[user.Id] = user
	}

	if err := rows.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Rows iteration error: %v\n", err)
		return users, err
	}
	return users, nil
}

func (m *MemStore) Get(id int64) (User, error) {
	var user User

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return user, err
	}
	defer conn.Close(context.Background())

	err = conn.QueryRow(
		context.Background(),
		"SELECT * FROM users",
	).Scan(&user.Id, &user.Name, &user.Status)

	return user, err
}

func (m *MemStore) Add(id int64, user User) error {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return err
	}
	defer conn.Close(context.Background())

	var newID int64
	err = conn.QueryRow(
		context.Background(),
		"INSERT INTO users (name, status) VALUES ($1, $2) RETURNING id",
		user.Name, user.Status,
	).Scan(&newID)
	if err != nil {
		return err
	}

	fmt.Printf("New ID: %d\n", newID)
	return nil
}

func (m *MemStore) Update(id int64, user User) error {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return err
	}
	defer conn.Close(context.Background())

	cmdTag, err := conn.Exec(
		context.Background(),
		"UPDATE users SET name = $1, status = $2 WHERE id = $3",
		user.Name, user.Status, id,
	)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() < 1 {
		return fmt.Errorf("Nothing updated")
	}
	fmt.Printf("Updated user with id=%d (new name=%s, status=%v) \n", id, user.Name, user.Id)
	return nil
}

func (m *MemStore) Remove(id int64) error {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return err
	}
	defer conn.Close(context.Background())

	cmdTag, err := conn.Exec(
		context.Background(),
		"DELETE FROM users WHERE id=$1",
		id,
	)
	if err != nil {
		fmt.Println("Unable to delete:", err)
		return err
	}
	fmt.Println("Deleted:", cmdTag)
	return nil
}
