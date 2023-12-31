package xianorm

import (
	"fmt"
	"strconv"
	"strings"
)

// Limit 用于设置查询的 LIMIT 条件
// 参数limit可以接收一个或两个int类型的参数：

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

// aggregateQuery 聚合查询函数，根据传入的聚合函数名称和参数进行对应的操作
func (d *DB) aggregateQuery(functionName, param string) (interface{}, error) {
	// 拼接SQL语句，使用传入的聚合函数和参数
	d.Prepare = fmt.Sprintf("SELECT %s(%s) AS result FROM %s", functionName, param, d.GetTable())

	// 如果WhereParam或OrWhereParam不为空，添加查询条件
	if d.WhereParam != "" {
		d.Prepare += " WHERE " + d.WhereParam
	}

	// 如果LimitParam不为空，添加限制条件
	if d.LimitParam != "" {
		d.Prepare += " LIMIT " + d.LimitParam
	}

	// 执行查询
	var result interface{}
	err := d.Db.QueryRow(d.Prepare, d.WhereExec...).Scan(&result)
	if err != nil {
		return nil, d.setErrorInfo(err)
	}

	return result, nil
}

// Count 统计数量
func (d *DB) Count() (int64, error) {
	count, err := d.aggregateQuery("count", "*")
	if err != nil {
		return 0, d.setErrorInfo(err)
	}
	c, err := strconv.ParseInt(string(count.([]byte)), 10, 64)
	if err != nil {
		return 0, d.setErrorInfo(err)
	}
	return c, err
}

// Max 最大值
func (d *DB) Max(param string) (string, error) {
	max, err := d.aggregateQuery("max", param)
	if err != nil {
		return "0", d.setErrorInfo(err)
	}
	return string(max.([]byte)), nil
}

// Min 最小值
func (d *DB) Min(param string) (string, error) {
	min, err := d.aggregateQuery("min", param)
	if err != nil {
		return "0", d.setErrorInfo(err)
	}

	return string(min.([]byte)), nil
}

// Avg 平均值
func (d *DB) Avg(param string) (float64, error) {
	avg, err := d.aggregateQuery("avg", param)
	if err != nil {
		return 0, d.setErrorInfo(err)
	}
	// 将[]byte转换为float64
	f, err := strconv.ParseFloat(string(avg.([]byte)), 64)
	return f, nil
}

// Sum 总和
func (d *DB) Sum(param string) (string, error) {
	sum, err := d.aggregateQuery("sum", param)
	if err != nil {
		return "0", d.setErrorInfo(err)
	}
	return string(sum.([]byte)), nil
}

// Order 排序
func (d *DB) Order(order ...string) *DB {
	orderLen := len(order)
	if orderLen%2 != 0 {
		panic("order by参数错误，请保证参数个数为偶数个")
	}

	// 排序的个数
	orderNum := orderLen / 2

	// 多次调用的情况下，如果已经有排序参数则在末尾添加逗号
	if d.OrderParam != "" {
		d.OrderParam += ","
	}

	// 用于存储拼接后的排序参数列表
	var orderParams []string

	// 循环处理每个排序字段和排序方式
	for i := 0; i < orderNum; i++ {
		keyString := strings.ToLower(order[i*2+1])
		if keyString != "desc" && keyString != "asc" {
			panic("排序关键字只能为：desc和asc")
		}

		// 将排序字段和排序方式拼接成一个字符串，并添加到列表中
		orderParams = append(orderParams, order[i*2]+" "+order[i*2+1])
	}

	// 将排序参数列表连接成一个字符串，以逗号分隔
	d.OrderParam += strings.Join(orderParams, ",")

	return d
}

// Group group分组
func (d *DB) Group(group ...string) *DB {
	if len(group) != 0 {
		d.GroupParam = strings.Join(group, ",")
	}
	return d
}

// Having Having过滤，跟where的实现一样
func (d *DB) Having(having string, args ...interface{}) *DB {
	// 创建一个新的DB结构体，复制原有的字段和值

	// 初始化HavingExec切片
	if d.HavingExec == nil {
		d.HavingExec = []interface{}{}
	}

	// 添加新的过滤条件和参数
	d.HavingParam += having + " "
	d.HavingExec = append(d.HavingExec, args...)

	return d
}
