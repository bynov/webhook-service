package config

import (
	"fmt"
	"os"
)

const (
	envKeyMode         = "MODE"
	envKeySlaveAddr    = "SLAVE_ADDR"
	envKeyDatabaseAddr = "DATABASE_ADDR"
)

const (
	ModeMaster = "master"
	ModeSlave  = "slave"
)

type Config struct {
	Mode         string
	SlaveAddr    string
	DatabaseAddr string
}

func (c Config) IsMaster() bool {
	return c.Mode == ModeMaster
}

// Parse is a simple function that parses config from envs.
func Parse() (*Config, error) {
	var cfg = &Config{
		Mode:         os.Getenv(envKeyMode),
		DatabaseAddr: os.Getenv(envKeyDatabaseAddr),
		SlaveAddr:    os.Getenv(envKeySlaveAddr),
	}

	switch cfg.Mode {
	case ModeMaster:
		if cfg.SlaveAddr == "" {
			return nil, fmt.Errorf("env key %q is empty or not provided", envKeySlaveAddr)
		}

		return cfg, nil
	case ModeSlave:
		return cfg, nil
	default:
		return nil, fmt.Errorf("invalid value in env key %q, got: %s", envKeyMode, cfg.Mode)
	}
}
