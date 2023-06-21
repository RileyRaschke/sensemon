package collector

import (
	"encoding/json"
	"fmt"
	"sensemon/db"
	"time"

	"sensemon/sensor"

	log "github.com/sirupsen/logrus"
)

type CollectorServiceOptions struct {
	PollingInverval string
	Sensors         []*sensor.Sensor
}

type CollectorService struct {
	dbc            *db.Connection
	pollInterval   time.Duration
	opts           *CollectorServiceOptions
	quit           chan int
	lastCollection time.Time
}

func NewCollectorService(db *db.Connection, opts *CollectorServiceOptions) *CollectorService {
	d, err := time.ParseDuration(opts.PollingInverval)
	if err != nil {
		panic(err)
	}
	return &CollectorService{dbc: db, pollInterval: d, opts: opts, quit: make(chan int)}
}

func (s *CollectorService) Run() {
	for {
		select {
		case <-s.quit:
			return
		default:
			if time.Since(s.lastCollection) >= s.pollInterval {
				s.CollectData()
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (s *CollectorService) CollectData() {
	log.Info("Logging temps...")
	s.lastCollection = time.Now()
	for _, sensor := range s.opts.Sensors {
		d := sensor.GetData()
		jsonData, _ := json.MarshalIndent(d, "", "  ")
		fmt.Printf("%s\n", string(jsonData))
	}
}

func (s *CollectorService) Stop() {
	s.quit <- 1
}
