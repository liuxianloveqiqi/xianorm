package xianorm

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
)

// Insert 插入数据
func (d *DB) Insert(data interface{}) (int64, error) {
	//return d.insertOrReplaceData(data, "insert")
	// 利用反射判断是批量还是单个插入
	getValue := reflect.ValueOf(data).Kind()
	if getValue == reflect.Struct {
		return d.insertOrReplaceData(data, "insert")
	} else if getValue == reflect.Slice || getValue == reflect.Array {
		return d.batchInsertOrReplaceData(data, "insert")
	} else {
		return 0, errors.New("插入的数据格式不正确，单个插入格式为: struct，批量插入格式为: []struct")
	}
}

// Replace 替换插入数据
func (d *DB) Replace(data interface{}) (int64, error) {
	// 利用反射判断是批量还是单个替换
	getValue := reflect.ValueOf(data).Kind()
	if getValue == reflect.Struct {
		return d.insertOrReplaceData(data, "insert")
	} else if getValue == reflect.Slice || getValue == reflect.Array {
		return d.batchInsertOrReplaceData(data, "insert")
	} else {
		return 0, errors.New("插入的数据格式不正确，单个插入格式为: struct，批量插入格式为: []struct")
	}
}

//// BatchInsert 批量插入数据
//func (d *DB) BatchInsert(data interface{}) (int64, error) {
//	return d.batchInsertData(data, "insert")
//}
//
//// BatchReplace 批量替换插入数据
//func (d *DB) BatchReplace(data interface{}) (int64, error) {
//	return d.batchInsertData(data, "replace")
//}

// 插入或者替换数据的真正方法
func (d *DB) insertOrReplaceData(data interface{}, insertType string) (int64, error) {
	// 反射type和value
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	// 字段名和
	var fieldName []string
	// 占位符
	var placeHolder []string

	// 循环判断结构体字段
	for i := 0; i < t.NumField(); i++ {
		// 跳过小写开头的字段
		if !v.Field(i).CanInterface() {
			continue
		}

		// 解析tag，找出真实的sql字段名
		tag := t.Field(i).Tag.Get("xianorm")
		parseFieldAndPlaceHolder(0, tag, &fieldName, &placeHolder, t, i)

		// 保存字段的值
		d.AllExec = append(d.AllExec, v.Field(i).Interface())
	}

	// 拼接表名、字段名和占位符，生成最终的插入SQL语句
	//d.Prepare = insertType + " into " + d.GetTable() + " (" + strings.Join(fieldName, ",") + ") values(" + strings.Join(placeHolder, ",") + ")"
	d.Prepare = fmt.Sprintf("%s into %s (%s) values (%s)", insertType, d.GetTable(), strings.Join(fieldName, ","), strings.Join(placeHolder, ","))

	// 判断是否是事务
	var stmt *sql.Stmt
	var err error
	fmt.Println("iii", d.TransStatus)
	// 准备SQL语句，返回一个预处理语句
	if d.TransStatus == 1 {
		stmt, err = d.Tx.Prepare(d.Prepare)
	} else {
		stmt, err = d.Db.Prepare(d.Prepare)
	}
	defer stmt.Close()

	// 来执行预处理语句
	result, err := stmt.Exec(d.AllExec...)
	if err != nil {
		return 0, d.setErrorInfo(err)
	}

	// 获取刚刚插入数据ID
	id, _ := result.LastInsertId()
	return id, nil
}

// 批量插入数据的真正方法
func (d *DB) batchInsertOrReplaceData(batchData interface{}, insertType string) (int64, error) {
	getValue := reflect.ValueOf(batchData)
	l := getValue.Len()

	// 字段名
	var fieldName []string
	// 占位符
	var placeHolder []string

	// 循环处理每个子元素
	for i := 0; i < l; i++ {
		v := getValue.Index(i)
		t := v.Type()
		if t.Kind() != reflect.Struct {
			log.Fatal("批量插入的子元素必须是结构体类型")
		}

		num := v.NumField()

		// 当前子元素的占位符
		var subPlaceHolder []string
		// 循环遍历子元素的字段
		for j := 0; j < num; j++ {
			// 跳过小写开头的字段
			if !v.Field(j).CanInterface() {
				continue
			}

			// 解析tag，找出真实的sql字段名，并生成占位符
			tag := t.Field(j).Tag.Get("xianorm")
			parseFieldAndPlaceHolder(i, tag, &fieldName, &subPlaceHolder, t, j)

			// 字段值
			d.AllExec = append(d.AllExec, v.Field(j).Interface())
		}

		// 子元素拼接成多个()括号后的值
		placeHolder = append(placeHolder, "("+strings.Join(subPlaceHolder, ",")+")")
	}

	// 拼接表名、字段名和占位符，生成最终的批量插入SQL语句
	d.Prepare = fmt.Sprintf("%s into %s (%s) values %s", insertType, d.GetTable(), strings.Join(fieldName, ","), strings.Join(placeHolder, ","))
	fmt.Println(d.Prepare)
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

	// 来执行预处理语句
	result, err := stmt.Exec(d.AllExec...)
	if err != nil {
		return 0, d.setErrorInfo(err)
	}

	// 获取刚刚插入数据ID
	id, _ := result.LastInsertId()
	return id, nil
}

// 解析tag，找出真实的sql字段名，并生成占位符
func parseFieldAndPlaceHolder(i int, tag string, fieldName *[]string, subPlaceHolder *[]string, t reflect.Type, j int) {
	// 字段名只记录第一个的
	if i == 0 {
		if tag != "" {
			// 跳过自增字段
			fmt.Println(tag)
			if !strings.Contains(strings.ToLower(tag), "auto_increment") {
				// 获取真实的sql字段名
				// 将标签值按逗号分割，取第一个部分作为真实的SQL字段名，并将其添加到
				*fieldName = append(*fieldName, strings.Split(tag, ",")[0])

			}
		} else {
			// 若字段没有tag，则使用字段名作为sql字段名
			*fieldName = append(*fieldName, t.Field(j).Name)
		}
	}
	// 在placeholder切片中添加一个问号 ?，用作占位符
	*subPlaceHolder = append(*subPlaceHolder, "?")
}
