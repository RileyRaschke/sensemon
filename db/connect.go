package db

import (
	_ "github.com/godror/godror"
	"github.com/jmoiron/sqlx"
	_ "github.com/sijms/go-ora/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Connection struct {
	*sqlx.DB
}

func Connect(args *ConnectArgs) (*Connection, error) {
	sqlx.BindDriver(args.Driver, sqlx.NAMED)
	c, err := sqlx.Open(args.Driver, args.ToConnectionString())
	if err != nil {
		log.Fatalf("error in sql.Open: %v", err)
	}
	db := &Connection{c}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging db: %v", err)
	}
	return db, nil
}

func FromViper(v *viper.Viper) (*Connection, error) {
	return Connect(
		&ConnectArgs{
			Driver:          v.GetString("Driver"),
			Username:        v.GetString("Username"),
			Password:        v.GetString("Password"),
			PasswordCommand: v.GetString("PasswordCommand"),
			Server:          v.GetString("Server"),
			Port:            v.GetInt("Port"),
			Service:         v.GetString("Service"),
			Opts:            v.Get("Options").(map[string]any),
		})
}
