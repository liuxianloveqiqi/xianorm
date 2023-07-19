package xianorm

// Where 构造原生SQL查询条件
func (d *DB) Where(query string, args ...interface{}) *DB {
	// 创建一个新的DB结构体，复制原有的字段和值
	newDB := *d

	// 初始化WhereExec切片
	if newDB.WhereExec == nil {
		newDB.WhereExec = []interface{}{}
	}

	// 添加新的查询条件和参数
	newDB.WhereParam += query + " "
	newDB.WhereExec = append(newDB.WhereExec, args...)

	return &newDB
}

// Select 构造原生SQL查询字段
func (d *DB) Select(fields string) *DB {
	// 创建一个新的DB结构体，复制原有的字段和值
	newDB := *d
	// 设置查询字段
	newDB.FieldParam = fields
	return &newDB
}
