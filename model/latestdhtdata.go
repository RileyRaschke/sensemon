package model

import "time"

type LatestDhtSensorData struct {
	SensorName    string    `db:"SENSOR_NAME" json:"sensor"`
	Fahrenheit    string    `db:"FAHRENHEIT" json:"fahrenheit"`
	Humidity      string    `db:"HUMIDITY" json:"humidity"`
	LastEntryDate time.Time `db:"LAST_ENTRY_DATE" json:"as_of"`
}
