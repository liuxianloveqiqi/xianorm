package xianorm

// Table 设置表名
func (d *DB) Table(name string) *DB {
	d.TableName = name

	//重置DB
	d.resetDB()
	return d
}

// 重置DB
func (d *DB) resetDB() {
	d.Prepare = ""
	d.AllExec = []interface{}{}
	d.Sql = ""
	d.WhereParam = ""
	d.LimitParam = ""
	d.OrderParam = ""
	d.OrWhereParam = ""
	d.WhereExec = []interface{}{}
	d.UpdateParam = ""
	d.UpdateExec = []interface{}{}
	d.FieldParam = "*"
	d.TransStatus = 0
	d.Tx = nil
	d.GroupParam = ""
	d.HavingParam = ""
}

// GetTable 获取表名
func (d *DB) GetTable() string {
	return d.TableName
}
