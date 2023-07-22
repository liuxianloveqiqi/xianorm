package test

import (
	"fmt"
	"log"
	"testing"
	"xianorm"
)

func TestAggregateFunctions(t *testing.T) {
	xdb, err := xianorm.NewDB("root", "root", "localhost:3306", "orm")
	if err != nil {
		t.Fatal(err)
	}

	// 测试Count函数
	count, err := xdb.Table("person").Count()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Count:", count)

	// 测试Max函数
	max, err := xdb.Table("person").Max("age")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Max Age:", max)

	// 测试Min函数
	min, err := xdb.Table("person").Min("age")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Min Age:", min)

	// 测试Avg函数
	avg, err := xdb.Table("person").Avg("age")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Average Age:", avg)

	// 测试Sum函数
	sum, err := xdb.Table("person").Sum("age")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Sum of Ages:", sum)
}
func TestOrder(t *testing.T) {
	xdb, err := xianorm.NewDB("root", "root", "localhost:3306", "orm")
	if err != nil {
		t.Fatal(err)
	}
	p := make([]Penson, 0)
	// 测试Order函数
	err = xdb.Table("person").Order("id", "desc", "age", "desc").Find(&p)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Order Param:", xdb.OrderParam)
	fmt.Println("perons：", p)
}

func TestGroup(t *testing.T) {
	xdb, err := xianorm.NewDB("root", "root", "localhost:3306", "orm")
	if err != nil {
		t.Fatal(err)
	}

	// 测试Group函数
	maps, err := xdb.Table("person").
		Select("CASE WHEN gender = 0 THEN '男' WHEN gender = 1 THEN '女' ELSE '未知' END as 'gender', COUNT(*) as 'count'").
		Group("gender").
		Query()

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("person的gender :", maps)
}

func TestHaving(t *testing.T) {
	xdb, err := xianorm.NewDB("root", "root", "localhost:3306", "orm")
	if err != nil {
		t.Fatal(err)
	}

	// 测试Having函数
	xdb.Having("age > ?", 30)
	fmt.Println("Having Param:", xdb.HavingParam)
	fmt.Println("Having Exec:", xdb.HavingExec)
}
