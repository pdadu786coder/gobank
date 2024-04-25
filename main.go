package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/pdadu/learn/api"
	db "github.com/pdadu/learn/db/sqlc"
	"github.com/pdadu/learn/db/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("unable to read config file", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}
