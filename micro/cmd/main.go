package main

import (
	"fmt"
	"time"

	"relorer/test_go_microservice/config"
	"relorer/test_go_microservice/internal/db"
	"relorer/test_go_microservice/internal/server"
)

func reindexerConnectWithRetry(params *db.ReindexerParams, delay time.Duration) *db.ReindexerDB {
	for {
		db, err := db.NewReindexerDB(params)

		if err != nil {
			fmt.Println(err)
		} else {
			return db
		}

		time.Sleep(delay)
	}
}

func main() {

	conf, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	params := &db.ReindexerParams{
		Host:     conf.Reindexer.Host,
		Port:     conf.Reindexer.Port,
		Database: conf.Reindexer.Database,
	}

	reindexerDb := reindexerConnectWithRetry(params, 5*time.Second)

	serverHandler := server.NewHandler(reindexerDb)

	server.StartServer(serverHandler, conf.App.Port)

}
