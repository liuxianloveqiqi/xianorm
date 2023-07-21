package test

import (
	"fmt"
	"log"
	"testing"
	"xianorm"
)

func TestQuery(t *testing.T) {
	xdb, err := xianorm.NewDB("root", "root", "localhost:3306", "orm")
	if err != nil {
		log.Fatal(err)
	}
	ms, err := xdb.Table("person").Query()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ms)

	m, err := xdb.Table("person").QueryOne()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(m)
}
func TestFind(t *testing.T) {
	xdb, err := xianorm.NewDB("root", "root", "localhost:3306", "orm")
	if err != nil {
		log.Fatal(err)
	}
	p := make([]Penson, 0)
	err = xdb.Table("person").Find(&p)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(p)
}
func TestFirst(t *testing.T) {
	xdb, err := xianorm.NewDB("root", "root", "localhost:3306", "orm")
	if err != nil {
		log.Fatal(err)
	}
	p := Penson{}
	err = xdb.Table("person").First(&p)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(p)
}
