package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "runtime"
    "time"
)

var ready = make(chan int)

type generic_payload struct {
    Status  string `json:"status"`
    Code    int    `json:"code"`
    Message string `json:"message,omitempty"`
}

func ServiceStatus(res http.ResponseWriter, req *http.Request) {
    res.Header().Set("Content-Type", "application/json")
    res.WriteHeader(http.StatusOK)

    type system struct {
        Pid      int `json:"pid"`
        CPUCount int `json:"cpu_count"`
    }

    type stats struct {
        IdleWorkers int `json:"idle_workers"`
        QueuedJobs  int `json:"queued_jobs"`
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
                QueuedJobs:  len(queue),
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

    incoming_job := InitJobTemplate()
    incoming_message := json.NewDecoder(req.Body)
    incoming_message.Decode(&incoming_job)

    if incoming_job.Message == "" {
        payload := generic_payload{
            Status:  STATUS_ERROR,
            Code:    CODE_INVALID_DATA,
            Message: "missing message",
        }

        json_response, _ := json.MarshalIndent(payload, "", "  ")

        res.WriteHeader(http.StatusBadRequest)
        res.Write(json_response)
        return
    }

    final_job := InitJob()
    final_job.Message = incoming_job.Message
    final_job.Priority = incoming_job.Priority

    InsertJob(final_job)

    payload := generic_payload{
        Status: STATUS_OK,
        Code:   CODE_SUCCESS,
    }

    json_response, _ := json.MarshalIndent(payload, "", "  ")

    res.WriteHeader(http.StatusBadRequest)
    res.Write(json_response)
}

func ServiceJobs(res http.ResponseWriter, req *http.Request) {
    res.Header().Set("Content-Type", "application/json")
    res.WriteHeader(http.StatusOK)

    json_response, _ := json.MarshalIndent(queue, "", "  ")

    res.Write(json_response)
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
