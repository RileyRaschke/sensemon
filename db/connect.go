package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const DB_TYPE = "pgx"

type Connection struct {
	*sqlx.DB
}

func Connect(args *ConnectArgs) (*Connection, error) {
	sqlx.BindDriver(DB_TYPE, sqlx.NAMED)
	c, err := sqlx.Open("postgres", args.ToConnectionString())
	log.Tracef("%s", args)
	if err != nil {
		panic(fmt.Errorf("error in sql.Open: %w", err))
	}
	db := &Connection{c}
	/*
		err = db.Ping()
		if err != nil {
			log.Fatalf("Error pinging db: %w", err)
		}
	*/
	return db, nil
}

func FromViper(v *viper.Viper) (*Connection, error) {
	return Connect(
		&ConnectArgs{
			Database:        v.GetString("Database"),
			Username:        v.GetString("Username"),
			Password:        v.GetString("Password"),
			PasswordCommand: v.GetString("PasswordCommand"),
			Server:          v.GetString("Server"),
			Port:            v.GetInt("Port"),
			Opts:            v.Get("Options").(map[string]any),
		})
}
