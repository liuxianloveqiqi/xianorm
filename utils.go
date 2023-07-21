package xianorm

import (
	"fmt"
	"runtime"
)

func (d *DB) BuildQuerySql() {
	// 拼接SQL语句
	d.Sql = "SELECT " + d.FieldParam + " FROM " + d.GetTable()

	// 如果 WhereParam 或 OrWhereParam 不为空，添加查询条件
	if d.WhereParam != "" || d.OrWhereParam != "" {
		d.Sql += " WHERE " + d.WhereParam + d.OrWhereParam
	}

	// 如果 GroupParam 不为空，添加 GROUP BY 条件
	if d.GroupParam != "" {
		d.Sql += " GROUP BY " + d.GroupParam
	}

	// 如果 HavingParam 不为空，添加 HAVING 条件
	if d.HavingParam != "" {
		d.Sql += " HAVING " + d.HavingParam
	}

	// 如果 OrderParam 不为空，添加排序条件
	if d.OrderParam != "" {
		d.Sql += " ORDER BY " + d.OrderParam
	}

	// 如果 LimitParam 不为空，添加限制条件
	if d.LimitParam != "" {
		d.Sql += " LIMIT " + d.LimitParam
	}

}

// 自定义错误格式
func (d *DB) setErrorInfo(err error) error {
	// 用于获取当前调用栈的信息
	_, file, line, _ := runtime.Caller(1)
	// 返回文件名和行号和错误信息
	return fmt.Errorf("file: %s:%d, %w", file, line, err)
}
