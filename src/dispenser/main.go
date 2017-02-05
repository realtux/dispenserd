package main

import (
    "fmt"
    "net/http"
    "os"
    "path/filepath"
    "sync"
)

var mu = &sync.Mutex{}
var ROOT string

func main() {
    ROOT, _ = filepath.Abs(filepath.Dir(os.Args[0]))

    // init config
    ConfigLoad()

    // handle service endpoints
    http.HandleFunc("/", ServiceStatus)
    http.HandleFunc("/jobs", ServiceJobs)
    http.HandleFunc("/schedule", ServiceSchedule)
    http.HandleFunc("/receive_block", ServiceReceiveBlock)
    http.HandleFunc("/receive_noblock", ServiceReceiveNoBlock)

    fmt.Println("server started on port", 8282)

    http.ListenAndServe(":8282", nil)
}
