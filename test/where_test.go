package test

import (
	"fmt"
	"testing"
	"xianorm"
)

func TestWhereAndSelect(t *testing.T) {
	// 创建数据库连接
	xdb, err := xianorm.NewDB("root", "root", "43.139.195.17:3301", "orm")
	if err != nil {
		t.Fatal(err)
	}

	// 测试Where函数
	xdb = xdb.Where("age > ?", 30).Where("name like ?", "%John%")

	// 测试Select函数
	xdb = xdb.Select("id, name, age")

	// 执行查询

	if err != nil {
		t.Fatal(err)
	}

	// 打印查询结果
	fmt.Println("Query Results:", results)
}
