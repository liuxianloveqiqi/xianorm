package xianorm

import (
	"database/sql"
	"log"
)

type DB struct {
	Db           *sql.DB
	TableName    string
	Prepare      string        // 生成的插入 SQL 语句
	AllExec      []interface{} // 存储所有字段的值
	Sql          string
	WhereParam   string
	LimitParam   string
	OrderParam   string
	OrWhereParam string
	WhereExec    []interface{}
	UpdateParam  string
	UpdateExec   []interface{}
	FieldParam   string
	TransStatus  int
	Tx           *sql.Tx
	GroupParam   string
	HavingParam  string
}

// NewDB 新建Mysql连接
func NewDB(Username string, Password string, Address string, Dbname string) (*DB, error) {
	dsn := Username + ":" + Password + "@tcp(" + Address + ")/" + Dbname + "?charset=utf8&timeout=5s&readTimeout=6s"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	//最大连接数等配置，先占个位
	//db.SetMaxOpenConns(3)
	//db.SetMaxIdleConns(3)
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	return &DB{
		Db:         db,
		FieldParam: "*",
	}, nil
}