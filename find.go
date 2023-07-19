package xianorm

// Find 执行原生 SQL 查询操作并返回结果
func (d *DB) Find() ([]map[string]interface{}, error) {
	// 拼接查询 SQL 语句
	d.Sql = "select " + d.FieldParam + " from " + d.GetTable()

	// 如果 WhereParam 不为空，添加查询条件
	if d.WhereParam != "" {
		d.Sql += " where " + d.WhereParam
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

	// 获取查询结果字段列表
	columns, err := rows.Columns()
	if err != nil {
		return nil, d.setErrorInfo(err)
	}

	// 构造返回结果
	var results []map[string]interface{}
	for rows.Next() {
		// 构造一个和查询字段一样大小的 interface{} 切片，用于存储每条记录的值
		values := make([]interface{}, len(columns))

		// 将每个字段的指针添加到 interface{} 切片中
		for i := range values {
			values[i] = new(interface{})
		}

		// 执行 Scan 操作将每条记录的值存储到 values 切片中
		if err := rows.Scan(values...); err != nil {
			return nil, d.setErrorInfo(err)
		}

		// 构造 map[string]interface{} 对象，将字段名和对应的值组成键值对
		result := make(map[string]interface{})
		for i, col := range columns {
			// 获取字段对应的值，并将其转换为真实类型
			val := *(values[i].(*interface{}))
			result[col] = val
		}

		// 将本次查询结果添加到结果集中
		results = append(results, result)
	}

	return results, nil
}
