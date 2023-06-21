package sensortypes

type SensorType int

const (
	DHT SensorType = iota
	Unknown
)

var (
	UnknownSensorType = "Unknown Sensor Type"
)

func (s SensorType) String() string {
	switch s {
	case DHT:
		return "Temperature and Humidity"
	default:
		return UnknownSensorType
	}
}

func ValidSensorType(s SensorType) bool {
	if s.String() == UnknownSensorType {
		return false
	}
	return true
}

func ParseType(s string) SensorType {
	switch s {
	case "DHT":
		return DHT
	default:
		return Unknown
	}
}
