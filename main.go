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
	"net"
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

// @title mend microservice example
// @version 1.0
// @host localhost:8010
// @contact.email piotr@beesoft.pl
func main() {
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
		dir, _ := os.Getwd()
		address := ipAddress() + ":8010"
		log.Infof("Server starts on %v", address)

		if err := router.RunTLS(address, dir+"/cert/localhost.crt", dir+"/cert/localhost.key"); err != nil {
			log.Error(err)
			cancel()
		}
	}()

	<-ctx.Done()
	cancel()
	log.Info("shutdown the server....")
}

func ipAddress() string {
	hostname, _ := os.Hostname()
	addresses, _ := net.LookupIP(hostname)
	for _, ip := range addresses {
		if ip.To4() != nil && !ip.IsLoopback() {
			return ip.String()
		}
	}
	return ""
}

func logFileHandle() (*os.File, error) {
	fpath, err := src.Config().LogFilePath()
	if err != nil {
		return nil, err
	}

	return os.OpenFile(fpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
}
