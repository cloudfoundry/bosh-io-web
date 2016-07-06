package index

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	boshlog "github.com/cloudfoundry/bosh-agent/logger"
)

type PostgresAdapterPool struct {
	conn *sql.DB

	logTag string
	logger boshlog.Logger
}

type PostgresAdapter struct {
	tableName string
	conn      *sql.DB

	logTag string
	logger boshlog.Logger
}

func NewPostgresAdapterPool(url string, logger boshlog.Logger) (PostgresAdapterPool, error) {
	conn, err := sql.Open("postgres", url)
	if err != nil {
		return PostgresAdapterPool{}, err
	}

	return PostgresAdapterPool{conn: conn, logTag: "PostgresAdapterPool", logger: logger}, nil
}

func (p PostgresAdapterPool) NewAdapter(tableName string) (PostgresAdapter, error) {
	logTag := "PostgresAdapter"

	q := fmt.Sprintf("CREATE TABLE %s (Key BYTEA UNIQUE, Value BYTEA)", tableName)

	p.logger.Debug(logTag, "Executing '%s'", q)

	_, err := p.conn.Exec(q)
	if err != nil {
		if typedErr, ok := err.(*pq.Error); ok {
			if typedErr.Code.Name() == "duplicate_table" {
				// todo check to make sure schema is the same
				goto returnAdapter
			}
		}

		return PostgresAdapter{}, err
	}

returnAdapter:
	adapter := PostgresAdapter{
		tableName: tableName,
		conn:      p.conn,

		logTag: logTag,
		logger: p.logger,
	}

	return adapter, nil
}

func (a PostgresAdapter) Clear() (sql.Result, error) {
	q := fmt.Sprintf("DELETE FROM %s", a.tableName)

	a.logger.Debug(a.logTag, "Executing '%s'", q)

	res, err := a.conn.Exec(q)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a PostgresAdapter) List() (*sql.Rows, error) {
	q := fmt.Sprintf("SELECT Value FROM %s", a.tableName)

	a.logger.Debug(a.logTag, "Executing '%s'", q)

	rows, err := a.conn.Query(q)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (a PostgresAdapter) ListKeys() (*sql.Rows, error) {
	q := fmt.Sprintf("SELECT Key FROM %s", a.tableName)

	a.logger.Debug(a.logTag, "Executing '%s'", q)

	rows, err := a.conn.Query(q)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (a PostgresAdapter) Find(keyBytes []byte) (*sql.Rows, error) {
	q := fmt.Sprintf("SELECT Value FROM %s WHERE Key = $1", a.tableName)

	a.logger.Debug(a.logTag, "Executing '%s' $1 = '%s'", q, string(keyBytes))

	rows, err := a.conn.Query(q, keyBytes)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (a PostgresAdapter) Save(keyBytes, valueBytes []byte) (sql.Result, error) {
	q := fmt.Sprintf("INSERT INTO %s (Key, Value) VALUES($1, $2)", a.tableName)

	a.logger.Debug(a.logTag, "Executing '%s' $1 = '%s' $2 = skip", q, string(keyBytes))

	res, err := a.conn.Exec(q, keyBytes, valueBytes)
	if err != nil {
		if typedErr, ok := err.(*pq.Error); ok {
			if typedErr.Code.Name() == "unique_violation" {
				return a.saveByUpdating(keyBytes, valueBytes)
			}
		}

		return nil, err
	}

	return res, nil
}

func (a PostgresAdapter) Insert(keyBytes, valueBytes []byte) (sql.Result, error) {
	q := fmt.Sprintf("INSERT INTO %s (Key, Value) VALUES($1, $2)", a.tableName)

	a.logger.Debug(a.logTag, "Executing '%s' $1 = '%s' $2 = skip", q, string(keyBytes))

	res, err := a.conn.Exec(q, keyBytes, valueBytes)
	if err != nil {
		if typedErr, ok := err.(*pq.Error); ok {
			if typedErr.Code.Name() == "unique_violation" {
				return nil, ErrExists
			}
		}

		return nil, err
	}

	return res, nil
}

func (a PostgresAdapter) Remove(keyBytes []byte) error {
	q := fmt.Sprintf("DELETE FROM %s WHERE Key = $1", a.tableName)

	a.logger.Debug(a.logTag, "Executing '%s' $1 = '%s'", q, string(keyBytes))

	_, err := a.conn.Exec(q, keyBytes)
	if err != nil {
		return err
	}

	return nil
}

func (a PostgresAdapter) FindLocked(keyBytes []byte) (*sql.Tx, *sql.Rows, error) {
	tx, err := a.conn.Begin()
	if err != nil {
		return nil, nil, err
	}

	q := fmt.Sprintf("SELECT Value FROM %s WHERE Key = $1 FOR UPDATE", a.tableName)

	a.logger.Debug(a.logTag, "Executing '%s' $1 = '%s'", q, string(keyBytes))

	rows, err := tx.Query(q, keyBytes)
	if err != nil {
		// todo error check
		tx.Rollback()
		return nil, nil, err
	}

	return tx, rows, nil
}

func (a PostgresAdapter) ReleaseLocked(tx *sql.Tx) error {
	return tx.Rollback()
}

func (a PostgresAdapter) InsertLocked(tx *sql.Tx, keyBytes, valueBytes []byte) (sql.Result, error) {
	q := fmt.Sprintf("INSERT INTO %s (Key, Value) VALUES($1, $2)", a.tableName)

	a.logger.Debug(a.logTag, "Executing '%s' $1 = '%s' $2 = skip", q, string(keyBytes))

	res, err := tx.Exec(q, keyBytes, valueBytes)
	if err != nil {
		// todo error check
		tx.Rollback()

		if typedErr, ok := err.(*pq.Error); ok {
			if typedErr.Code.Name() == "unique_violation" {
				return nil, ErrExists
			}
		}

		return nil, err
	}

	return res, tx.Commit()
}

func (a PostgresAdapter) UpdateLocked(tx *sql.Tx, keyBytes, valueBytes []byte) (sql.Result, error) {
	q := fmt.Sprintf("UPDATE %s SET Value = $2 WHERE Key = $1", a.tableName)

	a.logger.Debug(a.logTag, "Executing '%s' $1 = '%s' $2 = skip", q, string(keyBytes))

	res, err := tx.Exec(q, keyBytes, valueBytes)
	if err != nil {
		// todo error check
		tx.Rollback()
		return nil, err
	}

	return res, tx.Commit()
}

func (a PostgresAdapter) RemoveLocked(tx *sql.Tx, keyBytes []byte) error {
	q := fmt.Sprintf("DELETE FROM %s WHERE Key = $1", a.tableName)

	a.logger.Debug(a.logTag, "Executing '%s' $1 = '%s' $2 = skip", q, string(keyBytes))

	_, err := tx.Exec(q, keyBytes)
	if err != nil {
		// todo error check
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (a PostgresAdapter) saveByUpdating(keyBytes, valueBytes []byte) (sql.Result, error) {
	q := fmt.Sprintf("UPDATE %s SET Value = $2 WHERE Key = $1", a.tableName)

	a.logger.Debug(a.logTag, "Executing '%s' $1 = '%s' $2 = skip", q, string(keyBytes))

	res, err := a.conn.Exec(q, keyBytes, valueBytes)
	if err != nil {
		return nil, err
	}

	return res, nil
}
