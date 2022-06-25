// Mend API
//
// This is example of microservice.
// The example shows CRUD operations with REST-api.
//
// Schemes: https
// Host: localhost:12345
// BasePath: /
// Contact: Piotr Pszczółkowski <piotr@beesoft.pl> http://beesoft.pl
//
// Consumes:
// - application/json
//
// Produces:
// - appplication/json
// swagger:meta
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"mend/db"
	"mend/middleware"
	"mend/routes"
	"mend/src"
)

func main() {
	// u := model.NewUser()
	// u.ID = 1
	// u.FisrtName = "Piotr"
	// u.LastName = "Pszczółkowski"
	// u.Age = 100
	// fmt.Println(u)
	// return

	db.Use(db.InMemory)

	if src.Config().Gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	if fd, err := logFileHandle(); err != nil {
		log.Error(err)
	} else {
		// adds middleware to save log to a file
		logger := log.New()
		logger.Out = fd
		logger.SetLevel(log.DebugLevel)
		logger.SetFormatter(&log.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"})
		router.Use(middleware.LogToFile(logger))
	}
	// add middleware for recovery from panics
	router.Use(gin.Recovery())

	routes.Setup(router)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		if err := router.RunTLS(":12345", "/Users/piotr/.mend/cert/localhost.crt", "/Users/piotr/.mend/cert/localhost.key"); err != nil {
			log.Error(err)
			cancel()
		}
	}()

	<-ctx.Done()
	cancel()
	log.Info("shutdown the server....")
}

func logFileHandle() (*os.File, error) {
	fpath, err := src.Config().LogFilePath()
	if err != nil {
		return nil, err
	}

	return os.OpenFile(fpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
}
