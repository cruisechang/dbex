package config

import (
	"io/ioutil"
	"encoding/json"
	"path/filepath"
	"os"
	"bytes"
)

//Config config main struct
type Configurer struct {
	data confData
}

//NewConfig make a new config struct
func NewConfigurer(configFileName string) (*Configurer, error) {

	configFilePath := getFilePosition(configFileName)

	cf := &Configurer{
		data: confData{},
	}
	if err := cf.loadConfig(configFilePath); err != nil {
		return nil, err
	}
	return cf, nil
}

//LoadConfig loads config file.
func (c *Configurer) loadConfig(filePath string) error {

	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	//unmarshal to struct
	if err := json.Unmarshal(b, &c.data); err != nil {
		return err
	}

	return nil
}

func getFilePosition(fileName string) string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
	}

	var buf bytes.Buffer
	buf.WriteString(dir)
	buf.WriteString("/")
	buf.WriteString(fileName)

	return buf.String()
}

type confData struct {
	Version    string `json:"Version"`
	HTTPServer *httpServerConf `json:"HttpServer"`
	DefaultDB  *dbConf `json:"DefaultDB"`
	Logger     *loggerConf `json:"Logger"`
}

type httpServerConf struct {
	Address            string
	Port               string
	ReadTimeoutSecond  int
	WriteTimeoutSecond int
	IdleTimeoutSecond  int
	MaxHeaderBytes     int
}
type dbConf struct {
	DriverName     string
	DataSourceName string
	User           string
	Password       string
	Net            string
	Address        string
	DBName         string
	Timeout        string
	ReadTimeout    string
	WriteTimeout   string
}
type loggerConf struct {
	LogFileNamePrefix string
}

//http
//func (c *Configurer) HttpServerAddress() (addr string) {
//	return c.data.HttpServer.Address
//}
//func (c *Configurer) HttpServerPort() (port string) {
//	return c.data.HttpServer.Port
//}
//
//func (c *Configurer) HttpServerReadTimeout() time.Duration {
//	return time.Second * time.Duration(c.data.HttpServer.ReadTimeoutSecond)
//}
//func (c *Configurer) HttpServerWriteTimeout() time.Duration {
//	return time.Second * time.Duration(c.data.HttpServer.WriteTimeoutSecond)
//}
//func (c *Configurer) HttpServerIdleTimeout() time.Duration {
//	return time.Second * time.Duration(c.data.HttpServer.IdleTimeoutSecond)
//}
//func (c *Configurer) HttpServerMaxHeaderBytes() int {
//	return c.data.HttpServer.MaxHeaderBytes
//}

func (c *Configurer) Version() string {
	return c.data.Version
}

//http
func (c *Configurer) GetHTTPConfig() *httpServerConf {
	return c.data.HTTPServer
}

//logger
func (c *Configurer) GetLoggerConfig() *loggerConf {
	return c.data.Logger
}

//db
func (c *Configurer) GetDBConfig() *dbConf {
	return c.data.DefaultDB
}

//func (c *Configurer) DBDriverName() string {
//	return c.data.DB.DriverName
//}
//func (c *Configurer) DBDataSourceName() string {
//	return c.data.DB.DataSourceName
//
//}
//func (c *Configurer) User() string {
//	return c.data.DB.User
//
//}
//
//func (c *Configurer) DBName() string {
//	return c.data.DB.DBName
//
//}
//func (c *Configurer) Password() string {
//	return c.data.DB.pas
//
//}
