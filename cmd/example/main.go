package main

import (
	"fmt"
	"log"

	"github.com/ClavinJune/bloomfilter"
)

func main() {
	bf, err := bloomfilter.New(100, 0.001)

	if err != nil {
		log.Fatal(err)
	}

	bf.Add("aku mau makan")
	bf.Add("aku mau mandi")

	fmt.Println(bf.Check("aku mau makan"))
	fmt.Println(bf.Check("aku mau"))
}
