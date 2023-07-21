package xianorm

import "database/sql"

// Delete 执行删除操作
func (d *DB) Delete() (int64, error) {
	// 构建删除SQL语句
	d.Prepare = "delete from " + d.GetTable()

	// 如果有Where条件，将Where条件拼接到SQL语句中
	if d.WhereParam != "" {
		d.Prepare += " where " + d.WhereParam
	}

	// 如果有限制条件，将限制条件拼接到SQL语句中
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
	// 执行预处理语句
	result, err := stmt.Exec(d.WhereExec...)
	if err != nil {
		return 0, d.setErrorInfo(err)
	}

	// 获取受影响的行数
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, d.setErrorInfo(err)
	}

	return rowsAffected, nil
}
