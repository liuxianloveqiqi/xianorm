package xianorm

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Query 执行原生SQL查询操作并返回map切片

func (d *DB) Query() ([]map[string]string, error) {
	// 拼接SQL语句
	d.Sql = "select " + d.FieldParam + " from " + d.GetTable()

	// 如果WhereParam或OrWhereParam不为空，添加查询条件
	if d.WhereParam != "" || d.OrWhereParam != "" {
		d.Sql += " where " + d.WhereParam + d.OrWhereParam
	}

	// 如果 OrderParam 不为空，添加排序条件
	if d.OrderParam != "" {
		d.Sql += " order by " + d.OrderParam
	}

	// 如果 LimitParam 不为空，添加限制条件
	if d.LimitParam != "" {
		d.Sql += " limit " + d.LimitParam
	}
	// 执行查询
	rows, err := d.Db.Query(d.Sql, d.WhereExec...)
	if err != nil {
		return nil, d.setErrorInfo(err)
	}
	defer rows.Close()

	// 读取查询出的列字段名
	columns, err := rows.Columns()
	if err != nil {
		return nil, d.setErrorInfo(err)
	}

	// 构造返回结果
	var results []map[string]string
	for rows.Next() {
		// 准备一个新的 map 用于存储当前行的数据
		row := make(map[string]string)

		// 将当前行数据按列名存储到 map 中
		scanArgs := make([]interface{}, len(columns))
		scanValues := make([]interface{}, len(columns))
		// 把scanValues每个元素的值，都是scanArgs里的每个值的地址,2个进行了深度绑定。
		for i := range columns {
			scanValues[i] = &scanArgs[i]
		}

		if err := rows.Scan(scanValues...); err != nil {
			return nil, d.setErrorInfo(err)
		}

		for i, col := range columns {
			row[col] = fmt.Sprintf("%s", scanArgs[i])
		}

		// 将当前行数据的map添加到结果切片中
		results = append(results, row)
	}

	return results, nil
}

// QueryOne 查询返回单个map，直接调用Query再加上Limit
func (d *DB) QueryOne() (map[string]string, error) {
	//limit 1 单个查询
	results, err := d.Limit(1).Query()
	if err != nil {
		return nil, d.setErrorInfo(err)
	}
	//判断是否为空
	if len(results) == 0 {
		return nil, nil
	} else {
		return results[0], nil
	}
}

// Find 执行原生SQL查询操作并将结果映射到struct切片中
func (d *DB) Find(result interface{}) error {
	// 判断传入的参数是否为指针类型
	if reflect.ValueOf(result).Kind() != reflect.Ptr {
		return d.setErrorInfo(errors.New("参数请传指针变量！"))
	}

	// 判断传入的指针是否为nil
	if reflect.ValueOf(result).IsNil() {
		return d.setErrorInfo(errors.New("参数不能是空指针！"))
	}

	d.Sql = "select " + d.FieldParam + " from " + d.GetTable()

	// 如果WhereParam或OrWhereParam不为空，添加查询条件
	if d.WhereParam != "" || d.OrWhereParam != "" {
		d.Sql += " where " + d.WhereParam + d.OrWhereParam
	}

	// 如果 OrderParam 不为空，添加排序条件
	if d.OrderParam != "" {
		d.Sql += " order by " + d.OrderParam
	}

	// 如果 LimitParam 不为空，添加限制条件
	if d.LimitParam != "" {
		d.Sql += " limit " + d.LimitParam
	}

	// 执行查询
	rows, err := d.Db.Query(d.Sql, d.WhereExec...)
	if err != nil {
		return d.setErrorInfo(err)
	}
	defer rows.Close()

	// 读取查询出的列字段名
	columns, err := rows.Columns()
	if err != nil {
		return d.setErrorInfo(err)
	}

	// values是每个列的值，这里获取到byte里
	values := make([][]byte, len(columns))

	// 因为每次查询出来的列是不定长的，用len(column)定住当次查询的长度
	scans := make([]interface{}, len(columns))

	// 原始struct的切片值
	destSlice := reflect.ValueOf(result).Elem()

	// 原始单个struct的类型
	destType := destSlice.Type().Elem()

	for i := range values {
		scans[i] = &values[i]
	}

	// 循环遍历
	for rows.Next() {
		// 创建一个新的struct实例
		dest := reflect.New(destType).Elem()

		if err := rows.Scan(scans...); err != nil {
			// query.Scan查询出来的不定长值放到scans[i] = &values[i]，也就是每行都放在values里
			return d.setErrorInfo(err)
		}

		// 遍历一行数据的各个字段
		for k, v := range values {
			// 每行数据是放在values里面，现在把它挪到row里
			key := columns[k]
			value := string(v)

			// 遍历结构体
			for i := 0; i < destType.NumField(); i++ {
				// 看下是否有sql别名
				tag := destType.Field(i).Tag.Get("xianorm")
				var fieldName string
				// 如果标签不为空则跳过自增主键找的字段值
				if tag != "" && !strings.Contains(strings.ToLower(tag), "auto_increment") {
					fieldName = strings.Split(tag, ",")[0]

				} else {
					// 没有标签默认将结构体字段小写后取字段名
					fieldName = strings.ToLower(destType.Field(i).Name)

				}

				// struct里没有这个key
				if key != fieldName {
					continue
				}

				// 反射赋值
				if err := d.reflectSet(dest, i, value); err != nil {
					return err
				}
			}
		}
		// 赋值
		destSlice.Set(reflect.Append(destSlice, dest))
	}

	return nil
}

// reflectSet 反射赋值函数，根据字段类型将字符串值转换为相应的类型，并设置到结构体的字段中
func (d *DB) reflectSet(dest reflect.Value, i int, value string) error {
	// 获取字段的类型
	fieldKind := dest.Field(i).Kind()

	// 根据字段类型进行转换和赋值
	switch fieldKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// 将字符串值转换为int64类型
		res, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return d.setErrorInfo(err)
		}
		// 设置字段的值
		dest.Field(i).SetInt(res)
	case reflect.String:
		// 直接将字符串值设置为字段的值
		dest.Field(i).SetString(value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// 将字符串值转换为uint64类型
		res, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return d.setErrorInfo(err)
		}
		// 设置字段的值
		dest.Field(i).SetUint(res)
	case reflect.Float32, reflect.Float64:
		// 将字符串值转换为float32类型
		res, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return d.setErrorInfo(err)
		}
		// 设置字段的值
		dest.Field(i).SetFloat(res)

	case reflect.Bool:
		// 将字符串值转换为bool类型
		res, err := strconv.ParseBool(value)
		if err != nil {
			return d.setErrorInfo(err)
		}
		// 设置字段的值
		dest.Field(i).SetBool(res)
	default:
		// 如果不支持的类型，则返回错误
		return fmt.Errorf("不支持的字段类型：%v", fieldKind.String())
	}

	return nil
}
