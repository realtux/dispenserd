package main

import (
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "path/filepath"
    "sync"
    "syscall"
)

var mu = &sync.Mutex{}
var ROOT string

func cleanup() {
    fmt.Println("cleaning up...")
    if config.PersistQueue {
        fmt.Println("queue persistence enabled, saving queue...")
    }
}

func main() {
    ROOT, _ = filepath.Abs(filepath.Dir(os.Args[0]))

    // init config
    ConfigLoad()

    // handle sigterm
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        cleanup()
        os.Exit(1)
    }()

    // handle service endpoints
    http.HandleFunc("/", ServiceStatus)
    http.HandleFunc("/jobs", ServiceJobs)
    http.HandleFunc("/schedule", ServiceSchedule)
    http.HandleFunc("/receive_block", ServiceReceiveBlock)
    http.HandleFunc("/receive_noblock", ServiceReceiveNoBlock)

    fmt.Println("server started on port", 8282)

    http.ListenAndServe(":8282", nil)
}
