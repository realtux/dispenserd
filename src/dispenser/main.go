package main

import (
    "fmt"
    "net/http"
)

func main() {
    // init config
    ConfigLoad()

    // handle service endpoints
    http.HandleFunc("/", ServiceStatus)
    http.HandleFunc("/schedule", ServiceSchedule)
    http.HandleFunc("/receive_block", ServiceReceiveBlock)
    http.HandleFunc("/receive_noblock", ServiceReceiveNoBlock)

    fmt.Println("server started on port", 8282)

    http.ListenAndServe(":8282", nil)
}
