package dbex

import (
	"database/sql"
	goMysql "github.com/go-sql-driver/mysql"
	"time"
)

type dbParameter struct {
	DriverName   string
	User         string        // Username
	Password       string        // Password (requires User)
	Net          string        // Network type
	Addr         string        // Network address (requires Net)
	DBName       string        // Database name
	Timeout      time.Duration // Dial timeout
	ReadTimeout  time.Duration // I/O read timeout
	WriteTimeout time.Duration // I/O write timeout
}
type DB struct {
	sqlDB  *sql.DB
	config *dbParameter
}


type DBStats struct {
	MaxOpenConnections int // Maximum number of open connections to the database; added in Go 1.11

	// Pool Status
	OpenConnections int // The number of established connections both in use and idle.
	InUse           int // The number of connections currently in use; added in Go 1.11
	Idle            int // The number of idle connections; added in Go 1.11

	// Counters
	WaitCount         int64         // The total number of connections waited for; added in Go 1.11
	WaitDuration      time.Duration // The total time blocked waiting for a new connection; added in Go 1.11
	MaxIdleClosed     int64         // The total number of connections closed due to SetMaxIdleConns; added in Go 1.11
	MaxLifetimeClosed int64
}

//NewDB
//readTimeout I/O read timeout 30s, 0.5m, 1m30s
//writeTimeout I/O write timeout
//timeout timeout for establishing connections

func newDB(dbConfig *dbParameter) (*DB, error) {
	conf := goMysql.NewConfig()
	conf.User = dbConfig.User
	conf.Passwd = dbConfig.Password
	conf.Net=dbConfig.Net
	conf.Addr=dbConfig.Addr
	conf.DBName = dbConfig.DBName
	conf.WriteTimeout = dbConfig.WriteTimeout
	conf.ReadTimeout = dbConfig.ReadTimeout
	conf.Timeout = dbConfig.Timeout


	formatStr := conf.FormatDSN()

	//db, err := sql.Open(driverName, driverSourceName)
	sdb, err := sql.Open(dbConfig.DriverName, formatStr)
	if err != nil {
		return nil, err
	}

	err=sdb.Ping()
	if err != nil {
		return nil, err
	}

	return &DB{
		sqlDB:  sdb,
		config: dbConfig,
	}, nil
}

func(db *DB)GetSQLDB()*sql.DB {
	return db.sqlDB
}

func (db *DB) Close() {
	db.sqlDB.Close()
}
/*

func (db *db) Ping() error {
	return db.sqlDB.Ping()
}
func (db *db) Stats() *DBStats {
	s := db.sqlDB.Stats()
	return &DBStats{
		MaxOpenConnections: s.MaxOpenConnections,

		// Pool Status
		OpenConnections: s.OpenConnections,
		InUse:           s.InUse,
		Idle:            s.Idle,
		// Counters
		WaitCount:         s.WaitCount,
		WaitDuration:      s.WaitDuration,
		MaxIdleClosed:     s.MaxIdleClosed,
		MaxLifetimeClosed: s.MaxLifetimeClosed,
	}
}
*/
func (db *DB) SelectTableNames() ([]string, error) {
	res := make([]string, 0)

	var tableName string
	//rows, err := db.sqlDB.Query("SELECT table_name FROM information_schema.tables where table_schema  = '" + db.config.DBName + "'")
	rows, err := db.sqlDB.Query("SELECT table_name FROM information_schema.tables where table_schema  = ?","live")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&tableName)
		if err != nil {
			return nil, err
		}
		res = append(res, tableName)
	}
	return res, nil
}
