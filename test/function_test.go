package test

import (
	"fmt"
	"testing"
	"xianorm"
)

func TestAggregateFunctions(t *testing.T) {
	xdb, err := xianorm.NewDB("root", "root", "43.139.195.17:3301", "orm")
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
