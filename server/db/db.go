package db

import (
	"os"
	"strconv"
	"time"

	"bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/config"
	"bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/constant"
	"bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/lib/database"
	databasewrapper "bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/lib/database/wrapper"
	servicelogger "bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/lib/service-logger"
)

type Db struct {
	appConfig *config.Config
	logger    *servicelogger.AddonsLogrus
	DbSql     *database.DbSql
}

func NewDatabase(appConfig *config.Config, logger *servicelogger.AddonsLogrus) *Db {
	return &Db{
		appConfig: appConfig,
		logger:    logger,
	}
}

func (d *Db) StartDBConnection() {
	d.logger.Info("Starting Db Connections...")

	d.logger.Info("Main Db - Connecting")
	var err error
	maxRetry, convErr := strconv.Atoi(d.appConfig.DbMaxRetry)
	if convErr != nil {
		maxRetry = constant.DefaultDbMaxRetry
		d.logger.Infof("Failed to convert database max retry, set to default: %d", maxRetry)
	}

	dbTimeout, convErr := strconv.Atoi(d.appConfig.DbTimeout)
	if convErr != nil {
		dbTimeout = constant.DefaultDbTimeout
		d.logger.Infof("Failed to convert database Timeout, set to default: %ds", dbTimeout)
	}

	d.DbSql = database.InitConnectionDB("postgres", database.Config{
		Host:         d.appConfig.DbHost,
		Port:         d.appConfig.DbPort,
		User:         d.appConfig.DbUser,
		Password:     d.appConfig.DbPassword,
		DatabaseName: d.appConfig.DbName,
		SslMode:      d.appConfig.DbSslmode,
		TimeZone:     d.appConfig.DbTimezone,
		MaxRetry:     maxRetry,
		Timeout:      time.Duration(dbTimeout) * time.Second,
	}, &databasewrapper.DatabaseWrapper{})

	err = d.DbSql.Connect()

	if err != nil {
		d.logger.Fatalf("Failed connect to DB main: %v", err)
		os.Exit(1)
		return
	}

	d.DbSql.SetMaxIdleConns(constant.DefaultMaxIdleConns)
	d.DbSql.SetMaxOpenConns(constant.DefaultMaxOpenConns)

	d.logger.Info("Main Db - Connected")
}

func (d *Db) CloseDBConnections() {
	d.logger.Info("Closing DB Main Connection ... ")
	if err := d.DbSql.ClosePmConnection(); err != nil {
		d.logger.Fatalf("Error on disconnection with DB Main : %v", err)
	}
	d.logger.Info("Closing DB Main Success")
}
