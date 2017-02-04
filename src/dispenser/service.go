package main

import (
	"encoding/json"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"
)

var mu sync.Mutex

func ServiceStatus(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	type system struct {
		Pid      int `json:"pid"`
		CPUCount int `json:"cpu_count"`
	}

	type stats struct {
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
}
