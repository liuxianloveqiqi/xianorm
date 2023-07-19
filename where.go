package xianorm

//func (d *DB) Where(data interface{}) *DB {
//	// 反射获取数据的类型和值
//	t := reflect.TypeOf(data)
//	v := reflect.ValueOf(data)
//
//	// 切片用于存储字段名和条件
//	var conditions []string
//
//	// 遍历结构体字段并解析标签
//	for i := 0; i < t.NumField(); i++ {
//		// 跳过未导出的字段（小写开头的字段）
//		if !v.Field(i).CanInterface() {
//			continue
//		}
//
//		// 解析标签以获取真实的 SQL 字段名
//		tag := t.Field(i).Tag.Get("xianorm")
//		fieldName := t.Field(i).Name
//		if tag != "" {
//			fieldName = strings.Split(tag, ",")[0]
//		}
//
//		// 将字段名和条件添加到切片中
//		conditions = append(conditions, fmt.Sprintf("%s=?", fieldName))
//
//		// 反射获取字段值并添加到 WhereExec 切片中
//		d.WhereExec = append(d.WhereExec, v.Field(i).Interface())
//	}
//
//	// 使用" and "连接条件并添加到 WhereParam 中
//	conditionStr := strings.Join(conditions, " and ")
//
//	// 多次调用判断，用来构建复杂的 WHERE 子句
//	if d.WhereParam != "" {
//		// 如果不为空，则说明这是第二次调用了，用 "and (" 来做隔离
//		d.WhereParam += " and (" + conditionStr + ") "
//	} else {
//		d.WhereParam += "(" + conditionStr + ") "
//	}
//
//	return d
//}

// Where 构造原生 SQL 查询条件
func (d *DB) Where(query string, args ...interface{}) *DB {
	// 创建一个新的 DB 结构体，复制原有的字段和值
	newDB := *d

	// 初始化 WhereExec 切片
	if newDB.WhereExec == nil {
		newDB.WhereExec = []interface{}{}
	}

	// 添加新的查询条件和参数
	newDB.WhereParam += query + " "
	newDB.WhereExec = append(newDB.WhereExec, args...)

	return &newDB
}

// Select 构造原生 SQL 查询字段
func (d *DB) Select(fields string) *DB {
	// 创建一个新的 DB 结构体，复制原有的字段和值
	newDB := *d

	// 设置查询字段
	newDB.FieldParam = fields

	return &newDB
}
