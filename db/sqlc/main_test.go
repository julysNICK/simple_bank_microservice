// package db

// import (
// 	"database/sql"
// 	"log"
// 	"os"
// 	"testing"

// 	_ "github.com/lib/pq"
// )

// const (
// 	dbDriver = "postgres"
// 	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
// )

// var testQueries *Queries
// var testDB *sql.DB

// func TestMain(m *testing.M) {
// 	var err error
// 	testDB, err = sql.Open(dbDriver, dbSource)

// 	if err != nil {
// 		log.Fatal("cannot connect to db: ", err)
// 	}

// 	testQueries = New(testDB)

// 	os.Exit(m.Run())
// }

package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/julysNICK/simplebank/utils"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {

	config, err := utils.LoadConfig("../../")

	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	testDB, err = sql.Open(config.DBDrive, config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
