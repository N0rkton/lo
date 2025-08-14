package config

import (
	"log"
	"net"
	"os"
)

type Config struct {
	HostPort string
}

func NewConfig() *Config {
	hostPort := os.Getenv("HOST_PORT")
	if hostPort != "" {
		if err := validateHostPort(hostPort); err != nil {
			log.Println("invalid host port")
		} else {
			return &Config{
				HostPort: hostPort,
			}
		}
	}

	return &Config{
		HostPort: ":8080",
	}
}

func validateHostPort(hostPort string) error {
	_, _, err := net.SplitHostPort(hostPort)
	return err
}
