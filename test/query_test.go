package test

import (
	"fmt"
	"log"
	"testing"
	"xianorm"
)

func TestQuery(t *testing.T) {
	xdb, err := xianorm.NewDB("root", "root", "43.139.195.17:3301", "orm")
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
	xdb, err := xianorm.NewDB("root", "root", "43.139.195.17:3301", "orm")
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
