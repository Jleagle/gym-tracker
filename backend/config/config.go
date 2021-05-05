package config

import (
	"os"
)

var (
	Environment = os.Getenv("PURE_ENV")
	PortBackend = os.Getenv("PURE_PORT_BACKEND")

	GoogleProject = os.Getenv("PURE_GOOGLE_PROJECT")

	InfluxURL       = os.Getenv("PURE_INFLUX_URL")
	InfluxUser      = os.Getenv("PURE_INFLUX_USER")
	InfluxPass      = os.Getenv("PURE_INFLUX_PASS")
	InfluxDatabase  = os.Getenv("PURE_INFLUX_DATABASE")
	InfluxRetention = os.Getenv("PURE_INFLUX_RETENTION")
)
