package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "runtime"
    "sync"
    "time"
)

var mu sync.Mutex

var ready = make(chan int)

func ServiceStatus(res http.ResponseWriter, req *http.Request) {
    res.Header().Set("Content-Type", "application/json")
    res.WriteHeader(http.StatusOK)

    type system struct {
        Pid      int `json:"pid"`
        CPUCount int `json:"cpu_count"`
    }

    type stats struct {
        IdleWorkers int `json:"idle_workers"`
        QueuedJobs int `json:"queued_jobs"`
    }

    type payload struct {
        System system `json:"system"`
        Stats  stats  `json:"stats"`
    }

    type info struct {
        Name      string  `json:"name"`
        Version   string  `json:"version"`
        Timestamp string  `json:"timestamp"`
        Status    string  `json:"status"`
        Payload   payload `json:"payload"`
    }

    response := info{
        Name:      NAME,
        Version:   VERSION,
        Timestamp: time.Now().Format(time.RFC3339),
        Status:    STATUS_OK,
        Payload: payload{
            Stats: stats{
                IdleWorkers: idle_workers,
                QueuedJobs: len(queue),
            },
            System: system{
                Pid:      os.Getpid(),
                CPUCount: runtime.NumCPU(),
            },
        },
    }

    json_response, _ := json.MarshalIndent(response, "", "  ")

    res.Write(json_response)
}

func ServiceSchedule(res http.ResponseWriter, req *http.Request) {
    res.Header().Set("Content-Type", "application/json")
    res.WriteHeader(http.StatusOK)

    mu.Lock()

    incoming_job := job{}
    incoming_message := json.NewDecoder(req.Body)
    incoming_message.Decode(&incoming_job)

    queue = append(queue, incoming_job)

    if idle_workers > 0 {
        ready <- 1
    }

    mu.Unlock()
}

func ServiceReceiveBlock(res http.ResponseWriter, req *http.Request) {
    mu.Lock()

    if len(queue) == 0 {
        idle_workers += 1
        mu.Unlock()

        cn, _ := res.(http.CloseNotifier)

        select {
        case <-ready:
            mu.Lock()
            idle_workers -= 1
            mu.Unlock()
        case <-cn.CloseNotify():
            fmt.Println("client hung up")

            mu.Lock()
            idle_workers -= 1
            mu.Unlock()

            return
        }
    } else {
        mu.Unlock()
    }

    // take the next job off the front
    mu.Lock()
    next_job := queue[0]
    queue = queue[1:]
    mu.Unlock()

    res.Header().Set("Content-Type", "text/plain")
    res.WriteHeader(http.StatusOK)
    res.Write([]byte(next_job.Message))
}

func ServiceReceiveNoBlock(res http.ResponseWriter, req *http.Request) {
    mu.Lock()

    if len(queue) == 0 {
        mu.Unlock()

        // nothing in queue means return immediately
        res.Header().Set("Content-Type", "text/plain")
        res.WriteHeader(http.StatusOK)
        res.Write([]byte(`nothing ready`))
        return
    }

    // take the next job off the front
    next_job := queue[0]
    queue = queue[1:]
    mu.Unlock()

    res.Header().Set("Content-Type", "text/plain")
    res.WriteHeader(http.StatusOK)
    res.Write([]byte(next_job.Message))
}
