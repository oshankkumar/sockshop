package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/oshankkumar/sockshop/api"
	"github.com/oshankkumar/sockshop/api/handlers"
	"github.com/oshankkumar/sockshop/api/middleware"
	"github.com/oshankkumar/sockshop/api/router"
	"github.com/oshankkumar/sockshop/api/router/catalogue"
	"github.com/oshankkumar/sockshop/api/router/user"
	"github.com/oshankkumar/sockshop/internal/app"
	"github.com/oshankkumar/sockshop/internal/db/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type AppConfig struct {
	MySQLConnString string
	ImagePath       string
	Domain          string
}

func main() {
	var conf AppConfig
	flag.StringVar(&conf.MySQLConnString, "mysql-conn-str", "admin:password@tcp(mysql:3306)/socksdb", "MySQL connection string")
	flag.StringVar(&conf.ImagePath, "image-path", "assets/images", "Image path")
	flag.StringVar(&conf.Domain, "link-domain", "127.0.0.1:9090", "HATEAOS link domain")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := run(ctx, conf); err != nil {
		log.Fatalf("failed running app: %v", err)
	}

	log.Println("app closed")
}

func run(ctx context.Context, conf AppConfig) error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return fmt.Errorf("run logger initialization: %w", err)
	}

	db, err := sqlx.Open("mysql", conf.MySQLConnString)
	if err != nil {
		return fmt.Errorf("db open: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("db ping: %w", err)
	}

	routers := router.Routers{
		HealthCheckRouter(doHealthCheck(db)),
		catalogue.ImageRouter(conf.ImagePath),
	}

	{
		sockStore := mysql.NewSockStore(db)
		catalogueSvc := app.NewCatalogueService(sockStore)
		routers = append(routers, catalogue.NewRouter(catalogueSvc, sockStore))
	}

	{
		userStore := mysql.NewUserStore(db)
		userService := app.NewUserService(userStore, conf.Domain)
		routers = append(routers, user.NewRouter(userService))
	}

	var mux router.Mux = router.NewMux()
	mux = router.NewInstrumentedMux(mux, middleware.WithLog(logger))

	routers.InstallRoutes(mux)

	apiServer := &api.Server{
		Addr:   ":9090",
		Mux:    mux,
		Logger: logger,
	}

	return apiServer.Start(ctx)
}

func doHealthCheck(db *sqlx.DB) api.HealthCheckerFunc {
	return func(ctx context.Context) ([]api.Health, error) {
		if err := db.PingContext(ctx); err != nil {
			return nil, fmt.Errorf("db ping: %w", err)
		}

		var i int
		if err := db.Get(&i, "SELECT 1"); err != nil {
			return nil, fmt.Errorf("db read: %w", err)
		}

		return []api.Health{
			{Service: "sockshop", Status: "OK", Time: time.Now().Local().String()},
			{Service: "sockshop-db", Status: "OK", Time: time.Now().Local().String(), Details: db.Stats()},
		}, nil
	}
}

func HealthCheckRouter(hc api.HealthChecker) router.Router {
	return router.RouterFunc(func(m router.Mux) {
		m.Method(http.MethodGet, "/health", handlers.HealthCheckHandler(hc))
	})
}
