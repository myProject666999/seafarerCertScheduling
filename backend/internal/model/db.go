package model

import (
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(dsn string, log *logrus.Logger) error {
	gormLogger := logger.New(
		&gormLogWriter{log},
		logger.Config{
			SlowThreshold:             time.Second * 2,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return nil
}

type gormLogWriter struct {
	log *logrus.Logger
}

func (w *gormLogWriter) Printf(format string, args ...interface{}) {
	w.log.Debugf(format, args...)
}
