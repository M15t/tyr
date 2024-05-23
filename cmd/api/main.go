package main

import (
	"embed"
	"log/slog"
	"os"

	"tyr/config"

	"tyr/internal/api/root"
	"tyr/internal/api/v1/app/document"
	"tyr/internal/api/v1/auth"
	"tyr/internal/db"
	"tyr/internal/rbac"
	"tyr/internal/repo"
	"tyr/third_party/azure"

	"github.com/M15t/gram/pkg/server"
	"github.com/M15t/gram/pkg/server/middleware/jwt"
	"github.com/M15t/gram/pkg/server/middleware/secure"
	"github.com/M15t/gram/pkg/server/middleware/slogger"
	"github.com/M15t/gram/pkg/util/crypter"
	"github.com/M15t/gram/pkg/util/prettylog"

	"github.com/labstack/echo/v4"

	contextutil "tyr/internal/api/context"
)

// To embed SwaggerUI into api server using go:build tag
var (
	enableSwagger = false
	swaggerui     embed.FS
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	cfg, err := config.LoadAll()
	checkErr(err)

	db, sqldb, err := db.New(cfg.DB)
	checkErr(err)
	defer sqldb.Close()

	// Initialize HTTP server
	e := server.New(&server.Config{
		Port:              cfg.Server.Port,
		ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
		ReadTimeout:       cfg.Server.ReadTimeout,
		WriteTimeout:      cfg.Server.WriteTimeout,
		AllowOrigins:      cfg.Server.AllowOrigins,
		Debug:             cfg.General.Debug,
	})

	// Create a slog logger, which:
	//   - Logs to stdout.
	var logger *slog.Logger
	// formatLog = prettylog.JSONFormat
	logger = slog.New(prettylog.NewHandler(nil, prettylog.TextFormat))
	filters := make([]slogger.Filter, 0)
	filters = append(filters, slogger.IgnorePathContains("swagger"))

	loggerConfig := slogger.Config{
		Filters: filters,
	}

	if config.IsLambda() {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

		loggerConfig.WithUserAgent = true
	}

	e.Use(slogger.NewWithConfig(logger, loggerConfig))

	if enableSwagger {
		// Static page for SwaggerUI
		e.GET("/swagger-ui/*", echo.StaticDirectoryHandler(echo.MustSubFS(swaggerui, "swagger-ui"), false), secure.DisableCache())
	}

	// Initialize core services
	crypterSvc := crypter.New()
	repoSvc := repo.New(db)
	rbacSvc := rbac.New(cfg.General.Debug)
	jwtSvc := jwt.New(cfg.JWT.Algorithm, cfg.JWT.Secret, cfg.JWT.DurationAccessToken, cfg.JWT.DurationRefreshToken)

	azureSvc := azure.New(cfg.Azure, repoSvc)

	// Initialize services
	authSvc := auth.New(repoSvc, jwtSvc, crypterSvc)
	// sessionSvc := session.New(repoSvc, rbacSvc)
	// userSvc := user.New(repoSvc, rbacSvc, crypterSvc)

	documentSvc := document.New(repoSvc, rbacSvc, crypterSvc, azureSvc)

	// Initialize root API
	root.NewHTTP(e)

	v1router := e.Group("/v1")

	auth.NewHTTP(authSvc, v1router.Group("/auth"))

	// Initialize admin APIs
	// v1adminRouter := v1router.Group("/admin")
	v1appRouter := v1router.Group("/app")
	// v1adminRouter.Use(jwtSvc.MWFunc(), contextutil.MWContext())
	// session.NewHTTP(sessionSvc, v1adminRouter.Group("/sessions"))
	// user.NewHTTP(userSvc, v1adminRouter.Group("/users"))

	v1appRouter.Use(jwtSvc.MWFunc(), contextutil.MWContext())
	document.NewHTTP(documentSvc, v1appRouter.Group("/documents"))

	server.Start(e, config.IsLambda())
}
