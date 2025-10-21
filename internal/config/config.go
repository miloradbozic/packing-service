package config

import (
	"fmt"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Packs    PacksConfig    `yaml:"packs"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type DatabaseConfig struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	DBName          string `yaml:"dbname"`
	SSLMode         string `yaml:"sslmode"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	ConnMaxLifetime string `yaml:"conn_max_lifetime"`
}

type PacksConfig struct {
	Sizes []int `yaml:"sizes"`
}

func Load(path string) (*Config, error) {
	var config Config
	
	// Try to load from file first
	file, err := os.Open(path)
	if err == nil {
		defer file.Close()
		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(&config); err != nil {
			return nil, fmt.Errorf("failed to decode config: %w", err)
		}
	}

	// Override with environment variables if they exist
	overrideWithEnvVars(&config)

	// Sort pack sizes in ascending order
	sort.Ints(config.Packs.Sizes)

	return &config, nil
}

func overrideWithEnvVars(config *Config) {
	// Server configuration
	if port := os.Getenv("PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Server.Port = p
		}
	}
	if host := os.Getenv("HOST"); host != "" {
		config.Server.Host = host
	}

	// Database configuration from DATABASE_URL (Heroku format)
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		if err := parseDatabaseURL(dbURL, config); err != nil {
			// Log error but continue with other env vars
			fmt.Printf("Warning: Failed to parse DATABASE_URL: %v\n", err)
		}
	}
	
	// Individual database environment variables
	if host := os.Getenv("DB_HOST"); host != "" {
		config.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Database.Port = p
		}
	}
	if user := os.Getenv("DB_USER"); user != "" {
		config.Database.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		config.Database.Password = password
	}
	if dbname := os.Getenv("DB_NAME"); dbname != "" {
		config.Database.DBName = dbname
	}
	if sslmode := os.Getenv("DB_SSLMODE"); sslmode != "" {
		config.Database.SSLMode = sslmode
	}
}

func parseDatabaseURL(dbURL string, config *Config) error {
	u, err := url.Parse(dbURL)
	if err != nil {
		return fmt.Errorf("invalid database URL: %w", err)
	}

	// Extract host and port
	host := u.Hostname()
	port := u.Port()
	if port == "" {
		port = "5432" // Default PostgreSQL port
	}
	if portInt, err := strconv.Atoi(port); err == nil {
		config.Database.Port = portInt
	}
	config.Database.Host = host

	// Extract user and password
	if u.User != nil {
		config.Database.User = u.User.Username()
		if password, ok := u.User.Password(); ok {
			config.Database.Password = password
		}
	}

	// Extract database name
	config.Database.DBName = strings.TrimPrefix(u.Path, "/")

	// Extract SSL mode from query parameters
	query := u.Query()
	if sslmode := query.Get("sslmode"); sslmode != "" {
		config.Database.SSLMode = sslmode
	} else {
		config.Database.SSLMode = "require" // Default for Heroku
	}

	return nil
}
