package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"

	"github.com/subash68/ate/ate_menu_service/configuration"
	"github.com/subash68/ate/ate_menu_service/pkg/logger"
	"github.com/subash68/ate/ate_menu_service/pkg/protocol/grpc"
	"github.com/subash68/ate/ate_menu_service/pkg/protocol/rest"
	v1 "github.com/subash68/ate/ate_menu_service/pkg/service/v1"
)

type Config struct {
	GRPCPort string
	HTTPPort string

	LogLevel      int
	LogTimeFormat string

	DatabaseHost     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseSchema   string
}

func RunServer() error {
	ctx := context.Background()

	var cfg Config

	configuration.LoadConfig()

	cfg.GRPCPort = strconv.Itoa(configuration.PortConfig().GRPCPort)
	cfg.HTTPPort = strconv.Itoa(configuration.PortConfig().HTTPPort)

	cfg.LogLevel = configuration.LogConfig().LogLevel
	cfg.LogTimeFormat = configuration.LogConfig().LogTimeFormat

	cfg.DatabaseHost = os.Getenv("DB_HOSTNAME") //configuration.DbConfig().DatabaseHost
	cfg.DatabaseUser = configuration.DbConfig().DatabaseUser
	cfg.DatabasePassword = configuration.DbConfig().DatabasePassword
	cfg.DatabaseSchema = configuration.DbConfig().DatabaseSchema

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid tcp port for grpc server: '%s'", cfg.GRPCPort)
	}

	if len(cfg.HTTPPort) == 0 {
		return fmt.Errorf("invalid http port for http gateway: '%s'", cfg.HTTPPort)
	}

	if err := logger.Init(cfg.LogLevel, cfg.LogTimeFormat); err != nil {
		return fmt.Errorf("failed to initialize logger: %v", err)
	}

	param := "parseTime=true"
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseHost,
		cfg.DatabaseSchema,
		param)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("error while connecting with database %v", err)
	}

	defer db.Close()

	v1API := v1.NewMenuServiceServer(db)
	// v1API := v1.NewOnboardServiceServer()

	go func() {
		_ = rest.RunServer(ctx, cfg.GRPCPort, cfg.HTTPPort)
	}()

	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}
