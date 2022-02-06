package server

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sayan-biswas/file-store/pkg/database"
	"github.com/sayan-biswas/file-store/pkg/server/config"
	"github.com/sayan-biswas/file-store/pkg/server/logger"
	"github.com/sayan-biswas/file-store/pkg/server/middleware"
	"github.com/sayan-biswas/file-store/pkg/server/router"
	"go.uber.org/zap"
)

var httpServer *http.Server
var log = logger.Log
var store database.Store

func Start() {

	defer log.Sync()

	// load config
	err := config.Load()
	if err != nil {
		log.Error(err.Error())
	}
	config, err := config.Get()
	if err != nil {
		log.Error(err.Error())
	}

	// initialize database
	db, err := database.New(&database.Config{Path: config.Database.Path})
	if err != nil {
		log.Fatal(err.Error())
	}
	store = db
	defer store.Close()

	if config.Server.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// initialize router
	engine := gin.Default()

	// middleware chain
	engine.Use(middleware.CORS(config))

	// router chain
	router.Root(engine)
	router.Store(engine, store)

	// define server
	httpServer = &http.Server{
		Addr:    config.Server.Host + ":" + strconv.Itoa(config.Server.Port),
		Handler: engine,
	}

	// start listener
	log.Info("starting server", zap.Int("port", config.Server.Port))
	if config.Server.TLS {
		err = httpServer.ListenAndServeTLS(config.Server.Certificate, config.Server.PrivateKey)
	} else {
		err = httpServer.ListenAndServe()
	}
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err.Error())
	}

}

func Stop() {
	log.Info("shutting down server")
	if err := store.Close(); err != nil {
		log.Fatal(err.Error())
	}
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(context); err != nil {
		log.Fatal(err.Error())
	}
	log.Info("server successfully shutdown")
}
