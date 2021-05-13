package config

import (
	"os"
)

const EnvProduction = "PRODUCTION"

var (
	Environment = os.Getenv("GYMTRACKER_ENV")
	PortBackend = os.Getenv("GYMTRACKER_PORT_BACKEND")

	GoogleProject = os.Getenv("GYMTRACKER_GOOGLE_PROJECT")

	InfluxURL       = os.Getenv("GYMTRACKER_INFLUX_URL")
	InfluxUser      = os.Getenv("GYMTRACKER_INFLUX_USER")
	InfluxPass      = os.Getenv("GYMTRACKER_INFLUX_PASS")
	InfluxDatabase  = "GymTracker"
	InfluxRetention = "alltime"
)
