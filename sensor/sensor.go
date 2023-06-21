package sensor

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sensemon/sensor/sensortypes"
	"time"

	log "github.com/sirupsen/logrus"
)

type Sensor struct {
	Endpoint   string                 `mapstructure:"endpoint"`
	SensorType sensortypes.SensorType `mapstructure:"type"`
}

type SensorData interface {
	GetDeviceId() string
}

type DhtSensorData struct {
	rawData   string
	dataStore map[string]string
	DeviceID  string    `json:"DeviceID" db:"SR_DEVICE_ID"`
	Date      time.Time `json:"ts" db:"SR_DATE"`
	Farenheit float32   `json:"Fahrenheit" db:"SR_FARENHEIT"`
	Humidity  float32   `json:"Humidity" db:"SR_HUMIDITY"`
}

func (s *Sensor) GetData() (SensorData, error) {
	switch s.SensorType {
	case sensortypes.DHT:
		d := &DhtSensorData{}
		resp, err := http.Get(s.Endpoint)
		if err != nil {
			log.Errorf("Failed to refresh sensor data with error: %s", err)
			return nil, err
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), d)
		d.Date = time.Now()
		return d, nil
	default:
		return nil, errors.New("Unknown sensor type configured")
	}
}

func (sd *DhtSensorData) GetDeviceId() string {
	return sd.DeviceID
}
