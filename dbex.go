package dbex

import (
	"time"
)
type dbex struct {
	Configure  *config
	DB         *DB
	Logger     *Logger
	HttpServer *httpServer
}

func NewDBEX(configFilePath string) (*dbex, error) {

	conf, err := newConfigurer(configFilePath)
	if err != nil {
		return nil, err
	}
	logFile := conf.GetLoggerConfig().LogFileNamePrefix
	logger, err := newLogger(logFile)

	if err != nil {
		return nil, err
	}

	dbc := conf.GetDBConfig()
	tm, _ := time.ParseDuration(dbc.Timeout)
	rm, _ := time.ParseDuration(dbc.ReadTimeout)
	wm, _ := time.ParseDuration(dbc.WriteTimeout)

	dbConf := &dbParameter{
		DriverName:   dbc.DriverName,
		User:         dbc.User,
		Password:     dbc.Password,
		Net:          dbc.Net,
		Addr:         dbc.Address,
		DBName:       dbc.DBName,
		Timeout:      tm,
		ReadTimeout:  rm,
		WriteTimeout: wm,
	}
	db, err := newDB(dbConf)

	if err != nil {
		return nil, err
	}

	httpConf := conf.GetHTTPConfig()

	httpParams := &httpParameter{
		Address:        httpConf.Address,
		Port:           httpConf.Port,
		ReadTimeout:    time.Second * time.Duration(httpConf.ReadTimeoutSecond),
		WriteTimeout:   time.Second * time.Duration(httpConf.WriteTimeoutSecond),
		IdleTimeout:    time.Second * time.Duration(httpConf.IdleTimeoutSecond),
		MaxHeaderBytes: httpConf.MaxHeaderBytes,
	}

	hs, err := NewHTTPServer(httpParams)

	if err != nil {
		return nil, err
	}

	return &dbex{
		Logger:     logger,
		Configure:  conf,
		DB:         db,
		HttpServer: hs,
	}, nil
}
