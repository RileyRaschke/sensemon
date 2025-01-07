package db

import (
	"fmt"
	"net/url"
)

type ConnectArgs struct {
	Username         string
	Password         string
	PasswordCommand  string
	Server           string
	Port             int
	Database         string
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
	values := url.Values{}
	for key, value := range args.Opts {
		values.Add(key, value.(string))
	}
	urloptions := values.Encode()
	conUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s", args.Username, url.QueryEscape(args.Password), args.Server, args.Port, args.Database, urloptions)
	return conUrl
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
