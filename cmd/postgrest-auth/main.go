package main

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"time"

	"github.com/alexandrevilain/postgrest-auth/pkg/model"

	"github.com/alexandrevilain/postgrest-auth/pkg/api"
	"github.com/alexandrevilain/postgrest-auth/pkg/config"
	"github.com/alexandrevilain/postgrest-auth/pkg/mail"
	"github.com/labstack/gommon/log"

	_ "github.com/lib/pq"
)

func main() {
	logger := log.New("")
	config, err := config.LoadFromEnv()
	if err != nil {
		logger.Fatalf("Unable to load config file: %v", err.Error())
	}

	emailQueue := make(chan mail.EmailSendRequest, 100)
	worker := mail.NewSenderWorker(emailQueue, &config.Email, logger)

	db, err := sql.Open("postgres", config.DB.ConnectionString)
	if err != nil {
		logger.Fatalf("Unable to connect to database: %v", err.Error())
	}
	err = db.Ping()
	if err != nil {
		logger.Fatalf("Unable to connect to database: %v", err.Error())
	}

	err = model.EnsureDBElementsExists(db, &config.DB, logger)
	if err != nil {
		logger.Fatalf("Unable to create base elements on database: %v", err.Error())
	}

	logger.Info("Starting postgrest-auth server ...")
	api.Run(&config, db, emailQueue, logger)

	logger.Info("Stating email worker ...")
	worker.Start()

	// Wait for SIGINT (Ctrl+C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	// Stop API
	api.Stop(ctx)
	// Stop the Worker
	worker.Stop()
	// Close Database connection
	db.Close()
	logger.Info("Shutting down postgrest-auth server ...")
	os.Exit(0)
}
