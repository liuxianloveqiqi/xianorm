package xianorm

import (
	"strconv"
)

// Limit 用于设置查询的 LIMIT 条件
// 参数limit可以接收一个或两个int64类型的参数：

func (d *DB) Limit(limit ...int) *DB {
	if len(limit) == 1 {
		// 如果参数只有一个，表示限制查询结果为最多limit条记录
		d.LimitParam = strconv.Itoa(limit[0])
	} else if len(limit) == 2 {
		// 如果参数有两个，表示限制查询结果为从第一个参数开始，查询第二个参数数量的记录
		d.LimitParam = strconv.Itoa(limit[0]) + "," + strconv.Itoa(limit[2])
	} else {
		// 参数个数错误，抛出panic
		panic("参数个数错误")
	}
	return d
}
