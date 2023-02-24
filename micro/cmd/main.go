package main

import (
	"fmt"
	"time"

	"relorer/test_go_microservice/config"
	"relorer/test_go_microservice/internal/cache"
	"relorer/test_go_microservice/internal/database"
	"relorer/test_go_microservice/internal/server"
	"relorer/test_go_microservice/internal/server/handler"

	"github.com/gin-gonic/gin"
)

func main() {

	conf, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	reindexerParams := &database.ReindexerParams{
		Host:             conf.Reindexer.Host,
		Port:             conf.Reindexer.Port,
		Database:         conf.Reindexer.Database,
		GenerateTestData: conf.App.Mode == "debug",
	}

	reindexerConnection := database.ReindexerConnectWithRetry(reindexerParams, 5*time.Second)

	documentRepo := database.NewDocumentReindexerRepository(reindexerConnection)
	authorRepo := database.NewAuthorReindexerRepository(reindexerConnection)
	commentRepo := database.NewCommentReindexerRepository(reindexerConnection)

	simpleCache := cache.NewSimpleCache(time.Duration(conf.App.TTL)*time.Minute, time.Duration(conf.App.CleanupInterval)*time.Minute)
	documentCacheRepo := cache.NewDocumentCacheRepository(simpleCache, documentRepo)

	documentCRUDGroup := server.CRUDGroup{
		Handler:   handler.NewDocumentHandler(documentCacheRepo),
		Namespace: database.DocumentsNamespace,
	}

	authorCRUDGroup := server.CRUDGroup{
		Handler:   handler.NewAuthorHandler(authorRepo),
		Namespace: database.AuthorsNamespace,
	}

	commentCRUDGroup := server.CRUDGroup{
		Handler:   handler.NewCommentHandler(commentRepo),
		Namespace: database.CommentsNamespace,
	}

	server.ConfigureCRUDRoutes(gin.Default(), documentCRUDGroup, authorCRUDGroup, commentCRUDGroup).
		Run(fmt.Sprintf(":%d", conf.App.Port))

}
