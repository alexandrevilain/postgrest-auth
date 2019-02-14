package api

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/alexandrevilain/postgrest-auth/pkg/config"
	"github.com/alexandrevilain/postgrest-auth/pkg/mail"
)

var server *echo.Echo

// Run starts the API server
func Run(config *config.Config, db *sql.DB, emailQueue chan mail.EmailSendRequest, logger *log.Logger) {
	server = echo.New()
	server.HideBanner = true
	server.Logger = logger
	server.Use(middleware.Recover())
	server.Use(middleware.Logger())
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	h := handler{
		db:         db,
		config:     config,
		emailQueue: emailQueue,
		emails:     mail.NewEmailGenerator(&config.App),
	}

	server.POST("/signin", h.signin)
	server.POST("/signup", h.signup)
	server.GET("/confirm/:id", h.confirmAccount)
	server.POST("/reset", h.sendPasswordReset)
	server.POST("/reset/:token", h.resetPassword)
	server.POST("/provider/:provider", h.signinWithProvider)

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		listen := fmt.Sprintf("0.0.0.0:%v", config.API.Port)
		if err := server.Start(listen); err != nil {
			logger.Error(err)
		}
	}()
}

// Stop stops the API Server
func Stop(ctx context.Context) {
	server.Shutdown(ctx)
}
