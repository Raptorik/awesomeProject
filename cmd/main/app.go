package main

import (
	"awesomeProject/internal/config"
	course2 "awesomeProject/internal/course"
	course "awesomeProject/internal/course/db"
	"awesomeProject/internal/lesson"
	"awesomeProject/pkg/client/postrgresql"
	"awesomeProject/pkg/logging"
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("create router")
	router := httprouter.New()

	cfg := config.GetConfig()

	postgreSQLClient, err := postrgresql.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		logger.Fatalf("%v", err)
	}
	repository := course.NewRepository(postgreSQLClient, logger)

	logger.Info("register course handler")
	courseHandler := course2.NewHandler(repository, logger)
	courseHandler.Register(router)

	logger.Info("register create course handler")
	createHandler := course2.NewHandler(repository, logger)
	createHandler.Register(router)

	logger.Info("register update course handler")
	updateHandler := course2.NewHandler(repository, logger)
	updateHandler.Register(router)

	logger.Info("register delete course handler")
	deleteHandler := course2.NewHandler(repository, logger)
	deleteHandler.Register(router)

	logger.Info("register translate lesson name handler")
	lessonHandler := lesson.NewHandler(repository, logger)
	lessonHandler.Register(router)

	start(router, cfg)
}
func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start application")

	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == "sock" {
		logger.Info("detect app path")
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("create socket")
		socketPath := path.Join(appDir, "app.sock")

		logger.Info("listen unix socket")
		listener, listenErr = net.Listen("unix", socketPath)
		logger.Infof("server is listening unix socket: %s", socketPath)
	} else {
		logger.Info("listen tcp socket")
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		logger.Infof("server is listening port %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
	}
	if listenErr != nil {
		logger.Fatal(listenErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Fatal(server.Serve(listener))
}
