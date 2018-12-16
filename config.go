package dbex

import (
	"io/ioutil"
	"encoding/json"
	"path/filepath"
	"os"
	"bytes"
)

//Config config main struct
type config struct {
	data confData
}

//NewConfig make a new config struct
func newConfigurer(configFileName string) (*config, error) {

	configFilePath := getFilePosition(configFileName)

	cf := &config{
		data: confData{},
	}
	if err := cf.loadConfig(configFilePath); err != nil {
		return nil, err
	}
	return cf, nil
}

//LoadConfig loads config file.
func (c *config) loadConfig(filePath string) error {

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


func (c *config) Version() string {
	return c.data.Version
}

//http
func (c *config) GetHTTPConfig() *httpServerConf {
	return c.data.HTTPServer
}

//logger
func (c *config) GetLoggerConfig() *loggerConf {
	return c.data.Logger
}

//db
func (c *config) GetDBConfig() *dbConf {
	return c.data.DefaultDB
}

