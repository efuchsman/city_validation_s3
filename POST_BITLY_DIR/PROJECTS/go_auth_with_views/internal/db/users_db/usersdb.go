package usersdb

import (
	"database/sql"
	"fmt"

	"github.com/DATA-DOG/go-txdb"
	_ "github.com/lib/pq"
)

type DB struct {
	Conn  *sql.DB
	TxDB  bool   // Flag to indicate whether to use txdb (only use for testing)
	TxDrv string // Unique name for txdb registration
}

func NewDB(connStr string, useTxDB bool, TxDrv string) (*DB, error) {
	var db *sql.DB
	var err error

	if useTxDB {
		txdb.Register(TxDrv, "postgres", connStr)
		db, err = sql.Open(TxDrv, "")
	} else {
		db, err = sql.Open("postgres", connStr)
	}

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to the database")
	return &DB{
		Conn:  db,
		TxDB:  useTxDB,
		TxDrv: TxDrv,
	}, nil
}

func (db *DB) Close() {
	db.Conn.Close()
	fmt.Println("Closed the database connection")
}
