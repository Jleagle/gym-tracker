package config

import (
	"os"
)

var (
	Environment = os.Getenv("PURE_ENV")
	PortBackend = os.Getenv("PURE_PORT_BACKEND")

	User = os.Getenv("PURE_USER")
	Pass = os.Getenv("PURE_PASS")

	InfluxURL       = os.Getenv("PURE_INFLUX_URL")
	InfluxUser      = os.Getenv("PURE_INFLUX_USER")
	InfluxPass      = os.Getenv("PURE_INFLUX_PASS")
	InfluxDatabase  = os.Getenv("PURE_INFLUX_DATABASE")
	InfluxRetention = os.Getenv("PURE_INFLUX_RETENTION")
)
