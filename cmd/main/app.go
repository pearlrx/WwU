package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"techpoint/internal/config"
	"techpoint/internal/user"
	"techpoint/internal/user/db"
	"techpoint/pkg/client/mongodb"
	"techpoint/pkg/logging"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("create router")
	router := httprouter.New()

	cfg := config.GetConfig()
	cfgMongo := cfg.MongoDB

	mongoDBClient, err := mongodb.NewClient(context.Background(), cfgMongo.Host, cfgMongo.Port, cfgMongo.Username,
		cfgMongo.Password, cfgMongo.Database, cfgMongo.AuthDB)
	if err != nil {
		panic(err)
	}
	storage := db.NewStorage(mongoDBClient, cfg.MongoDB.Collection, logger)
	user1 := user.User{
		ID:           "",
		Email:        "pearl@gmail.com",
		Username:     "pearl",
		PasswordHash: "12345",
	}
	userID1, err := storage.Create(context.Background(), user1)
	if err != nil {
		panic(err)
	}
	logger.Info(userID1)

	logger.Info("register user handler")
	handler := user.NewHandler()
	handler.Register(router)

	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start application")

	var listener net.Listener
	var ListenError error

	if cfg.Listen.Type == "sock" {
		logger.Info("detect app path")
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("create socket")
		socketPath := path.Join(appDir, "app.sock")

		logger.Info("listen unix socket")
		listener, err = net.Listen("unix", socketPath)
		logger.Infof("server is listening unix socket: %s", socketPath)
		if err != nil {
			logger.Fatal(err)
		}
	} else {
		logger.Info("listen tcp")
		listener, ListenError = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		logger.Infof("Server listen and serve on port %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
	}

	if ListenError != nil {
		logger.Fatal(ListenError)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	logger.Fatal(server.Serve(listener))
}
