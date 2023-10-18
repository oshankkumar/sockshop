package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/oshankkumar/sockshop/api"
	"github.com/oshankkumar/sockshop/internal/app"
	"github.com/oshankkumar/sockshop/internal/db/mysql"
	"github.com/oshankkumar/sockshop/internal/transport/http"

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

	routers := http.Routers{
		http.HealthCheckRouter(doHealthCheck(db)),
		http.ImageServeRouter(conf.ImagePath),
	}

	{
		sockStore := mysql.NewSockStore(db)
		catalogueSvc := app.NewCatalogueService(sockStore)
		routers = append(routers, &http.CatalogueRoutes{SockLister: catalogueSvc, SockStore: sockStore})
	}

	{
		userStore := mysql.NewUserStore(db)
		userService := app.NewUserService(userStore, conf.Domain)
		routers = append(routers, &http.UserRoutes{UserService: userService})
	}

	var mux http.Mux = chi.NewMux()
	mux = http.NewInstrumentedMux(mux, http.LogginMiddleware(logger))

	routers.InstallRoutes(mux)

	apiServer := &http.APIServer{
		Addr:   ":9090",
		Mux:    mux,
		Logger: logger,
	}

	return apiServer.Start(ctx)
}

func doHealthCheck(db *sqlx.DB) http.HealthCheckerFunc {
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
