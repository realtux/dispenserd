package main

import (
    "crypto/sha1"
    "fmt"
    "io"
    "math/rand"
    "sort"
    "strconv"
    "time"
)

type job struct {
    JobNum    uint64  `json:"job_num"`
    Hash      string  `json:"hash"`
    Timestamp string  `json:"timestamp"`
    Priority  uint    `json:"priority"`
    Message   *string `json:"message"`
}

type job_set []job

var queue = job_set{}
var idle_workers = 0
var total_jobs uint64 = 0

func (a job_set) Len() int      { return len(a) }
func (a job_set) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a job_set) Less(i, j int) bool {
    if a[i].Priority < a[j].Priority {
        return true
    }
    if a[i].Priority > a[j].Priority {
        return false
    }
    return a[i].JobNum < a[j].JobNum
}

func InitJob() job {
    mu.Lock()
    total_jobs += 1
    var job_num uint64 = total_jobs
    mu.Unlock()

    rand_number := rand.Intn(999999999)
    timestamp := time.Now().Format(time.RFC3339)

    hash := sha1.New()
    io.WriteString(hash, strconv.Itoa(rand_number)+timestamp)

    return job{
        JobNum:    job_num,
        Hash:      fmt.Sprintf("%x", hash.Sum(nil)),
        Timestamp: timestamp,
        Priority:  JOB_DEFAULT_PRIORITY,
        Message:   nil,
    }
}

func InitJobTemplate() job {
    return job{
        JobNum:    0,
        Hash:      "",
        Timestamp: "",
        Priority:  0,
        Message:   nil,
    }
}

func InsertJob(job job) {
    mu.Lock()

    queue = append(queue, job)

    // this kills perf, consider replacing with periodic sort
    SortQueue()

    if idle_workers > 0 {
        ready <- 1
    }

    mu.Unlock()
}

func SortQueue() {
    //mu.Lock()
    sort.Sort(job_set(queue))
    //mu.Unlock()
}
