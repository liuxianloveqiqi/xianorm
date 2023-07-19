package test

import (
	"fmt"
	"log"
	"testing"
	"xianorm"
)

type Penson struct {
	ID   int `xianorm:auto_increment`
	Name string
	Age  int
}

func TestInsert(t *testing.T) {
	xdb, err := xianorm.NewDB("root", "root", "43.139.195.17:3301", "orm")
	if err != nil {
		log.Fatal(err)
	}
	p := Penson{
		Name: "ggr",
		Age:  23,
	}
	_, err = xdb.Table("person").Insert(p)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(xdb.Prepare)
}
func TestIBatchInsert(t *testing.T) {
	xdb, err := xianorm.NewDB("root", "root", "43.139.195.17:3301", "orm")
	if err != nil {
		log.Fatal(err)
	}
	ps := make([]Penson, 0)
	ps = append(ps, Penson{
		Name: "oop",
		Age:  22,
	})
	ps = append(ps, Penson{
		Name: "xxz",
		Age:  32,
	})
	_, err = xdb.Table("person").Insert(ps)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(xdb.Prepare)
}
