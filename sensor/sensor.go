package sensor

import (
	"encoding/json"
	"io"
	"net/http"
	st "sensemon/sensor/sensortype"

	log "github.com/sirupsen/logrus"
)

type Sensor struct {
	Endpoint   string        `mapstructure:"endpoint"`
	SensorType st.SensorType `mapstructure:"type"`
}

type SensorData interface {
	GetDeviceId() string
}

type DhtSensorData struct {
	rawData   string
	dataStore map[string]string
	DeviceID  string `json:"DeviceID", db:"SR_DEVICE_ID"`
	Farenheit string `json:"Fahrenheit", db:"SR_FARENHEIT"`
	Humidity  string `json:"Humidity", db"SR_HUMIDITY"`
}

func (s *Sensor) GetData() SensorData {
	switch s.SensorType {
	case st.DHT:
		d := &DhtSensorData{}
		resp, err := http.Get(s.Endpoint)
		if err != nil {
			log.Errorf("Failed to refresh sensor data with error: %s", err)
			return nil
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), d)
		return d
	default:
		return nil
	}
}

func (sd *DhtSensorData) GetDeviceId() string {
	return sd.DeviceID
}
