package main

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/entities/course"
	"awesomeProject/internal/entities/course/repository"
	"awesomeProject/internal/entities/lesson"
	repository2 "awesomeProject/internal/entities/lesson/repository"
	"awesomeProject/pkg/client/postrgresql"
	"awesomeProject/pkg/logging"
	translation_lesson "awesomeProject/pkg/translation"
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
	router := httprouter.New()
	cfg := config.GetConfig()

	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", cfg.GoogleCredentialsFile)
	if err != nil {
		return
	}

	postgreSQLClient, err := postrgresql.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		logger.Fatalf("%v", err)
	}
	translator, err := translation_lesson.NewGoogleTranslator(cfg.GoogleCredentialsFile)
	if err != nil {
		logger.Fatalf("failed to initialize translator: %v", err)
	}

	repositoryCourse := repository.NewRepositoryCourse(postgreSQLClient, logger)
	repositoryLesson := repository2.NewRepositoryLesson(postgreSQLClient, logger)

	course.NewHandler(repositoryCourse, logger).Register(router)
	lesson.NewHandler(repositoryLesson, translator, logger).Register(router)

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
