package db

import (
	go_ora "github.com/sijms/go-ora/v2"
)

type ConnectArgs struct {
	Username         string
	Password         string
	PasswordCommand  string
	Server           string
	Port             int
	Service          string
	ConnectionString string
	Opts             map[string]interface{}
}

func (args *ConnectArgs) String() string {
	cleanArgs := *args
	cleanArgs.Password = "XXXXXXXX"
	return (&cleanArgs).ToConnectionString()
}

func (args *ConnectArgs) ToConnectionString() string {
	args.GetPass()
	urloptions := make(map[string]string)
	for key, val := range args.Opts {
		urloptions[key] = val.(string)
	}
	url := go_ora.BuildUrl(args.Server, args.Port, args.Service, args.Username, args.Password, urloptions)
	return url
}

func (args *ConnectArgs) GetPass() {
	var err error
	if args.Password == "" {
		if args.PasswordCommand != "" {
			args.Password, err = passwordFromCommand(args.PasswordCommand)
			if err != nil {
				panic(err)
			}
		}
	}
	if args.Password == "" {
		args.Password, err = passwordFromShell()
		if err != nil {
			panic(err)
		}
	}
}
