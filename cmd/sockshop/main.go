package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/oshankkumar/sockshop/api"
	"github.com/oshankkumar/sockshop/api/router"
	"github.com/oshankkumar/sockshop/api/router/catalogue"
	"github.com/oshankkumar/sockshop/api/router/user"
	"github.com/oshankkumar/sockshop/internal/app"
	"github.com/oshankkumar/sockshop/internal/db/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	conf := NewConfigFromFlags()
	if err := mainE(ctx, conf); err != nil {
		log.Fatalf("failed running app: %v", err)
	}

	log.Println("app closed")
}

func mainE(ctx context.Context, conf AppConfig) error {
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

	sockStore := mysql.NewSockStore(db)
	catalogueSvc := app.NewCatalogueService(sockStore)

	userService := &app.UserService{
		UserStore:    mysql.NewUserStore(db),
		CardStore:    mysql.NewCardStore(db),
		AddressStore: mysql.NewAddressStore(db),
		TxBeginner:   db,
		Domain:       conf.Domain,
	}

	routers := router.Routers{
		catalogue.ImageRouter(conf.ImagePath),
		catalogue.NewRouter(catalogueSvc, sockStore),
		user.NewRouter(userService),
	}

	apiServer := &api.Server{
		Addr:          ":9090",
		Logger:        logger,
		HealthChecker: doHealthCheck(db),
		Router:        routers,
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
