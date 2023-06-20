package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/sijms/go-ora/v2"
	log "github.com/sirupsen/logrus"
)

type Connection struct {
	*sqlx.DB
}

func Connect(args *ConnectArgs) (*Connection, error) {
	var err error
	if args.Password == "" {
		if args.PasswordCommand != "" {
			args.Password, err = passwordFromCommand(args.PasswordCommand)
			if err != nil {
				return nil, err
			}
		}
	}
	if args.Password == "" {
		args.Password, err = passwordFromShell()
		if err != nil {
			return nil, err
		}
	}

	c, err := sqlx.Open(args.DbType, args.ToConnectionString())
	if err != nil {
		panic(fmt.Errorf("error in sql.Open: %s", err))
	}
	db := &Connection{c}
	err = db.Ping()
	if err != nil {
		//panic(fmt.Errorf("error pinging db: %s", err))
		log.Fatalf("Error pinging db: %s", err)
	}
	return db, nil
}
