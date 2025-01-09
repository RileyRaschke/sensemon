package model

import "time"

type LatestDhtSensorData struct {
	SensorName    string    `db:"sensor_name" json:"sensor"`
	Fahrenheit    string    `db:"fahrenheit" json:"fahrenheit"`
	Humidity      string    `db:"humidity" json:"humidity"`
	LastEntryDate time.Time `db:"last_entry_date" json:"as_of"`
}
