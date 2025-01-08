package sensor

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sensemon/sensor/sensortypes"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Sensor struct {
	DeviceID   string                 `json:"deviceID" db:"sensor_device_id"`
	Endpoint   string                 `json:"-" db:"sensor_address"`
	SensorType sensortypes.SensorType `json:"type" db:"sensor_type_id"`
	Name       string                 `json:"name" db:"sensor_name"`
}

type SensorData interface {
	GetDeviceId() string
}

type DhtSensorData struct {
	rawData   string
	dataStore map[string]string
	DeviceID  string    `json:"deviceID" db:"sr_device_id"`
	Date      time.Time `json:"ts" db:"sr_date"`
	Farenheit float32   `json:"fahrenheit" db:"sr_farenheit"`
	Humidity  float32   `json:"humidity" db:"sr_humidity"`
}

func SensorsFromViper() []*Sensor {
	c := viper.Get("sensors")
	sensors := make([]*Sensor, len(c.([]interface{})))
	for idx, v := range c.([]interface{}) {
		m := v.(map[string]interface{})
		sensors[idx] = &Sensor{
			Endpoint:   m["endpoint"].(string),
			SensorType: sensortypes.ParseType(m["sensor_type"].(string)),
		}
	}
	return sensors
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
