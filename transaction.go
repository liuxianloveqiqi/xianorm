package xianorm

// Begin 开启事务
func (d *DB) Begin() error {
	// 调用原生的开启事务方法
	tx, err := d.Db.Begin()
	if err != nil {
		return d.setErrorInfo(err)
	}
	d.TransStatus = 1
	d.Tx = tx
	return nil
}

// Rollback 事务回滚
func (d *DB) Rollback() error {
	d.TransStatus = 0
	return d.Tx.Rollback()
}

// Commit 事务提交
func (d *DB) Commit() error {
	d.TransStatus = 0
	return d.Tx.Commit()
}
