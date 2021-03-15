package config

import (
	"fmt"
	"os"
)

const (
	envKeyMode      = "MODE"
	envKeySlaveAddr = "SLAVE_ADDR"
)

const (
	ModeMaster = "master"
	ModeSlave  = "slave"
)

type Config struct {
	Mode      string
	SlaveAddr string
}

func (c Config) IsMaster() bool {
	return c.Mode == ModeMaster
}

func Parse() (*Config, error) {
	mode := os.Getenv(envKeyMode)

	switch mode {
	case ModeMaster:
		return &Config{
			Mode: mode,
		}, nil
	case ModeSlave:
		slaveAddr := os.Getenv(envKeySlaveAddr)

		if slaveAddr == "" {
			return nil, fmt.Errorf("env key %q is empty or not provided", envKeySlaveAddr)
		}

		return &Config{
			Mode:      mode,
			SlaveAddr: slaveAddr,
		}, nil
	default:
		return nil, fmt.Errorf("invalid value in env key %q, got: %s", envKeyMode, mode)
	}
}
