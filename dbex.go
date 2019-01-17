package dbex

import (
	"time"
)

//DBEX is main struct, do not init this struct.
//using NewDBEX function
type DBEX struct {
	Configure  *config
	DB         *DB
	Logger     *Logger
	HTTPServer *httpServer
	HTTPClient *httpClient
}

//NewDBEX returns a dbex structure with config file path parameter passed in.
func NewDBEX(configFilePath string) (*DBEX, error) {

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

	return &DBEX{
		Logger:     logger,
		Configure:  conf,
		DB:         db,
		HTTPServer: hs,
	}, nil
}
