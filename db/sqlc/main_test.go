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
var testDB *sql.DB

func TestMain(m *testing.M) {
	var error error

	testDB, error = sql.Open(driverName, dataSourceName)

	if error != nil {
		log.Fatal("Cannon't connect to the database.")
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
