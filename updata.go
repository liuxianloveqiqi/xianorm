package xianorm

import (
	"database/sql"
	"errors"
	"reflect"
	"strings"
)

func (d *DB) Update(data ...interface{}) (int64, error) {
	// 判断传入参数的类型
	switch len(data) {
	case 1:
		// 单个参数为结构体
		return d.updateWithStruct(data[0])
	case 2:
		// 两个参数直接拼接SQL
		return d.updateWithSQL(data[0].(string), data[1])
	default:
		return 0, errors.New("参数个数错误")
	}
}

// updateWithStruct 根据结构体生成更新SQL
func (d *DB) updateWithStruct(data interface{}) (int64, error) {
	// 反射获取结构体字段信息
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	var fieldNameArray []string
	for i := 0; i < t.NumField(); i++ {
		// 跳过小写开头的字段
		if !v.Field(i).CanInterface() {
			continue
		}

		// 解析tag，找出真实的sql字段名
		tag := t.Field(i).Tag.Get("xianorm")
		if tag != "" {
			fieldNameArray = append(fieldNameArray, strings.Split(tag, ",")[0]+"=?")
		} else {
			fieldNameArray = append(fieldNameArray, t.Field(i).Name+"=?")
		}

		// 添加字段的值到UpdateExec
		d.UpdateExec = append(d.UpdateExec, v.Field(i).Interface())
	}
	d.UpdateParam += strings.Join(fieldNameArray, ",")

	// 拼接SQL
	return d.executeUpdate()
}

// updateWithSQL 根据传入的SQL字符串执行更新
func (d *DB) updateWithSQL(sql string, args interface{}) (int64, error) {
	// 使用传入的SQL字符串和参数直接执行更新
	d.UpdateParam += sql + "=?"
	d.UpdateExec = append(d.UpdateExec, args)

	// 拼接SQL
	return d.executeUpdate()
}

// executeUpdate 执行更新SQL
func (d *DB) executeUpdate() (int64, error) {
	// 拼接SQL
	d.Prepare = "update " + d.GetTable() + " set " + d.UpdateParam

	// 如果WhereParam不为空，添加查询条件
	if d.WhereParam != "" {
		d.Prepare += " where " + d.WhereParam
	}

	// 如果LimitParam不为空，添加限制条件
	if d.LimitParam != "" {
		d.Prepare += " limit " + d.LimitParam
	}

	// 判断是否是事务
	var stmt *sql.Stmt
	var err error
	// 准备SQL语句，返回一个预处理语句
	if d.TransStatus == 1 {
		stmt, err = d.Tx.Prepare(d.Prepare)
	} else {
		stmt, err = d.Db.Prepare(d.Prepare)
	}
	defer stmt.Close()

	// 合并UpdateExec和WhereExec
	d.AllExec = append(d.UpdateExec, d.WhereExec...)

	// 执行更新操作
	result, err := stmt.Exec(d.AllExec...)
	if err != nil {
		return 0, d.setErrorInfo(err)
	}

	// 获取影响的行数
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, d.setErrorInfo(err)
	}

	return rowsAffected, nil
}
