package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strconv"
	"syscall"

	"github.com/Hymiside/hezzl-api/pkg/handler"
	"github.com/Hymiside/hezzl-api/pkg/models"
	"github.com/Hymiside/hezzl-api/pkg/queue"
	"github.com/Hymiside/hezzl-api/pkg/repository/clickhouse"
	"github.com/Hymiside/hezzl-api/pkg/repository/postgres"
	"github.com/Hymiside/hezzl-api/pkg/repository/redis"
	"github.com/Hymiside/hezzl-api/pkg/server"
	"github.com/Hymiside/hezzl-api/pkg/service"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			_, filename := path.Split(f.File)
			filename = fmt.Sprintf("%s:%d", filename, f.Line)
			return "", filename
		},
	})

	if err := godotenv.Load(); err != nil {
		log.Panicf("error to load .env file: %v", err)
	}

	dbPostgres, err := postgres.NewPostgresDB(ctx, models.ConfigPostgres{
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Database: os.Getenv("POSTGRES_DB"),
	})
	if err != nil {
		log.Fatalf("error to connect postgres: %v", err)
	}

	dbClickhouse, err := clickhouse.NewClickhouseDB(ctx, models.ConfigClickhouse{
		Host:     os.Getenv("CLICKHOUSE_HOST"),
		Port:     os.Getenv("CLICKHOUSE_PORT"),
		Database: os.Getenv("CLICKHOUSE_DATABASE"),
	})
	if err != nil {
		log.Fatalf("error to connect clickhouse: %v", err)
	}

	rdb, err := redis.NewRedisDB(ctx, models.ConfigRedis{
		Host: os.Getenv("REDIS_HOST"),
		Port: os.Getenv("REDIS_PORT"),
	})
	if err != nil {
		log.Fatalf("error to connect redis: %v", err)
	}

	qu, err := queue.NewNats(ctx, models.ConfigNats{
		Host: os.Getenv("NATS_HOST"),
		Port: os.Getenv("NATS_PORT"),
	})
	if err != nil {
		log.Fatalf("error to connect nats: %v", err)
	}

	numOfLogs, err := strconv.Atoi(os.Getenv("NUM_OF_LOGS"))
	if err != nil {
		numOfLogs = 25
	}

	repoPostgres := postgres.NewRepositoryPostgres(dbPostgres)
	repoClickhouse := clickhouse.NewRepositoryClickhouse(dbClickhouse)
	repoRedis := redis.NewRepositoryRedis(rdb)
	quNats := queue.NewQueue(qu, repoClickhouse, numOfLogs)
	services := service.NewService(repoPostgres, repoRedis, quNats)
	handlers := handler.NewHandler(services)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
		select {
		case <-quit:
			cancel()
		case <-ctx.Done():
			return
		}
	}()

	if err = server.StartServer(ctx, handlers.NewRoutes(), models.ConfigServer{
		Host: os.Getenv("SERVER_HOST"),
		Port: os.Getenv("SERVER_PORT"),
	}); err != nil {
		log.Fatalf("error to start server: %v", err)
	}
}
