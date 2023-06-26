package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/pborman/getopt/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	me                = filepath.Base(os.Args[0])
	yamlFile          = fmt.Sprintf("%s.yaml", me)
	envPrefix         = "SENSEMON"
	configSearchPaths = []string{".", "./etc", "$HOME/etc", "/etc"}
	genConfig         = getopt.BoolLong("genconfig", 'x', "Write example config to \"./"+yamlFile+"\"")
	webOnly           = getopt.BoolLong("webOnly", 0x00, "Disable ecollector service")
)

func init() {
	os_user, err := user.Current()
	if err != nil {
		panic(err)
	}
	viper.SetConfigName(yamlFile)
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	for _, p := range configSearchPaths {
		viper.AddConfigPath(p)
	}
	viper.SetDefault("log.level", "TRACE")

	viper.SetDefault("app.port", "3001")
	viper.SetDefault("app.proxyHostsCidr", []string{})
	viper.SetDefault("app.realIpHeader", "X-Forwarded-For")

	viper.SetDefault("collector.polling_interval", "5s")

	viper.SetDefault("db.Type", "oracle")
	viper.SetDefault("db.Username", strings.ToUpper(os_user.Username))
	viper.SetDefault("db.Password", "")
	viper.SetDefault("db.PasswordCommand", "")
	viper.SetDefault("db.WalletLocation", "")
	viper.SetDefault("db.Server", "localhost")
	viper.SetDefault("db.Port", "1521")
	viper.SetDefault("db.Service", "")
	viper.SetDefault("db.SID", "")
	viper.SetDefault("db.Options", make(map[string]string))

	getopt.SetUsage(func() { usage() })
	getopt.Parse()

	if *genConfig {
		ConfigWrite()
		os.Exit(0)
		return
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Errorf("%v\n", err)
			usage(fmt.Sprintf("\nTry: %s --genconfig\n", me))
			os.Exit(0)
		} else {
			log.Fatalf("Failed to parse config: %v\n", err)
			os.Exit(0)
		}
	}

	initLogger(viper.GetString("log.level"))
}

func usage(msg ...string) {
	if len(msg) > 0 {
		fmt.Fprintf(os.Stderr, "%s\n", msg[0])
	}
	// strip off the first line of generated usage
	b := &strings.Builder{}
	getopt.PrintUsage(b)
	u := strings.SplitAfterN(b.String(), "\n", 2)
	fmt.Printf(`Usage: %s

OPTIONS
%s
`, me, u[1])

	//os.Exit(1)
}

func ConfigWrite() {
	viper.SafeWriteConfigAs(fmt.Sprintf("./%s", yamlFile))
	log.Printf("Wrote example config to: \"./%s\", feel free to move to: %v", yamlFile, configSearchPaths[1:])
}
