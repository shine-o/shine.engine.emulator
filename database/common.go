package database

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v9"
)

type ConnectionParams struct {
	User string
	Password string
	Host string
	Port string
	Database string
	Schema string
}

func CreateSchema(db *pg.DB, schema string) error {
	_, err := db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %v;", schema))
	return err
}

func Connection(ctx context.Context, cp ConnectionParams) *pg.DB {
	return pg.Connect(&pg.Options{
		Addr:            fmt.Sprintf("%v:%v", cp.Host, cp.Port),
		User:            cp.User,
		Password:        cp.Password,
		Database:        cp.Database,
		ApplicationName: cp.Schema,
		TLSConfig:       nil,
		PoolSize:    5,
		PoolTimeout: 5,
	}).WithParam(cp.Schema, pg.Safe(cp.Schema)).WithContext(ctx)
}

