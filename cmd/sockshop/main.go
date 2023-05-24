package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	gohttp "net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/oshankkumar/sockshop/api"
	"github.com/oshankkumar/sockshop/internal/app"
	"github.com/oshankkumar/sockshop/internal/db/mysql"
	"github.com/oshankkumar/sockshop/transport/http"

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

	var catalogueRouter http.Router
	{
		sockStore := mysql.NewSockStore(db)
		catalogueSvc := app.NewCatalogueService(sockStore)
		catalogueRouter = &http.CatalogueRouter{SockLister: catalogueSvc, SockStore: sockStore}
	}

	var userRouter http.Router
	{
		userStore := mysql.NewUserStore(db)
		userService := app.NewUserService(userStore, conf.Domain)
		userRouter = &http.UserRouter{UserService: userService}
	}

	apiServer := &http.APIServer{
		ImagePath:     conf.ImagePath,
		HealthChecker: doHealthCheck(db),
		Middleware:    http.ChainMiddleware(http.WithLogging(logger)),
	}

	serveMux := apiServer.CreateMux(userRouter, catalogueRouter)

	return startHTTPServer(ctx, ":9090", serveMux, logger)
}

func startHTTPServer(ctx context.Context, addr string, handler gohttp.Handler, logger *zap.Logger) error {
	srv := &gohttp.Server{Addr: ":9090", Handler: handler}
	errc := make(chan error, 1)

	logger.Info("starting app http server :9090")

	go func(errc chan<- error) {
		if err := srv.ListenAndServe(); err != nil && err != gohttp.ErrServerClosed {
			errc <- fmt.Errorf("http server shutdonwn: %w", err)
		}
	}(errc)

	shutdown := func(timeout time.Duration) error {
		logger.Info("received context cancellation; shutting down server")

		shutCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := srv.Shutdown(shutCtx); err != nil {
			return fmt.Errorf("http server shutdonwn: %w", err)
		}
		return nil
	}

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
		return shutdown(time.Second * 5)
	}
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
