package xianorm

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

// DB 管理数据库连接和执行数据库操作
type DB struct {
	Db           *sql.DB       // 数据库连接
	TableName    string        // 表名
	Prepare      string        // 准备的SQL查询语句
	AllExec      []interface{} // 存储所有字段的值
	Sql          string        // 最终生成的SQL语句
	WhereParam   string        // WHERE 条件
	WhereExec    []interface{} // 存储 WHERE 条件中的参数值
	LimitParam   string        // LIMIT 条件
	OrderParam   string        // ORDER BY 条件
	OrWhereParam string        // OR WHERE 条件
	UpdateParam  string        // UPDATE 条件
	UpdateExec   []interface{} // 存储UPDATE 条件中的参数值
	FieldParam   string        // SELECT 字段
	TransStatus  int           // 事务状态，0表示无事务，1表示有事务
	Tx           *sql.Tx       // 事务对象
	GroupParam   string        // GROUP BY 条件
	HavingParam  string        // HAVING 条件
	HavingExec   []interface{} // 存储HAVING 条件中的参数值
}

// NewDB 新建Mysql连接
func NewDB(Username string, Password string, Address string, Dbname string) (*DB, error) {
	dsn := Username + ":" + Password + "@tcp(" + Address + ")/" + Dbname + "?charset=utf8&timeout=5s&readTimeout=6s"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	return &DB{
		Db:         db,
		FieldParam: "*",
	}, nil
}
