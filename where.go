package xianorm

// Where 构造原生SQL查询条件
func (d *DB) Where(query string, args ...interface{}) *DB {

	// 初始化WhereExec切片
	if d.WhereExec == nil {
		d.WhereExec = []interface{}{}
	}

	// 添加新的查询条件和参数
	d.WhereParam += query + " "
	d.WhereExec = append(d.WhereExec, args...)

	return d
}

// Select 构造原生SQL查询字段
func (d *DB) Select(fields string) *DB {
	// 设置查询字段
	d.FieldParam = fields
	return d
}
