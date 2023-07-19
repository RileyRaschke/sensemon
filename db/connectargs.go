package db

import (
	"fmt"
)

type ConnectArgs struct {
	Driver          string
	Username        string
	Password        string
	PasswordCommand string
	Server          string
	Port            int
	Service         string
	Opts            map[string]interface{}
}

func (args *ConnectArgs) String() string {
	cleanArgs := *args
	cleanArgs.Password = "XXXXXXXX"
	return (&cleanArgs).ToConnectionString()
}

func (args *ConnectArgs) ToConnectionString() string {
	opts := ""
	for key, val := range args.Opts {
		switch val.(type) {
		case int:
			opts = opts + fmt.Sprintf(`%s=%d`, key, val.(int))
			break
		case bool:
			opts = opts + fmt.Sprintf(`%s=%b`, key, val.(string))
			break
		default: // string
			opts = opts + fmt.Sprintf(`%s="%s"`, key, val.(string))
		}
	}

	connStr := fmt.Sprintf(`user="%s" password="%s" connectString="%s" %s`,
		args.Username,
		args.GetPass(),
		args.DSN(),
		opts,
	)
	return connStr
}

func (args *ConnectArgs) DSN() string {
	return fmt.Sprintf("%s:%d/%s",
		args.Server,
		args.Port,
		args.Service,
	)
}

func (args *ConnectArgs) GetPass() string {
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
	return args.Password
}
