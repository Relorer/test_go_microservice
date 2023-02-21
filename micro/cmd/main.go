package main

import (
	"fmt"
	"log"
	"time"

	"relorer/test_go_microservice/config"
	"relorer/test_go_microservice/internal/repository"
	"relorer/test_go_microservice/internal/server"
	"relorer/test_go_microservice/internal/util"
)

func reindexerConnectWithRetry(params *repository.ReindexerParams, delay time.Duration) *repository.ReindexerRepository {
	for {
		db, err := repository.NewReindexerRepository(params)

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

	reindexerRepo := reindexerConnectWithRetry(params, 5*time.Second)

	cache := util.NewCache(15*time.Minute, 15*time.Minute)

	cacheRepo := repository.NewCacheRepository(cache, reindexerRepo)

	serverHandler := server.NewHandler(cacheRepo)

	server.StartServer(serverHandler, conf.App.Port)

}
