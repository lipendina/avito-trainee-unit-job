package db

import (
	"../config"
	"context"
	"fmt"
	"github.com/jackc/pgx/pgxpool"
	"log"
)

type PgxSource interface {
	NewConnect()
}

type ConnDB struct {
	DB *pgxpool.Pool
	ctx context.Context
}

func NewConnectToPG(dbConfig *config.DBConfig, ctx context.Context) ConnDB {
	poolConfig, err := pgxpool.ParseConfig(fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName))
	if err != nil {
		log.Fatalf("Cannot parse config", err)
	}
	poolConfig.ConnConfig.RuntimeParams["standard_conforming_strings"] = "on";
	poolConfig.ConnConfig.PreferSimpleProtocol = true

	db, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalf("Unable to create connection pool: ", err)
	}

	return ConnDB{
		DB: db,
		ctx: ctx,
	}
}
