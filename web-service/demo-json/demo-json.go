package main

import (
	"encoding/json"
	"fmt"
)

type employee struct {
	ID           int
	EmployeeName string
	Tel          string
	Email        string
}

func main() {
	data, _ := json.Marshal(&employee{101, "Phanthakarn Khumphai", "0983405579", "izephanthakarn@hotmail.com"})
	fmt.Println(string(data))
}
