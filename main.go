package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"

	"sensemon/api"
	"sensemon/collector"
	database "sensemon/db"
	"sensemon/sensor"
	"sensemon/sensor/sensortypes"

	log "github.com/sirupsen/logrus"
)

var (
	Version string
	signals chan os.Signal
	dbc     *database.Connection
	router  *chi.Mux
)

func init() {
	signals = make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
}

func main() {
	log.Infof("Starting %s", Version)

	dbc, err := dbConnect()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := dbc.Close()
		if err != nil {
			panic(fmt.Errorf("Failed to close connections: %w", err))
		}
	}()

	// Build and Start the collection service
	collectorService := collector.NewCollectorService(dbc,
		&collector.CollectorServiceOptions{
			PollingInverval: "5s",
			Sensors: []*sensor.Sensor{
				&sensor.Sensor{Endpoint: "http://10.1.1.50/", SensorType: sensortypes.ParseType("DHT")},
			},
		},
	)

	go collectorService.Run()

	// Create routers
	router = chi.NewRouter()

	// init middleware
	router.Use(middleware.RealIP)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.GetHead)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// bind routes
	router.Get("/favicon.ico", noData)
	router.Mount("/api/", api.NewApiController(dbc).Handler())
	router.Mount("/", StaticServerChroot())

	// Run the server
	port := ":" + viper.GetString("app.port")
	srv := &http.Server{Addr: port, Handler: router}
	go func() {
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed && err != nil {
			panic(err)
		}
	}()
	log.Infof("Server listening on port %s", port)

	// block until a signal is recieved from the OS
	i := <-signals

	log.Infof("Received signal: %v, exiting...\n", i)
	err = srv.Shutdown(context.Background())
	if err != nil {
		panic(err)
	}
	srv.Close()
	collectorService.Stop()
}

func noData(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "")
}
func root(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/app/", 302)
}

func dbConnect() (*database.Connection, error) {
	return database.Connect(
		&database.ConnectArgs{
			DbType:          viper.GetString("db.Type"),
			Username:        viper.GetString("db.Username"),
			Password:        viper.GetString("db.Password"),
			PasswordCommand: viper.GetString("db.PasswordCommand"),
			WalletLocation:  viper.GetString("db.WalletLocation"),
			Server:          viper.GetString("db.Server"),
			Port:            viper.GetInt("db.Port"),
			Service:         viper.GetString("db.Service"),
			SID:             viper.GetString("db.SID"),
			Opts:            viper.Get("db.Options").(map[string]any),
		})
}
