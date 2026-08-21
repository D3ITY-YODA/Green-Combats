package normalization

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUnsupportedSource = errors.New("unsupported source for normalization")
	ErrMalformedPayload  = errors.New("malformed raw payload")
	ErrInvalidValue      = errors.New("observation value out of physical bounds")
)

const (
	VarTemperature      = "temperature_celsius"
	VarPrecipitation    = "precipitation_mm"
	VarRelativeHumidity = "relative_humidity_percent"
	VarWindSpeed        = "wind_speed_ms"
	VarSoilMoisture     = "soil_moisture_m3m3"
	VarSolarRadiation   = "solar_radiation_mj"
	VarWeatherCode      = "weather_code"

	UnitCelsius         = "celsius"
	UnitMillimeter      = "mm"
	UnitPercent         = "percent"
	UnitMetersPerSecond = "m/s"
	UnitM3PerM3         = "m3/m3"
	UnitMegajoulesPerM2 = "mj/m2"
	UnitWMOCode         = "wmo_code"
)

type CanonicalObservation struct {
	ID          uuid.UUID
	RawRecordID uuid.UUID
	SourceID    uuid.UUID
	PlaceID     uuid.UUID
	SourceCode  string
	Variable    string
	Value       float64
	Unit        string
	ObservedAt  time.Time
	Valid       bool
}
