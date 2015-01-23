package index

import (
	"database/sql"
	"errors"
)

var (
	ErrExists = errors.New("Record already exists")
)

type DBAdapter interface {
	List() (*sql.Rows, error)
	ListKeys() (*sql.Rows, error)

	Find([]byte) (*sql.Rows, error)
	Save([]byte, []byte) (sql.Result, error)
	Remove([]byte) error

	FindLocked([]byte) (*sql.Tx, *sql.Rows, error)
	ReleaseLocked(*sql.Tx) error
	InsertLocked(*sql.Tx, []byte, []byte) (sql.Result, error)
	UpdateLocked(*sql.Tx, []byte, []byte) (sql.Result, error)
	RemoveLocked(*sql.Tx, []byte) error
}
