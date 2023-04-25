package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	"github.com/oshankkumar/sockshop/internal/app"
	"github.com/oshankkumar/sockshop/internal/pkg/mysql"
	"github.com/oshankkumar/sockshop/transport/http"

	"go.uber.org/zap"
)

type AppConfig struct {
	MySQLConnString string
}

func main() {
	var conf AppConfig
	flag.StringVar(&conf.MySQLConnString, "mysql-conn-str", "admin:password@tcp(mysql:3306)/socksdb", "MySQL connection string")
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

	sockStore := mysql.NewSockStore(db)
	catalogueSvc := &app.CatalogueService{SockStore: sockStore}

	log.Println("starting app http server :9090")

	return http.NewServer(
		catalogueSvc,
		logger,
	).Start(ctx, ":9090")
}
