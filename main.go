package main

import (
	"database/sql"
	"log"

	"github.com/julysNICK/simplebank/api"
	db "github.com/julysNICK/simplebank/db/sqlc"
	"github.com/julysNICK/simplebank/utils"
	_ "github.com/lib/pq"
)

func main() {

	config, err := utils.LoadConfig(".") // load config from .env file

	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn, err := sql.Open(config.DBDrive, config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)

	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

}
