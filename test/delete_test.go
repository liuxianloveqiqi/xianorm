package test

import (
	"fmt"
	"log"
	"testing"
	"xianorm"
)

func TestDelete(t *testing.T) {
	// 初始化数据库连接
	xdb, err := xianorm.NewDB("root", "root", "localhost:3306", "orm")
	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, err := xdb.Table("person").Where("id > ?", 5).Delete()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(xdb.Prepare)

	// 输出受影响的行数
	fmt.Printf("删除了 %d 行\n", rowsAffected)
}
