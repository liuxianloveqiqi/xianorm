package test

import (
	"fmt"
	"log"
	"testing"
	"xianorm"
)

func TestUpdata(t *testing.T) {
	xdb, err := xianorm.NewDB("root", "root", "43.139.195.17:3301", "orm")
	if err != nil {
		log.Fatal(err)
	}
	//p := Penson{
	//	Name: "xxz",
	//	Age:  15,
	//}
	//r, err := xdb.Table("person").Update(p)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("更新了 %d 行\n", r)

	rowsAffected, err := xdb.Table("person").Where("id = ?", 24).Update("name", "ggo")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("更新了 %d 行\n", rowsAffected)
}
