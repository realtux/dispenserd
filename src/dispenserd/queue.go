package main

import (
    "crypto/sha1"
    "fmt"
    "io"
    "math/rand"
    //"sort"
    "strconv"
    "time"
)

type job struct {
    JobNum    uint64  `json:"job_num"`
    Hash      string  `json:"hash"`
    Timestamp string  `json:"timestamp"`
    Lane      *string `json:"lane"`
    Priority  *uint   `json:"priority"`
    Message   *string `json:"message"`
}

type job_set []job

var queue = []job_set{}
var lanes = []string{"main"}

var idle_workers = 0
var total_jobs uint64 = 0
var current_jobs uint64 = 0

func (a job_set) Len() int      { return len(a) }
func (a job_set) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a job_set) Less(i, j int) bool {
    if *a[i].Priority < *a[j].Priority {
        return true
    }
    if *a[i].Priority > *a[j].Priority {
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

    var default_lane string = JOB_DEFAULT_LANE
    var default_priority uint = JOB_DEFAULT_PRIORITY

    return job{
        JobNum:    job_num,
        Hash:      fmt.Sprintf("%x", hash.Sum(nil)),
        Timestamp: timestamp,
        Lane:      &default_lane,
        Priority:  &default_priority,
        Message:   nil,
    }
}

func InitJobTemplate() job {
    return job{
        JobNum:    0,
        Hash:      "",
        Timestamp: "",
        Lane:      nil,
        Priority:  nil,
        Message:   nil,
    }
}

func InsertJob(job job) {
    mu.Lock()

    lane_index := LaneIndex(*job.Lane)
fmt.Println(lane_index)
    queue[lane_index] = append(queue[lane_index], job)

    current_jobs += 1

    // this kills perf, consider replacing with periodic sort
    //SortQueue()

    if idle_workers > 0 {
        ready <- 1
    }

    mu.Unlock()
}

func SortQueue() {
    //mu.Lock()
    //sort.Sort(job_set(queue))
    //mu.Unlock()
}

func LaneIndex(search_lane string) int {
    index := -1
fmt.Println("searching for lane", search_lane)
    if search_lane == "" || search_lane == "main" {
fmt.Println("found main lane")
        return 0
    }

    for i, lane := range lanes {
        if search_lane == lane {
            index = i
fmt.Println("found lane at", i)
            break
        }
    }

    if index == -1 {
fmt.Println("creating new lane", search_lane, "at", len(lanes))
        lanes = append(lanes, search_lane)
        queue = append(queue, job_set{})

        new_index := len(lanes) - 1

        return new_index
    } else {
        return index
    }
}
