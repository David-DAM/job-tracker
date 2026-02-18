package bootstrap

import (
	"log"

	"github.com/grafana/pyroscope-go"
)

func InitProfile(config *Config) (*pyroscope.Profiler, error) {

	pyroscopeConfig := pyroscope.Config{
		ApplicationName: config.AppName,
		ServerAddress:   "http://localhost:9999",
		ProfileTypes:    pyroscope.DefaultProfileTypes,
	}

	profiler, err := pyroscope.Start(pyroscopeConfig)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return profiler, nil
}
