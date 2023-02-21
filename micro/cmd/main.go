package main

import (
	"fmt"
	"log"
	"time"

	"relorer/test_go_microservice/config"
	"relorer/test_go_microservice/internal/repository"
	"relorer/test_go_microservice/internal/server"
)

func reindexerConnectWithRetry(params *repository.ReindexerParams, delay time.Duration) *repository.ReindexerRepository {
	for {
		db, err := repository.NewReindexerDB(params)

		if err != nil {
			log.Printf("Error connecting: %s. Retry in %s", err, delay.String())
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

	params := &repository.ReindexerParams{
		Host:     conf.Reindexer.Host,
		Port:     conf.Reindexer.Port,
		Database: conf.Reindexer.Database,
	}

	reindexerDb := reindexerConnectWithRetry(params, 5*time.Second)

	serverHandler := server.NewHandler(reindexerDb)

	server.StartServer(serverHandler, conf.App.Port)

}
