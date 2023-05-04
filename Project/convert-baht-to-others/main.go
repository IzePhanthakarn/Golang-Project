package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

type exchangeRates struct {
    Rates map[string]float64 `json:"rates"`
}

func main() {
    var thb float64
    var targetCurrency string

    // รับค่าจำนวนเงินในหน่วยบาทและสกุลเงินปลายทางจากผู้ใช้งาน
    fmt.Print("Enter amount in THB: ")
    fmt.Scanln(&thb)
    fmt.Print("Enter target currency (USD, JPY, CNY): ")
    fmt.Scanln(&targetCurrency)

    // ส่ง request ไปยัง API เพื่อขอ exchange rates ล่าสุด
    res, err := http.Get("https://api.exchangerate-api.com/v4/latest/" + targetCurrency)
    if err != nil {
        fmt.Println("Error while fetching exchange rates:", err)
        return
    }
    defer res.Body.Close()

    // อ่าน response body และแปลงเป็น struct
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        fmt.Println("Error while reading response body:", err)
        return
    }
    var exchange exchangeRates
    err = json.Unmarshal(body, &exchange)
    if err != nil {
        fmt.Println("Error while unmarshalling response body:", err)
        return
    }

    // คำนวณค่าสกุลเงินตามสกุลเงินปลายทาง
    targetRate, ok := exchange.Rates["THB"]
    if !ok {
        fmt.Println("Error: target currency not found.")
        return
    }
    targetAmount := thb / targetRate

    // แสดงผลลัพธ์
    fmt.Printf("%.2f THB = %.2f %s\n", thb, targetAmount, targetCurrency)
}

// >go build