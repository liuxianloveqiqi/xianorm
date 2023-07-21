package test

import (
	"log"
	"testing"
	"xianorm"
)

func TestTransaction(t *testing.T) {

	xdb, err := xianorm.NewDB("root", "root", "localhost:3306", "orm")
	if err != nil {
		log.Fatal(err)
	}
	// 执行事务操作
	err = xdb.Begin()
	if err != nil {
		t.Fatalf("开启事务失败: %v", err)
	}

	p := Penson{
		ID:   82,
		Name: "yyh",
		Age:  34,
	}
	_, err = xdb.Table("person").Insert(p)
	if err != nil {
		t.Fatalf("insert失败: %v", err)
	}

	err = xdb.Commit()
	if err != nil {
		t.Fatalf("提交事务失败: %v", err)
	}

}
