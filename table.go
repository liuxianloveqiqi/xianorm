package xianorm

// 设置表名
func (d *DB) Table(name string) *DB {
	d.TableName = name

	//重置引擎
	// d.resetSmallormEngine()
	return d
}

// 获取表名
func (d *DB) GetTable() string {
	return d.TableName
}
