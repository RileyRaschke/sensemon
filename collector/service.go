package collector

import (
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
	dbc                *db.Connection
	pollInterval       time.Duration
	opts               *CollectorServiceOptions
	quit               chan int
	lastCollectionTime time.Time
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
			if time.Since(s.lastCollectionTime) >= s.pollInterval {
				s.CollectData()
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (s *CollectorService) CollectData() {
	log.Trace("Logging temps...")
	s.lastCollectionTime = time.Now()
	for _, sen := range s.opts.Sensors {
		d, err := sen.GetData()
		//jsonData, _ := json.MarshalIndent(d, "", "  ")
		//fmt.Printf("%s\n", string(jsonData))
		if err != nil {
			log.Errorf("Skipping data insert due to sensor error: %s", err)
			continue
		}
		if d == nil {
			log.Errorf("No data returned from sensor at endpoint: %s", sen.Endpoint)
			continue
		}
		switch d.(type) {
		case *sensor.DhtSensorData:
			err := s.dbc.InsertDhtData(d.(*sensor.DhtSensorData))
			if err != nil {
				log.Errorf("Data insert failed with error: %s", err)
			}
			continue
		default:
			continue
		}
	}
}

func (s *CollectorService) Stop() {
	s.quit <- 1
}
