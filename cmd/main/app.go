package main

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/entities/course"
	"awesomeProject/internal/entities/course/repository"
	"awesomeProject/internal/entities/lesson"
	repository2 "awesomeProject/internal/entities/lesson/repository"
	"awesomeProject/pkg/client/postgresql"
	"awesomeProject/pkg/logging"
	translationlesson "awesomeProject/pkg/translation"
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
	logger, router, cfg := logging.GetLogger(), httprouter.New(), config.GetConfig()

	postgreSQLClient, err := postgresql.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		logger.Fatalf("%v", err)
	}
	translator, err := translationlesson.NewGoogleTranslator(cfg.GoogleCredentialsFile)
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
	var listener net.Listener
	var err error

	if cfg.Listen.Type == "sock" {
		appDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		socketPath := path.Join(appDir, "app.sock")
		listener, err = net.Listen("unix", socketPath)
		logger.Infof("server is listening unix socket: %s", socketPath)
	} else {
		addr := fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
		listener, err = net.Listen("tcp", addr)
		logger.Infof("server is listening port %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
	}

	if err != nil {
		logger.Fatal(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Fatal(server.Serve(listener))
}
