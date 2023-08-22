package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	driverName     = "postgres"
	dataSourceName = "postgresql://root:PaSSw0rD@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(driverName, dataSourceName)

	if err != nil {
		log.Fatal("Cannon't connect to the database.")
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
