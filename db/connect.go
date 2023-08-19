package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/sijms/go-ora/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const DB_TYPE = "oracle"

type Connection struct {
	*sqlx.DB
}

func Connect(args *ConnectArgs) (*Connection, error) {
	sqlx.BindDriver(DB_TYPE, sqlx.NAMED)
	c, err := sqlx.Open(DB_TYPE, args.ToConnectionString())
	log.Tracef("%s", args)
	if err != nil {
		panic(fmt.Errorf("error in sql.Open: %w", err))
	}
	db := &Connection{c}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging db: %w", err)
	}
	return db, nil
}

func FromViper(v *viper.Viper) (*Connection, error) {
	return Connect(
		&ConnectArgs{
			Username:        v.GetString("Username"),
			Password:        v.GetString("Password"),
			PasswordCommand: v.GetString("PasswordCommand"),
			Server:          v.GetString("Server"),
			Port:            v.GetInt("Port"),
			Service:         v.GetString("Service"),
			Opts:            v.Get("Options").(map[string]any),
		})
}
