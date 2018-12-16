package dbex

import (
	"reflect"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"time"
)


func TestNewDB(t *testing.T){

	conf, err := NewConfigurer("dbexConfig.json")
	if err != nil {
		t.Fatalf("TestNewDB err=%s", err.Error())
	}

	dbc := conf.GetDBConfig()

	t.Logf("timeout %s,read %s, write %s",dbc.Timeout,dbc.ReadTimeout,dbc.WriteTimeout)

	tm, err := time.ParseDuration(dbc.Timeout)
	if err != nil {
		t.Fatalf("TestNewDB parse timeout duration err=%s", err.Error())
	}
	rm, err := time.ParseDuration(dbc.ReadTimeout)
	if err != nil {
		t.Fatalf("TestNewDB parse read timeout duration err=%s", err.Error())
	}
	wm, err := time.ParseDuration(dbc.WriteTimeout)
	if err != nil {
		t.Fatalf("TestNewDB parse write timeout duration err=%s", err.Error())
	}

	dbConf := &DBParameter{
		DriverName:dbc.DriverName,
		User:         dbc.User,
		Password:       dbc.Password,
		Net:          dbc.Net,
		Addr:         dbc.Address,
		DBName:       dbc.DBName,
		Timeout:      tm,
		ReadTimeout:  rm,
		WriteTimeout: wm,

	}
	db, err := newDB(dbConf)
	if err != nil {
		panic(err.Error())
	}

	db.Close()

}

func getDB()*db{
	conf, _ := NewConfigurer("dbexConfig.json")

	dbc := conf.GetDBConfig()


	tm, _ := time.ParseDuration(dbc.Timeout)
	rm, _ := time.ParseDuration(dbc.ReadTimeout)
	wm, _ := time.ParseDuration(dbc.WriteTimeout)

	dbConf := &DBParameter{
		DriverName:dbc.DriverName,
		User:         dbc.User,
		Password:       dbc.Password,
		Net:          dbc.Net,
		Addr:         dbc.Address,
		DBName:       dbc.DBName,
		Timeout:      tm,
		ReadTimeout:  rm,
		WriteTimeout: wm,

	}
	db, _ := newDB(dbConf)
	return db
}

//func TestDB_Close(t *testing.T) {
//	type fields struct {
//		sqlDB *sql.DB
//	}
//	tests := []struct {
//		name   string
//		fields fields
//	}{
//	// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			db := &DB{
//				sqlDB: tt.fields.sqlDB,
//			}
//			db.Close()
//		})
//	}
//}

//func TestDB_Ping(t *testing.T) {
//	type fields struct {
//		sqlDB *sql.DB
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		wantErr bool
//	}{
//	// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			db := &DB{
//				sqlDB: tt.fields.sqlDB,
//			}
//			if err := db.Ping(); (err != nil) != tt.wantErr {
//				t.Errorf("DB.Ping() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}

func TestDB_Stats(t *testing.T) {

	db:=getDB()
	sdb:=db.GetSQLDB()
	dbStats := sdb.Stats()
	defer db.Close()
	t.Logf("db stats =%+v", dbStats)

}

//func TestDB_Query(t *testing.T) {
//	type fields struct {
//		sqlDB *sql.DB
//	}
//	type args struct {
//		query string
//		args  []interface{}
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    *sql.Rows
//		wantErr bool
//	}{
//	// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			db := &DB{
//				sqlDB: tt.fields.sqlDB,
//			}
//			got, err := db.Query(tt.args.query, tt.args.args...)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("DB.Query() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("DB.Query() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

//func TestDB_QueryToJSON(t *testing.T) {
//	type fields struct {
//		sqlDB *sql.DB
//	}
//	tests := []struct {
//		name   string
//		fields fields
//	}{
//	// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			db := &DB{
//				sqlDB: tt.fields.sqlDB,
//			}
//			db.QueryToJSON()
//		})
//	}
//}

func TestDB_SelectTableNames(t *testing.T) {

	db:=getDB()
	defer db.Close()

	tests := []struct {
		name    string
		want    []string
		errorTable bool
	}{
		{
			name:    "0",
			want:    []string{"partner", "user"},
			errorTable :false,
		},
		{
			name:    "0",
			want:    []string{"partner", "xxxx"},
			errorTable:true,

		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.SelectTableNames()
			if (err != nil) {
				t.Errorf("DB.SelectTableNames() error = %v", err)
				return
			}
			if (!reflect.DeepEqual(got, tt.want)) !=tt.errorTable{
				t.Errorf("DB.SelectTableNames() = %v, want %v", got, tt.want)
			}
		})
	}

}



