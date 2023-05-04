package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type employee struct {
	ID           int
	EmployeeName string
	Tel          string
	Email        string
}

func main() {
	
	e := employee{}
	err := json.Unmarshal([]byte(`{"ID":101,"EmployeeName": "Phanthakarn Khumphai","Tel" : "0983405579","Email" :"izephanthakarn@hotmail.com"}`), &e)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(e.EmployeeName)
}
