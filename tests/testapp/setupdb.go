package testapp

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"path"
	"runtime"

	"github.com/Edgar200021/netowork-server-go/internal/config"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/redis/go-redis/v9"
	tc "github.com/testcontainers/testcontainers-go"
	tcRedis "github.com/testcontainers/testcontainers-go/modules/redis"
)

func setupDb(config *config.PostgresConfig) (*sql.DB, func()) {
	config.Database = fmt.Sprintf("test-%s", uuid.New().String())

	db, err := sql.Open("postgres", config.ConnectionStringWithoutDb())
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(
		fmt.Sprintf(
			"CREATE DATABASE %s",
			pq.QuoteIdentifier(config.Database),
		),
	); err != nil {
		log.Fatal(err)
	}

	connectedDb, err := sql.Open("postgres", config.ConnectionString())
	if err != nil {
		log.Fatal(err)
	}

	if err := connectedDb.Ping(); err != nil {
		log.Fatal(err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	goose.SetLogger(
		log.New(io.Discard, "", 0),
	)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Cannot get current file path")
	}

	if err := goose.Up(connectedDb, path.Join(path.Dir(filename), "../../migrations")); err != nil {
		log.Fatal(err)
	}

	return connectedDb, func() {
		if err := connectedDb.Close(); err != nil {
			log.Printf("failed to close connected database: %s", err)
		}

		if _, err := db.Exec(
			fmt.Sprintf(
				"DROP DATABASE %s",
				pq.QuoteIdentifier(config.Database),
			),
		); err != nil {
			log.Printf("failed to remove database %s\n", err.Error())
		}

		if err := db.Close(); err != nil {
			log.Printf("failed to close database: %s", err)
		}
	}
}

func setupRedis(config *config.RedisConfig) (*redis.Client, func()) {
	ctx := context.Background()

	redisContainer, err := tcRedis.Run(
		ctx, "redis:8", tc.WithLogger(log.New(io.Discard, "", 0)),
	)

	if err != nil {
		log.Fatal(err)
	}

	host, err := redisContainer.Host(ctx)
	if err != nil {
		log.Fatal(err)
	}

	port, err := redisContainer.MappedPort(
		ctx, nat.Port(
			fmt.Sprintf(
				"%d/tcp",
				config.Port,
			),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	config.Host = host
	config.Port = port.Int()

	client := redis.NewClient(config.ConnectOptions())

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatal(err)
	}

	return client, func() {

		if err := client.Close(); err != nil {
			log.Printf("failed to close redis: %s", err)
		}

		if err := tc.TerminateContainer(redisContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}
}
