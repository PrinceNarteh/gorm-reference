// Package config provides application configuration structures.
package config

import "time"

type Config struct {
	App appConfig
	DB  dbConfig
}

type appConfig struct {
	Port    int
	Env     string
	Version string
}

type dbConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxIdleConns    int
	MaxOpenConns    int
	MaxIdleTime     time.Duration
	MaxConnLifetime time.Duration
}
