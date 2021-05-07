package main

import (
	"fmt"
	"log"

	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/exam_marker/v2/utils"
)

func main() {
	fmt.Println("vim-go")
	db := utils.InitDB()
	log.Println(db)
}
