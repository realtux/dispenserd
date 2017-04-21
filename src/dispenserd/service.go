package main

import (
    "encoding/json"
    "net/http"
    "os"
    "runtime"
    "time"
)

var listeners = map[string]chan int{"main": make(chan int)}

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
    }

    type payload struct {
        System      system          `json:"system"`
        QueuedJobs  map[string]int  `json:"queued_jobs"`
        IdleWorkers map[string]uint `json:"idle_workers"`
    }

    type info struct {
        Name      string  `json:"name"`
        Version   string  `json:"version"`
        Timestamp string  `json:"timestamp"`
        Status    string  `json:"status"`
        Payload   payload `json:"payload"`
    }

    mu.Lock()

    qj := make(map[string]int)

    for k, v := range queue {
        qj[k] = len(v)
    }

    iw := make(map[string]uint)

    for k, v := range idle_workers {
        iw[k] = v
    }

    mu.Unlock()

    response := info{
        Name:      NAME,
        Version:   VERSION,
        Timestamp: time.Now().Format(time.RFC3339),
        Status:    STATUS_OK,
        Payload: payload{
            QueuedJobs:  qj,
            IdleWorkers: iw,
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

    if incoming_job.Message == nil {
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

    if incoming_job.Message != nil {
        final_job.Message = incoming_job.Message
    }

    if incoming_job.Priority != nil {
        final_job.Priority = incoming_job.Priority
    }

    if incoming_job.Lane != nil && *incoming_job.Lane != "" {
        final_job.Lane = incoming_job.Lane
    }

    InsertJob(final_job)

    payload := generic_payload{
        Status: STATUS_OK,
        Code:   CODE_SUCCESS,
    }

    json_response, _ := json.MarshalIndent(payload, "", "  ")

    if config.ThrottleSchedule {
        time.Sleep(time.Duration(CONFIG_DEFAULT_THROTTLE_MS) * time.Millisecond)
    }

    res.WriteHeader(http.StatusOK)
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

    if config.ThrottleReceive {
        time.Sleep(time.Duration(CONFIG_DEFAULT_THROTTLE_MS) * time.Millisecond)
    }

    var current_lane string

    send_job := func() {
        // take the next job off the front
        next_job := queue[current_lane][0]

        if len(queue[current_lane]) == 1 {
            queue[current_lane] = job_set{}

            // queue empty, good time to reset indexes
            ResetIndexes(*next_job.Lane)
        } else {
            queue[current_lane] = queue[current_lane][1:]

            // queue not empty, decrement
            DecrementIndexes(*next_job.Lane)
        }

        mu.Unlock()

        payload := generic_payload{
            Status:  STATUS_OK,
            Code:    CODE_SUCCESS,
            Message: *next_job.Message,
        }

        json_response, _ := json.MarshalIndent(payload, "", "  ")

        res.Header().Set("Content-Type", "application/json")
        res.WriteHeader(http.StatusOK)
        res.Write(json_response)
    }

    type request struct {
        Lane *string `json:"lane"`
    }

    incoming_request := request{
        Lane: nil,
    }

    incoming_data := json.NewDecoder(req.Body)
    incoming_data.Decode(&incoming_request)

    if incoming_request.Lane == nil || *incoming_request.Lane == "" {
        current_lane = "main"
    } else {
        CreateIndex(*incoming_request.Lane)
        current_lane = *incoming_request.Lane
    }

    if len(queue[current_lane]) == 0 {
        idle_workers[current_lane] += 1

        mu.Unlock()

        cn, _ := res.(http.CloseNotifier)

        for {
            select {
            case <-listeners[current_lane]:
                mu.Lock()

                // empty queue or no workers means do nothing
                if len(queue[current_lane]) == 0 || idle_workers[current_lane] == 0 {
                    mu.Unlock()
                    continue
                }

                idle_workers[current_lane] -= 1

                send_job()

                return
            case <-cn.CloseNotify():
                mu.Lock()
                if idle_workers[current_lane] > 0 {
                    idle_workers[current_lane] -= 1
                }
                mu.Unlock()

                return
            }
        }
    } else {
        send_job()
        return
    }
}

func ServiceReceiveNoBlock(res http.ResponseWriter, req *http.Request) {
    var current_lane string

    type request struct {
        Lane *string `json:"lane"`
    }

    incoming_request := request{
        Lane: nil,
    }

    incoming_data := json.NewDecoder(req.Body)
    incoming_data.Decode(&incoming_request)

    if incoming_request.Lane == nil || *incoming_request.Lane == "" {
        current_lane = "main"
    } else {
        CreateIndex(*incoming_request.Lane)
        current_lane = *incoming_request.Lane
    }

    mu.Lock()

    if len(queue[current_lane]) == 0 {
        mu.Unlock()

        // nothing in queue means return immediately
        payload := generic_payload{
            Status:  STATUS_OK,
            Code:    CODE_NO_DATA_AVAILABLE,
            Message: "empty queue",
        }

        json_response, _ := json.MarshalIndent(payload, "", "  ")

        res.Header().Set("Content-Type", "application/json")
        res.WriteHeader(http.StatusOK)
        res.Write(json_response)
        return
    }

    // take the next job off the front
    next_job := queue[current_lane][0]

    if len(queue[current_lane]) == 1 {
        queue[current_lane] = job_set{}

        // queue empty, good time to reset indexes
        ResetIndexes(*next_job.Lane)
    } else {
        queue[current_lane] = queue[current_lane][1:]

        // queue not empty, decrement
        DecrementIndexes(*next_job.Lane)
    }

    mu.Unlock()

    payload := generic_payload{
        Status:  STATUS_OK,
        Code:    CODE_SUCCESS,
        Message: *next_job.Message,
    }

    json_response, _ := json.MarshalIndent(payload, "", "  ")

    res.Header().Set("Content-Type", "application/json")
    res.WriteHeader(http.StatusOK)
    res.Write(json_response)
}
