package main

import (
    "fmt"
    "time"

    "github.com/pquerna/otp/totp"
)

func main() {
    key := "JBSWY3DPEHPK3PXP"

    // Generate a mock TOTP code with 6 digits for current time
    totpCode, err := totp.GenerateCodeCustom(key, time.Now(), totp.ValidateOpts{
        Digits:    6,
        Period:    1,
        Skew:      1,

    })

    if err != nil {
        fmt.Println("Error generating TOTP code:", err)
        return
    }

    fmt.Println("Mock TOTP code:", totpCode)
}
