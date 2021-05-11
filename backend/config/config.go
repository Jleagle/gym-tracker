package config

import (
	"os"
)

var (
	Environment = os.Getenv("GYMTRACKER_ENV")
	PortBackend = os.Getenv("GYMTRACKER_PORT_BACKEND")

	GoogleProject = os.Getenv("GYMTRACKER_GOOGLE_PROJECT")

	InfluxURL       = os.Getenv("GYMTRACKER_INFLUX_URL")
	InfluxUser      = os.Getenv("GYMTRACKER_INFLUX_USER")
	InfluxPass      = os.Getenv("GYMTRACKER_INFLUX_PASS")
	InfluxDatabase  = os.Getenv("GYMTRACKER_INFLUX_DATABASE")
	InfluxRetention = os.Getenv("GYMTRACKER_INFLUX_RETENTION")
)
