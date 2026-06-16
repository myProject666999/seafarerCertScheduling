package config

import "github.com/sirupsen/logrus"

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: ":3000",
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     "3306",
			User:     "root",
			Password: "123456",
			DBName:   "seafarer_cert_scheduling",
		},
	}
}

func (c *DatabaseConfig) DSN() string {
	return c.User + ":" + c.Password + "@tcp(" + c.Host + ":" + c.Port + ")/" + c.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
}

func InitLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	return logger
}
