package index

// import (
//   "database/sql"
//   "fmt"
//
//   _ "code.google.com/p/go-sqlite/go1/sqlite3"
// )
//
// type SqliteAdapter struct {
//   tableName string
//   conn      *sql.DB
// }
//
// func NewMemorySqliteAdapter(tableName string) (SqliteAdapter, error) {
//   conn, err := sql.Open("sqlite3", ":memory:")
//   if err != nil {
//     return SqliteAdapter{}, err
//   }
//
//   q := fmt.Sprintf("CREATE TABLE %s (Key BLOB UNIQUE, Value BLOB)", tableName)
//
//   _, err = conn.Exec(q)
//   if err != nil {
//     return SqliteAdapter{}, err
//   }
//
//   return SqliteAdapter{tableName: tableName, conn: conn}, nil
// }
//
// func (a SqliteAdapter) List() (*sql.Rows, error) {
//   q := fmt.Sprintf("SELECT Value FROM %s", a.tableName)
//
//   rows, err := a.conn.Query(q)
//   if err != nil {
//     return nil, err
//   }
//
//   return rows, nil
// }
//
// func (a SqliteAdapter) ListKeys() (*sql.Rows, error) {
//   q := fmt.Sprintf("SELECT Key FROM %s", a.tableName)
//
//   rows, err := a.conn.Query(q)
//   if err != nil {
//     return nil, err
//   }
//
//   return rows, nil
// }
//
// func (a SqliteAdapter) Find(keyBytes []byte) (*sql.Rows, error) {
//   q := fmt.Sprintf("SELECT Value FROM %s WHERE Key = ?", a.tableName)
//
//   rows, err := a.conn.Query(q, keyBytes)
//   if err != nil {
//     return nil, err
//   }
//
//   return rows, nil
// }
//
// func (a SqliteAdapter) Save(keyBytes, valueBytes []byte) (sql.Result, error) {
//   q := fmt.Sprintf("INSERT OR REPLACE INTO %s (Key, Value) VALUES(?, ?)", a.tableName)
//
//   res, err := a.conn.Exec(q, keyBytes, valueBytes)
//   if err != nil {
//     return nil, err
//   }
//
//   return res, nil
// }
//
// func (a SqliteAdapter) Remove(keyBytes []byte) error {
//   q := fmt.Sprintf("DELETE FROM %s WHERE Key = ?", a.tableName)
//
//   _, err := a.conn.Exec(q, keyBytes)
//   if err != nil {
//     return err
//   }
//
//   return nil
// }
