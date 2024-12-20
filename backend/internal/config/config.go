package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type DatabaseConfig struct {
	Username     string
	Password     string
	Host         string
	Port         uint16
	DBName       string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

type Config struct {
	Port    string
	Env     string
	DB      DatabaseConfig
	Limiter struct {
		RPS     float64
		Burst   int
		Enabled bool
	}
}

func NewConfig() *Config {
	return &Config{}
}

func (cfg *Config) ParseFlags() error {
	username, ok := os.LookupEnv("POSTGRES_USER")
	if !ok {
		return fmt.Errorf("no POSTGRES_USER env variable set")
	}

	password, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if !ok {
		return fmt.Errorf("no POSTGRES_PASSWORD env variable set")
	}

	host, ok := os.LookupEnv("POSTGRES_HOST")
	if !ok {
		return fmt.Errorf("no POSTGRES_HOST env variable set")
	}

	portStr, ok := os.LookupEnv("POSTGRES_PORT")
	if !ok {
		return fmt.Errorf("no POSTGRES_PORT env variable set")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("failed to convert port to int: %w", err)
	}

	dbname, ok := os.LookupEnv("POSTGRES_DB")
	if !ok {
		return fmt.Errorf("no POSTGRES_DB env variable set")
	}

	sslmode, ok := os.LookupEnv("POSTGRES_SSLMODE")
	if !ok {
		return fmt.Errorf("no SSLMode env variable set")
	}

	cfg.DB = DatabaseConfig{
		Username: username,
		Password: password,
		Host:     host,
		Port:     uint16(port),
		DBName:   dbname,
		SSLMode:  sslmode,
	}

	flag.StringVar(&cfg.Port, "port", os.Getenv("PORT"), "API server port")
	flag.StringVar(&cfg.Env, "env", os.Getenv("ENV"), "Environment (development|staging|production)")

	flag.IntVar(&cfg.DB.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.DB.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.DB.MaxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Float64Var(&cfg.Limiter.RPS, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.Limiter.Burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.Limiter.Enabled, "limiter-enabled", false, "Enable rate limiter")

	return nil
}

func (cfg *Config) URL() string {
  // "host=localhost user=admin password=admin dbname=swim-results port=5432 sslmode=disable TimeZone=Europe/Vienna"
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Europe/Vienna",
    cfg.DB.Host,
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.DBName,
    cfg.DB.Port,
		cfg.DB.SSLMode,
	)
}
