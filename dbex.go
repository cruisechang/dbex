package dbex

import (
	"github.com/cruisechang/dbex/log"
	"github.com/cruisechang/dbex/http"
	"github.com/cruisechang/dbex/db"
	"github.com/cruisechang/dbex/config"
	"time"
)

type dbex struct {
	Configurer *config.Configurer
	DB         *db.DB
	Logger     *log.Logger
	HttpServer *http.Server
}

func NewDBEX(configFilePath string) (*dbex, error) {

	conf, err := config.NewConfigurer(configFilePath)
	if err != nil {
		return nil, err
	}
	logFile := conf.GetLoggerConfig().LogFileNamePrefix
	logger, err := log.NewLogger(logFile)

	if err != nil {
		return nil, err
	}



	dbc := conf.GetDBConfig()
	tm, _ := time.ParseDuration(dbc.Timeout)
	rm, _ := time.ParseDuration(dbc.ReadTimeout)
	wm, _ := time.ParseDuration(dbc.WriteTimeout)

	dbConf := &db.DBParameter{
		DriverName:   dbc.DriverName,
		User:         dbc.User,
		Passwd:       dbc.Password,
		Net:          dbc.Net,
		Addr:         dbc.Address,
		DBName:       dbc.DBName,
		Timeout:      tm,
		ReadTimeout:  rm,
		WriteTimeout: wm,
	}
	db, err := db.NewDB(dbConf)

	if err != nil {
		return nil, err
	}

	httpConf:=conf.GetHTTPConfig()

	httpParams := &http.ServerParameters{
		Address:        httpConf.Address,
		Port:           httpConf.Port,
		ReadTimeout:    time.Second * time.Duration(httpConf.ReadTimeoutSecond),
		WriteTimeout:   time.Second * time.Duration(httpConf.WriteTimeoutSecond),
		IdleTimeout:    time.Second * time.Duration(httpConf.IdleTimeoutSecond),
		MaxHeaderBytes: httpConf.MaxHeaderBytes,
	}

	hs,err:=http.NewServer(httpParams)

	if err != nil {
		return nil, err
	}


	return &dbex{
		Logger:     logger,
		Configurer: conf,
		DB:         db,
		HttpServer:hs,
	}, nil
}
