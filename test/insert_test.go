package test

import (
	"fmt"
	"log"
	"testing"
	"xianorm"
)

type Penson struct {
	ID     int    `xianorm:"id"`
	Name   string `xianorm:"name"`
	Age    int
	Gender int
}

func TestInsert(t *testing.T) {
	xdb, err := xianorm.NewDB("root", "root", "localhost:3306", "orm")
	if err != nil {
		log.Fatal(err)
	}
	p := Penson{
		ID:   99,
		Name: "jjj",
		Age:  56,
	}
	_, err = xdb.Table("person").Insert(p)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(xdb.Prepare)
}
func TestIBatchInsert(t *testing.T) {
	xdb, err := xianorm.NewDB("root", "root", "localhost:3306", "orm")
	if err != nil {
		log.Fatal(err)
	}
	ps := make([]Penson, 0)
	ps = append(ps, Penson{
		ID:   289,
		Name: "ver",
		Age:  43,
	})
	ps = append(ps, Penson{
		ID:   70,
		Name: "x532",
		Age:  42,
	})
	_, err = xdb.Table("person").Insert(ps)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(xdb.Prepare)
}
