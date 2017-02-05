package main

import (
    "fmt"
    "time"
)

func Persist() {
    for {
        fmt.Println("replace with sync code")
        time.Sleep(time.Duration(config.PersistInterval) * time.Second)
    }
}
