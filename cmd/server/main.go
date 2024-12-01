package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/mixailo/go-training-metrics/internal/repository/storage"
	"github.com/mixailo/go-training-metrics/internal/service/database"
	"github.com/mixailo/go-training-metrics/internal/service/logger"
)

func newMux(sa *storageAware) *chi.Mux {
	router := chi.NewRouter()

	router.Use(gzipMiddleware)
	router.Use(logger.RequestResponseLogger)

	router.Post("/update/{type}/{name}/{value}", sa.updateItemValue)
	router.Get("/value/{type}/{name}", sa.getItemValue)
	router.Post("/update/", sa.update)
	router.Post("/value/", sa.value)
	router.Get("/", sa.getAllValues)
	router.Get("/ping", sa.ping)

	return router
}

func gracefulShutdownCatcher(c *config) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		logger.Log.Debug("shutdown signal caught")
		shutdown(c)
	}()
}

func persistenceTicker(c *config) {
	ticker := time.NewTicker(time.Duration(c.storeInterval) * time.Second)
	defer ticker.Stop()

	logger.Log.Debug("init persistence ticker", zap.Int64("interval", c.storeInterval))

	for range ticker.C {
		err := sa.store(c.fileStoragePath)
		if err != nil {
			logger.Log.Error("persistence ticker error", zap.Error(err), zap.String("path", c.fileStoragePath))
		} else {
			logger.Log.Debug("data stored by ticking timer")
		}
	}
}

func shutdown(c *config) {
	logger.Log.Info("shutting down gracefully, let's save data to disk", zap.String("path", c.fileStoragePath))
	err := sa.store(c.fileStoragePath)
	if err != nil {
		logger.Log.Error("graceful shutdown error", zap.Error(err))
	}
	os.Exit(0)
}

func storingMiddleware(cnf *config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			err := sa.store(cnf.fileStoragePath)
			if err != nil {
				logger.Log.Error(err.Error())
			}
		})
	}
}

var sa *storageAware

func main() {
	// init logging
	serverConf, err := initConfig()
	if err != nil {
		panic(err) // logger has not been initialized yet
	}
	if err := logger.Initialize(serverConf.logLevel); err != nil {
		panic(err) // cannot log without logger
	}

	gracefulShutdownCatcher(&serverConf)

	// init storage
	sa = newStorageAware(storage.NewMemStorage())
	if serverConf.doRestoreValues {
		sa.restore(serverConf.fileStoragePath)
	}

	//inject DB connection
	db, err := database.NewConnection(database.Config{DSN: serverConf.dsn})
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Cannot connect to database '%s': '%s'", serverConf.dsn, err.Error()))
	}
	sa.DB = db

	logger.Log.Info(fmt.Sprintf("Starting server at %s:%d", serverConf.endpoint.host, serverConf.endpoint.port))

	chiMux := newMux(sa)
	if serverConf.storeInterval == 0 {
		logger.Log.Info("will save data to disk immediately")
		chiMux.Use(storingMiddleware(&serverConf))
	} else {
		logger.Log.Info("will save data to disk periodically", zap.Int64("interval", serverConf.storeInterval))
		go persistenceTicker(&serverConf)
	}

	err = http.ListenAndServe(serverConf.endpoint.String(), chiMux)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
}
