package main

import (
	"fmt"
	"time"

	"relorer/test_go_microservice/config"
	"relorer/test_go_microservice/internal/repository"
	"relorer/test_go_microservice/internal/server"
	"relorer/test_go_microservice/internal/util"
)

func main() {

	conf, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	reindexerParams := &repository.ReindexerParams{
		Host:     conf.Reindexer.Host,
		Port:     conf.Reindexer.Port,
		Database: conf.Reindexer.Database,
	}

	reindexerRepo := repository.ReindexerConnectWithRetry(reindexerParams, 5*time.Second)

	cache := util.NewCache(time.Duration(conf.App.TTL)*time.Minute, time.Duration(conf.App.CleanupInterval)*time.Minute)
	cacheRepo := repository.NewCacheRepository(cache, reindexerRepo)

	serverHandler := server.NewDefaultHandler(cacheRepo)

	server.StartServer(serverHandler, conf.App.Port)

}
